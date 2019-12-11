/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"errors"
	"sort"
)

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

// ListNodes retrieves a list of nodes.
func (c *VirtualEnvironmentClient) ListNodes() ([]*VirtualEnvironmentNodeListResponseData, error) {
	resBody := &VirtualEnvironmentNodeListResponseBody{}
	err := c.DoRequest(hmGET, "nodes", nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Name < resBody.Data[j].Name
	})

	return resBody.Data, nil
}
