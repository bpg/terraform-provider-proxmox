/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package files

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Schema defines the schema for the files data source.
func (d *Datasource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a list of files available in a datastore on a specific Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"node_name": schema.StringAttribute{
				Description: "The name of the node.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"datastore_id": schema.StringAttribute{
				Description: "The identifier of the datastore.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"content_type": schema.StringAttribute{
				Description: "The content type to filter by. When set, only files of this type " +
					"are returned. Valid values are `backup`, `images`, `import`, `iso`, " +
					"`rootdir`, `snippets`, `vztmpl`.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"backup",
						"images",
						"import",
						"iso",
						"rootdir",
						"snippets",
						"vztmpl",
					),
				},
			},
			"files": schema.ListNestedAttribute{
				Description: "The list of files in the datastore.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the file (volume ID), " +
								"e.g. `local:iso/ubuntu.iso`.",
							Computed: true,
						},
						"content_type": schema.StringAttribute{
							Description: "The content type of the file.",
							Computed:    true,
						},
						"file_name": schema.StringAttribute{
							Description: "The name of the file.",
							Computed:    true,
						},
						"file_format": schema.StringAttribute{
							Description: "The format of the file.",
							Computed:    true,
						},
						"file_size": schema.Int64Attribute{
							Description: "The size of the file in bytes.",
							Computed:    true,
						},
						"vmid": schema.Int64Attribute{
							Description: "The VM ID associated with the file, if applicable.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
