/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"fmt"
	"net/url"
	"strconv"
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
	ID                 *string                         `json:"storage" url:"storage"`
	Type               *string                         `json:"type" url:"type"`
	ContentTypes       *types.CustomCommaSeparatedList `json:"content,omitempty" url:"content,omitempty,comma"`
	Path               *string                         `json:"path,omitempty" url:"path,omitempty"`
	Nodes              *types.CustomCommaSeparatedList `json:"nodes,omitempty" url:"nodes,omitempty,comma"`
	Disable            *types.CustomBool               `json:"disable,omitempty" url:"disable,omitempty,int"`
	Shared             *types.CustomBool               `json:"shared,omitempty" url:"shared,omitempty,int"`
	Server             *string                         `json:"server,omitempty" url:"server,omitempty"`
	Export             *string                         `json:"export,omitempty" url:"export,omitempty"`
	Options            *string                         `json:"options,omitempty" url:"options,omitempty"`
	Preallocation      *string                         `json:"preallocation,omitempty" url:"preallocation,omitempty"`
	Datastore          *string                         `json:"datastore,omitempty" url:"datastore,omitempty"`
	Username           *string                         `json:"username,omitempty" url:"username,omitempty"`
	Password           *string                         `json:"password,omitempty" url:"password,omitempty"`
	Namespace          *string                         `json:"namespace,omitempty" url:"namespace,omitempty"`
	Fingerprint        *string                         `json:"fingerprint,omitempty" url:"fingerprint,omitempty"`
	EncryptionKey      *string                         `json:"keyring,omitempty" url:"keyring,omitempty"`
	ZFSPool            *string                         `json:"pool,omitempty" url:"pool,omitempty"`
	ThinProvision      *types.CustomBool               `json:"sparse,omitempty" url:"sparse,omitempty,int"`
	Blocksize          *string                         `json:"blocksize,omitempty" url:"blocksize,omitempty"`
	VolumeGroup        *string                         `json:"vgname,omitempty" url:"vgname,omitempty"`
	WipeRemovedVolumes *types.CustomBool               `json:"saferemove,omitempty" url:"saferemove,omitempty,int"`
	ThinPool           *string                         `json:"thinpool,omitempty" url:"thinpool,omitempty"`
	Share              *string                         `json:"share,omitempty" url:"share,omitempty"`
	Domain             *string                         `json:"domain,omitempty" url:"domain,omitempty"`
	SubDirectory       *string                         `json:"subdir,omitempty" url:"subdir,omitempty"`
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
	Shared       *bool                           `json:"shared,omitempty" url:"shared,omitempty,int"`
}

// DataStoreWithBackups holds optional retention settings for backups.
type DataStoreWithBackups struct {
	MaxProtectedBackups *types.CustomInt64 `json:"max-protected-backups,omitempty" url:"max,omitempty"`
	KeepAll             *types.CustomBool  `json:"-" url:"-"`
	KeepDaily           *int               `json:"-" url:"-"`
	KeepHourly          *int               `json:"-" url:"-"`
	KeepLast            *int               `json:"-" url:"-"`
	KeepMonthly         *int               `json:"-" url:"-"`
	KeepWeekly          *int               `json:"-" url:"-"`
	KeepYearly          *int               `json:"-" url:"-"`
}

// String serializes DataStoreWithBackups into the Proxmox "key=value,key=value" format.
// Only defined (non-nil) fields will be included.
func (b *DataStoreWithBackups) String() string {
	var parts []string

	if b.KeepAll != nil {
		return fmt.Sprintf("keep-all=1")
	}

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

func (b *DataStoreWithBackups) EncodeValues(key string, v *url.Values) error {
	if b.MaxProtectedBackups != nil {
		v.Set("max-protected-backups", strconv.FormatInt(int64(*b.MaxProtectedBackups), 10))
	}

	backupString := b.String()
	if backupString != "" {
		v.Set("prune-backups", backupString)
	}

	return nil
}
