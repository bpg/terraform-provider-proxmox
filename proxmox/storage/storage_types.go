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

// DatastoreListRequestBody contains the body for a datastore list request.
type DatastoreListRequestBody struct {
	ContentTypes types.CustomCommaSeparatedList `json:"content,omitempty" url:"content,omitempty,comma"`
	Enabled      *types.CustomBool              `json:"enabled,omitempty" url:"enabled,omitempty,int"`
	Format       *types.CustomBool              `json:"format,omitempty"  url:"format,omitempty,int"`
	ID           *string                        `json:"storage,omitempty" url:"storage,omitempty"`
	Target       *string                        `json:"target,omitempty"  url:"target,omitempty"`
}

// DatastoreListResponseBody contains the body from a datastore list response.
type DatastoreListResponseBody struct {
	Data []*DatastoreListResponseData `json:"data,omitempty"`
}

// DatastoreListResponseData contains the data from a datastore list response.
type DatastoreListResponseData struct {
	Active              *types.CustomBool               `json:"active,omitempty"`
	ContentTypes        *types.CustomCommaSeparatedList `json:"content,omitempty"`
	Enabled             *types.CustomBool               `json:"enabled,omitempty"`
	ID                  string                          `json:"storage,omitempty"`
	Shared              *types.CustomBool               `json:"shared,omitempty"`
	SpaceAvailable      *types.CustomInt64              `json:"avail,omitempty"`
	SpaceTotal          *types.CustomInt64              `json:"total,omitempty"`
	SpaceUsed           *types.CustomInt64              `json:"used,omitempty"`
	SpaceUsedPercentage *types.CustomFloat64            `json:"used_fraction,omitempty"`
	Type                string                          `json:"type,omitempty"`
}

// DataStoreBase contains the common fields for all storage types.
type DataStoreBase struct {
	Storage string `json:"storage"`
	Nodes   string `json:"nodes,omitempty"`
	Content string `json:"content,omitempty"`
	Enable  bool   `json:"enable,omitempty"`
	Shared  bool   `json:"shared,omitempty"`
}

// DataStoreWithBackups holds optional retention settings for backups.
type DataStoreWithBackups struct {
	MaxProtectedBackups *types.CustomInt64 `json:"max-protected-backups,omitempty"`
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
		PruneBackups string `json:"prune-backups,omitempty"`
	}{
		Alias:        (*Alias)(&b),
		PruneBackups: str,
	}
	return json.Marshal(aux)
}

// DirectoryStorageRequestBody defines options for 'dir' type storage.
type DirectoryStorageRequestBody struct {
	DataStoreBase
	DataStoreWithBackups
	Path                   string `json:"path"`
	Preallocation          string `json:"preallocation,omitempty"`
	SnapshotsAsVolumeChain bool   `json:"snapshot,omitempty"`
}

// LVMStorageRequestBody defines options for 'lvm' type storage.
type LVMStorageRequestBody struct {
	DataStoreBase
	VolumeGroup        string `json:"volume_group"`
	WipeRemovedVolumes bool   `json:"wipe_removed_volumes,omitempty"`
}

// LVMThinStorageRequestBody defines options for 'lvmthin' type storage.
type LVMThinStorageRequestBody struct {
	DataStoreBase
	VolumeGroup string `json:"volume_group"`
	ThinPool    string `json:"thin_pool,omitempty"`
}

// BTRFSStorageRequestBody defines options for 'btrfs' type storage.
type BTRFSStorageRequestBody struct {
	DataStoreBase
	DataStoreWithBackups
	Path          string `json:"path"`
	Preallocation string `json:"preallocation,omitempty"`
}

// NFSStorageRequestBody defines specific options for 'nfs' type storage.
type NFSStorageRequestBody struct {
	DataStoreBase
	Export                 string `json:"export"`
	NFSVersion             string `json:"nfs_version,omitempty"`
	Server                 string `json:"server"`
	Preallocation          string `json:"preallocation,omitempty"`
	SnapshotsAsVolumeChain bool   `json:"snapshot-as-volume-chain,omitempty"`
}

// SMBStorageRequestBody defines specific options for 'smb'/'cifs' type storage.
type SMBStorageRequestBody struct {
	DataStoreBase
	DataStoreWithBackups
	Username               string `json:"username"`
	Password               string `json:"password"`
	Share                  string `json:"share"`
	Domain                 string `json:"domain,omitempty"`
	Subdirectory           string `json:"subdirectory,omitempty"`
	Server                 string `json:"server"`
	Preallocation          string `json:"preallocation,omitempty"`
	SnapshotsAsVolumeChain bool   `json:"snapshot-as-volume-chain,omitempty"`
}

// ISCSIStorageRequestBody defines options for 'iscsi' type storage.
type ISCSIStorageRequestBody struct {
	DataStoreBase
	Portal          string `json:"portal"`
	Target          string `json:"target"`
	UseLUNsDirectly bool   `json:"use_luns_directly,omitempty"`
}

// CephFSStorageRequestBody defines options for 'cephfs' type storage.
type CephFSStorageRequestBody struct {
	DataStoreBase
	DataStoreWithBackups
	Monitors  string `json:"monhost"`
	Username  string `json:"username,omitempty"`
	FSName    string `json:"fs_name,omitempty"`
	SecretKey string `json:"keyring,omitempty"`
	Managed   bool   `json:"managed,omitempty"`
}

// RBDStorageRequestBody defines options for 'rbd' type storage.
type RBDStorageRequestBody struct {
	DataStoreBase
	Pool      string `json:"pool"`
	Monitors  string `json:"monhost"`
	Username  string `json:"username,omitempty"`
	KRBD      bool   `json:"krbd,omitempty"`
	SecretKey string `json:"keyring"`
	Managed   bool   `json:"managed,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

// ZFSStorageRequestBody defines options for 'zfs' type storage.
type ZFSStorageRequestBody struct {
	DataStoreBase
	ZFSPool       string `json:"zfs_pool"`
	ThinProvision bool   `json:"thin_provision,omitempty"`
	Blocksize     string `json:"blocksize,omitempty"`
}

// ZFSOverISCSIOptions defines options for 'zfs over iscsi' type storage.
type ZFSOverISCSIOptions struct {
	DataStoreBase
	Portal            string `json:"portal"`
	Pool              string `json:"pool"`
	Blocksize         string `json:"blocksize,omitempty"`
	Target            string `json:"target"`
	TargetGroup       string `json:"target_group,omitempty"`
	ISCSIProvider     string `json:"iscsi_provider"`
	ThinProvision     bool   `json:"thin_provision,omitempty"`
	WriteCache        bool   `json:"write_cache,omitempty"`
	HostGroup         string `json:"host_group,omitempty"`
	TargetPortalGroup string `json:"target_portal_group,omitempty"`
}

// PBSStorageRequestBody defines options for 'pbs' (Proxmox Backup Server) type storage.
type PBSStorageRequestBody struct {
	DataStoreBase
	DataStoreWithBackups
	Server      string `json:"server"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Datastore   string `json:"datastore"`
	Namespace   string `json:"namespace,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	Encryption  string `json:"encryption-key,omitempty"`
}

// ESXiStorageRequestBody defines options for 'esxi' type storage.
type ESXiStorageRequestBody struct {
	DataStoreBase
	Server               string `json:"server"`
	Username             string `json:"username"`
	Password             string `json:"password"`
	SkipCertVerification bool   `json:"skip_cert_verification,omitempty"`
}
