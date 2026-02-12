/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package backup

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// stringPtrOr converts a *string to a types.String, returning def when nil.
func stringPtrOr(p *string, def types.String) types.String {
	if p != nil {
		return types.StringPointerValue(p)
	}

	return def
}

// customBoolPtrOr converts a *CustomBool to a types.Bool, returning BoolValue(def) when nil.
func customBoolPtrOr(p *proxmoxtypes.CustomBool, def bool) types.Bool {
	if p != nil {
		return types.BoolValue(bool(*p))
	}

	return types.BoolValue(def)
}

// intPtrOr converts a *int to a types.Int64, returning def when nil.
func intPtrOr(p *int, def types.Int64) types.Int64 {
	if p != nil {
		return types.Int64Value(int64(*p))
	}

	return def
}

// BackupJobModel maps the resource schema data.
type BackupJobModel struct {
	ID                     types.String `tfsdk:"id"`
	Enabled                types.Bool   `tfsdk:"enabled"`
	Schedule               types.String `tfsdk:"schedule"`
	Storage                types.String `tfsdk:"storage"`
	Node                   types.String `tfsdk:"node"`
	VMID                   types.String `tfsdk:"vmid"`
	All                    types.Bool   `tfsdk:"all"`
	Mode                   types.String `tfsdk:"mode"`
	Compress               types.String `tfsdk:"compress"`
	StartTime              types.String `tfsdk:"starttime"`
	MaxFiles               types.Int64  `tfsdk:"maxfiles"`
	MailTo                 types.String `tfsdk:"mailto"`
	MailNotification       types.String `tfsdk:"mailnotification"`
	BwLimit                types.Int64  `tfsdk:"bwlimit"`
	IONice                 types.Int64  `tfsdk:"ionice"`
	Pigz                   types.Int64  `tfsdk:"pigz"`
	Zstd                   types.Int64  `tfsdk:"zstd"`
	PruneBackups           types.String `tfsdk:"prune_backups"`
	Remove                 types.Bool   `tfsdk:"remove"`
	NotesTemplate          types.String `tfsdk:"notes_template"`
	Protected              types.Bool   `tfsdk:"protected"`
	RepeatMissed           types.Bool   `tfsdk:"repeat_missed"`
	Script                 types.String `tfsdk:"script"`
	StdExcludes            types.Bool   `tfsdk:"stdexcludes"`
	ExcludePath            types.List   `tfsdk:"exclude_path"`
	Pool                   types.String `tfsdk:"pool"`
	Fleecing               types.Object `tfsdk:"fleecing"`
	Performance            types.Object `tfsdk:"performance"`
	PBSChangeDetectionMode types.String `tfsdk:"pbs_change_detection_mode"`
	LockWait               types.Int64  `tfsdk:"lockwait"`
	StopWait               types.Int64  `tfsdk:"stopwait"`
	TmpDir                 types.String `tfsdk:"tmpdir"`
}

// FleecingModel maps the fleecing nested object.
type FleecingModel struct {
	Enabled types.Bool   `tfsdk:"enabled"`
	Storage types.String `tfsdk:"storage"`
}

// PerformanceModel maps the performance nested object.
type PerformanceModel struct {
	MaxWorkers    types.Int64 `tfsdk:"max_workers"`
	PBSEntriesMax types.Int64 `tfsdk:"pbs_entries_max"`
}
