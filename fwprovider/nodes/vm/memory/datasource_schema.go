/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package memory

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// DataSourceSchema defines the schema for the memory datasource block.
func DataSourceSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		Description: "Memory configuration for the VM.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"size": schema.Int64Attribute{
				Description: "Total memory available to the VM in MiB.",
				Computed:    true,
			},
			"balloon": schema.Int64Attribute{
				Description: "Minimum guaranteed memory in MiB via balloon device. 0 disables the balloon driver.",
				Computed:    true,
			},
			"shares": schema.Int64Attribute{
				Description: "CPU scheduler priority for memory ballooning.",
				Computed:    true,
			},
			"hugepages": schema.StringAttribute{
				Description: "Use hugepages for VM memory. Options: '2' (2 MiB), '1024' (1 GiB), 'any'.",
				Computed:    true,
			},
			"keep_hugepages": schema.BoolAttribute{
				Description: "Keep hugepages allocated when VM is stopped.",
				Computed:    true,
			},
		},
	}
}
