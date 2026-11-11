/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package pool

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	poolapi "github.com/bpg/terraform-provider-proxmox/proxmox/nodes/ceph/pool"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// Proxmox creates an erasure-coded pool as two pools named after the requested name; none has the name itself.
const (
	ecDataPoolSuffix     = "-data"
	ecMetadataPoolSuffix = "-metadata"
)

// cephPoolModel maps the schema data for proxmox_ceph_pool.
type cephPoolModel struct {
	ID               types.String  `tfsdk:"id"`
	Name             types.String  `tfsdk:"name"`
	NodeName         types.String  `tfsdk:"node_name"`
	Application      types.String  `tfsdk:"application"`
	CrushRule        types.String  `tfsdk:"crush_rule"`
	ErasureCoding    types.Object  `tfsdk:"erasure_coding"`
	MinSize          types.Int64   `tfsdk:"min_size"`
	PGAutoscaleMode  types.String  `tfsdk:"pg_autoscale_mode"`
	PGNum            types.Int64   `tfsdk:"pg_num"`
	PGNumMin         types.Int64   `tfsdk:"pg_num_min"`
	Size             types.Int64   `tfsdk:"size"`
	TargetSize       types.String  `tfsdk:"target_size"`
	TargetSizeRatio  types.Float64 `tfsdk:"target_size_ratio"`
	PoolID           types.Int64   `tfsdk:"pool_id"`
	DataPool         types.Object  `tfsdk:"data_pool"`
	MetadataPoolName types.String  `tfsdk:"metadata_pool_name"`
	MetadataPoolID   types.Int64   `tfsdk:"metadata_pool_id"`
	ForceDestroy     types.Bool    `tfsdk:"force_destroy"`
}

// erasureCodingModel maps the erasure_coding attribute, the fields of Proxmox's `erasure-coding` property string.
type erasureCodingModel struct {
	K             types.Int64  `tfsdk:"k"`
	M             types.Int64  `tfsdk:"m"`
	FailureDomain types.String `tfsdk:"failure_domain"`
	DeviceClass   types.String `tfsdk:"device_class"`
	Profile       types.String `tfsdk:"profile"`
}

// dataPoolModel maps the read-only data_pool attribute of an erasure-coded pool.
type dataPoolModel struct {
	Name      types.String `tfsdk:"name"`
	PoolID    types.Int64  `tfsdk:"pool_id"`
	Size      types.Int64  `tfsdk:"size"`
	MinSize   types.Int64  `tfsdk:"min_size"`
	CrushRule types.String `tfsdk:"crush_rule"`
}

func dataPoolAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":       types.StringType,
		"pool_id":    types.Int64Type,
		"size":       types.Int64Type,
		"min_size":   types.Int64Type,
		"crush_rule": types.StringType,
	}
}

// erasureCoded reports whether the model describes an erasure-coded pool: configured as one, or read as one (an
// import, where erasure_coding is still null).
func (m *cephPoolModel) erasureCoded() bool {
	return !m.ErasureCoding.IsNull() || (!m.DataPool.IsNull() && !m.DataPool.IsUnknown())
}

// erasureCoding returns the configured erasure coding, or nil when it is null or unknown.
func (m *cephPoolModel) erasureCoding(ctx context.Context) (*erasureCodingModel, diag.Diagnostics) {
	if m.ErasureCoding.IsNull() || m.ErasureCoding.IsUnknown() {
		return nil, nil
	}

	var ec erasureCodingModel

	diags := m.ErasureCoding.As(ctx, &ec, basetypes.ObjectAsOptions{})

	return &ec, diags
}

// String builds Proxmox's erasure-coding property string, e.g. `k=2,m=2,failure-domain=osd,device-class=hdd`.
func (e *erasureCodingModel) String() string {
	parts := []string{
		fmt.Sprintf("k=%d", e.K.ValueInt64()),
		fmt.Sprintf("m=%d", e.M.ValueInt64()),
	}

	if v := e.FailureDomain.ValueString(); v != "" {
		parts = append(parts, "failure-domain="+v)
	}

	if v := e.DeviceClass.ValueString(); v != "" {
		parts = append(parts, "device-class="+v)
	}

	if v := e.Profile.ValueString(); v != "" {
		parts = append(parts, "profile="+v)
	}

	return strings.Join(parts, ",")
}

// toCreateBody builds the POST body for creating a pool. For an erasure-coded pool, Proxmox applies size, min_size
// and crush_rule to the metadata pool and the other settings to the data pool (and at creation also to the
// metadata pool, except pg_num).
func (m *cephPoolModel) toCreateBody(ctx context.Context) (*poolapi.CreateRequestBody, diag.Diagnostics) {
	body := &poolapi.CreateRequestBody{
		Name: m.Name.ValueString(),
		// Storages on a pool are separate resources; Proxmox adds one to erasure-coded pools unless told not to.
		AddStorages:     proxmoxtypes.CustomBool(false).Pointer(),
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

	ec, diags := m.erasureCoding(ctx)
	if ec != nil {
		s := ec.String()
		body.ErasureCoding = &s
	}

	return body, diags
}

// toUpdateBodies builds the PUT bodies, keyed by the Ceph pool they go to. An erasure-coded pool splits them as at
// creation: size, min_size and crush_rule go to the metadata pool, the other settings to the data pool.
func (m *cephPoolModel) toUpdateBodies() map[string]*poolapi.UpdateRequestBody {
	name := m.Name.ValueString()

	if !m.erasureCoded() {
		return map[string]*poolapi.UpdateRequestBody{
			name: {
				Application:     attribute.StringPtrFromValue(m.Application),
				CrushRule:       attribute.StringPtrFromValue(m.CrushRule),
				MinSize:         attribute.Int64PtrFromValue(m.MinSize),
				PGAutoscaleMode: attribute.StringPtrFromValue(m.PGAutoscaleMode),
				PGNum:           attribute.Int64PtrFromValue(m.PGNum),
				PGNumMin:        attribute.Int64PtrFromValue(m.PGNumMin),
				Size:            attribute.Int64PtrFromValue(m.Size),
				TargetSize:      attribute.StringPtrFromValue(m.TargetSize),
				TargetSizeRatio: attribute.Float64PtrFromValue(m.TargetSizeRatio),
			},
		}
	}

	return map[string]*poolapi.UpdateRequestBody{
		name + ecMetadataPoolSuffix: {
			CrushRule: attribute.StringPtrFromValue(m.CrushRule),
			MinSize:   attribute.Int64PtrFromValue(m.MinSize),
			Size:      attribute.Int64PtrFromValue(m.Size),
		},
		name + ecDataPoolSuffix: {
			Application:     attribute.StringPtrFromValue(m.Application),
			PGAutoscaleMode: attribute.StringPtrFromValue(m.PGAutoscaleMode),
			PGNum:           attribute.Int64PtrFromValue(m.PGNum),
			PGNumMin:        attribute.Int64PtrFromValue(m.PGNumMin),
			TargetSize:      attribute.StringPtrFromValue(m.TargetSize),
			TargetSizeRatio: attribute.Float64PtrFromValue(m.TargetSizeRatio),
		},
	}
}

// toDeleteParams builds the DELETE query params. remove_storages is never sent: storages on a pool are separate
// resources. The erasure code profile goes with the data pool only when Proxmox generated it for this pool; one
// named in erasure_coding.profile belongs to the user, and after an import its origin is unknown.
func (m *cephPoolModel) toDeleteParams(ctx context.Context) (*poolapi.DeleteRequestParams, diag.Diagnostics) {
	ec, diags := m.erasureCoding(ctx)
	generatedProfile := ec != nil && ec.Profile.ValueString() == ""

	return &poolapi.DeleteRequestParams{
		Force:           attribute.CustomBoolPtrFromValue(m.ForceDestroy),
		RemoveECProfile: proxmoxtypes.CustomBool(generatedProfile).Pointer(),
	}, diags
}

// fromAPI populates the model from the /status response of a replicated pool. Local-only attributes
// (erasure_coding, force_destroy) and target_size / target_size_ratio are not round-tripped and are left untouched
// (target_size has a representation mismatch — PVE accepts a unit-suffixed string but returns bytes integer;
// target_size_ratio is omitted for consistency with target_size).
func (m *cephPoolModel) fromAPI(data *poolapi.StatusResponseData) {
	m.ID = types.StringValue(m.NodeName.ValueString() + "/" + data.Name)
	m.Name = types.StringValue(data.Name)
	m.PoolID = types.Int64Value(data.ID)
	m.DataPool = types.ObjectNull(dataPoolAttrTypes())
	m.MetadataPoolName = types.StringNull()
	m.MetadataPoolID = types.Int64Null()

	m.setApplication(data.ApplicationList)

	m.CrushRule = types.StringValue(data.CrushRule)
	m.MinSize = types.Int64Value(data.MinSize)
	m.Size = types.Int64Value(data.Size)
	m.setPGSettings(data)
}

// fromAPIErasureCoded populates the model from the /status responses of an erasure-coded pool's two pools, mapping
// them as at creation: size, min_size and crush_rule from the metadata pool, the other settings from the data pool.
func (m *cephPoolModel) fromAPIErasureCoded(ctx context.Context, data, metadata *poolapi.StatusResponseData) diag.Diagnostics {
	m.ID = types.StringValue(m.NodeName.ValueString() + "/" + m.Name.ValueString())
	m.PoolID = types.Int64Null()
	m.MetadataPoolName = types.StringValue(metadata.Name)
	m.MetadataPoolID = types.Int64Value(metadata.ID)

	dataPool, diags := types.ObjectValueFrom(ctx, dataPoolAttrTypes(), dataPoolModel{
		Name:      types.StringValue(data.Name),
		PoolID:    types.Int64Value(data.ID),
		Size:      types.Int64Value(data.Size),
		MinSize:   types.Int64Value(data.MinSize),
		CrushRule: types.StringValue(data.CrushRule),
	})
	m.DataPool = dataPool

	m.setApplication(data.ApplicationList)

	m.CrushRule = types.StringValue(metadata.CrushRule)
	m.MinSize = types.Int64Value(metadata.MinSize)
	m.Size = types.Int64Value(metadata.Size)
	m.setPGSettings(data)

	return diags
}

// setApplication sets application from application_list. PVE returns application_list (when verbose=1 is set on
// /status) for a fully provisioned pool. If the API response is transiently empty, prefer the existing state value
// over an unconditional fallback so a user-set value (e.g. "cephfs") is never clobbered to the server default. Fall
// back to "rbd" only when the state has no value yet (first read after a create where the user didn't specify the
// attribute).
func (m *cephPoolModel) setApplication(apps []string) {
	if app := applicationFromList(apps); app != "" {
		m.Application = types.StringValue(app)
	} else if !attribute.IsDefined(m.Application) {
		m.Application = types.StringValue("rbd")
	}
}

// setPGSettings sets the placement group settings.
func (m *cephPoolModel) setPGSettings(data *poolapi.StatusResponseData) {
	m.PGAutoscaleMode = types.StringValue(data.PGAutoscaleMode)
	m.PGNum = types.Int64Value(data.PGNum)

	// pg_num_min is Computed and may be unset on the server (null in /status). Set an
	// explicit null in state so the value is known after Read; UseStateForUnknown
	// handles the planning side.
	if data.PGNumMin != nil {
		m.PGNumMin = types.Int64Value(*data.PGNumMin)
	} else {
		m.PGNumMin = types.Int64Null()
	}
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
