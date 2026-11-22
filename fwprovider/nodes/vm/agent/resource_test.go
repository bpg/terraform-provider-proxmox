/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL can be obtained
 * at http://mozilla.org/MPL/2.0/.
 */

package agent

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

func TestResourceSchemaFreezeFS(t *testing.T) {
	t.Parallel()

	attribute, ok := ResourceSchema().(schema.SingleNestedAttribute)
	require.True(t, ok)

	freezeFS, ok := attribute.Attributes["freeze_fs"].(schema.BoolAttribute)
	require.True(t, ok)
	require.True(t, freezeFS.Optional)
	require.True(t, freezeFS.Computed)
	require.NotNil(t, freezeFS.Default)
}

func TestFillUpdateBodyRemovesAgentBlock(t *testing.T) {
	t.Parallel()

	plan := types.ObjectNull(map[string]attr.Type{
		"enabled":   types.BoolType,
		"trim":      types.BoolType,
		"type":      types.StringType,
		"freeze_fs": types.BoolType,
	})

	stateModel := Model{
		Enabled:  types.BoolValue(true),
		Trim:     types.BoolValue(false),
		Type:     types.StringValue("virtio"),
		FreezeFS: types.BoolValue(false),
	}
	state, d := types.ObjectValueFrom(context.Background(), attributeTypes(), stateModel)
	require.False(t, d.HasError())

	body := &vms.UpdateRequestBody{}
	FillUpdateBody(context.Background(), plan, state, body, &diag.Diagnostics{})

	require.Contains(t, body.Delete, "agent")
}
