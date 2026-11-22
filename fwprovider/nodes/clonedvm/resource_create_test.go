/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package clonedvm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// TestCreateKeepsClonedVMInStateWhenALaterStepFails checks that a VM that was cloned stays in the state when a later
// create step fails, so Terraform taints it instead of losing track of it.
func TestCreateKeepsClonedVMInStateWhenALaterStepFails(t *testing.T) {
	t.Parallel()

	const upid = "UPID:pve:0000ABCD:0000ABCD:6543210A:qmclone:100:root@pam:"

	write := func(w http.ResponseWriter, body string) {
		if _, err := w.Write([]byte(body)); err != nil {
			t.Errorf("writing response: %v", err)
		}
	}

	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/nodes/pve/qemu/100/clone"):
			write(w, `{"data":"`+upid+`"}`)
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/status") &&
			strings.Contains(r.URL.Path, "/nodes/pve/tasks/"):
			write(w, `{"data":{"status":"stopped","exitstatus":"OK"}}`)
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/nodes/pve/qemu/101/config"):
			w.WriteHeader(http.StatusInternalServerError)
			write(w, `{"data":null,"message":"config read failed"}`)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotImplemented)
		}
	}))
	t.Cleanup(srv.Close)

	conn, err := api.NewConnection(srv.URL, true, "", nil)
	require.NoError(t, err)

	creds, err := api.NewCredentials("", "", "", "root@pam!test=00000000-0000-0000-0000-000000000000", "", "")
	require.NoError(t, err)

	apiClient, err := api.NewClient(creds, conn)
	require.NoError(t, err)

	r := &Resource{client: proxmox.NewClient(apiClient, nil, "")}
	ctx := context.Background()

	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	require.False(t, schemaResp.Diagnostics.HasError(), schemaResp.Diagnostics)

	s := schemaResp.Schema
	nullObject := tftypes.NewValue(s.Type().TerraformType(ctx), nil)

	plan := tfsdk.Plan{Schema: s, Raw: nullObject}
	for p, v := range map[string]any{
		"id":                                   types.Int64Value(101),
		"node_name":                            types.StringValue("pve"),
		"started":                              types.BoolValue(false),
		"stop_on_destroy":                      types.BoolValue(true),
		"purge_on_destroy":                     types.BoolValue(true),
		"delete_unreferenced_disks_on_destroy": types.BoolValue(false),
		"clone": CloneModel{
			SourceVMID:      types.Int64Value(100),
			SourceNodeName:  types.StringNull(),
			Full:            types.BoolValue(true),
			TargetDatastore: types.StringNull(),
			TargetFormat:    types.StringNull(),
			SnapshotName:    types.StringNull(),
			PoolID:          types.StringNull(),
			Retries:         types.Int64Value(1),
			BandwidthLimit:  types.Int64Null(),
		},
	} {
		diags := plan.SetAttribute(ctx, path.Root(p), v)
		require.False(t, diags.HasError(), diags)
	}

	resp := resource.CreateResponse{State: tfsdk.State{Schema: s, Raw: nullObject}}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)

	require.True(t, resp.Diagnostics.HasError(), "create should fail on the config read")
	require.False(t, resp.State.Raw.IsNull(), "the cloned VM must stay in the state")

	var state Model

	diags := resp.State.Get(ctx, &state)
	require.False(t, diags.HasError(), diags)
	require.Equal(t, int64(101), state.ID.ValueInt64())
	require.Equal(t, "pve", state.NodeName.ValueString())
	require.True(t, state.StopOnDestroy.ValueBool())
	require.True(t, state.PurgeOnDestroy.ValueBool())
	require.Equal(t, int64(100), state.Clone.SourceVMID.ValueInt64())
}
