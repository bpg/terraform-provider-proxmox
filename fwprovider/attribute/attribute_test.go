/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package attribute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
)

// fakeBody is a minimal DeleteAppender used by CheckDeleteBody tests.
type fakeBody struct {
	deletes []string
}

func (b *fakeBody) AppendDelete(apiName string) {
	b.deletes = append(b.deletes, apiName)
}

func newStringSet(t *testing.T, elements []string) stringset.Value {
	t.Helper()

	var diags diag.Diagnostics

	v := stringset.NewValueList(elements, &diags)
	require.False(t, diags.HasError(), "stringset.NewValueList diags: %v", diags)

	return v
}

func TestCheckDelete(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		plan       attr.Value
		state      attr.Value
		apiName    string
		wantDelete []string
	}{
		{
			name:       "null plan, known state adds delete",
			plan:       types.StringNull(),
			state:      types.StringValue("foo"),
			apiName:    "alias",
			wantDelete: []string{"alias"},
		},
		{
			name:       "null plan, null state does nothing",
			plan:       types.StringNull(),
			state:      types.StringNull(),
			apiName:    "alias",
			wantDelete: nil,
		},
		{
			name:       "known plan, known state does nothing (population handled elsewhere)",
			plan:       types.StringValue("bar"),
			state:      types.StringValue("foo"),
			apiName:    "alias",
			wantDelete: nil,
		},
		{
			name:       "unknown plan, known state does nothing (unknown != deletion)",
			plan:       types.StringUnknown(),
			state:      types.StringValue("foo"),
			apiName:    "alias",
			wantDelete: nil,
		},
		{
			name:       "empty stringset plan, non-empty state adds delete",
			plan:       newStringSet(t, nil),
			state:      newStringSet(t, []string{"a", "b"}),
			apiName:    "tags",
			wantDelete: []string{"tags"},
		},
		{
			name:       "non-empty stringset plan, non-empty state does nothing",
			plan:       newStringSet(t, []string{"c"}),
			state:      newStringSet(t, []string{"a", "b"}),
			apiName:    "tags",
			wantDelete: nil,
		},
		{
			name:       "empty list plan, non-empty state adds delete",
			plan:       types.ListValueMust(types.StringType, nil),
			state:      types.ListValueMust(types.StringType, []attr.Value{types.StringValue("x")}),
			apiName:    "vmid",
			wantDelete: []string{"vmid"},
		},
		{
			name:       "empty map plan, non-empty state adds delete",
			plan:       types.MapValueMust(types.StringType, nil),
			state:      types.MapValueMust(types.StringType, map[string]attr.Value{"k": types.StringValue("v")}),
			apiName:    "ide0",
			wantDelete: []string{"ide0"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var toDelete []string

			attribute.CheckDelete(tc.plan, tc.state, &toDelete, tc.apiName)
			assert.Equal(t, tc.wantDelete, toDelete, "slice form")

			body := &fakeBody{}
			attribute.CheckDeleteBody(tc.plan, tc.state, body, tc.apiName)
			assert.Equal(t, tc.wantDelete, body.deletes, "body form")
		})
	}
}

func TestCheckDeleteBody_AppendsInOrder(t *testing.T) {
	t.Parallel()

	body := &fakeBody{}

	attribute.CheckDeleteBody(types.StringNull(), types.StringValue("a"), body, "first")
	attribute.CheckDeleteBody(types.StringValue("kept"), types.StringValue("kept"), body, "second") // no-op
	attribute.CheckDeleteBody(types.StringNull(), types.StringValue("b"), body, "third")

	assert.Equal(t, []string{"first", "third"}, body.deletes)
}
