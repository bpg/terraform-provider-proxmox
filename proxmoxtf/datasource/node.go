/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datasource

import (
	"context"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	mkDataSourceVirtualEnvironmentNodeCPUCores        = "cpu_count"
	mkDataSourceVirtualEnvironmentNodeCPUSockets      = "cpu_sockets"
	mkDataSourceVirtualEnvironmentNodeCPUModel        = "cpu_model"
	mkDataSourceVirtualEnvironmentNodeMemoryAvailable = "memory_available"
	mkDataSourceVirtualEnvironmentNodeMemoryUsed      = "memory_used"
	mkDataSourceVirtualEnvironmentNodeMemoryTotal     = "memory_total"
	mkDataSourceVirtualEnvironmentNodeUptime          = "uptime"
	mkDataSourceVirtualEnvironmentNodeName            = "node_name"
)

// Node returns a resource for the Proxmox node.
func Node() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentNodeCPUCores: {
				Type:        schema.TypeInt,
				Description: "The CPU count on the node",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentNodeCPUSockets: {
				Type:        schema.TypeInt,
				Description: "The CPU sockets on the node",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentNodeCPUModel: {
				Type:        schema.TypeString,
				Description: "The CPU model on the node",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentNodeMemoryAvailable: {
				Type:        schema.TypeInt,
				Description: "The available memory in bytes on the node",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentNodeMemoryUsed: {
				Type:        schema.TypeInt,
				Description: "The used memory in bytes on the node",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentNodeMemoryTotal: {
				Type:        schema.TypeInt,
				Description: "The total memory in bytes on the node",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentNodeUptime: {
				Type:        schema.TypeInt,
				Description: "The uptime in seconds on the node",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentNodeName: {
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
			},
		},
		ReadContext: nodeRead,
	}
}

func nodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)

	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeID := d.Get(mkDataSourceVirtualEnvironmentNodeName).(string)

	node, err := api.Node(nodeID).GetInfo(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(nodeID)

	if node.CPUInfo.CPUCores != nil {
		err = d.Set(mkDataSourceVirtualEnvironmentNodeCPUCores, *node.CPUInfo.CPUCores)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentNodeCPUCores, 1)
	}

	diags = append(diags, diag.FromErr(err)...)

	if node.CPUInfo.CPUSockets != nil {
		err = d.Set(mkDataSourceVirtualEnvironmentNodeCPUSockets, *node.CPUInfo.CPUSockets)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentNodeCPUSockets, 1)
	}

	diags = append(diags, diag.FromErr(err)...)

	if node.CPUInfo.CPUModel != nil {
		err = d.Set(mkDataSourceVirtualEnvironmentNodeCPUModel, *node.CPUInfo.CPUModel)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentNodeCPUModel, "")
	}

	diags = append(diags, diag.FromErr(err)...)

	if node.MemoryInfo.Total != nil {
		err = d.Set(mkDataSourceVirtualEnvironmentNodeMemoryAvailable, node.MemoryInfo.Free)
		diags = append(diags, diag.FromErr(err)...)
		err = d.Set(mkDataSourceVirtualEnvironmentNodeMemoryUsed, node.MemoryInfo.Used)
		diags = append(diags, diag.FromErr(err)...)
		err = d.Set(mkDataSourceVirtualEnvironmentNodeMemoryTotal, node.MemoryInfo.Total)
		diags = append(diags, diag.FromErr(err)...)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentNodeMemoryAvailable, 0)
		diags = append(diags, diag.FromErr(err)...)
		err = d.Set(mkDataSourceVirtualEnvironmentNodeMemoryUsed, 0)
		diags = append(diags, diag.FromErr(err)...)
		err = d.Set(mkDataSourceVirtualEnvironmentNodeMemoryTotal, 0)
		diags = append(diags, diag.FromErr(err)...)
	}

	if node.Uptime != nil {
		err = d.Set(mkDataSourceVirtualEnvironmentNodeUptime, *node.Uptime)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentNodeUptime, 0)
	}

	diags = append(diags, diag.FromErr(err)...)

	return diags
}
