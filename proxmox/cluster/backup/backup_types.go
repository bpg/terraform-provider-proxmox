/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package backup

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// PruneBackupsString is a custom type that handles prune-backups.
// The API accepts a string like "keep-last=3,keep-daily=7" but returns an object.
type PruneBackupsString string

// UnmarshalJSON handles both string and object formats from the API.
func (p *PruneBackupsString) UnmarshalJSON(data []byte) error {
	// try string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*p = PruneBackupsString(str)
		return nil
	}

	// try object format with int values: {"keep-last": 3, "keep-daily": 7}
	var objInt map[string]int
	if err := json.Unmarshal(data, &objInt); err == nil {
		*p = formatPruneBackups(objInt)
		return nil
	}

	// try object format with string values: {"keep-last": "3", "keep-daily": "7"}
	var objStr map[string]string
	if err := json.Unmarshal(data, &objStr); err == nil {
		objInt = make(map[string]int, len(objStr))

		for k, v := range objStr {
			var val int
			if _, err := fmt.Sscanf(v, "%d", &val); err != nil {
				return fmt.Errorf("could not parse value %q for key %q in prune-backups: %w", v, k, err)
			}

			objInt[k] = val
		}

		*p = formatPruneBackups(objInt)

		return nil
	}

	return fmt.Errorf("prune-backups must be a string or object, got: %s", string(data))
}

func formatPruneBackups(obj map[string]int) PruneBackupsString {
	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%d", k, obj[k]))
	}

	return PruneBackupsString(strings.Join(parts, ","))
}

// String returns the string value.
func (p PruneBackupsString) String() string {
	return string(p)
}

// Pointer returns a pointer to the string value, or nil if empty.
func (p PruneBackupsString) Pointer() *string {
	if p == "" {
		return nil
	}

	s := string(p)

	return &s
}

// ListResponseBody contains the response body for listing backup jobs.
type ListResponseBody struct {
	Data []*GetResponseData `json:"data,omitempty"`
}

// GetResponseBody contains the response body for getting a backup job.
type GetResponseBody struct {
	Data *GetResponseData `json:"data,omitempty"`
}

// GetResponseData contains the data from a backup job response.
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

// RequestBodyCommon contains optional fields shared between create and update request bodies.
type RequestBodyCommon struct {
	Enabled                *types.CustomBool               `json:"enabled,omitempty"                   url:"enabled,omitempty,int"`
	Node                   *string                         `json:"node,omitempty"                      url:"node,omitempty"`
	VMID                   *string                         `json:"vmid,omitempty"                      url:"vmid,omitempty"`
	All                    *types.CustomBool               `json:"all,omitempty"                       url:"all,omitempty,int"`
	Mode                   *string                         `json:"mode,omitempty"                      url:"mode,omitempty"`
	Compress               *string                         `json:"compress,omitempty"                  url:"compress,omitempty"`
	StartTime              *string                         `json:"starttime,omitempty"                 url:"starttime,omitempty"`
	MaxFiles               *int                            `json:"maxfiles,omitempty"                  url:"maxfiles,omitempty"`
	MailTo                 *string                         `json:"mailto,omitempty"                    url:"mailto,omitempty"`
	MailNotification       *string                         `json:"mailnotification,omitempty"          url:"mailnotification,omitempty"`
	BwLimit                *int                            `json:"bwlimit,omitempty"                   url:"bwlimit,omitempty"`
	IONice                 *int                            `json:"ionice,omitempty"                    url:"ionice,omitempty"`
	Pigz                   *int                            `json:"pigz,omitempty"                      url:"pigz,omitempty"`
	Zstd                   *int                            `json:"zstd,omitempty"                      url:"zstd,omitempty"`
	PruneBackups           *string                         `json:"prune-backups,omitempty"             url:"prune-backups,omitempty"`
	Remove                 *types.CustomBool               `json:"remove,omitempty"                    url:"remove,omitempty,int"`
	NotesTemplate          *string                         `json:"notes-template,omitempty"            url:"notes-template,omitempty"`
	Protected              *types.CustomBool               `json:"protected,omitempty"                 url:"protected,omitempty,int"`
	RepeatMissed           *types.CustomBool               `json:"repeat-missed,omitempty"             url:"repeat-missed,omitempty,int"`
	Script                 *string                         `json:"script,omitempty"                    url:"script,omitempty"`
	StdExcludes            *types.CustomBool               `json:"stdexcludes,omitempty"               url:"stdexcludes,omitempty,int"`
	ExcludePath            *types.CustomCommaSeparatedList `json:"exclude-path,omitempty"              url:"exclude-path,omitempty,comma"`
	Pool                   *string                         `json:"pool,omitempty"                      url:"pool,omitempty"`
	Fleecing               *FleecingConfig                 `json:"fleecing,omitempty"                  url:"fleecing,omitempty"`
	Performance            *PerformanceConfig              `json:"performance,omitempty"               url:"performance,omitempty"`
	PBSChangeDetectionMode *string                         `json:"pbs-change-detection-mode,omitempty" url:"pbs-change-detection-mode,omitempty"`
	LockWait               *int                            `json:"lockwait,omitempty"                  url:"lockwait,omitempty"`
	StopWait               *int                            `json:"stopwait,omitempty"                  url:"stopwait,omitempty"`
	TmpDir                 *string                         `json:"tmpdir,omitempty"                    url:"tmpdir,omitempty"`
}

// CreateRequestBody contains the request body for creating a backup job.
type CreateRequestBody struct {
	RequestBodyCommon

	ID       string `json:"id"       url:"id"`
	Schedule string `json:"schedule" url:"schedule"`
	Storage  string `json:"storage"  url:"storage"`
}

// UpdateRequestBody contains the request body for updating a backup job.
type UpdateRequestBody struct {
	RequestBodyCommon

	Schedule *string `json:"schedule,omitempty" url:"schedule,omitempty"`
	Storage  *string `json:"storage,omitempty"  url:"storage,omitempty"`
	Delete   *string `json:"delete,omitempty"   url:"delete,omitempty"`
}

// FleecingConfig contains backup fleecing configuration.
type FleecingConfig struct {
	Enabled *types.CustomBool `json:"enabled,omitempty" url:"enabled,omitempty,int"`
	Storage *string           `json:"storage,omitempty" url:"storage,omitempty"`
}

// PerformanceConfig contains performance tuning configuration.
type PerformanceConfig struct {
	MaxWorkers    *int `json:"max-workers,omitempty"     url:"max-workers,omitempty"`
	PBSEntriesMax *int `json:"pbs-entries-max,omitempty" url:"pbs-entries-max,omitempty"`
}
