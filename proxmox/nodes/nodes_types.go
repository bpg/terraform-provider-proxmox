/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/bpg/terraform-provider-proxmox/internal/types"
)

// CustomCommands contains an array of commands to execute.
type CustomCommands []string

// ExecuteRequestBody contains the data for a node execute request.
type ExecuteRequestBody struct {
	Commands CustomCommands `json:"commands" url:"commands"`
}

// GetTimeResponseBody contains the body from a node time zone get response.
type GetTimeResponseBody struct {
	Data *GetTimeResponseData `json:"data,omitempty"`
}

// GetTimeResponseData contains the data from a node list response.
type GetTimeResponseData struct {
	LocalTime types.CustomTimestamp `json:"localtime"`
	TimeZone  string                `json:"timezone"`
	UTCTime   types.CustomTimestamp `json:"time"`
}

// ListResponseBody contains the body from a node list response.
type ListResponseBody struct {
	Data []*ListResponseData `json:"data,omitempty"`
}

// ListResponseData contains the data from a node list response.
type ListResponseData struct {
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

// UpdateTimeRequestBody contains the body for a node time update request.
type UpdateTimeRequestBody struct {
	TimeZone string `json:"timezone" url:"timezone"`
}

// EncodeValues converts a CustomCommands array to a JSON encoded URL value.
func (r CustomCommands) EncodeValues(key string, v *url.Values) error {
	jsonArrayBytes, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("error marshalling CustomCommands array: %w", err)
	}

	v.Add(key, string(jsonArrayBytes))

	return nil
}
