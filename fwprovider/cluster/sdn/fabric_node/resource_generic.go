/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabric_node

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/fabric_nodes"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type genericModel struct {
	ID             types.String `tfsdk:"id"`
	NodeID         types.String `tfsdk:"node_id"`
	FabricID       types.String `tfsdk:"fabric_id"`
	InterfaceNames types.Set    `tfsdk:"interface_names"`
}

func (m *genericModel) fromAPI(name string, data *fabric_nodes.FabricNodeData, diags *diag.Diagnostics) {
	m.ID = types.StringValue(name)

	parts := strings.SplitN(name, "/", 2)
	if len(parts) != 2 {
		diags.AddError(
			"Unexpected SDN Fabric Node ID Format",
			fmt.Sprintf("Expected SDN Fabric Node ID to be in the format <fabric_id>/<node_id>, got: %s", name),
		)

		return
	}

	fabricID := parts[0]
	nodeID := parts[1]
	m.NodeID = types.StringValue(nodeID)
	m.FabricID = types.StringValue(fabricID)
	interfaceNameSet, d := m.toInterfaceNamesFromInterfaces(data.Interfaces)
	diags.Append(d...)

	if diags.HasError() {
		return
	}

	m.InterfaceNames = interfaceNameSet
}

func (m *genericModel) toAPI(ctx context.Context, diags *diag.Diagnostics) *fabric_nodes.FabricNode {
	data := &fabric_nodes.FabricNode{}
	data.NodeID = m.NodeID.ValueString()
	data.FabricID = m.FabricID.ValueString()

	if m.InterfaceNames.IsNull() {
		data.Interfaces = []string{}
	} else {
		var interfaces []string
		diags.Append(m.InterfaceNames.ElementsAs(ctx, &interfaces, false)...)

		if diags.HasError() {
			return nil
		}

		for i, iface := range interfaces {
			interfaces[i] = fmt.Sprintf("name=%s", iface)
		}

		data.Interfaces = interfaces
	}

	return data
}

func (m *genericModel) getID() string {
	return fmt.Sprintf("%s/%s", m.FabricID.ValueString(), m.NodeID.ValueString())
}

func (m *genericModel) toInterfaceNamesFromInterfaces(value []string) (types.Set, diag.Diagnostics) {
	if value == nil {
		return types.SetValue(types.StringType, []attr.Value{})
	}

	// Convet filtered slice to types.List
	interfaces := make([]attr.Value, 0)

	for _, iface := range value {
		parts := strings.SplitN(iface, "=", 2)
		if len(parts) != 2 {
			emptySet, diags := types.SetValue(types.StringType, []attr.Value{})
			diags.AddError(
				"Unexpected SDN Fabric Node ID Format",
				fmt.Sprintf("Expected SDN Fabric Node ID to be in the format <fabric_id>/<node_id>, got: %s", iface),
			)

			return emptySet, diags
		}

		k := parts[0]
		v := parts[1]
		// currently, only "name" key is relevant for interfaces
		if k == "name" {
			ifaceStringValue := types.StringValue(v)
			interfaces = append(interfaces, ifaceStringValue)
		}
	}

	return types.SetValue(types.StringType, interfaces)
}

func (m *genericModel) handleDeletedIPAddrValue(value *string) customtypes.IPAddrValue {
	if value == nil {
		return customtypes.IPAddrValue{
			StringValue: types.StringNull(),
		}
	}

	return customtypes.NewIPAddrPointerValue(value)
}

func checkDeletedFields(state, plan *genericModel) []string {
	var toDelete []string

	if plan.InterfaceNames.IsNull() && !state.InterfaceNames.IsNull() {
		toDelete = append(toDelete, "interfaces")
	}

	return toDelete
}

func genericAttributesWith(extraAttributes map[string]schema.Attribute) map[string]schema.Attribute {
	// Start with generic attributes as the base
	result := map[string]schema.Attribute{
		"node_id": schema.StringAttribute{
			Description: "The unique identifier of the SDN fabric node.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"fabric_id": schema.StringAttribute{
			Description: "The unique identifier of the SDN fabric.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: validators.SDNID(),
		},
		"interface_names": schema.SetAttribute{
			Description: "Set of interfaces associated with the fabric node.",
			ElementType: types.StringType,
			Required:    true,
		},
		"id": schema.StringAttribute{
			Description: "The ID of the SDN fabric node, in the format <fabric_id>/<node_id>.",
			Computed:    true,
		},
	}

	// Add extra attributes, allowing them to override generic ones if needed
	if extraAttributes != nil {
		maps.Copy(result, extraAttributes)
	}

	return result
}

type fabricNodeModel interface {
	fromAPI(name string, data *fabric_nodes.FabricNodeData, diags *diag.Diagnostics)
	toAPI(ctx context.Context, diags *diag.Diagnostics) *fabric_nodes.FabricNode
	getID() string
	getGenericModel() *genericModel
}

type fabricNodeResourceConfig struct {
	typeNameSuffix string
	// fabricID       string
	fabricProtocol string
	modelFunc      func() fabricNodeModel
}

type genericFabricNodeResource struct {
	client *cluster.Client
	config fabricNodeResourceConfig
}

func newGenericFabricNodeResource(cfg fabricNodeResourceConfig) resource.Resource {
	return &genericFabricNodeResource{config: cfg}
}

func (r *genericFabricNodeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.config.typeNameSuffix
}

func (r *genericFabricNodeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.Resource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected config.Resource, got: %T", req.ProviderData),
		)

		return
	}

	r.client = cfg.Client.Cluster()
}

func (r *genericFabricNodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := r.config.modelFunc()
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	newFabricNode := plan.toAPI(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	newFabricNode.Protocol = ptr.Ptr(r.config.fabricProtocol)

	client := r.client.SDNFabricNodes(plan.getGenericModel().FabricID.ValueString(), r.config.fabricProtocol)

	tflog.Debug(ctx, fmt.Sprintf("Create plan %+v", plan))

	if err := client.CreateFabricNode(ctx, newFabricNode); err != nil {
		resp.Diagnostics.AddError("Unable to Create SDN Fabric Node", err.Error())
		return
	}

	r.readAndSetState(ctx, client, plan.getID(), &resp.State, &resp.Diagnostics)
}

func (r *genericFabricNodeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := r.config.modelFunc()
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := r.client.SDNFabricNodes(state.getGenericModel().FabricID.ValueString(), r.config.fabricProtocol)
	tflog.Debug(ctx, fmt.Sprintf("Getting fabric node with ID %s", state.getGenericModel().NodeID.ValueString()))

	fabricNode, err := client.GetFabricNodeWithParams(
		ctx,
		state.getGenericModel().NodeID.ValueString(),
		&sdn.QueryParams{Pending: proxmoxtypes.CustomBool(true).Pointer()},
	)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Unable to Read SDN Fabric Node", err.Error())

		return
	}

	r.setModelFromFabricNode(ctx, fabricNode, &resp.State, &resp.Diagnostics)
}

func (r *genericFabricNodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := r.config.modelFunc()
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)

	state := r.config.modelFunc()
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateFabric := plan.toAPI(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	toDelete := checkDeletedFields(state.getGenericModel(), plan.getGenericModel())
	update := &fabric_nodes.FabricNodeUpdate{
		FabricNode: *updateFabric,
		Delete:     toDelete,
	}

	client := r.client.SDNFabricNodes(plan.getGenericModel().FabricID.ValueString(), r.config.fabricProtocol)
	if err := client.UpdateFabricNode(ctx, update); err != nil {
		resp.Diagnostics.AddError("Unable to Update SDN Fabric Node", err.Error())
		return
	}

	r.readAndSetState(ctx, client, plan.getID(), &resp.State, &resp.Diagnostics)
}

func (r *genericFabricNodeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := r.config.modelFunc()
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := r.client.SDNFabricNodes(state.getGenericModel().FabricID.ValueString(), r.config.fabricProtocol)
	if err := client.DeleteFabricNode(ctx, state.getGenericModel().NodeID.ValueString()); err != nil &&
		!errors.Is(err, api.ErrResourceDoesNotExist) {
		resp.Diagnostics.AddError("Unable to Delete SDN Fabric Node", err.Error())
	}
}

func (r *genericFabricNodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected SDN Fabric Node ID Format",
			fmt.Sprintf("Expected SDN Fabric Node ID to be in the format <fabric_id>/<node_id>, got: %s", req.ID),
		)

		return
	}

	fabricID := parts[0]
	nodeID := parts[1]
	client := r.client.SDNFabricNodes(fabricID, r.config.fabricProtocol)

	fabricNode, err := client.GetFabricNode(ctx, nodeID)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError(fmt.Sprintf("Fabric node %s does not exist", req.ID), err.Error())
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Import SDN Fabric node %s", req.ID), err.Error())

		return
	}

	r.setModelFromFabricNode(ctx, fabricNode, &resp.State, &resp.Diagnostics)
}

// Schema is required to satisfy the resource.Resource interface. It should be implemented by the specific resource.
func (r *genericFabricNodeResource) Schema(_ context.Context, _ resource.SchemaRequest, _ *resource.SchemaResponse) {
	// Intentionally left blank. Should be set by the specific resource.
}

func (r *genericFabricNodeResource) readAndSetState(
	ctx context.Context,
	client *fabric_nodes.Client,
	fabricNodeID string,
	state *tfsdk.State,
	diags *diag.Diagnostics,
) {
	parts := strings.SplitN(fabricNodeID, "/", 2)
	if len(parts) != 2 {
		diags.AddError(
			"Unexpected SDN Fabric Node ID Format",
			fmt.Sprintf("Expected SDN Fabric Node ID to be in the format <fabric_id>/<node_id>, got: %s", fabricNodeID),
		)

		return
	}

	nodeID := parts[1]

	fabricNode, err := client.GetFabricNodeWithParams(ctx, nodeID, &sdn.QueryParams{Pending: proxmoxtypes.CustomBool(true).Pointer()})
	if err != nil {
		diags.AddError("Unable to Read SDN Fabric Node", err.Error())
		return
	}

	r.setModelFromFabricNode(ctx, fabricNode, state, diags)
}

func (r *genericFabricNodeResource) setModelFromFabricNode(
	ctx context.Context,
	fabricNode *fabric_nodes.FabricNodeData,
	state *tfsdk.State,
	diags *diag.Diagnostics,
) {
	readModel := r.config.modelFunc()
	id := fmt.Sprintf("%s/%s", fabricNode.FabricID, fabricNode.NodeID)
	readModel.fromAPI(id, fabricNode, diags)

	if diags.HasError() {
		return
	}

	diags.Append(state.Set(ctx, readModel)...)
}
