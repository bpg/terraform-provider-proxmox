package sdn

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/subnets"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &sdnSubnetResource{}
	_ resource.ResourceWithConfigure   = &sdnSubnetResource{}
	_ resource.ResourceWithImportState = &sdnSubnetResource{}
)

type sdnSubnetResource struct {
	client *subnets.Client
}

func NewSDNSubnetResource() resource.Resource {
	return &sdnSubnetResource{}
}

func (r *sdnSubnetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdn_subnet"
}

func (r *sdnSubnetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = cfg.Client.Cluster().SDNSubnets()
}

func (r *sdnSubnetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages SDN Subnets in Proxmox VE.",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ResourceID(),
			"subnet": schema.StringAttribute{
				Required:    true,
				Description: "The name/ID of the subnet.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"canonical_name": schema.StringAttribute{
				Computed:    true,
				Description: "Canonical name of the subnet (e.g. zoneM-10.10.0.0-24).",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "Subnet type (set default at 'subnet')",
				Default:     stringdefault.StaticString("subnet"),
			},
			"vnet": schema.StringAttribute{
				Required:    true,
				Description: "The VNet to which this subnet belongs.",
			},
			"dhcp_dns_server": schema.StringAttribute{
				Optional:    true,
				Description: "The DNS server used for DHCP.",
			},
			"dhcp_range": schema.ListNestedAttribute{
				Optional:    true,
				Description: "List of DHCP ranges (start and end IPs).",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"start_address": schema.StringAttribute{
							Required:    true,
							Description: "Start of the DHCP range.",
						},
						"end_address": schema.StringAttribute{
							Required:    true,
							Description: "End of the DHCP range.",
						},
					},
				},
			},
			"dnszoneprefix": schema.StringAttribute{
				Optional:    true,
				Description: "Prefix used for DNS zone delegation.",
			},
			"gateway": schema.StringAttribute{
				Optional:    true,
				Description: "The gateway address for the subnet.",
			},
			"snat": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether SNAT is enabled for the subnet.",
			},
		},
	}
}

func (r *sdnSubnetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sdnSubnetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Vnet.IsNull() || plan.Vnet.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("vnet"),
			"missing required field",
			"Missing the parent vnet's ID attribute, which is required to define a subnet")
		return
	}
	err := r.client.CreateSubnet(ctx, plan.Vnet.ValueString(), plan.toAPIRequestBody())
	if err != nil {
		resp.Diagnostics.AddError("Error creating subnet", err.Error())
		return
	}

	tflog.Debug(ctx, "Created object's ID", map[string]any{"plan name:": plan.Subnet})
	plan.ID = plan.Subnet

	// Because proxmox API doesn't return the created object's properties and the subnet's name gets modified by proxmox internally
	// Read it back to get the canonical-ID from proxmox
	canonicalID, err := resolveCanonicalSubnetID(ctx, r.client, plan.Vnet.ValueString(), plan.Subnet.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error resolving canonical subnet ID", err.Error())
		return
	}

	plan.ID = types.StringValue(canonicalID)
	plan.CanonicalName = types.StringValue(canonicalID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *sdnSubnetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sdnSubnetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subnet, err := r.client.GetSubnet(ctx, state.Vnet.ValueString(), state.ID.ValueString())
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Error reading subnet", err.Error())
		return
	}

	readModel := &sdnSubnetModel{}
	readModel.Subnet = state.Subnet
	readModel.importFromAPI(state.ID.ValueString(), subnet)

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *sdnSubnetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sdnSubnetModel
	// var state sdnSubnetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	reqData := plan.toAPIRequestBody()
	// reqData.Delete = toDelete

	if plan.Vnet.IsNull() || plan.Vnet.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("vnet"),
			"missing required field",
			"Missing the parent vnet's ID attribute, which is required to define a subnet")
		return
	}
	err := r.client.UpdateSubnet(ctx, plan.Vnet.ValueString(), reqData)
	if err != nil {
		resp.Diagnostics.AddError("Error updating subnet", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *sdnSubnetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sdnSubnetModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSubnet(ctx, state.Vnet.ValueString(), state.ID.ValueString())
	if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
		resp.Diagnostics.AddError("Error deleting subnet", err.Error())
	}
}

func (r *sdnSubnetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expect ID format: "vnet/subnet"
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			"Expected import identifier in format 'vnet-id/subnet-id'.",
		)
		return
	}
	vnetID := parts[0]
	subnetID := parts[1]
	subnet, err := r.client.GetSubnet(ctx, vnetID, subnetID)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError("Subnet does not exist", err.Error())
			return
		}

		resp.Diagnostics.AddError("Unable to import subnet", err.Error())
		return
	}

	readModel := &sdnSubnetModel{}
	readModel.importFromAPI(req.ID, subnet)
	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func resolveCanonicalSubnetID(ctx context.Context, client *subnets.Client, vnet string, originalID string) (string, error) {
	subnets, err := client.GetSubnets(ctx, vnet)
	if err != nil {
		return "", fmt.Errorf("failed to list subnets for canonical name resolution: %w", err)
	}

	for _, subnet := range subnets {
		if subnet.ID == originalID {
			return subnet.ID, nil // Already canonical
		}

		// Proxmox canonical format is usually zone-prefixed:
		// e.g., zoneM-10-10-0-0-24 instead of 10.10.0.0/24
		if strings.HasSuffix(subnet.ID, strings.ReplaceAll(originalID, "/", "-")) {
			return subnet.ID, nil
		}
	}

	return "", fmt.Errorf("could not resolve canonical subnet ID for %s", originalID)
}

// ValidateConfig checks that the subnet's field are correctly set. Particularly that gateway, dhcp and dns are within CIDR
func (r *sdnSubnetResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config sdnSubnetModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, ipnet, err := net.ParseCIDR(config.Subnet.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("subnet"),
			"Invalid Subnet",
			fmt.Sprintf("Could not parse subnet: %s", err),
		)
		return
	}

	checkIPInCIDR := func(attrName string, ipVal types.String) {
		if !ipVal.IsNull() {
			ip := net.ParseIP(ipVal.ValueString())
			if ip == nil {
				resp.Diagnostics.AddAttributeError(
					path.Root(attrName),
					"Invalid IP Address",
					fmt.Sprintf("Could not parse IP address: %s", ipVal.ValueString()),
				)
				return
			}

			if !ipnet.Contains(ip) {
				resp.Diagnostics.AddAttributeError(
					path.Root(attrName),
					"Invalid IP for Subnet",
					fmt.Sprintf("%s must be within the subnet %s", ipVal.ValueString(), config.Subnet.ValueString()),
				)
			}
		}
	}

	checkIPInCIDR("gateway", config.Gateway)
	checkIPInCIDR("dhcp_dns_server", config.DhcpDnsServer)

	for i, r := range config.DhcpRange {
		if !r.StartAddress.IsNull() {
			ip := net.ParseIP(r.StartAddress.ValueString())
			if !ipnet.Contains(ip) {
				resp.Diagnostics.AddAttributeError(
					path.Root("dhcp_range").AtListIndex(i).AtMapKey("start_address"),
					"Invalid DHCP Range Start Address",
					fmt.Sprintf("Start address %s must be within the subnet %s", ip, config.Subnet.ValueString()),
				)
			}
		}

		if !r.EndAddress.IsNull() {
			ip := net.ParseIP(r.EndAddress.ValueString())
			if !ipnet.Contains(ip) {
				resp.Diagnostics.AddAttributeError(
					path.Root("dhcp_range").AtListIndex(i).AtMapKey("end_address"),
					"Invalid DHCP Range End Address",
					fmt.Sprintf("End address %s must be within the subnet %s", ip, config.Subnet.ValueString()),
				)
			}
		}
	}
}
