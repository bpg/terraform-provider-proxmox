/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type DatastoreGetRequest struct {
	ID *string `json:"storage" url:"storage"`
}

type DatastoreGetResponseBody struct {
	Data *DatastoreGetResponseData `json:"data,omitempty"`
}

// DatastoreListRequest contains the body for a datastore list request.
type DatastoreListRequest struct {
	Type *string `json:"type,omitempty" url:"type,omitempty,omitempty"`
}

// DatastoreListResponse contains the body from a datastore list response.
type DatastoreListResponse struct {
	Data []*DatastoreGetResponseData `json:"data,omitempty"`
}

type DatastoreGetResponseData struct {
	ID           *string                         `json:"storage" url:"storage"`
	Type         *string                         `json:"type" url:"type"`
	ContentTypes *types.CustomCommaSeparatedList `json:"content,omitempty" url:"content,omitempty,comma"`
	Path         *string                         `json:"path,omitempty" url:"path,omitempty"`
	Nodes        *types.CustomCommaSeparatedList `json:"nodes,omitempty" url:"nodes,omitempty,comma"`
	Disable      *types.CustomBool               `json:"disable,omitempty" url:"disable,omitempty,int"`
	Shared       *types.CustomBool               `json:"shared,omitempty" url:"shared,omitempty,int"`
}

type DataStoreCommonImmutableFields struct {
	ID   *string `json:"storage" url:"storage"`
	Type *string `json:"type,omitempty" url:"type,omitempty"`
}

type DataStoreCommonMutableFields struct {
	ContentTypes *types.CustomCommaSeparatedList `json:"content,omitempty" url:"content,omitempty,comma"`
	Nodes        *types.CustomCommaSeparatedList `json:"nodes,omitempty" url:"nodes,omitempty,comma"`
	Disable      *types.CustomBool               `json:"disable,omitempty" url:"disable,omitempty,int"`
	Shared       *types.CustomBool               `json:"shared,omitempty" url:"shared,omitempty,int"`
}

// DataStoreWithBackups holds optional retention settings for backups.
type DataStoreWithBackups struct {
	MaxProtectedBackups *types.CustomInt64 `json:"max-protected-backups,omitempty" url:"max,omitempty"`
	KeepDaily           *int               `json:"-"`
	KeepHourly          *int               `json:"-"`
	KeepLast            *int               `json:"-"`
	KeepMonthly         *int               `json:"-"`
	KeepWeekly          *int               `json:"-"`
	KeepYearly          *int               `json:"-"`
}

// String serializes DataStoreWithBackups into the Proxmox "key=value,key=value" format.
// Only defined (non-nil) fields will be included.
func (b DataStoreWithBackups) String() string {
	var parts []string

	if b.KeepLast != nil {
		parts = append(parts, fmt.Sprintf("keep-last=%d", *b.KeepLast))
	}
	if b.KeepHourly != nil {
		parts = append(parts, fmt.Sprintf("keep-hourly=%d", *b.KeepHourly))
	}
	if b.KeepDaily != nil {
		parts = append(parts, fmt.Sprintf("keep-daily=%d", *b.KeepDaily))
	}
	if b.KeepWeekly != nil {
		parts = append(parts, fmt.Sprintf("keep-weekly=%d", *b.KeepWeekly))
	}
	if b.KeepMonthly != nil {
		parts = append(parts, fmt.Sprintf("keep-monthly=%d", *b.KeepMonthly))
	}
	if b.KeepYearly != nil {
		parts = append(parts, fmt.Sprintf("keep-yearly=%d", *b.KeepYearly))
	}

	return strings.Join(parts, ",")
}

// MarshalJSON ensures DataStoreWithBackups is encoded into a JSON field "prune-backups".
func (b DataStoreWithBackups) MarshalJSON() ([]byte, error) {
	str := b.String()

	// Special case; nothing defined so we omit the field
	if str == "" && b.MaxProtectedBackups == nil {
		return []byte(`{}`), nil
	}

	type Alias DataStoreWithBackups
	aux := struct {
		*Alias
		PruneBackups string `json:"prune-backups,omitempty" url:"prune,omitempty"`
	}{
		Alias:        (*Alias)(&b),
		PruneBackups: str,
	}
	return json.Marshal(aux)
}

// DirectoryStorageMutableFields defines the mutable attributes for 'dir' type storage.
type DirectoryStorageMutableFields struct {
	DataStoreCommonMutableFields
	Preallocation          *string          `json:"preallocation,omitempty" url:"preallocation,omitempty"`
	SnapshotsAsVolumeChain types.CustomBool `json:"snapshot-as-volume-chain,omitempty" url:"snapshot-as-volume-chain,omitempty"`
}

// DirectoryStorageImmutableFields defines the immutable attributes for 'dir' type storage.
type DirectoryStorageImmutableFields struct {
	Path *string `json:"path,omitempty" url:"path,omitempty"`
}

// DirectoryStorageCreateRequest defines options for 'dir' type storage.
type DirectoryStorageCreateRequest struct {
	DataStoreCommonImmutableFields
	DirectoryStorageMutableFields
	DirectoryStorageImmutableFields
}

type DirectoryStorageUpdateRequest struct {
	DirectoryStorageMutableFields
}
