package sdn

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/vnets"
)

var (
	_ resource.Resource                = &sdnVnetResource{}
	_ resource.ResourceWithConfigure   = &sdnVnetResource{}
	_ resource.ResourceWithImportState = &sdnVnetResource{}
)

type sdnVnetResource struct {
	client *vnets.Client
}

func NewSDNVnetResource() resource.Resource {
	return &sdnVnetResource{}
}

func (r *sdnVnetResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_sdn_vnet"
}

func (r *sdnVnetResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
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

	r.client = cfg.Client.Cluster().SDNVnets()
}

func (r *sdnVnetResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages Proxmox VE SDN vnet.",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ResourceID(),
			"name": schema.StringAttribute{
				Description: "Unique identifier for the vnet.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"zonetype": schema.StringAttribute{
				Required:    true,
				Description: "Parent's zone type. MUST be specified.",
			},
			"zone": schema.StringAttribute{
				Description: "The zone to which this vnet belongs.",
				Required:    true,
			},
			"alias": schema.StringAttribute{
				Optional:    true,
				Description: "An optional alias for this vnet.",
			},
			"isolate_ports": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to isolate ports within this vnet.",
			},
			"tag": schema.Int64Attribute{
				Optional:    true,
				Description: "Tag value for VLAN/VXLAN (depends on zone type).",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "Type of vnet (e.g. 'vnet').",
				Default:     stringdefault.StaticString("vnet"),
			},
			"vlanaware": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether this vnet is VLAN aware.",
			},
		},
	}
}

func (r *sdnVnetResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan sdnVnetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CreateVnet(ctx, plan.toAPIRequestBody())
	if err != nil {
		resp.Diagnostics.AddError("Error creating vnet", err.Error())
		return
	}

	plan.ID = plan.Name
	tflog.Info(ctx, "ZONETYPE value", map[string]any{"zonetype": plan.ZoneType.ValueString()})
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *sdnVnetResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state sdnVnetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.GetVnet(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Error reading vnet", err.Error())
		return
	}

	readModel := &sdnVnetModel{}
	readModel.importFromAPI(state.ID.ValueString(), data)
	// Preserve provider-only field
	readModel.ZoneType = state.ZoneType
	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *sdnVnetResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan sdnVnetModel
	var state sdnVnetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var toDelete []string
	checkDelete(plan.Alias, state.Alias, &toDelete, "alias")
	checkDelete(plan.IsolatePorts, state.IsolatePorts, &toDelete, "isolate-ports")
	checkDelete(plan.Tag, state.Tag, &toDelete, "tag")
	checkDelete(plan.Type, state.Type, &toDelete, "type")
	checkDelete(plan.VlanAware, state.VlanAware, &toDelete, "vlanaware")

	reqData := plan.toAPIRequestBody()
	reqData.Delete = toDelete

	err := r.client.UpdateVnet(ctx, reqData)
	if err != nil {
		resp.Diagnostics.AddError("Error updating vnet", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *sdnVnetResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state sdnVnetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteVnet(ctx, state.ID.ValueString())
	if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
		resp.Diagnostics.AddError("Error deleting vnet", err.Error())
	}
}

func (r *sdnVnetResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	data, err := r.client.GetVnet(ctx, req.ID)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError("Resource does not exist", err.Error())
			return
		}
		resp.Diagnostics.AddError("Failed to import resource", err.Error())
		return
	}

	readModel := &sdnVnetModel{}
	readModel.importFromAPI(req.ID, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func checkDelete(planField, stateField attr.Value, toDelete *[]string, apiName string) {
	if planField.IsNull() && !stateField.IsNull() {
		*toDelete = append(*toDelete, apiName)
	}
}

func (r *sdnVnetResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data sdnVnetModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Zone.IsNull() || data.Zone.IsUnknown() {
		return
	}

	if data.ZoneType.IsNull() || data.ZoneType.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("zonetype"),
			"Missing Required Field",
			"No Zone linked to this Vnet, please set the 'zonetype' property. \nEither from a created zone or a datasource import.")
		return
	}

	zoneType := data.ZoneType.ValueString()

	required := map[string][]string{
		"simple": {"name", "zone"},
		"vlan":   {"name", "zone", "tag"},
		"qinq":   {"name", "zone"},
		"vxlan":  {"name", "zone", "tag"},
		"evpn":   {"name", "zone", "tag"},
	}

	authorized := map[string]map[string]bool{
		"simple": {"name": true, "alias": true, "zone": true, "isolate_ports": true, "vlanaware": true},
		"vlan":   {"name": true, "alias": true, "zone": true, "tag": true, "isolate_ports": true, "vlanaware": true},
		"qinq":   {"name": true, "alias": true, "zone": true, "tag": true, "isolate_ports": true, "vlanaware": true},
		"vxlan":  {"name": true, "alias": true, "zone": true, "tag": true, "isolate_ports": true, "vlanaware": true},
		"evpn":   {"name": true, "alias": true, "zone": true, "tag": true, "isolate_ports": true},
	}

	fieldMap := map[string]attr.Value{
		"name":          data.Name,
		"zone":          data.Zone,
		"alias":         data.Alias,
		"tag":           data.Tag,
		"isolate_ports": data.IsolatePorts,
		"vlanaware":     data.VlanAware,
		"type":          data.Type,
	}

	// Check required fields
	for _, field := range required[zoneType] {
		if val, ok := fieldMap[field]; ok {
			if val.IsNull() || val.IsUnknown() {
				resp.Diagnostics.AddAttributeError(
					path.Root(field),
					"Missing Required Attribute",
					fmt.Sprintf("The attribute %q is required for SDN VNETs in a %q zone.", field, zoneType),
				)
			}
		}
	}

	for fieldName, val := range fieldMap {
		if !authorized[zoneType][fieldName] && !val.IsNull() && !val.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root(fieldName),
				"Unauthorized Attribute for Zone Type",
				fmt.Sprintf("The attribute %q is not allowed in VNETs under a %q zone.", fieldName, zoneType),
			)
		}
	}

}
