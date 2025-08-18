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

type DatastoreGetResponse struct {
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
	ID            *string                         `json:"storage" url:"storage"`
	Type          *string                         `json:"type" url:"type"`
	ContentTypes  *types.CustomCommaSeparatedList `json:"content,omitempty" url:"content,omitempty,comma"`
	Path          *string                         `json:"path,omitempty" url:"path,omitempty"`
	Nodes         *types.CustomCommaSeparatedList `json:"nodes,omitempty" url:"nodes,omitempty,comma"`
	Disable       *types.CustomBool               `json:"disable,omitempty" url:"disable,omitempty,int"`
	Shared        *types.CustomBool               `json:"shared,omitempty" url:"shared,omitempty,int"`
	Server        *string                         `json:"server,omitempty" url:"server,omitempty"`
	Export        *string                         `json:"export,omitempty" url:"export,omitempty"`
	Options       *string                         `json:"options,omitempty" url:"options,omitempty"`
	Preallocation *string                         `json:"preallocation,omitempty" url:"preallocation,omitempty"`
	Datastore     *string                         `json:"datastore,omitempty" url:"datastore,omitempty"`
	Username      *string                         `json:"username,omitempty" url:"username,omitempty"`
	Password      *string                         `json:"password,omitempty" url:"password,omitempty"`
	Namespace     *string                         `json:"namespace,omitempty" url:"namespace,omitempty"`
	Fingerprint   *string                         `json:"fingerprint,omitempty" url:"fingerprint,omitempty"`
	EncryptionKey *string                         `json:"keyring,omitempty" url:"keyring,omitempty"`
	ZFSPool       *string                         `json:"pool,omitempty" url:"pool,omitempty"`
	ThinProvision *types.CustomBool               `json:"sparse,omitempty" url:"sparse,omitempty"`
	Blocksize     *string                         `json:"blocksize,omitempty" url:"blocksize,omitempty"`
}

type DatastoreCreateResponse struct {
	Data *DatastoreCreateResponseData `json:"data,omitempty" url:"data,omitempty"`
}

type DatastoreCreateResponseData struct {
	Type    *string                           `json:"type" url:"type"`
	Storage *string                           `json:"storage,omitempty" url:"storage,omitempty"`
	Config  DatastoreCreateResponseConfigData `json:"config,omitempty" url:"config,omitempty"`
}

type DatastoreCreateResponseConfigData struct {
	EncryptionKey *string `json:"encryption-key,omitempty" url:"encryption-key,omitempty"`
}

type DataStoreCommonImmutableFields struct {
	ID   *string `json:"storage" url:"storage"`
	Type *string `json:"type,omitempty" url:"type,omitempty"`
}

type DataStoreCommonMutableFields struct {
	ContentTypes *types.CustomCommaSeparatedList `json:"content,omitempty" url:"content,omitempty,comma"`
	Nodes        *types.CustomCommaSeparatedList `json:"nodes,omitempty" url:"nodes,omitempty,comma"`
	Disable      *types.CustomBool               `json:"disable,omitempty" url:"disable,omitempty,int"`
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
