/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package pool

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	poolapi "github.com/bpg/terraform-provider-proxmox/proxmox/nodes/ceph/pool"
)

// cephPoolModel maps the schema data for proxmox_ceph_pool.
type cephPoolModel struct {
	ID              types.String  `tfsdk:"id"`
	Name            types.String  `tfsdk:"name"`
	NodeName        types.String  `tfsdk:"node_name"`
	Application     types.String  `tfsdk:"application"`
	CrushRule       types.String  `tfsdk:"crush_rule"`
	ErasureCoding   types.String  `tfsdk:"erasure_coding"`
	MinSize         types.Int64   `tfsdk:"min_size"`
	PGAutoscaleMode types.String  `tfsdk:"pg_autoscale_mode"`
	PGNum           types.Int64   `tfsdk:"pg_num"`
	PGNumMin        types.Int64   `tfsdk:"pg_num_min"`
	Size            types.Int64   `tfsdk:"size"`
	TargetSize      types.String  `tfsdk:"target_size"`
	TargetSizeRatio types.Float64 `tfsdk:"target_size_ratio"`
	AddStorages     types.Bool    `tfsdk:"add_storages"`
	ForceDestroy    types.Bool    `tfsdk:"force_destroy"`
	RemoveStorages  types.Bool    `tfsdk:"remove_storages"`
	RemoveECProfile types.Bool    `tfsdk:"remove_ecprofile"`
}

// toCreateBody builds the POST body for creating a pool.
func (m *cephPoolModel) toCreateBody() *poolapi.CreateRequestBody {
	return &poolapi.CreateRequestBody{
		Name:            m.Name.ValueString(),
		AddStorages:     attribute.CustomBoolPtrFromValue(m.AddStorages),
		Application:     attribute.StringPtrFromValue(m.Application),
		CrushRule:       attribute.StringPtrFromValue(m.CrushRule),
		ErasureCoding:   attribute.StringPtrFromValue(m.ErasureCoding),
		MinSize:         attribute.Int64PtrFromValue(m.MinSize),
		PGAutoscaleMode: attribute.StringPtrFromValue(m.PGAutoscaleMode),
		PGNum:           attribute.Int64PtrFromValue(m.PGNum),
		PGNumMin:        attribute.Int64PtrFromValue(m.PGNumMin),
		Size:            attribute.Int64PtrFromValue(m.Size),
		TargetSize:      attribute.StringPtrFromValue(m.TargetSize),
		TargetSizeRatio: attribute.Float64PtrFromValue(m.TargetSizeRatio),
	}
}

// toUpdateBody builds the PUT body for updating a pool. erasure_coding and add_storages
// are not updatable and therefore omitted.
func (m *cephPoolModel) toUpdateBody() *poolapi.UpdateRequestBody {
	return &poolapi.UpdateRequestBody{
		Application:     attribute.StringPtrFromValue(m.Application),
		CrushRule:       attribute.StringPtrFromValue(m.CrushRule),
		MinSize:         attribute.Int64PtrFromValue(m.MinSize),
		PGAutoscaleMode: attribute.StringPtrFromValue(m.PGAutoscaleMode),
		PGNum:           attribute.Int64PtrFromValue(m.PGNum),
		PGNumMin:        attribute.Int64PtrFromValue(m.PGNumMin),
		Size:            attribute.Int64PtrFromValue(m.Size),
		TargetSize:      attribute.StringPtrFromValue(m.TargetSize),
		TargetSizeRatio: attribute.Float64PtrFromValue(m.TargetSizeRatio),
	}
}

// toDeleteParams builds the DELETE query params from local-only state attributes.
func (m *cephPoolModel) toDeleteParams() *poolapi.DeleteRequestParams {
	return &poolapi.DeleteRequestParams{
		Force:           attribute.CustomBoolPtrFromValue(m.ForceDestroy),
		RemoveECProfile: attribute.CustomBoolPtrFromValue(m.RemoveECProfile),
		RemoveStorages:  attribute.CustomBoolPtrFromValue(m.RemoveStorages),
	}
}

// fromAPI populates the model from a pool list-response entry. Local-only attributes
// (add_storages, force_destroy, remove_storages, remove_ecprofile) and target_size are
// not round-tripped and are left untouched.
func (m *cephPoolModel) fromAPI(data *poolapi.ListResponseData) {
	app := applicationFromMetadata(data.ApplicationMetadata)
	if app == "" {
		// PVE always returns the application in metadata once the pool is fully provisioned;
		// the empty-map fallback guards against transient list responses during create.
		app = "rbd"
	}

	m.ID = types.StringValue(m.NodeName.ValueString() + "/" + data.PoolName)
	m.Name = types.StringValue(data.PoolName)
	m.Application = types.StringValue(app)
	m.CrushRule = types.StringValue(data.CrushRuleName)
	m.MinSize = types.Int64Value(data.MinSize)
	m.PGAutoscaleMode = types.StringValue(data.PGAutoscaleMode)
	m.PGNum = types.Int64Value(data.PGNum)
	m.Size = types.Int64Value(data.Size)
	// pg_num_min, target_size, and target_size_ratio are write-only on the PVE side: the list
	// endpoint omits them so we deliberately do not touch the model values here. The user's
	// configured value remains in state.
}

// applicationFromMetadata extracts the first key from application_metadata. A pool has at
// most one application, so this is the application name; empty string when absent.
func applicationFromMetadata(meta map[string]any) string {
	for k := range meta {
		return k
	}

	return ""
}
