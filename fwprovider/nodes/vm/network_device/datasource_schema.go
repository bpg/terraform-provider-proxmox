/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package network_device

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSourceSchema defines the schema for a list of network devices on a VM datasource.
func DataSourceSchema() schema.Attribute {
	return schema.ListNestedAttribute{
		Description: "Network device configurations.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"bridge":       schema.StringAttribute{Computed: true},
				"disconnected": schema.BoolAttribute{Computed: true},
				"firewall":     schema.BoolAttribute{Computed: true},
				"mac_address":  schema.StringAttribute{Computed: true},
				"model":        schema.StringAttribute{Computed: true},
				"mtu":          schema.Int64Attribute{Computed: true},
				"queues":       schema.Int64Attribute{Computed: true},
				"rate_limit":   schema.Float64Attribute{Computed: true},
				"trunks":       schema.ListAttribute{Computed: true, ElementType: types.Int64Type},
				"vlan_id":      schema.Int64Attribute{Computed: true},
			},
		},
	}
}
