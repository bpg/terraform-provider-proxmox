/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

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
