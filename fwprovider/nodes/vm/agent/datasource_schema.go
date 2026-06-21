/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package agent

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// DataSourceSchema defines the schema for the QEMU guest agent datasource block.
func DataSourceSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		Description: "QEMU guest agent configuration.",
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				Description: "Whether the QEMU guest agent is enabled.",
				Computed:    true,
			},
			"trim": schema.BoolAttribute{
				Description: "Whether fstrim runs after cloning or moving a disk.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Guest agent channel type.",
				Computed:    true,
			},
		},
	}
}
