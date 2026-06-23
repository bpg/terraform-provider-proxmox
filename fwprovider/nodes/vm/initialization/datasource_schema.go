/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package initialization

import (
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSourceSchema returns the schema for the initialization block in a datasource.
// The password field is omitted: Proxmox never returns the plaintext password via the API.
func DataSourceSchema() dsschema.Attribute {
	return dsschema.SingleNestedAttribute{
		Description: "The cloud-init initialization configuration.",
		Computed:    true,
		Attributes: map[string]dsschema.Attribute{
			"dns": dsschema.SingleNestedAttribute{
				Description: "DNS configuration applied via cloud-init.",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"domain": dsschema.StringAttribute{
						Description: "The DNS search domain.",
						Computed:    true,
					},
					"servers": dsschema.ListAttribute{
						Description: "List of DNS server IP addresses.",
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			"ip_config": dsschema.ListNestedAttribute{
				Description: "IP configuration per network interface.",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"ipv4_address": dsschema.StringAttribute{
							Description: `IPv4 address in CIDR notation, or "dhcp".`,
							Computed:    true,
						},
						"ipv4_gateway": dsschema.StringAttribute{
							Description: "Default IPv4 gateway.",
							Computed:    true,
						},
						"ipv6_address": dsschema.StringAttribute{
							Description: `IPv6 address in CIDR notation, "dhcp", or "auto".`,
							Computed:    true,
						},
						"ipv6_gateway": dsschema.StringAttribute{
							Description: "Default IPv6 gateway.",
							Computed:    true,
						},
					},
				},
			},
			"meta_data_file_id": dsschema.StringAttribute{
				Description: "The file ID of a custom cloud-init meta data snippet.",
				Computed:    true,
			},
			"network_data_file_id": dsschema.StringAttribute{
				Description: "The file ID of a custom cloud-init network configuration snippet.",
				Computed:    true,
			},
			"type": dsschema.StringAttribute{
				Description: "The cloud-init configuration format.",
				Computed:    true,
			},
			"upgrade": dsschema.BoolAttribute{
				Description: "Whether to run package upgrades on the first boot.",
				Computed:    true,
			},
			"user_account": dsschema.SingleNestedAttribute{
				Description: "Cloud-init user account configuration.",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"keys": dsschema.ListAttribute{
						Description: "SSH public keys for the default user.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"username": dsschema.StringAttribute{
						Description: "The default user.",
						Computed:    true,
					},
				},
			},
			"user_data_file_id": dsschema.StringAttribute{
				Description: "The file ID of a custom cloud-init user data snippet.",
				Computed:    true,
			},
			"vendor_data_file_id": dsschema.StringAttribute{
				Description: "The file ID of a custom cloud-init vendor data snippet.",
				Computed:    true,
			},
		},
	}
}
