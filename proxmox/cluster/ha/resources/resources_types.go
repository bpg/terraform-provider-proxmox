/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resources

import "github.com/bpg/terraform-provider-proxmox/internal/types"

// HAResourceListResponseBody contains the body from a HA resource list response.
type HAResourceListResponseBody struct {
	Data []*HAResourceListResponseData `json:"data,omitempty"`
}

// HAResourceListResponseData contains the data from a HA resource list response.
type HAResourceListResponseData struct {
	ID types.HAResourceID `json:"sid"`
}

// HAResourceGetResponseBody contains the body from a HA resource get response.
type HAResourceGetResponseBody struct {
	Data *HAResourceGetResponseData `json:"data,omitempty"`
}

// HAResourceGetResponseData contains data received from the HA resource API when requesting information about a single
// HA resource.
type HAResourceGetResponseData struct {
	ID          types.HAResourceID    `json:"sid"`               // Identifier of this resource
	Type        types.HAResourceType  `json:"type"`              // Type of this resource
	Comment     *string               `json:"comment,omitempty"` // Resource comment, if defined
	Digest      *string               `json:"digest,omitempty"`  // SHA-1 digest of the resources' configuration.
	Group       *string               `json:"group,omitempty"`   // HA group identifier, if the resource is part of one.
	MaxRelocate *int64                `json:"max_relocate"`      // Maximal number of service relocation attempts.
	MaxRestart  *int64                `json:"max_restart"`       // Maximal number of service restart attempts.
	State       types.HAResourceState `json:"state"`             // Requested resource state.
}
