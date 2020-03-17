/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"math"

	"github.com/hashicorp/terraform/helper/schema"
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

func dataSourceVirtualEnvironmentNodes() *schema.Resource {
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
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Read: dataSourceVirtualEnvironmentNodesRead,
	}
}

func dataSourceVirtualEnvironmentNodesRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	list, err := veClient.ListNodes()

	if err != nil {
		return err
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

	d.Set(mkDataSourceVirtualEnvironmentNodesCPUCount, cpuCount)
	d.Set(mkDataSourceVirtualEnvironmentNodesCPUUtilization, cpuUtilization)
	d.Set(mkDataSourceVirtualEnvironmentNodesMemoryAvailable, memoryAvailable)
	d.Set(mkDataSourceVirtualEnvironmentNodesMemoryUsed, memoryUsed)
	d.Set(mkDataSourceVirtualEnvironmentNodesNames, name)
	d.Set(mkDataSourceVirtualEnvironmentNodesOnline, online)
	d.Set(mkDataSourceVirtualEnvironmentNodesSSLFingerprints, sslFingerprints)
	d.Set(mkDataSourceVirtualEnvironmentNodesSupportLevels, supportLevels)
	d.Set(mkDataSourceVirtualEnvironmentNodesUptime, uptime)

	return nil
}
