/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package initialization

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceSchema returns the schema for the initialization block in a resource.
func ResourceSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		Description: "The cloud-init initialization configuration.",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			"dns": schema.SingleNestedAttribute{
				Description: "DNS configuration applied via cloud-init.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"domain": schema.StringAttribute{
						Description: "The DNS search domain.",
						Optional:    true,
					},
					"servers": schema.ListAttribute{
						Description: "List of DNS server IP addresses.",
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
			"ip_config": schema.ListNestedAttribute{
				Description: "IP configuration per network interface (up to 8). " +
					"The first entry maps to `ipconfig0`, the second to `ipconfig1`, etc.",
				Optional: true,
				Validators: []validator.List{
					listvalidator.SizeAtMost(8),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ipv4_address": schema.StringAttribute{
							Description: `IPv4 address in CIDR notation, or "dhcp".`,
							Optional:    true,
						},
						"ipv4_gateway": schema.StringAttribute{
							Description: "Default IPv4 gateway.",
							Optional:    true,
						},
						"ipv6_address": schema.StringAttribute{
							Description: `IPv6 address in CIDR notation, "dhcp", or "auto".`,
							Optional:    true,
						},
						"ipv6_gateway": schema.StringAttribute{
							Description: "Default IPv6 gateway.",
							Optional:    true,
						},
					},
				},
			},
			"meta_data_file_id": schema.StringAttribute{
				Description: "The file ID of a custom cloud-init meta data snippet.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"network_data_file_id": schema.StringAttribute{
				Description: "The file ID of a custom cloud-init network configuration snippet.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The cloud-init configuration format. " +
					"Defaults to the format inferred from the OS type.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("configdrive2", "nocloud", "openstacklatests"),
				},
			},
			"upgrade": schema.BoolAttribute{
				Description: "Whether to run package upgrades on the first boot.",
				Optional:    true,
			},
			"user_account": schema.SingleNestedAttribute{
				Description: "Cloud-init user account configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"keys": schema.ListAttribute{
						Description: "SSH public keys to add to the default user's authorized_keys file.",
						Optional:    true,
						ElementType: types.StringType,
					},
					"password": schema.StringAttribute{
						Description: "The login password for the default user. " +
							"This is a write-only field: the value is applied to Proxmox during apply " +
							"but is never stored in Terraform state.",
						Optional:  true,
						WriteOnly: true,
					},
					"username": schema.StringAttribute{
						Description: "The default user to configure.",
						Optional:    true,
					},
				},
			},
			"user_data_file_id": schema.StringAttribute{
				Description: "The file ID of a custom cloud-init user data snippet.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vendor_data_file_id": schema.StringAttribute{
				Description: "The file ID of a custom cloud-init vendor data snippet.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}
