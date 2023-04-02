/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datasource

import (
	"context"
	"math"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	mkDataSourceVirtualEnvironmentNodesCPUCount        = "cpu_count"
	mkDataSourceVirtualEnvironmentNodesCPUUtilization  = "cpu_utilization"
	mkDataSourceVirtualEnvironmentNodesMemoryAvailable = "memory_available"
	mkDataSourceVirtualEnvironmentNodesMemoryUsed      = "memory_used"
	mkDataSourceVirtualEnvironmentNodesNames           = "names"
	mkDataSourceVirtualEnvironmentNodesOnline          = "online"
	mkDataSourceVirtualEnvironmentNodesSSLFingerprints = "ssl_fingerprints"
	mkDataSourceVirtualEnvironmentNodesSupportLevels   = "support_levels"
	mkDataSourceVirtualEnvironmentNodesUptime          = "uptime"
)

func Nodes() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentNodesCPUCount: {
				Type:        schema.TypeList,
				Description: "The CPU count for each node",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			mkDataSourceVirtualEnvironmentNodesCPUUtilization: {
				Type:        schema.TypeList,
				Description: "The CPU utilization on each node",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeFloat},
			},
			mkDataSourceVirtualEnvironmentNodesMemoryAvailable: {
				Type:        schema.TypeList,
				Description: "The available memory in bytes on each node",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			mkDataSourceVirtualEnvironmentNodesMemoryUsed: {
				Type:        schema.TypeList,
				Description: "The used memory in bytes on each node",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			mkDataSourceVirtualEnvironmentNodesNames: {
				Type:        schema.TypeList,
				Description: "The node names",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceVirtualEnvironmentNodesOnline: {
				Type:        schema.TypeList,
				Description: "Whether a node is online",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeBool},
			},
			mkDataSourceVirtualEnvironmentNodesSSLFingerprints: {
				Type:        schema.TypeList,
				Description: "The SSL fingerprint for each node",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceVirtualEnvironmentNodesSupportLevels: {
				Type:        schema.TypeList,
				Description: "The support level for each node",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceVirtualEnvironmentNodesUptime: {
				Type:        schema.TypeList,
				Description: "The uptime in seconds for each node",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
		ReadContext: nodesRead,
	}
}

func nodesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	list, err := veClient.ListNodes(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	cpuCount := make([]interface{}, len(list))
	cpuUtilization := make([]interface{}, len(list))
	memoryAvailable := make([]interface{}, len(list))
	memoryUsed := make([]interface{}, len(list))
	name := make([]interface{}, len(list))
	online := make([]interface{}, len(list))
	sslFingerprints := make([]interface{}, len(list))
	supportLevels := make([]interface{}, len(list))
	uptime := make([]interface{}, len(list))

	for i, v := range list {
		if v.CPUCount != nil {
			cpuCount[i] = v.CPUCount
		} else {
			cpuCount[i] = 0
		}

		if v.CPUUtilization != nil {
			cpuUtilization[i] = math.Round(*v.CPUUtilization*100) / 100
		} else {
			cpuUtilization[i] = 0
		}

		if v.MemoryAvailable != nil {
			memoryAvailable[i] = v.MemoryAvailable
		} else {
			memoryAvailable[i] = 0
		}

		if v.MemoryUsed != nil {
			memoryUsed[i] = v.MemoryUsed
		} else {
			memoryUsed[i] = 0
		}

		name[i] = v.Name

		if v.Status != nil {
			online[i] = *v.Status == "online"
		} else {
			online[i] = false
		}

		if v.SSLFingerprint != nil {
			sslFingerprints[i] = v.SSLFingerprint
		} else {
			sslFingerprints[i] = 0
		}

		if v.SupportLevel != nil {
			supportLevels[i] = v.SupportLevel
		} else {
			supportLevels[i] = ""
		}

		if v.Uptime != nil {
			uptime[i] = v.Uptime
		} else {
			uptime[i] = 0
		}
	}

	d.SetId("nodes")

	err = d.Set(mkDataSourceVirtualEnvironmentNodesCPUCount, cpuCount)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentNodesCPUUtilization, cpuUtilization)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentNodesMemoryAvailable, memoryAvailable)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentNodesMemoryUsed, memoryUsed)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentNodesNames, name)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentNodesOnline, online)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentNodesSSLFingerprints, sslFingerprints)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentNodesSupportLevels, supportLevels)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentNodesUptime, uptime)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}
