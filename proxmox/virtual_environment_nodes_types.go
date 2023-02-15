/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"encoding/json"
	"net/url"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// CustomNodeCommands contains an array of commands to execute.
type CustomNodeCommands []string

// VirtualEnvironmentNodeExecuteRequestBody contains the data for a node execute request.
type VirtualEnvironmentNodeExecuteRequestBody struct {
	Commands CustomNodeCommands `json:"commands" url:"commands"`
}

// VirtualEnvironmentNodeGetTimeResponseBody contains the body from a node time zone get response.
type VirtualEnvironmentNodeGetTimeResponseBody struct {
	Data *VirtualEnvironmentNodeGetTimeResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentNodeGetTimeResponseData contains the data from a node list response.
type VirtualEnvironmentNodeGetTimeResponseData struct {
	LocalTime types.CustomTimestamp `json:"localtime"`
	TimeZone  string                `json:"timezone"`
	UTCTime   types.CustomTimestamp `json:"time"`
}

// VirtualEnvironmentNodeGetTaskStatusResponseBody contains the body from a node get task status response.
type VirtualEnvironmentNodeGetTaskStatusResponseBody struct {
	Data *VirtualEnvironmentNodeGetTaskStatusResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentNodeGetTaskStatusResponseData contains the data from a node get task status response.
type VirtualEnvironmentNodeGetTaskStatusResponseData struct {
	PID      int    `json:"pid,omitempty"`
	Status   string `json:"status,omitempty"`
	ExitCode string `json:"exitstatus,omitempty"`
}

// VirtualEnvironmentNodeListResponseBody contains the body from a node list response.
type VirtualEnvironmentNodeListResponseBody struct {
	Data []*VirtualEnvironmentNodeListResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentNodeListResponseData contains the data from a node list response.
type VirtualEnvironmentNodeListResponseData struct {
	CPUCount        *int     `json:"maxcpu,omitempty"`
	CPUUtilization  *float64 `json:"cpu,omitempty"`
	MemoryAvailable *int     `json:"maxmem,omitempty"`
	MemoryUsed      *int     `json:"mem,omitempty"`
	Name            string   `json:"node"`
	SSLFingerprint  *string  `json:"ssl_fingerprint,omitempty"`
	Status          *string  `json:"status"`
	SupportLevel    *string  `json:"level,omitempty"`
	Uptime          *int     `json:"uptime"`
}

// VirtualEnvironmentNodeNetworkDeviceListResponseBody contains the body from a node network device list response.
type VirtualEnvironmentNodeNetworkDeviceListResponseBody struct {
	Data []*VirtualEnvironmentNodeNetworkDeviceListResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentNodeNetworkDeviceListResponseData contains the data from a node network device list response.
type VirtualEnvironmentNodeNetworkDeviceListResponseData struct {
	Active      *types.CustomBool `json:"active,omitempty"`
	Address     *string           `json:"address,omitempty"`
	Autostart   *types.CustomBool `json:"autostart,omitempty"`
	BridgeFD    *string           `json:"bridge_fd,omitempty"`
	BridgePorts *string           `json:"bridge_ports,omitempty"`
	BridgeSTP   *string           `json:"bridge_stp,omitempty"`
	CIDR        *string           `json:"cidr,omitempty"`
	Exists      *types.CustomBool `json:"exists,omitempty"`
	Families    *[]string         `json:"families,omitempty"`
	Gateway     *string           `json:"gateway,omitempty"`
	Iface       string            `json:"iface"`
	MethodIPv4  *string           `json:"method,omitempty"`
	MethodIPv6  *string           `json:"method6,omitempty"`
	Netmask     *string           `json:"netmask,omitempty"`
	Priority    int               `json:"priority"`
	Type        string            `json:"type"`
}

// VirtualEnvironmentNodeUpdateTimeRequestBody contains the body for a node time update request.
type VirtualEnvironmentNodeUpdateTimeRequestBody struct {
	TimeZone string `json:"timezone" url:"timezone"`
}

// EncodeValues converts a CustomNodeCommands array to a JSON encoded URL vlaue.
func (r CustomNodeCommands) EncodeValues(key string, v *url.Values) error {
	jsonArrayBytes, err := json.Marshal(r)
	if err != nil {
		return err
	}

	v.Add(key, string(jsonArrayBytes))

	return nil
}
