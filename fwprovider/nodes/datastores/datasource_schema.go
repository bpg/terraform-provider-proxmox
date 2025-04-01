/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datastores

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
)

// Schema defines the schema for the resource.
func (d *Datasource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about all the datastores available to a specific node.",
		Attributes: map[string]schema.Attribute{
			"node_name": schema.StringAttribute{
				Description: "The name of the node to retrieve the stores from.",
				Required:    true,
			},
			"filters": schema.SingleNestedAttribute{
				Description: "The filters to apply to the stores.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"content_types": stringset.DataSourceAttribute("Only list stores with the given content types.", "", true),
					"id": schema.StringAttribute{
						Description: "Only list stores with the given ID.",
						Optional:    true,
					},
					"target": schema.StringAttribute{
						Description: "If `target` is different to `node_name`, then only lists shared stores which " +
							"content is accessible on this node and the specified `target` node.",
						Optional: true,
					},
				},
			},
			"datastores": schema.ListNestedAttribute{
				Description: "The list of datastores.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"active": schema.BoolAttribute{
							Description: "Whether the store is active.",
							Optional:    true,
						},
						"content_types": stringset.DataSourceAttribute("Allowed store content types.", "", false),
						"enabled": schema.BoolAttribute{
							Description: "Whether the store is enabled.",
							Optional:    true,
						},
						"id": schema.StringAttribute{
							Description: "The ID of the store.",
							Required:    true,
						},
						"node_name": schema.StringAttribute{
							Description: "The name of the node the store is on.",
							Required:    true,
						},
						"shared": schema.BoolAttribute{
							Description: "Shared flag from store configuration.",
							Optional:    true,
						},
						"space_available": schema.Int64Attribute{
							Description: "Available store space in bytes.",
							Optional:    true,
						},
						"space_total": schema.Int64Attribute{
							Description: "Total store space in bytes.",
							Optional:    true,
						},
						"space_used": schema.Int64Attribute{
							Description: "Used store space in bytes.",
							Optional:    true,
						},
						"space_used_fraction": schema.Float64Attribute{
							Description: "Used fraction (used/total).",
							Optional:    true,
						},
						"type": schema.StringAttribute{
							Description: "Store type.",
							Required:    true,
						},
					},
				},
			},
		},
	}
}
