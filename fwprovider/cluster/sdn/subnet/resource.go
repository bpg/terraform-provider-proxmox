/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package subnet

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/subnets"
)

var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithConfigure   = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

type Resource struct {
	client *cluster.Client
}

func NewResource() resource.Resource {
	return &Resource{}
}

func (r *Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdn_subnet"
}

func (r *Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages SDN Subnets in Proxmox VE.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the SDN Subnet. Read only.",
				Required:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cidr": schema.StringAttribute{
				Description: "A CIDR network address, for example 10.0.0.0/8",
				Required:    true,
				CustomType:  customtypes.IPCIDRType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vnet": schema.StringAttribute{
				Description: "The VNet to which this subnet belongs.",
				Required:    true,
				Validators:  validators.SDNID(),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dhcp_dns_server": schema.StringAttribute{
				Description: "The DNS server used for DHCP.",
				Optional:    true,
			},
			"dhcp_range": schema.SingleNestedAttribute{
				Description: "DHCP range (start and end IPs).",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"start_address": schema.StringAttribute{
						Description: "Start of the DHCP range.",
						CustomType:  customtypes.IPAddrType{},
						Required:    true,
					},
					"end_address": schema.StringAttribute{
						Description: "End of the DHCP range.",
						CustomType:  customtypes.IPAddrType{},
						Required:    true,
					},
				},
			},
			"dns_zone_prefix": schema.StringAttribute{
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

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := r.client.SDNVnets(plan.VNet.ValueString()).Subnets()
	subnet := plan.toAPI()
	// this is a special case for subnet creation, use the CIDR as the ID
	// after creation, the ID will be set to the canonical ID
	subnet.ID = plan.CIDR.ValueString()

	err := client.CreateSubnet(ctx, subnet)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create SDN Subnet", err.Error())
		return
	}

	canonicalID, err := resolveCanonicalSubnetID(ctx, client, plan.CIDR.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Resolve SDN Subnet ID", err.Error())
		return
	}

	plan.ID = types.StringValue(canonicalID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := r.client.SDNVnets(state.VNet.ValueString()).Subnets()

	subnet, err := client.GetSubnet(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Unable to Read SDN Subnet", err.Error())

		return
	}

	readModel := &model{}
	if err := readModel.fromAPI(&subnet.Subnet); err != nil {
		resp.Diagnostics.AddError("Invalid Subnet Data", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	reqData := plan.toAPI()

	client := r.client.SDNVnets(plan.VNet.ValueString()).Subnets()

	err := client.UpdateSubnet(ctx, &subnets.SubnetUpdate{
		Subnet: *reqData,
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Update SDN Subnet", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := r.client.SDNVnets(state.VNet.ValueString()).Subnets()

	err := client.DeleteSubnet(ctx, state.ID.ValueString())
	if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
		resp.Diagnostics.AddError("Unable to Delete SDN Subnet", err.Error())
	}
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			"Expected import identifier in format 'vnet-id/subnet-id'.",
		)

		return
	}

	vnetID := parts[0]
	subnetID := parts[1]

	client := r.client.SDNVnets(vnetID).Subnets()

	subnet, err := client.GetSubnet(ctx, subnetID)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError("SDN Subnet Not Found", fmt.Sprintf("SDN Subnet with ID '%s' in VNet '%s' was not found", subnetID, vnetID))
			return
		}

		resp.Diagnostics.AddError("Unable to Import SDN Subnet", err.Error())

		return
	}

	readModel := &model{}
	if err := readModel.fromAPI(&subnet.Subnet); err != nil {
		resp.Diagnostics.AddError("Invalid Subnet Data", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func resolveCanonicalSubnetID(ctx context.Context, client *subnets.Client, cidr string) (string, error) {
	subnetList, err := client.GetSubnets(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list subnets for canonical name resolution: %w", err)
	}

	for _, subnet := range subnetList {
		// Proxmox canonical format is zone-prefixed.
		// e.g., zoneM-10.10.0.0-24 instead of 10.10.0.0/24.
		canonicalCIDR := strings.ReplaceAll(cidr, "/", "-")
		if strings.HasSuffix(subnet.ID, canonicalCIDR) {
			return subnet.ID, nil
		}
	}

	return "", fmt.Errorf("could not resolve canonical subnet ID for %s", cidr)
}

/*
ValidateConfig checks that the subnet's field are correctly set.
Particularly that gateway, dhcp and dns are within CIDR.
*/
func (r *Resource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var config model
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, ipnet, err := net.ParseCIDR(config.CIDR.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("cidr"),
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
					fmt.Sprintf("%s must be within the subnet %s", ipVal.ValueString(), config.CIDR.ValueString()),
				)
			}
		}
	}

	checkIPInCIDR("gateway", config.Gateway)
	checkIPInCIDR("dhcp_dns_server", config.DhcpDnsServer)

	if config.DhcpRange != nil {
		r := config.DhcpRange
		if !r.StartAddress.IsNull() {
			ip := net.ParseIP(r.StartAddress.ValueString())
			if !ipnet.Contains(ip) {
				resp.Diagnostics.AddAttributeError(
					path.Root("dhcp_range").AtName("start_address"),
					"Invalid DHCP Range Start Address",
					fmt.Sprintf("Start address %s must be within the subnet %s", ip, config.CIDR.ValueString()),
				)
			}
		}

		if !r.EndAddress.IsNull() {
			ip := net.ParseIP(r.EndAddress.ValueString())
			if !ipnet.Contains(ip) {
				resp.Diagnostics.AddAttributeError(
					path.Root("dhcp_range").AtName("end_address"),
					"Invalid DHCP Range End Address",
					fmt.Sprintf("End address %s must be within the subnet %s", ip, config.CIDR.ValueString()),
				)
			}
		}

		// Validate that start address is not after end address
		if !r.StartAddress.IsNull() && !r.EndAddress.IsNull() {
			startIP := net.ParseIP(r.StartAddress.ValueString())

			endIP := net.ParseIP(r.EndAddress.ValueString())
			if startIP != nil && endIP != nil {
				if bytes.Compare(startIP, endIP) > 0 {
					resp.Diagnostics.AddAttributeError(
						path.Root("dhcp_range"),
						"Invalid DHCP Range",
						fmt.Sprintf("Start address %s must be less than or equal to end address %s", r.StartAddress.ValueString(), r.EndAddress.ValueString()),
					)
				}
			}
		}
	}
}
