/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package rng

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// DataSourceSchema defines the schema for the RNG datasource.
func DataSourceSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		CustomType: basetypes.ObjectType{
			AttrTypes: attributeTypes(),
		},
		Description: "The RNG (Random Number Generator) configuration.",
		Optional:    true,
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"source": schema.StringAttribute{
				Description: "The entropy source for the RNG device.",
				Optional:    true,
				Computed:    true,
			},
			"max_bytes": schema.Int64Attribute{
				Description: "Maximum bytes of entropy allowed to get injected into the guest every period.",
				Optional:    true,
				Computed:    true,
			},
			"period": schema.Int64Attribute{
				Description: "Period in milliseconds to limit entropy injection to the guest.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}
