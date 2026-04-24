/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package memory_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/memory"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func customInt64Ptr(v int64) *proxmoxtypes.CustomInt64 {
	c := proxmoxtypes.CustomInt64(v)
	return &c
}

func customBoolPtr(v bool) *proxmoxtypes.CustomBool {
	c := proxmoxtypes.CustomBool(v)
	return &c
}

// planValue builds a memory.Value from a Model via ObjectValueFrom, asserting no diag errors.
func planValue(t *testing.T, m memory.Model) memory.Value {
	t.Helper()

	ctx := context.Background()
	attrTypes := memory.NullValue().AttributeTypes(ctx)

	v, diags := types.ObjectValueFrom(ctx, attrTypes, m)
	require.False(t, diags.HasError(), "planValue: diags: %+v", diags)

	return v
}

func TestNewValue_AllNil_ReturnsNull(t *testing.T) {
	t.Parallel()

	var diags diag.Diagnostics

	value := memory.NewValue(context.Background(), &vms.GetResponseData{}, &diags)

	require.False(t, diags.HasError())
	assert.True(t, value.IsNull(), "expected NullValue when all memory fields are nil on the API response")
}

func TestNewValue_SomeFieldsSet_ReturnsPopulated(t *testing.T) {
	t.Parallel()

	var diags diag.Diagnostics

	resp := &vms.GetResponseData{
		DedicatedMemory: customInt64Ptr(4096),
		FloatingMemory:  customInt64Ptr(1024),
	}

	value := memory.NewValue(context.Background(), resp, &diags)

	require.False(t, diags.HasError())
	require.False(t, value.IsNull())

	attrs := value.Attributes()
	assert.Equal(t, types.Int64Value(4096), attrs["size"])
	assert.Equal(t, types.Int64Value(1024), attrs["balloon"])
	assert.Equal(t, types.Int64Null(), attrs["shares"])
	assert.Equal(t, types.StringNull(), attrs["hugepages"])
	assert.Equal(t, types.BoolNull(), attrs["keep_hugepages"])
}

func TestNewValue_HugepagesAndKeep_PropagateThrough(t *testing.T) {
	t.Parallel()

	var diags diag.Diagnostics

	resp := &vms.GetResponseData{
		Hugepages:            new("1024"),
		KeepHugepages:        customBoolPtr(true),
		FloatingMemoryShares: new(2000),
	}

	value := memory.NewValue(context.Background(), resp, &diags)

	require.False(t, diags.HasError())

	attrs := value.Attributes()
	assert.Equal(t, types.StringValue("1024"), attrs["hugepages"])
	assert.Equal(t, types.BoolValue(true), attrs["keep_hugepages"])
	assert.Equal(t, types.Int64Value(2000), attrs["shares"])
}

func TestFillCreateBody_NullPlan_NoChanges(t *testing.T) {
	t.Parallel()

	var (
		body  vms.CreateRequestBody
		diags diag.Diagnostics
	)

	memory.FillCreateBody(context.Background(), memory.NullValue(), &body, &diags)

	require.False(t, diags.HasError())
	assert.Nil(t, body.DedicatedMemory)
	assert.Nil(t, body.FloatingMemory)
	assert.Nil(t, body.FloatingMemoryShares)
	assert.Nil(t, body.Hugepages)
	assert.Nil(t, body.KeepHugepages)
	assert.Empty(t, body.Delete)
}

func TestFillCreateBody_PopulatedPlan_SetsBody(t *testing.T) {
	t.Parallel()

	plan := planValue(t, memory.Model{
		Size:          types.Int64Value(4096),
		Balloon:       types.Int64Value(1024),
		Shares:        types.Int64Null(),
		Hugepages:     types.StringValue("2"),
		KeepHugepages: types.BoolValue(true),
	})

	var (
		body  vms.CreateRequestBody
		diags diag.Diagnostics
	)

	memory.FillCreateBody(context.Background(), plan, &body, &diags)

	require.False(t, diags.HasError())
	require.NotNil(t, body.DedicatedMemory)
	assert.Equal(t, 4096, *body.DedicatedMemory)
	require.NotNil(t, body.FloatingMemory)
	assert.Equal(t, 1024, *body.FloatingMemory)
	assert.Nil(t, body.FloatingMemoryShares)
	require.NotNil(t, body.Hugepages)
	assert.Equal(t, "2", *body.Hugepages)
	require.NotNil(t, body.KeepHugepages)
	assert.Equal(t, proxmoxtypes.CustomBool(true), *body.KeepHugepages)
	assert.Empty(t, body.Delete)
}

func TestFillUpdateBody_PlanRemovesBlock_DeletesPreviouslySetFields(t *testing.T) {
	t.Parallel()

	state := planValue(t, memory.Model{
		Size:          types.Int64Value(4096),
		Balloon:       types.Int64Value(1024),
		Shares:        types.Int64Null(),
		Hugepages:     types.StringValue("2"),
		KeepHugepages: types.BoolNull(),
	})

	var (
		body  vms.UpdateRequestBody
		diags diag.Diagnostics
	)

	memory.FillUpdateBody(context.Background(), memory.NullValue(), state, &body, &diags)

	require.False(t, diags.HasError())
	assert.ElementsMatch(t, []string{"memory", "balloon", "hugepages"}, body.Delete,
		"fields that had state values must be deleted; fields that were already null stay off the wire")
	assert.Nil(t, body.DedicatedMemory)
	assert.Nil(t, body.FloatingMemory)
	assert.Nil(t, body.Hugepages)
}

func TestFillUpdateBody_PlanRemovesSingleField_DeletesThatFieldOnly(t *testing.T) {
	t.Parallel()

	state := planValue(t, memory.Model{
		Size:          types.Int64Value(4096),
		Balloon:       types.Int64Value(1024),
		Shares:        types.Int64Null(),
		Hugepages:     types.StringNull(),
		KeepHugepages: types.BoolNull(),
	})
	plan := planValue(t, memory.Model{
		Size:          types.Int64Value(4096),
		Balloon:       types.Int64Null(),
		Shares:        types.Int64Null(),
		Hugepages:     types.StringNull(),
		KeepHugepages: types.BoolNull(),
	})

	var (
		body  vms.UpdateRequestBody
		diags diag.Diagnostics
	)

	memory.FillUpdateBody(context.Background(), plan, state, &body, &diags)

	require.False(t, diags.HasError())
	assert.Equal(t, []string{"balloon"}, body.Delete)
	require.NotNil(t, body.DedicatedMemory)
	assert.Equal(t, 4096, *body.DedicatedMemory)
	assert.Nil(t, body.FloatingMemory)
}

func TestFillUpdateBody_PlanEqualsState_NoChange(t *testing.T) {
	t.Parallel()

	m := memory.Model{
		Size:          types.Int64Value(2048),
		Balloon:       types.Int64Null(),
		Shares:        types.Int64Null(),
		Hugepages:     types.StringNull(),
		KeepHugepages: types.BoolNull(),
	}
	state := planValue(t, m)
	plan := planValue(t, m)

	var (
		body  vms.UpdateRequestBody
		diags diag.Diagnostics
	)

	memory.FillUpdateBody(context.Background(), plan, state, &body, &diags)

	require.False(t, diags.HasError())
	assert.Empty(t, body.Delete)
	assert.Nil(t, body.DedicatedMemory, "plan==state is a no-op on the wire")
}

func TestFillUpdateBody_PlanAddsFieldFromNullState_NoDelete(t *testing.T) {
	t.Parallel()

	plan := planValue(t, memory.Model{
		Size:          types.Int64Value(4096),
		Balloon:       types.Int64Null(),
		Shares:        types.Int64Null(),
		Hugepages:     types.StringNull(),
		KeepHugepages: types.BoolNull(),
	})

	var (
		body  vms.UpdateRequestBody
		diags diag.Diagnostics
	)

	memory.FillUpdateBody(context.Background(), plan, memory.NullValue(), &body, &diags)

	require.False(t, diags.HasError())
	assert.Empty(t, body.Delete)
	require.NotNil(t, body.DedicatedMemory)
	assert.Equal(t, 4096, *body.DedicatedMemory)
}
