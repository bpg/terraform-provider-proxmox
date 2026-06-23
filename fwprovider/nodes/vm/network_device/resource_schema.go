/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package network_device

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceSchema defines the schema for a list of network devices on a VM resource.
func ResourceSchema() schema.Attribute {
	return schema.ListNestedAttribute{
		Description: "Network device configurations.",
		Optional:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"bridge": schema.StringAttribute{
					Description: "The bridge to attach the network device to (e.g. `vmbr0`).",
					Optional:    true,
				},
				"disconnected": schema.BoolAttribute{
					Description: "Whether the network cable is disconnected.",
					Optional:    true,
				},
				"firewall": schema.BoolAttribute{
					Description: "Whether the Proxmox firewall is enabled on this network device.",
					Optional:    true,
				},
				"mac_address": schema.StringAttribute{
					Description: "The MAC address of the network device. PVE generates one when not provided.",
					Optional:    true,
					Computed:    true,
				},
				"model": schema.StringAttribute{
					Description: "The network device model.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.OneOf("e1000", "e1000e", "rtl8139", "virtio", "vmxnet3"),
					},
				},
				"mtu": schema.Int64Attribute{
					Description: "The MTU for the network device. Only valid for `virtio` model.",
					Optional:    true,
					Validators: []validator.Int64{
						int64validator.Between(1, 65536),
					},
				},
				"queues": schema.Int64Attribute{
					Description: "The number of packet queues. Only valid for `virtio` model.",
					Optional:    true,
					Validators: []validator.Int64{
						int64validator.Between(1, 64),
					},
				},
				"rate_limit": schema.Float64Attribute{
					Description: "The rate limit in megabytes per second.",
					Optional:    true,
					Validators: []validator.Float64{
						float64validator.AtLeast(0),
					},
				},
				"trunks": schema.ListAttribute{
					Description: "List of VLAN IDs passed through the network device.",
					Optional:    true,
					ElementType: types.Int64Type,
				},
				"vlan_id": schema.Int64Attribute{
					Description: "The VLAN identifier assigned to the network device.",
					Optional:    true,
					Validators: []validator.Int64{
						int64validator.Between(1, 4094),
					},
				},
			},
		},
	}
}
