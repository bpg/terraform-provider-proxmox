/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cdrom

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// DataSourceSchema defines the schema for the CD-ROM datasource.
func DataSourceSchema() schema.Attribute {
	return schema.MapNestedAttribute{
		Description: "The CD-ROM configuration.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"file_id": schema.StringAttribute{
					Description: "The file ID of the CD-ROM.",
					Computed:    true,
				},
			},
		},
	}
}
