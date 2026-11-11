/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cpumodels

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// ListResponseBody contains the body from a custom CPU model list response.
type ListResponseBody struct {
	Data []*GetResponseData `json:"data,omitempty"`
}

// GetResponseBody contains the body from a custom CPU model get response.
type GetResponseBody struct {
	Data *GetResponseData `json:"data,omitempty"`
}

// GetResponseData contains a custom CPU model definition.
type GetResponseData struct {
	// CPUType is the model name including the `custom-` prefix, as PVE always returns it.
	CPUType       string             `json:"cputype"`
	Digest        *string            `json:"digest,omitempty"`
	Flags         *string            `json:"flags,omitempty"`
	GuestPhysBits *types.CustomInt64 `json:"guest-phys-bits,omitempty"`
	Hidden        *types.CustomBool  `json:"hidden,omitempty"`
	HVVendorID    *string            `json:"hv-vendor-id,omitempty"`
	Level         *types.CustomInt64 `json:"level,omitempty"`
	PhysBits      *string            `json:"phys-bits,omitempty"`
	ReportedModel *string            `json:"reported-model,omitempty"`
}

// RequestBodyCommon contains fields shared by custom CPU model create and update requests.
type RequestBodyCommon struct {
	Flags         *string           `json:"flags,omitempty"           url:"flags,omitempty"`
	GuestPhysBits *int64            `json:"guest-phys-bits,omitempty" url:"guest-phys-bits,omitempty"`
	Hidden        *types.CustomBool `json:"hidden,omitempty"          url:"hidden,omitempty,int"`
	HVVendorID    *string           `json:"hv-vendor-id,omitempty"    url:"hv-vendor-id,omitempty"`
	Level         *int64            `json:"level,omitempty"           url:"level,omitempty"`
	PhysBits      *string           `json:"phys-bits,omitempty"       url:"phys-bits,omitempty"`
	ReportedModel *string           `json:"reported-model,omitempty"  url:"reported-model,omitempty"`
}

// CreateRequestBody contains the body for creating a custom CPU model.
type CreateRequestBody struct {
	RequestBodyCommon

	CPUType string `json:"cputype" url:"cputype"`
}

// UpdateRequestBody contains the body for updating a custom CPU model.
type UpdateRequestBody struct {
	RequestBodyCommon

	Delete []string `json:"delete,omitempty" url:"delete,omitempty,comma"`
}
