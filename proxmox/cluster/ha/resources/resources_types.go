/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resources

import (
	types2 "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// HAResourceListResponseBody contains the body from a HA resource list response.
type HAResourceListResponseBody struct {
	Data []*HAResourceListResponseData `json:"data,omitempty"`
}

// HAResourceListResponseData contains the data from a HA resource list response.
type HAResourceListResponseData struct {
	ID types2.HAResourceID `json:"sid"`
}

// HAResourceGetResponseBody contains the body from a HA resource get response.
type HAResourceGetResponseBody struct {
	Data *HAResourceGetResponseData `json:"data,omitempty"`
}

// HAResourceDataBase contains data common to all HA resource API calls.
type HAResourceDataBase struct {
	// Resource comment, if defined
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	// HA group identifier, if the resource is part of one.
	Group *string `json:"group,omitempty" url:"group,omitempty"`
	// Maximal number of service relocation attempts.
	MaxRelocate *int64 `json:"max_relocate,omitempty" url:"max_relocate,omitempty"`
	// Maximal number of service restart attempts.
	MaxRestart *int64 `json:"max_restart" url:"max_restart,omitempty"`
	// Requested resource state.
	State types2.HAResourceState `json:"state" url:"state"`
}

// HAResourceGetResponseData contains data received from the HA resource API when requesting information about a single
// HA resource.
type HAResourceGetResponseData struct {
	HAResourceDataBase
	// Identifier of this resource
	ID types2.HAResourceID `json:"sid"`
	// Type of this resource
	Type types2.HAResourceType `json:"type"`
	// SHA-1 digest of the resources' configuration.
	Digest *string `json:"digest,omitempty"`
}

// HAResourceCreateRequestBody contains data received from the HA resource API when creating a new HA resource.
type HAResourceCreateRequestBody struct {
	HAResourceDataBase
	// Identifier of this resource
	ID types2.HAResourceID `url:"sid"`
	// Type of this resource
	Type *types2.HAResourceType `url:"type,omitempty"`
	// SHA-1 digest of the resources' configuration.
	Digest *string `url:"comment,omitempty"`
}

// HAResourceUpdateRequestBody contains data received from the HA resource API when updating an existing HA resource.
type HAResourceUpdateRequestBody struct {
	HAResourceDataBase
	// Settings that must be deleted from the resource's configuration
	Delete []string `url:"delete,omitempty,comma"`
}
