/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	poolapi "github.com/bpg/terraform-provider-proxmox/proxmox/nodes/ceph/pool"
)

// cephFSModel maps the schema data for proxmox_ceph_fs.
type cephFSModel struct {
	ID           types.String `tfsdk:"id"`
	NodeName     types.String `tfsdk:"node_name"`
	Name         types.String `tfsdk:"name"`
	PGNum        types.Int64  `tfsdk:"pg_num"`
	MetadataPool types.Object `tfsdk:"metadata_pool"`
	DataPool     types.Object `tfsdk:"data_pool"`
}

// cephFSPoolModel maps one of the filesystem's two pools.
type cephFSPoolModel struct {
	Name      types.String `tfsdk:"name"`
	PoolID    types.Int64  `tfsdk:"pool_id"`
	CrushRule types.String `tfsdk:"crush_rule"`
	Size      types.Int64  `tfsdk:"size"`
	MinSize   types.Int64  `tfsdk:"min_size"`
}

func poolAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":       types.StringType,
		"pool_id":    types.Int64Type,
		"crush_rule": types.StringType,
		"size":       types.Int64Type,
		"min_size":   types.Int64Type,
	}
}

// poolSettings returns the replication settings planned for a pool. Settings that are unset or unknown are nil, so
// they aren't sent.
func poolSettings(ctx context.Context, obj types.Object) (*poolapi.UpdateRequestBody, diag.Diagnostics) {
	body := &poolapi.UpdateRequestBody{}

	if obj.IsNull() || obj.IsUnknown() {
		return body, nil
	}

	var m cephFSPoolModel

	diags := obj.As(ctx, &m, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

	body.CrushRule = attribute.StringPtrFromValue(m.CrushRule)
	body.Size = attribute.Int64PtrFromValue(m.Size)
	body.MinSize = attribute.Int64PtrFromValue(m.MinSize)

	return body, diags
}

// poolObject builds a pool's state from its name, its ID (nil when Proxmox doesn't report it) and its status.
func poolObject(name string, id *int64, status *poolapi.StatusResponseData) types.Object {
	return types.ObjectValueMust(poolAttrTypes(), map[string]attr.Value{
		"name":       types.StringValue(name),
		"pool_id":    types.Int64PointerValue(id),
		"crush_rule": types.StringValue(status.CrushRule),
		"size":       types.Int64Value(status.Size),
		"min_size":   types.Int64Value(status.MinSize),
	})
}
