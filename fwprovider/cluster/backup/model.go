/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package backup

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/backup"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type backupJobModel struct {
	ID                     types.String `tfsdk:"id"`
	Schedule               types.String `tfsdk:"schedule"`
	Storage                types.String `tfsdk:"storage"`
	Enabled                types.Bool   `tfsdk:"enabled"`
	Node                   types.String `tfsdk:"node"`
	VMIDs                  types.List   `tfsdk:"vmid"`
	All                    types.Bool   `tfsdk:"all"`
	Mode                   types.String `tfsdk:"mode"`
	Compress               types.String `tfsdk:"compress"`
	StartTime              types.String `tfsdk:"starttime"`
	MaxFiles               types.Int64  `tfsdk:"maxfiles"`
	MailTo                 types.List   `tfsdk:"mailto"`
	MailNotification       types.String `tfsdk:"mailnotification"`
	BwLimit                types.Int64  `tfsdk:"bwlimit"`
	IONice                 types.Int64  `tfsdk:"ionice"`
	Pigz                   types.Int64  `tfsdk:"pigz"`
	Zstd                   types.Int64  `tfsdk:"zstd"`
	PruneBackups           types.Map    `tfsdk:"prune_backups"`
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

type fleecingModel struct {
	Enabled types.Bool   `tfsdk:"enabled"`
	Storage types.String `tfsdk:"storage"`
}

type performanceModel struct {
	MaxWorkers    types.Int64 `tfsdk:"max_workers"`
	PBSEntriesMax types.Int64 `tfsdk:"pbs_entries_max"`
}

func fleecingAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled": types.BoolType,
		"storage": types.StringType,
	}
}

func performanceAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"max_workers":     types.Int64Type,
		"pbs_entries_max": types.Int64Type,
	}
}

// int64PtrToIntPtr converts *int64 to *int for API fields.
func int64PtrToIntPtr(v *int64) *int {
	if v == nil {
		return nil
	}

	i := int(*v)

	return &i
}

// intPtrToInt64Ptr converts *int to *int64 for Terraform state.
func intPtrToInt64Ptr(v *int) *int64 {
	if v == nil {
		return nil
	}

	i := int64(*v)

	return &i
}

// int64PtrToCustomIntPtr converts *int64 to *types.CustomInt for API fields.
func int64PtrToCustomIntPtr(v *int64) *proxmoxtypes.CustomInt {
	if v == nil {
		return nil
	}

	i := proxmoxtypes.CustomInt(*v)

	return &i
}

// customIntPtrToInt64Ptr converts *types.CustomInt to *int64 for Terraform state.
func customIntPtrToInt64Ptr(v *proxmoxtypes.CustomInt) *int64 {
	if v == nil {
		return nil
	}

	i := int64(*v)

	return &i
}

func (m *backupJobModel) toCreateAPI(ctx context.Context, diags *diag.Diagnostics) *backup.CreateRequestBody {
	body := &backup.CreateRequestBody{}

	body.ID = m.ID.ValueString()
	body.Schedule = m.Schedule.ValueString()
	body.Storage = m.Storage.ValueString()

	m.fillCommonFields(ctx, &body.RequestBodyCommon, diags)

	return body
}

func (m *backupJobModel) toUpdateAPI(
	ctx context.Context,
	state *backupJobModel,
	diags *diag.Diagnostics,
) *backup.UpdateRequestBody {
	body := &backup.UpdateRequestBody{}

	body.Schedule = m.Schedule.ValueStringPointer()
	body.Storage = m.Storage.ValueStringPointer()

	m.fillCommonFields(ctx, &body.RequestBodyCommon, diags)

	var toDelete []string

	attribute.CheckDelete(m.Node, state.Node, &toDelete, "node")
	attribute.CheckDelete(m.VMIDs, state.VMIDs, &toDelete, "vmid")

	attribute.CheckDelete(m.Mode, state.Mode, &toDelete, "mode")
	attribute.CheckDelete(m.Compress, state.Compress, &toDelete, "compress")
	attribute.CheckDelete(m.StartTime, state.StartTime, &toDelete, "starttime")
	attribute.CheckDelete(m.MaxFiles, state.MaxFiles, &toDelete, "maxfiles")
	attribute.CheckDelete(m.MailTo, state.MailTo, &toDelete, "mailto")
	attribute.CheckDelete(m.MailNotification, state.MailNotification, &toDelete, "mailnotification")
	attribute.CheckDelete(m.BwLimit, state.BwLimit, &toDelete, "bwlimit")
	attribute.CheckDelete(m.IONice, state.IONice, &toDelete, "ionice")
	attribute.CheckDelete(m.Pigz, state.Pigz, &toDelete, "pigz")
	attribute.CheckDelete(m.Zstd, state.Zstd, &toDelete, "zstd")
	attribute.CheckDelete(m.PruneBackups, state.PruneBackups, &toDelete, "prune-backups")
	attribute.CheckDelete(m.Remove, state.Remove, &toDelete, "remove")
	attribute.CheckDelete(m.NotesTemplate, state.NotesTemplate, &toDelete, "notes-template")
	attribute.CheckDelete(m.Protected, state.Protected, &toDelete, "protected")
	attribute.CheckDelete(m.RepeatMissed, state.RepeatMissed, &toDelete, "repeat-missed")
	attribute.CheckDelete(m.Script, state.Script, &toDelete, "script")
	attribute.CheckDelete(m.StdExcludes, state.StdExcludes, &toDelete, "stdexcludes")
	attribute.CheckDelete(m.ExcludePath, state.ExcludePath, &toDelete, "exclude-path")

	attribute.CheckDelete(m.Pool, state.Pool, &toDelete, "pool")
	attribute.CheckDelete(m.Fleecing, state.Fleecing, &toDelete, "fleecing")
	attribute.CheckDelete(m.Performance, state.Performance, &toDelete, "performance")
	attribute.CheckDelete(m.PBSChangeDetectionMode, state.PBSChangeDetectionMode, &toDelete, "pbs-change-detection-mode")
	attribute.CheckDelete(m.LockWait, state.LockWait, &toDelete, "lockwait")
	attribute.CheckDelete(m.StopWait, state.StopWait, &toDelete, "stopwait")
	attribute.CheckDelete(m.TmpDir, state.TmpDir, &toDelete, "tmpdir")
	attribute.CheckDelete(m.Enabled, state.Enabled, &toDelete, "enabled")
	attribute.CheckDelete(m.All, state.All, &toDelete, "all")

	if len(toDelete) > 0 {
		body.Delete = toDelete
	}

	return body
}

func (m *backupJobModel) fillCommonFields(
	ctx context.Context,
	common *backup.RequestBodyCommon,
	diags *diag.Diagnostics,
) {
	common.Enabled = attribute.CustomBoolPtrFromValue(m.Enabled)
	common.Node = attribute.StringPtrFromValue(m.Node)
	common.All = attribute.CustomBoolPtrFromValue(m.All)
	common.Mode = attribute.StringPtrFromValue(m.Mode)
	common.Compress = attribute.StringPtrFromValue(m.Compress)
	common.StartTime = attribute.StringPtrFromValue(m.StartTime)
	common.MaxFiles = int64PtrToIntPtr(attribute.Int64PtrFromValue(m.MaxFiles))
	// MailTo: convert types.List to comma-separated string
	if !m.MailTo.IsNull() && !m.MailTo.IsUnknown() {
		var emails []string

		d := m.MailTo.ElementsAs(ctx, &emails, false)
		diags.Append(d...)

		if !d.HasError() && len(emails) > 0 {
			emailStr := strings.Join(emails, ",")
			common.MailTo = &emailStr
		}
	}

	common.MailNotification = attribute.StringPtrFromValue(m.MailNotification)
	common.BwLimit = int64PtrToIntPtr(attribute.Int64PtrFromValue(m.BwLimit))
	common.IONice = int64PtrToIntPtr(attribute.Int64PtrFromValue(m.IONice))
	common.Pigz = int64PtrToIntPtr(attribute.Int64PtrFromValue(m.Pigz))
	common.Zstd = int64PtrToIntPtr(attribute.Int64PtrFromValue(m.Zstd))
	// PruneBackups: convert types.Map to comma-separated key=value string
	if !m.PruneBackups.IsNull() && !m.PruneBackups.IsUnknown() {
		pruneMap := make(map[string]string)

		d := m.PruneBackups.ElementsAs(ctx, &pruneMap, false)
		diags.Append(d...)

		if !d.HasError() && len(pruneMap) > 0 {
			keys := make([]string, 0, len(pruneMap))
			for k := range pruneMap {
				keys = append(keys, k)
			}

			sort.Strings(keys)

			parts := make([]string, 0, len(pruneMap))
			for _, k := range keys {
				parts = append(parts, fmt.Sprintf("%s=%s", k, pruneMap[k]))
			}

			pruneStr := strings.Join(parts, ",")
			common.PruneBackups = &pruneStr
		}
	}

	common.Remove = attribute.CustomBoolPtrFromValue(m.Remove)
	common.NotesTemplate = attribute.StringPtrFromValue(m.NotesTemplate)
	common.Protected = attribute.CustomBoolPtrFromValue(m.Protected)
	common.RepeatMissed = attribute.CustomBoolPtrFromValue(m.RepeatMissed)
	common.Script = attribute.StringPtrFromValue(m.Script)
	common.StdExcludes = attribute.CustomBoolPtrFromValue(m.StdExcludes)
	common.Pool = attribute.StringPtrFromValue(m.Pool)
	common.PBSChangeDetectionMode = attribute.StringPtrFromValue(m.PBSChangeDetectionMode)
	common.LockWait = int64PtrToIntPtr(attribute.Int64PtrFromValue(m.LockWait))
	common.StopWait = int64PtrToIntPtr(attribute.Int64PtrFromValue(m.StopWait))
	common.TmpDir = attribute.StringPtrFromValue(m.TmpDir)

	// VMID: convert types.List to comma-separated string
	if !m.VMIDs.IsNull() && !m.VMIDs.IsUnknown() {
		var vmids []string

		d := m.VMIDs.ElementsAs(ctx, &vmids, false)
		diags.Append(d...)

		if !d.HasError() && len(vmids) > 0 {
			vmidStr := strings.Join(vmids, ",")
			common.VMID = &vmidStr
		}
	}

	// ExcludePath: convert types.List to comma-separated string
	if !m.ExcludePath.IsNull() && !m.ExcludePath.IsUnknown() {
		var paths []string

		d := m.ExcludePath.ElementsAs(ctx, &paths, false)
		diags.Append(d...)

		if !d.HasError() && len(paths) > 0 {
			excludeStr := strings.Join(paths, ",")
			common.ExcludePath = &excludeStr
		}
	}

	// Fleecing: extract nested object
	if !m.Fleecing.IsNull() && !m.Fleecing.IsUnknown() {
		var fleecing fleecingModel

		d := m.Fleecing.As(ctx, &fleecing, basetypes.ObjectAsOptions{})
		diags.Append(d...)

		if !d.HasError() {
			common.Fleecing = &backup.FleecingConfig{
				Enabled: proxmoxtypes.CustomBoolPtr(fleecing.Enabled.ValueBoolPointer()),
				Storage: fleecing.Storage.ValueStringPointer(),
			}
		}
	}

	// Performance: extract nested object
	if !m.Performance.IsNull() && !m.Performance.IsUnknown() {
		var perf performanceModel

		d := m.Performance.As(ctx, &perf, basetypes.ObjectAsOptions{})
		diags.Append(d...)

		if !d.HasError() {
			common.Performance = &backup.PerformanceConfig{
				MaxWorkers:    int64PtrToCustomIntPtr(perf.MaxWorkers.ValueInt64Pointer()),
				PBSEntriesMax: int64PtrToCustomIntPtr(perf.PBSEntriesMax.ValueInt64Pointer()),
			}
		}
	}
}

func (m *backupJobModel) fromAPI(
	ctx context.Context,
	data *backup.GetResponseData,
) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(data.ID)
	m.Schedule = types.StringValue(data.Schedule)
	m.Storage = types.StringValue(data.Storage)
	m.Enabled = types.BoolPointerValue(data.Enabled.PointerBool())
	m.Node = types.StringPointerValue(data.Node)
	m.All = types.BoolPointerValue(data.All.PointerBool())

	// VMID: convert comma-separated string to list
	if data.VMID != nil && *data.VMID != "" {
		ids := strings.Split(*data.VMID, ",")
		vmidValues := make([]attr.Value, len(ids))

		for i, id := range ids {
			vmidValues[i] = types.StringValue(strings.TrimSpace(id))
		}

		listVal, d := types.ListValue(types.StringType, vmidValues)
		diags.Append(d...)

		m.VMIDs = listVal
	} else {
		m.VMIDs = types.ListNull(types.StringType)
	}

	m.Mode = types.StringPointerValue(data.Mode)
	m.Compress = types.StringPointerValue(data.Compress)
	m.StartTime = types.StringPointerValue(data.StartTime)
	m.MaxFiles = types.Int64PointerValue(intPtrToInt64Ptr(data.MaxFiles))

	// MailTo: convert comma-separated string to list
	if data.MailTo != nil && *data.MailTo != "" {
		emails := strings.Split(*data.MailTo, ",")
		emailValues := make([]attr.Value, len(emails))

		for i, e := range emails {
			emailValues[i] = types.StringValue(strings.TrimSpace(e))
		}

		listVal, d := types.ListValue(types.StringType, emailValues)
		diags.Append(d...)

		m.MailTo = listVal
	} else {
		m.MailTo = types.ListNull(types.StringType)
	}

	m.MailNotification = types.StringPointerValue(data.MailNotification)
	m.BwLimit = types.Int64PointerValue(intPtrToInt64Ptr(data.BwLimit))
	m.IONice = types.Int64PointerValue(intPtrToInt64Ptr(data.IONice))
	m.Pigz = types.Int64PointerValue(intPtrToInt64Ptr(data.Pigz))
	m.Zstd = types.Int64PointerValue(intPtrToInt64Ptr(data.Zstd))
	m.Remove = types.BoolPointerValue(data.Remove.PointerBool())
	m.NotesTemplate = types.StringPointerValue(data.NotesTemplate)
	m.Protected = types.BoolPointerValue(data.Protected.PointerBool())
	m.RepeatMissed = types.BoolPointerValue(data.RepeatMissed.PointerBool())
	m.Script = types.StringPointerValue(data.Script)
	m.StdExcludes = types.BoolPointerValue(data.StdExcludes.PointerBool())
	m.Pool = types.StringPointerValue(data.Pool)
	m.PBSChangeDetectionMode = types.StringPointerValue(data.PBSChangeDetectionMode)
	m.LockWait = types.Int64PointerValue(intPtrToInt64Ptr(data.LockWait))
	m.StopWait = types.Int64PointerValue(intPtrToInt64Ptr(data.StopWait))
	m.TmpDir = types.StringPointerValue(data.TmpDir)

	// PruneBackups: convert comma-separated key=value string to map
	if data.PruneBackups != nil {
		s := data.PruneBackups.Pointer()
		if s != nil && *s != "" {
			mapValues := make(map[string]attr.Value)

			for part := range strings.SplitSeq(*s, ",") {
				kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
				if len(kv) == 2 {
					mapValues[kv[0]] = types.StringValue(kv[1])
				}
			}

			mapVal, d := types.MapValue(types.StringType, mapValues)
			diags.Append(d...)

			m.PruneBackups = mapVal
		} else {
			m.PruneBackups = types.MapNull(types.StringType)
		}
	} else {
		m.PruneBackups = types.MapNull(types.StringType)
	}

	// ExcludePath: convert CustomCommaSeparatedList to types.List
	if data.ExcludePath != nil {
		paths := make([]attr.Value, len(*data.ExcludePath))
		for i, p := range *data.ExcludePath {
			paths[i] = types.StringValue(p)
		}

		listVal, d := types.ListValue(types.StringType, paths)
		diags.Append(d...)

		m.ExcludePath = listVal
	} else {
		m.ExcludePath = types.ListNull(types.StringType)
	}

	// Fleecing: convert to types.Object
	if data.Fleecing != nil {
		fleecingVal := fleecingModel{
			Enabled: types.BoolPointerValue(data.Fleecing.Enabled.PointerBool()),
			Storage: types.StringPointerValue(data.Fleecing.Storage),
		}

		obj, d := types.ObjectValueFrom(ctx, fleecingAttrTypes(), fleecingVal)
		diags.Append(d...)

		m.Fleecing = obj
	} else {
		m.Fleecing = types.ObjectNull(fleecingAttrTypes())
	}

	// Performance: convert to types.Object
	if data.Performance != nil {
		perfVal := performanceModel{
			MaxWorkers:    types.Int64PointerValue(customIntPtrToInt64Ptr(data.Performance.MaxWorkers)),
			PBSEntriesMax: types.Int64PointerValue(customIntPtrToInt64Ptr(data.Performance.PBSEntriesMax)),
		}

		obj, d := types.ObjectValueFrom(ctx, performanceAttrTypes(), perfVal)
		diags.Append(d...)

		m.Performance = obj
	} else {
		m.Performance = types.ObjectNull(performanceAttrTypes())
	}

	return diags
}

// backupJobDatasourceModel is a simplified model for the backup job data source.
type backupJobDatasourceModel struct {
	ID               types.String `tfsdk:"id"`
	Schedule         types.String `tfsdk:"schedule"`
	Storage          types.String `tfsdk:"storage"`
	Enabled          types.Bool   `tfsdk:"enabled"`
	Node             types.String `tfsdk:"node"`
	VMIDs            types.List   `tfsdk:"vmid"`
	All              types.Bool   `tfsdk:"all"`
	Mode             types.String `tfsdk:"mode"`
	Compress         types.String `tfsdk:"compress"`
	MailTo           types.List   `tfsdk:"mailto"`
	MailNotification types.String `tfsdk:"mailnotification"`
	NotesTemplate    types.String `tfsdk:"notes_template"`
	Pool             types.String `tfsdk:"pool"`
	PruneBackups     types.Map    `tfsdk:"prune_backups"`
	Protected        types.Bool   `tfsdk:"protected"`
}

func (m *backupJobDatasourceModel) fromAPI(data *backup.GetResponseData) {
	m.ID = types.StringValue(data.ID)
	m.Schedule = types.StringValue(data.Schedule)
	m.Storage = types.StringValue(data.Storage)
	m.Enabled = types.BoolPointerValue(data.Enabled.PointerBool())
	m.Node = types.StringPointerValue(data.Node)
	m.All = types.BoolPointerValue(data.All.PointerBool())
	m.Mode = types.StringPointerValue(data.Mode)
	m.Compress = types.StringPointerValue(data.Compress)

	// MailTo: convert comma-separated string to list
	if data.MailTo != nil && *data.MailTo != "" {
		emails := strings.Split(*data.MailTo, ",")
		emailValues := make([]attr.Value, len(emails))

		for i, e := range emails {
			emailValues[i] = types.StringValue(strings.TrimSpace(e))
		}

		m.MailTo = types.ListValueMust(types.StringType, emailValues)
	} else {
		m.MailTo = types.ListNull(types.StringType)
	}

	m.MailNotification = types.StringPointerValue(data.MailNotification)
	m.NotesTemplate = types.StringPointerValue(data.NotesTemplate)
	m.Pool = types.StringPointerValue(data.Pool)
	m.Protected = types.BoolPointerValue(data.Protected.PointerBool())

	// VMID: convert comma-separated string to list
	if data.VMID != nil && *data.VMID != "" {
		ids := strings.Split(*data.VMID, ",")
		vmidValues := make([]attr.Value, len(ids))

		for i, id := range ids {
			vmidValues[i] = types.StringValue(strings.TrimSpace(id))
		}

		m.VMIDs = types.ListValueMust(types.StringType, vmidValues)
	} else {
		m.VMIDs = types.ListNull(types.StringType)
	}

	// PruneBackups: convert comma-separated key=value string to map
	if data.PruneBackups != nil {
		s := data.PruneBackups.Pointer()
		if s != nil && *s != "" {
			mapValues := make(map[string]attr.Value)

			for part := range strings.SplitSeq(*s, ",") {
				kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
				if len(kv) == 2 {
					mapValues[kv[0]] = types.StringValue(kv[1])
				}
			}

			m.PruneBackups = types.MapValueMust(types.StringType, mapValues)
		} else {
			m.PruneBackups = types.MapNull(types.StringType)
		}
	} else {
		m.PruneBackups = types.MapNull(types.StringType)
	}
}
