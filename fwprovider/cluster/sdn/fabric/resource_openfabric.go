/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabric

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/fabrics"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

var (
	_ resource.ResourceWithConfigure   = &OpenFabricResource{}
	_ resource.ResourceWithImportState = &OpenFabricResource{}
)

type openFabricModel struct {
	genericModel

	IPv4Prefix    types.String `tfsdk:"ip_prefix"`
	IPv6Prefix    types.String `tfsdk:"ip6_prefix"`
	CsnpInterval  types.Int64  `tfsdk:"csnp_interval"`
	HelloInterval types.Int64  `tfsdk:"hello_interval"`
}

func (m *openFabricModel) fromAPI(name string, data *fabrics.FabricData, diags *diag.Diagnostics) {
	m.genericModel.fromAPI(name, data, diags)

	m.IPv4Prefix = m.handleDeletedStringValue(data.IPv4Prefix)
	m.IPv6Prefix = m.handleDeletedStringValue(data.IPv6Prefix)
	m.CsnpInterval = m.handleDeletedInt64Value(data.CsnpInterval)
	m.HelloInterval = m.handleDeletedInt64Value(data.HelloInterval)
}

func (m *openFabricModel) toAPI(ctx context.Context, diags *diag.Diagnostics) *fabrics.Fabric {
	data := m.genericModel.toAPI(ctx, diags)

	data.IPv4Prefix = m.IPv4Prefix.ValueStringPointer()
	data.IPv6Prefix = m.IPv6Prefix.ValueStringPointer()
	data.CsnpInterval = m.CsnpInterval.ValueInt64Pointer()
	data.HelloInterval = m.HelloInterval.ValueInt64Pointer()

	return data
}

func checkDeletedOpenFabricFields(state, plan *openFabricModel) []string {
	var toDelete []string

	if plan.IPv4Prefix.IsNull() && !state.IPv4Prefix.IsNull() {
		toDelete = append(toDelete, "ip_prefix")
	}

	if plan.IPv6Prefix.IsNull() && !state.IPv6Prefix.IsNull() {
		toDelete = append(toDelete, "ip6_prefix")
	}

	if plan.CsnpInterval.IsNull() && !state.CsnpInterval.IsNull() {
		toDelete = append(toDelete, "csnp_interval")
	}

	if plan.HelloInterval.IsNull() && !state.HelloInterval.IsNull() {
		toDelete = append(toDelete, "hello_interval")
	}

	toDelete = append(toDelete, checkDeletedFields(state.getGenericModel(), plan.getGenericModel())...)

	return toDelete
}

type OpenFabricResource struct {
	*genericFabricResource
}

func NewOpenFabricResource() resource.Resource {
	return &OpenFabricResource{
		genericFabricResource: newGenericFabricResource(fabricResourceConfig{
			typeNameSuffix: "_sdn_fabric_openfabric",
			fabricProtocol: fabrics.ProtocolOpenFabric,
			modelFunc:      func() fabricModel { return &openFabricModel{} },
		}).(*genericFabricResource),
	}
}

func (r *OpenFabricResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan openFabricModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state openFabricModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	toDelete := checkDeletedOpenFabricFields(&state, &plan)

	updateFabric := plan.toAPI(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	update := &fabrics.FabricUpdate{
		Fabric: *updateFabric,
		Delete: toDelete,
	}

	err := r.client.UpdateFabric(ctx, update)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating OpenFabric SDN Fabric",
			fmt.Sprintf("Could not update OpenFabric SDN Fabric %q: %v", plan.ID.ValueString(), err),
		)
		return
	}

	// Read updated state
	r.readAndSetState(ctx, plan.ID.ValueString(), &resp.State, &resp.Diagnostics)
}

func (r *OpenFabricResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "OpenFabric Fabric in Proxmox SDN. Fabrics in Proxmox VE SDN provide automated routing between nodes in a cluster.",
		MarkdownDescription: "OpenFabric Fabric in Proxmox SDN. Fabrics in Proxmox VE SDN provide automated routing between nodes in a cluster.",
		Attributes: genericAttributesWith(map[string]schema.Attribute{
			"ip_prefix": schema.StringAttribute{
				Description: "IPv4 prefix cidr for the fabric.",
				Optional:    true,
			},
			"ip6_prefix": schema.StringAttribute{
				Description: "IPv6 prefix cidr for the fabric.",
				Optional:    true,
			},
			"csnp_interval": schema.Int64Attribute{
				Description: "The csnp_interval property for OpenFabric.",
				Optional:    true,
			},
			"hello_interval": schema.Int64Attribute{
				Description: "The hello_interval property for OpenFabric.",
				Optional:    true,
			},
		}),
	}
}

func (r *OpenFabricResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("ip_prefix"),
			path.MatchRoot("ip6_prefix"),
		),
	}
}

func (m *openFabricModel) getGenericModel() *genericModel {
	return &m.genericModel
}
