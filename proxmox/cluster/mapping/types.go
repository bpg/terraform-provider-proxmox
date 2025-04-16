/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package mapping

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types/hardwaremapping"
)

const (
	// APIParamNamePCIMediatedDevices is the API attribute name of the Proxmox VE API "mediated devices" parameter for a
	// PCI hardware mapping.
	APIParamNamePCIMediatedDevices = "mdev"
)

type listQuery struct {
	// CheckNode is the name of the node those configuration should be checked for correctness.
	CheckNode string `url:"check-node,omitempty"`
}

// DataBase contains common data for hardware mapping API calls.
type DataBase struct {
	// Description is the optional key for the description for a hardware mapping that is omitted by the Proxmox VE API
	// when not set.
	// Note that even though the Proxmox VE API attribute is named "description" it is generally labeled as "comment"
	// cross the Proxmox VE web UI while only being named "description" in the Proxmox VE API and its documentation.
	Description *string `url:"description,omitempty"`

	// Map is the list of device mappings.
	Map []proxmoxtypes.Map `json:"map" url:"map"`

	// MediatedDevices is the indicator for the optional APIParamNamePCIMediatedDevices parameter.
	MediatedDevices types.CustomBool `json:"mdev" url:"mdev,omitempty,int"`
}

// CreateRequestBody contains the data which must be sent when creating a hardware mapping.
type CreateRequestBody struct {
	DataBase

	// ID is the hardware mappings identifier.
	ID string `url:"id"`
}

// GetResponseBody contains the body from a hardware mapping get response.
type GetResponseBody struct {
	// Data is the hardware mapping get response data.
	Data *GetResponseData `json:"data,omitempty"`
}

// ListResponseBody contains the body from a hardware mapping list response.
type ListResponseBody struct {
	// Data is the hardware mapping list response data.
	Data []*ListResponseData `json:"data,omitempty"`
}

// GetResponseData contains data received from the hardware mapping API when requesting information about a single
// mapping.
type GetResponseData struct {
	DataBase

	// Type is the required types of the hardware mapping.
	Type proxmoxtypes.Type `json:"type"`
}

// ListResponseData contains the data from a hardware mapping list response.
type ListResponseData struct {
	DataBase

	// Checks might contain relevant diagnostics about incorrect [typesHWM.TypePCI] configurations.
	// The name of the node must be passed to the Proxmox VE API call which maps to the "check-node" URL parameter.
	// Note that the Proxmox VE API, for whatever reason, only returns one error at a time, even though the field is an
	// array.
	Checks []NodeCheckDiag `json:"checks,omitempty"`

	// ChecksUSB might contain relevant diagnostics about incorrect [typesHWM.TypeUSB] configurations.
	// The name of the node must be passed to the Proxmox VE API call which maps to the "check-node" URL parameter.
	// Note that the actual JSON field name matches the Proxmox VE API, but the name of this variable has been adjusted
	// for clarity.
	// Also note that the Proxmox VE API, for whatever reason, only returns one error at a time, even though the field is
	// an array.
	ChecksUSB []NodeCheckDiag `json:"errors,omitempty"`

	// ID is the hardware mappings identifier.
	ID string `json:"id"`

	// Type is the required types of the hardware mapping.
	Type proxmoxtypes.Type `json:"type"`
}

// NodeCheckDiag is a hardware mapping configuration correctness diagnostic entry.
type NodeCheckDiag struct {
	// Message is the message of the node check diagnostic entry.
	Message *string `json:"message"`

	// Severity is the severity of the node check diagnostic entry.
	Severity *string `json:"severity"`
}

// UpdateRequestBody contains data received from the hardware mapping resource API when updating an existing hardware
// mapping resource.
type UpdateRequestBody struct {
	DataBase

	// Delete are settings that must be deleted from the resource's configuration.
	Delete []string `url:"delete,omitempty,comma"`
}
