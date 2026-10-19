/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

// RBDStorageMutableFields defines the updatable options for 'rbd' type storage.
type RBDStorageMutableFields struct {
	DataStoreCommonMutableFields

	KRBD *types.CustomBool `json:"krbd,omitempty" url:"krbd,omitempty,int"`
}

// RBDStorageImmutableFields defines the options for 'rbd' type storage that are only set on creation.
type RBDStorageImmutableFields struct {
	Pool      *string `json:"pool"                url:"pool"`
	DataPool  *string `json:"data-pool,omitempty" url:"data-pool,omitempty"`
	Namespace *string `json:"namespace,omitempty" url:"namespace,omitempty"`
}

// RBDStorageCreateRequest defines the request body for creating a new RBD storage.
type RBDStorageCreateRequest struct {
	DataStoreCommonImmutableFields
	RBDStorageMutableFields
	RBDStorageImmutableFields
}

// RBDStorageUpdateRequest defines the request body for updating an existing RBD storage.
type RBDStorageUpdateRequest struct {
	RBDStorageMutableFields
}
