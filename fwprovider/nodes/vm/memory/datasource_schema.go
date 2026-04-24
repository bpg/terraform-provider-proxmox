/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package memory

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// DataSourceSchema defines the schema for the memory datasource.
func DataSourceSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		CustomType: basetypes.ObjectType{
			AttrTypes: attributeTypes(),
		},
		Description: "Memory configuration. Controls total available RAM and minimum guaranteed RAM via ballooning.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"size": schema.Int64Attribute{
				Description: "Total memory available to the VM in MiB.",
				Computed:    true,
			},
			"balloon": schema.Int64Attribute{
				Description: "Minimum guaranteed memory in MiB via balloon device. `0` means ballooning is disabled.",
				Computed:    true,
			},
			"shares": schema.Int64Attribute{
				Description: "CPU scheduler priority for memory ballooning.",
				Computed:    true,
			},
			"hugepages": schema.StringAttribute{
				Description: "Hugepages setting for VM memory (`2`, `1024`, or `any`).",
				Computed:    true,
			},
			"keep_hugepages": schema.BoolAttribute{
				Description: "Whether hugepages are kept allocated when the VM is stopped.",
				Computed:    true,
			},
		},
	}
}
