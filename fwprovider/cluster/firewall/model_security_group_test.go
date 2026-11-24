/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"encoding/json"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

// replayRuleOps applies ops with PVE's list semantics (pve-firewall src/PVE/API2/Firewall/Rules.pm): create unshifts,
// moveto splices the rule in before the original index moveto, and delete removes the position.
func replayRuleOps(t *testing.T, rules []rule, ops []ruleOp) []rule {
	t.Helper()

	rules = slices.Clone(rules)

	for _, op := range ops {
		switch op.kind {
		case ruleOpCreate:
			rules = slices.Insert(rules, 0, op.rule)
		case ruleOpMove:
			require.Less(t, op.pos, len(rules), "move from a missing position")

			moved := rules[op.pos]
			next := make([]rule, 0, len(rules))

			for i, r := range rules {
				if i == op.pos {
					continue
				}

				if i == op.moveTo {
					next = append(next, moved)
				}

				next = append(next, r)
			}

			if op.moveTo >= len(rules) {
				next = append(next, moved)
			}

			rules = next
		case ruleOpUpdate:
			require.Less(t, op.pos, len(rules), "update of a missing position")
			rules[op.pos] = op.rule
		case ruleOpDelete:
			require.Less(t, op.pos, len(rules), "delete of a missing position")
			rules = slices.Delete(rules, op.pos, op.pos+1)
		}
	}

	return rules
}

func countOps(ops []ruleOp, kind ruleOpKind) int {
	n := 0

	for _, op := range ops {
		if op.kind == kind {
			n++
		}
	}

	return n
}

func TestPlanRuleOps(t *testing.T) {
	t.Parallel()

	ssh := rule{Type: "in", Action: "ACCEPT", Macro: "SSH", Enabled: true}
	http := rule{Type: "in", Action: "ACCEPT", Proto: "tcp", DPort: "80", Enabled: true}
	https := rule{Type: "in", Action: "ACCEPT", Proto: "tcp", DPort: "443", Enabled: true}
	drop := rule{Type: "in", Action: "DROP", Log: "info", Enabled: true}

	withComment := func(r rule, c string) rule {
		r.Comment = c
		return r
	}

	tests := []struct {
		name    string
		current []rule
		desired []rule
		creates int
		moves   int
		updates int
		deletes int
	}{
		{name: "empty to empty"},
		{name: "create into an empty group", desired: []rule{ssh, http, https}, creates: 3},
		{name: "delete everything", current: []rule{ssh, http, https}, deletes: 3},
		{name: "unchanged", current: []rule{ssh, http}, desired: []rule{ssh, http}},
		{name: "reverse order", current: []rule{ssh, http, https}, desired: []rule{https, http, ssh}, moves: 2},
		{name: "move last to top", current: []rule{ssh, http, https}, desired: []rule{https, ssh, http}, moves: 1},
		{name: "move top to bottom", current: []rule{ssh, http, https}, desired: []rule{http, https, ssh}, moves: 2},
		{name: "insert in the middle", current: []rule{ssh, https}, desired: []rule{ssh, http, https}, creates: 1, moves: 1},
		{name: "append", current: []rule{ssh}, desired: []rule{ssh, drop}, creates: 1, moves: 1},
		{name: "remove from the middle", current: []rule{ssh, http, https}, desired: []rule{ssh, https}, moves: 1, deletes: 1},
		{
			name:    "comment change updates in place",
			current: []rule{ssh, http},
			desired: []rule{ssh, withComment(http, "web")},
			updates: 1,
		},
		{
			name:    "identity change replaces the rule",
			current: []rule{ssh, http},
			desired: []rule{ssh, https},
			creates: 1, moves: 1, deletes: 1,
		},
		{
			name:    "duplicate signatures keep their own attributes",
			current: []rule{withComment(http, "a"), withComment(http, "b")},
			desired: []rule{withComment(http, "b"), withComment(http, "a")},
			moves:   1,
		},
		{
			name:    "macro case differences are not a change",
			current: []rule{{Type: "in", Action: "ACCEPT", Macro: "SSH", Enabled: true}},
			desired: []rule{{Type: "in", Action: "ACCEPT", Macro: "ssh", Enabled: true}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ops := planRuleOps(tt.current, tt.desired)
			got := replayRuleOps(t, tt.current, ops)

			require.Len(t, got, len(tt.desired))

			for i := range tt.desired {
				require.Truef(t, got[i].equal(tt.desired[i]), "position %d: got %+v, want %+v", i, got[i], tt.desired[i])
			}

			require.Equal(t, tt.creates, countOps(ops, ruleOpCreate), "creates")
			require.Equal(t, tt.moves, countOps(ops, ruleOpMove), "moves")
			require.Equal(t, tt.updates, countOps(ops, ruleOpUpdate), "updates")
			require.Equal(t, tt.deletes, countOps(ops, ruleOpDelete), "deletes")

			for _, op := range ops {
				if op.kind == ruleOpMove {
					require.Less(t, op.moveTo, op.pos, "only upward moves land exactly on moveto")
				}
			}
		})
	}
}

func TestRuleUpdateBodyDeletesClearedFields(t *testing.T) {
	t.Parallel()

	current := rule{Type: "in", Action: "ACCEPT", Comment: "old", Log: "info", Proto: "tcp", DPort: "22", Enabled: true}
	want := rule{Type: "in", Action: "ACCEPT", Proto: "tcp", DPort: "22", Enabled: false}

	body := want.updateBody(current)

	require.Equal(t, []string{"comment", "log"}, body.Delete)
	require.Nil(t, body.Comment)
	require.Equal(t, "tcp", *body.Proto)
	require.NotNil(t, body.Enable)
	require.False(t, bool(*body.Enable))
}

func TestMoveSecurityGroupFromSDK(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	r := &securityGroupResource{}

	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	require.False(t, schemaResp.Diagnostics.HasError())

	sdkState, err := json.Marshal(map[string]any{
		"id":           "webserver",
		"name":         "webserver",
		"comment":      "",
		"node_name":    nil,
		"vm_id":        nil,
		"container_id": nil,
		"rule": []map[string]any{{
			"pos": 0, "security_group": "", "action": "ACCEPT", "type": "in", "comment": "Allow HTTP",
			"dest": "", "dport": "80", "enabled": true, "iface": "", "log": "", "macro": "", "proto": "tcp",
			"source": "", "sport": "",
		}},
	})
	require.NoError(t, err)

	resp := &resource.MoveStateResponse{
		TargetState: tfsdk.State{
			Schema: schemaResp.Schema,
			Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil),
		},
	}

	moveSecurityGroupFromSDK(ctx, resource.MoveStateRequest{
		SourceTypeName: sdkSecurityGroupTypeName,
		SourceRawState: &tfprotov6.RawState{JSON: sdkState},
	}, resp)
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)

	var got securityGroupModel
	require.False(t, resp.TargetState.Get(ctx, &got).HasError())

	require.Equal(t, "webserver", got.ID.ValueString())
	require.Equal(t, "webserver", got.Name.ValueString())
	require.True(t, got.Comment.IsNull())

	diags := schemaResp.Diagnostics

	rules := got.ruleModels(ctx, &diags)
	require.False(t, diags.HasError())
	require.Len(t, rules, 1)
	require.Equal(t, int64(0), rules[0].Pos.ValueInt64())
	require.Equal(t, "Allow HTTP", rules[0].Comment.ValueString())
	require.Equal(t, "80", rules[0].DPort.ValueString())
	require.True(t, rules[0].Dest.IsNull())
	require.True(t, rules[0].Enabled.ValueBool())
}

func TestMoveSecurityGroupIgnoresOtherTypes(t *testing.T) {
	t.Parallel()

	resp := &resource.MoveStateResponse{TargetState: tfsdk.State{Schema: schema.Schema{}}}

	moveSecurityGroupFromSDK(context.Background(), resource.MoveStateRequest{
		SourceTypeName: "proxmox_virtual_environment_firewall_rules",
	}, resp)

	require.False(t, resp.Diagnostics.HasError())
	require.Nil(t, resp.TargetState.Raw.Type())
}
