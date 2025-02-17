/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/vm/cpu"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/vm/rng"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/vm/vga"
)

// Schema defines the schema for the resource.
func (d *Datasource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "This is an experimental implementation of a Proxmox VM datasource using Plugin Framework.",
		Attributes: map[string]schema.Attribute{
			"clone": schema.SingleNestedAttribute{
				Description: "The cloning configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "The ID of the VM to clone.",
						Required:    true,
					},
					"retries": schema.Int64Attribute{
						Description: "The number of retries to perform when cloning the VM (default: 3).",
						Optional:    true,
						Computed:    true,
					},
				},
			},
			"cpu": cpu.DataSourceSchema(),
			"description": schema.StringAttribute{
				Description: "The description of the VM.",
				Optional:    true,
			},
			"id": schema.Int64Attribute{
				Required:    true,
				Description: "The unique identifier of the VM in the Proxmox cluster.",
			},
			"name": schema.StringAttribute{
				Description: "The name of the VM.",
				Optional:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the node where the VM is provisioned.",
				Required:    true,
			},
			"rng":  rng.DataSourceSchema(),
			"tags": stringset.ResourceAttribute("The tags assigned to the VM.", ""),
			"template": schema.BoolAttribute{
				Description: "Whether the VM is a template.",
				Optional:    true,
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"vga": vga.DataSourceSchema(),
		},
	}
}
