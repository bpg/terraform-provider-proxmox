/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cluster

import (
	"github.com/bpg/terraform-provider-proxmox/internal/types"
)

// NextIDRequestBody contains the data for a cluster next id request.
type NextIDRequestBody struct {
	VMID *int `json:"vmid,omitempty" url:"vmid,omitempty"`
}

// NextIDResponseBody contains the body from a cluster next id response.
type NextIDResponseBody struct {
	Data *types.CustomInt `json:"data,omitempty"`
}

// ClusterResourcesListBody contains the body from a cluste resource list response.
type ClusterResourcesListBody struct {
	Data []*ClusterResourcesListResponseData `json:"data,omitempty"`
}

type ClusterResourcesListRequestBody struct {
	Type string `json:"type"  url:"type"`
}

// ClusterResourcesListResponseData contains the data from a cluster resource list body response
type ClusterResourcesListResponseData struct {
	Type       string  `json:"type"`
	ID         string  `json:"id"`
	CgroupMode int     `json:"cgroup-mode,omitempty"`
	Content    int     `json:"content,omitempty"`
	Cpu        float64 `json:"cpu,omitempty"`
	Disk       int     `json:"disk,omitempty"`
	HaState    string  `json:"hastate,omitempty"`
	Level      string  `json:"level,omitempty"`
	MaxCpu     float64 `json:"maxcpu,omitempty"`
	MaxDisk    int     `json:"maxdisk,omitempty"`
	MaxMem     int     `json:"maxmem,omitempty"`
	Mem        int     `json:"mem,omitempty"`
	Name       string  `json:"name,omitempty"`
	NodeName   string  `json:"node,omitempty"`
	PluginType string  `json:"plugintype,omitempty"`
	PoolName   string  `json:"poolname,omitempty"`
	Status     string  `json:"status,omitempty"`
	Storage    string  `json:"storage,omitempty"`
	Uptime     int     `json:"uptime,omitempty"`
	VMID       int     `json:"vmid,omitempty"`
}
