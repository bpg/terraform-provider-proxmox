/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package pool

import (
	"sort"

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

// fromAPI populates the model from a pool /status response. Local-only attributes
// (add_storages, force_destroy, remove_storages, remove_ecprofile) and target_size /
// target_size_ratio are not round-tripped and are left untouched (target_size has a
// representation mismatch — PVE accepts a unit-suffixed string but returns bytes
// integer; target_size_ratio is omitted for consistency with target_size).
func (m *cephPoolModel) fromAPI(data *poolapi.StatusResponseData) {
	m.ID = types.StringValue(m.NodeName.ValueString() + "/" + data.Name)
	m.Name = types.StringValue(data.Name)

	// application handling: PVE returns application_list (when verbose=1 is set on
	// /status) for a fully provisioned pool. If the API response is transiently empty,
	// prefer the existing state value over an unconditional fallback so a user-set
	// value (e.g. "cephfs") is never clobbered to the server default. Fall back to
	// "rbd" only when the state has no value yet (first read after a create where the
	// user didn't specify the attribute).
	if app := applicationFromList(data.ApplicationList); app != "" {
		m.Application = types.StringValue(app)
	} else if !attribute.IsDefined(m.Application) {
		m.Application = types.StringValue("rbd")
	}

	m.CrushRule = types.StringValue(data.CrushRule)
	m.MinSize = types.Int64Value(data.MinSize)
	m.PGAutoscaleMode = types.StringValue(data.PGAutoscaleMode)
	m.PGNum = types.Int64Value(data.PGNum)
	m.Size = types.Int64Value(data.Size)

	// pg_num_min is Computed and may be unset on the server (null in /status). Set an
	// explicit null in state so the value is known after Read; UseStateForUnknown
	// handles the planning side.
	if data.PGNumMin != nil {
		m.PGNumMin = types.Int64Value(*data.PGNumMin)
	} else {
		m.PGNumMin = types.Int64Null()
	}
	// target_size and target_size_ratio intentionally omitted — kept write-only for
	// consistency (see godoc above).
}

// applicationFromList returns the application name from application_list (returned
// by /status?verbose=1). PVE guarantees at most one entry per pool. If a future API
// version surfaces multiple entries, the lowest-sorted name wins so callers see
// deterministic output.
func applicationFromList(apps []string) string {
	if len(apps) == 0 {
		return ""
	}

	sorted := append([]string(nil), apps...)
	sort.Strings(sorted)

	return sorted[0]
}
