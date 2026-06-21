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

	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/agent"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/cdrom"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/cpu"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/initialization"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/memory"
	network_device "github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/network_device"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/rng"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/vga"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
)

// Schema defines the schema for the resource.
func (d *Datasource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		DeprecationMessage: migration.DeprecationMessage("proxmox_vm"),
		Description:        "Retrieves information about a specific VM.",
		Attributes: map[string]schema.Attribute{
			"agent":          agent.DataSourceSchema(),
			"cdrom":          cdrom.DataSourceSchema(),
			"cpu":            cpu.DataSourceSchema(),
			"initialization": initialization.DataSourceSchema(),
			"memory":         memory.DataSourceSchema(),
			"network_device": network_device.DataSourceSchema(),
			"description": schema.StringAttribute{
				Description: "The description of the VM.",
				Computed:    true,
			},
			"id": schema.Int64Attribute{
				Required:    true,
				Description: "The unique identifier of the VM in the Proxmox cluster.",
			},
			"name": schema.StringAttribute{
				Description: "The name of the VM.",
				Computed:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the node where the VM is provisioned.",
				Required:    true,
			},
			"rng": rng.DataSourceSchema(),
			"started": schema.BoolAttribute{
				Description: "Whether the VM is currently running.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "The status of the VM (e.g., `running`, `stopped`).",
				Computed:    true,
			},
			"tags": stringset.DataSourceAttribute("The tags assigned to the VM.", ""),
			"template": schema.BoolAttribute{
				Description: "Whether the VM is a template.",
				Computed:    true,
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"vga": vga.DataSourceSchema(),
		},
	}
}
