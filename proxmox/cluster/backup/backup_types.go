/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package backup

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// PruneBackupsString is a custom string type that handles the Proxmox API returning
// prune-backups as either a plain string (e.g. "keep-last=2,keep-weekly=1") or a JSON object
// (e.g. {"keep-last":2,"keep-weekly":"1"}). It always stores the value as a comma-separated string.
type PruneBackupsString string

// UnmarshalJSON handles string, object with int values, and object with string values.
func (p *PruneBackupsString) UnmarshalJSON(data []byte) error {
	// Try string first.
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*p = PruneBackupsString(s)
		return nil
	}

	// Try object with varying value types (int or string).
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return fmt.Errorf("failed to unmarshal PruneBackupsString: %w", err)
	}

	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	parts := make([]string, 0, len(obj))

	for _, k := range keys {
		raw := obj[k]

		// Try as integer.
		var intVal int
		if err := json.Unmarshal(raw, &intVal); err == nil {
			parts = append(parts, fmt.Sprintf("%s=%d", k, intVal))
			continue
		}

		// Try as string.
		var strVal string
		if err := json.Unmarshal(raw, &strVal); err == nil {
			parts = append(parts, fmt.Sprintf("%s=%s", k, strVal))
			continue
		}

		return fmt.Errorf("failed to unmarshal PruneBackupsString value for key %q", k)
	}

	*p = PruneBackupsString(strings.Join(parts, ","))

	return nil
}

// String returns the string representation.
func (p PruneBackupsString) String() string {
	return string(p)
}

// Pointer returns a pointer to the underlying string value, or nil if empty.
func (p PruneBackupsString) Pointer() *string {
	if p == "" {
		return nil
	}

	s := string(p)

	return &s
}

// ListResponseBody contains the body from a backup job list response.
type ListResponseBody struct {
	Data []*GetResponseData `json:"data,omitempty"`
}

// GetResponseBody contains the body from a backup job get response.
type GetResponseBody struct {
	Data *GetResponseData `json:"data,omitempty"`
}

// GetResponseData contains the data from a backup job get response.
type GetResponseData struct {
	ID                     string                          `json:"id"`
	Type                   *string                         `json:"type,omitempty"`
	Enabled                *types.CustomBool               `json:"enabled,omitempty"`
	Schedule               string                          `json:"schedule"`
	Storage                string                          `json:"storage"`
	Node                   *string                         `json:"node,omitempty"`
	VMID                   *string                         `json:"vmid,omitempty"`
	All                    *types.CustomBool               `json:"all,omitempty"`
	Mode                   *string                         `json:"mode,omitempty"`
	Compress               *string                         `json:"compress,omitempty"`
	StartTime              *string                         `json:"starttime,omitempty"`
	MaxFiles               *int                            `json:"maxfiles,omitempty"`
	MailTo                 *string                         `json:"mailto,omitempty"`
	MailNotification       *string                         `json:"mailnotification,omitempty"`
	BwLimit                *int                            `json:"bwlimit,omitempty"`
	IONice                 *int                            `json:"ionice,omitempty"`
	Pigz                   *int                            `json:"pigz,omitempty"`
	Zstd                   *int                            `json:"zstd,omitempty"`
	PruneBackups           *PruneBackupsString             `json:"prune-backups,omitempty"`
	Remove                 *types.CustomBool               `json:"remove,omitempty"`
	NotesTemplate          *string                         `json:"notes-template,omitempty"`
	Protected              *types.CustomBool               `json:"protected,omitempty"`
	RepeatMissed           *types.CustomBool               `json:"repeat-missed,omitempty"`
	Script                 *string                         `json:"script,omitempty"`
	StdExcludes            *types.CustomBool               `json:"stdexcludes,omitempty"`
	ExcludePath            *types.CustomCommaSeparatedList `json:"exclude-path,omitempty"`
	Pool                   *string                         `json:"pool,omitempty"`
	Fleecing               *FleecingConfig                 `json:"fleecing,omitempty"`
	Performance            *PerformanceConfig              `json:"performance,omitempty"`
	PBSChangeDetectionMode *string                         `json:"pbs-change-detection-mode,omitempty"`
	LockWait               *int                            `json:"lockwait,omitempty"`
	StopWait               *int                            `json:"stopwait,omitempty"`
	TmpDir                 *string                         `json:"tmpdir,omitempty"`
}

// RequestBodyCommon contains common fields for backup job create and update requests.
type RequestBodyCommon struct {
	Enabled                *types.CustomBool  `json:"enabled,omitempty"                   url:"enabled,omitempty,int"`
	Node                   *string            `json:"node,omitempty"                      url:"node,omitempty"`
	VMID                   *string            `json:"vmid,omitempty"                      url:"vmid,omitempty"`
	All                    *types.CustomBool  `json:"all,omitempty"                       url:"all,omitempty,int"`
	Mode                   *string            `json:"mode,omitempty"                      url:"mode,omitempty"`
	Compress               *string            `json:"compress,omitempty"                  url:"compress,omitempty"`
	StartTime              *string            `json:"starttime,omitempty"                 url:"starttime,omitempty"`
	MaxFiles               *int               `json:"maxfiles,omitempty"                  url:"maxfiles,omitempty"`
	MailTo                 *string            `json:"mailto,omitempty"                    url:"mailto,omitempty"`
	MailNotification       *string            `json:"mailnotification,omitempty"          url:"mailnotification,omitempty"`
	BwLimit                *int               `json:"bwlimit,omitempty"                   url:"bwlimit,omitempty"`
	IONice                 *int               `json:"ionice,omitempty"                    url:"ionice,omitempty"`
	Pigz                   *int               `json:"pigz,omitempty"                      url:"pigz,omitempty"`
	Zstd                   *int               `json:"zstd,omitempty"                      url:"zstd,omitempty"`
	PruneBackups           *string            `json:"prune-backups,omitempty"             url:"prune-backups,omitempty"`
	Remove                 *types.CustomBool  `json:"remove,omitempty"                    url:"remove,omitempty,int"`
	NotesTemplate          *string            `json:"notes-template,omitempty"            url:"notes-template,omitempty"`
	Protected              *types.CustomBool  `json:"protected,omitempty"                 url:"protected,omitempty,int"`
	RepeatMissed           *types.CustomBool  `json:"repeat-missed,omitempty"             url:"repeat-missed,omitempty,int"`
	Script                 *string            `json:"script,omitempty"                    url:"script,omitempty"`
	StdExcludes            *types.CustomBool  `json:"stdexcludes,omitempty"               url:"stdexcludes,omitempty,int"`
	ExcludePath            *string            `json:"exclude-path,omitempty"              url:"exclude-path,omitempty"`
	Pool                   *string            `json:"pool,omitempty"                      url:"pool,omitempty"`
	Fleecing               *FleecingConfig    `json:"fleecing,omitempty"                  url:"fleecing,omitempty"`
	Performance            *PerformanceConfig `json:"performance,omitempty"               url:"performance,omitempty"`
	PBSChangeDetectionMode *string            `json:"pbs-change-detection-mode,omitempty" url:"pbs-change-detection-mode,omitempty"`
	LockWait               *int               `json:"lockwait,omitempty"                  url:"lockwait,omitempty"`
	StopWait               *int               `json:"stopwait,omitempty"                  url:"stopwait,omitempty"`
	TmpDir                 *string            `json:"tmpdir,omitempty"                    url:"tmpdir,omitempty"`
}

// CreateRequestBody contains the body for creating a new backup job.
type CreateRequestBody struct {
	RequestBodyCommon

	ID       string `json:"id"       url:"id"`
	Schedule string `json:"schedule" url:"schedule"`
	Storage  string `json:"storage"  url:"storage"`
}

// UpdateRequestBody contains the body for updating an existing backup job.
type UpdateRequestBody struct {
	RequestBodyCommon

	Schedule *string  `json:"schedule,omitempty" url:"schedule,omitempty"`
	Storage  *string  `json:"storage,omitempty"  url:"storage,omitempty"`
	Delete   []string `json:"delete,omitempty"   url:"delete,omitempty,comma"`
}

// FleecingConfig contains the fleecing configuration for a backup job.
type FleecingConfig struct {
	Enabled *types.CustomBool `json:"enabled,omitempty" url:"enabled,omitempty,int"`
	Storage *string           `json:"storage,omitempty" url:"storage,omitempty"`
}

// EncodeValues encodes the FleecingConfig into URL values as a comma-separated key=value string.
func (f *FleecingConfig) EncodeValues(key string, v *url.Values) error {
	var parts []string

	if f.Enabled != nil {
		if *f.Enabled {
			parts = append(parts, "enabled=1")
		} else {
			parts = append(parts, "enabled=0")
		}
	}

	if f.Storage != nil {
		parts = append(parts, fmt.Sprintf("storage=%s", *f.Storage))
	}

	if len(parts) > 0 {
		v.Add(key, strings.Join(parts, ","))
	}

	return nil
}

// PerformanceConfig contains the performance configuration for a backup job.
type PerformanceConfig struct {
	MaxWorkers    *types.CustomInt `json:"max-workers,omitempty"     url:"max-workers,omitempty"`
	PBSEntriesMax *types.CustomInt `json:"pbs-entries-max,omitempty" url:"pbs-entries-max,omitempty"`
}

// EncodeValues encodes the PerformanceConfig into URL values as a comma-separated key=value string.
func (p *PerformanceConfig) EncodeValues(key string, v *url.Values) error {
	var parts []string

	if p.MaxWorkers != nil {
		parts = append(parts, fmt.Sprintf("max-workers=%d", int(*p.MaxWorkers)))
	}

	if p.PBSEntriesMax != nil {
		parts = append(parts, fmt.Sprintf("pbs-entries-max=%d", int(*p.PBSEntriesMax)))
	}

	if len(parts) > 0 {
		v.Add(key, strings.Join(parts, ","))
	}

	return nil
}
