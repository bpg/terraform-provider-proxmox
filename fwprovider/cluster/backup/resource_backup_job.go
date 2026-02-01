/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package backup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	backupAPI "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/backup"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

var (
	_ resource.Resource                = &backupJobResource{}
	_ resource.ResourceWithConfigure   = &backupJobResource{}
	_ resource.ResourceWithImportState = &backupJobResource{}
)

type backupJobResource struct {
	client *backupAPI.Client
}

// NewBackupJobResource creates a new backup job resource.
func NewBackupJobResource() resource.Resource {
	return &backupJobResource{}
}

func (r *backupJobResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_backup_job"
}

func (r *backupJobResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.Resource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected config.Resource, got: %T", req.ProviderData),
		)

		return
	}

	r.client = cfg.Client.Cluster().Backup()
}

func (r *backupJobResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE cluster backup job (vzdump).",
		MarkdownDescription: "Manages a Proxmox VE cluster backup job. Backup jobs are scheduled tasks that " +
			"automatically backup VMs and containers according to a defined schedule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "The unique identifier of the backup job.",
				MarkdownDescription: "The unique identifier of the backup job. This will be used as the job ID in Proxmox VE.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description:         "Enable or disable the backup job.",
				MarkdownDescription: "Enable or disable the backup job. When disabled, the job will not run according to its schedule.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"schedule": schema.StringAttribute{
				Description: "The schedule for the backup job in systemd calendar event format.",
				MarkdownDescription: "The schedule for the backup job in systemd calendar event format. Examples: `0 2 * * *` (daily at 2 AM), " +
					"`0 */4 * * *` (every 4 hours), `*:0/30` (every 30 minutes).",
				Required: true,
			},
			"storage": schema.StringAttribute{
				Description:         "The storage where backups will be stored.",
				MarkdownDescription: "The storage identifier where backups will be stored (e.g., `local`, `pbs`).",
				Required:            true,
			},
			"node": schema.StringAttribute{
				Description: "The node name where the backup job will run.",
				MarkdownDescription: "The node name where the backup job will run. If not specified, the job can run on any node. " +
					"Use this to restrict backup execution to a specific node.",
				Optional: true,
			},
			"vmid": schema.StringAttribute{
				Description:         "Comma-separated list of VM/container IDs to backup.",
				MarkdownDescription: "Comma-separated list of VM/container IDs to backup (e.g., `100,101,102`). Mutually exclusive with `all` and `pool`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("all"), path.MatchRoot("pool")),
				},
			},
			"all": schema.BoolAttribute{
				Description:         "Backup all VMs and containers.",
				MarkdownDescription: "Backup all VMs and containers on the node. Mutually exclusive with `vmid` and `pool`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(path.MatchRoot("vmid"), path.MatchRoot("pool")),
				},
			},
			"pool": schema.StringAttribute{
				Description:         "Backup all VMs in the specified pool.",
				MarkdownDescription: "Backup all VMs and containers in the specified pool. Mutually exclusive with `vmid` and `all`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("vmid"), path.MatchRoot("all")),
				},
			},
			"mode": schema.StringAttribute{
				Description: "The backup mode.",
				MarkdownDescription: "The backup mode. Options: `snapshot` (minimal downtime), `suspend` (suspend during backup), " +
					"`stop` (stop VM during backup for maximum consistency). Default is `snapshot`.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("snapshot"),
				Validators: []validator.String{
					stringvalidator.OneOf("snapshot", "suspend", "stop"),
				},
			},
			"compress": schema.StringAttribute{
				Description: "The compression algorithm.",
				MarkdownDescription: "The compression algorithm. Options: `0` (no compression), `1` (default compression), " +
					"`gzip`, `lzo`, `zstd`. Default is `0`.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("0"),
				Validators: []validator.String{
					stringvalidator.OneOf("0", "1", "gzip", "lzo", "zstd"),
				},
			},
			"starttime": schema.StringAttribute{
				Description:         "The start time in HH:MM format.",
				MarkdownDescription: "The start time in HH:MM format (24-hour). This works in conjunction with the schedule.",
				Optional:            true,
			},
			"maxfiles": schema.Int64Attribute{
				Description: "Maximum number of backup files per VM (deprecated, use prune_backups).",
				MarkdownDescription: "Maximum number of backup files per VM. **Deprecated:** Use `prune_backups` instead. " +
					"This option is maintained for backward compatibility.",
				Optional: true,
			},
			"mailto": schema.StringAttribute{
				Description: "Email recipients for notifications.",
				MarkdownDescription: "Email recipients for notifications. Can be comma-separated email addresses or Proxmox VE user names. " +
					"Example: `admin@example.com,root`",
				Optional: true,
			},
			"mailnotification": schema.StringAttribute{
				Description:         "When to send email notifications.",
				MarkdownDescription: "When to send email notifications. Options: `always` (send on success and failure), `failure` (send only on failure).",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("always"),
				Validators: []validator.String{
					stringvalidator.OneOf("always", "failure"),
				},
			},
			"bwlimit": schema.Int64Attribute{
				Description:         "Bandwidth limit in KiB/s.",
				MarkdownDescription: "Bandwidth limit in KiB/s. 0 means unlimited. Use this to prevent storage I/O bottlenecks.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"ionice": schema.Int64Attribute{
				Description: "I/O priority (0-8).",
				MarkdownDescription: "I/O priority when using BFQ scheduler (0-8). Higher values mean lower priority. " +
					"Default is 7.",
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(7),
				Validators: []validator.Int64{
					int64validator.Between(0, 8),
				},
			},
			"pigz": schema.Int64Attribute{
				Description:         "Use pigz for parallel gzip compression.",
				MarkdownDescription: "Use pigz instead of gzip for parallel compression. Specify the number of threads (0 means disabled).",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"zstd": schema.Int64Attribute{
				Description: "Zstd thread count.",
				MarkdownDescription: "Zstd compression thread count. 0 means half of available CPU cores. " +
					"Only relevant when compress is set to `zstd`.",
				Optional: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"prune_backups": schema.StringAttribute{
				Description: "Retention policy for backups.",
				MarkdownDescription: "Retention policy for backups in the format " +
					"`keep-last=N,keep-hourly=N,keep-daily=N,keep-weekly=N,keep-monthly=N,keep-yearly=N`. " +
					"Example: `keep-last=3,keep-daily=7,keep-weekly=4`. Default is `keep-all=1` (keep all backups).",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("keep-all=1"),
			},
			"remove": schema.BoolAttribute{
				Description: "Prune old backups according to retention policy.",
				MarkdownDescription: "Prune old backups according to the `prune_backups` retention policy. " +
					"When disabled, old backups are not automatically removed.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"notes_template": schema.StringAttribute{
				Description: "Template for backup notes.",
				MarkdownDescription: "Template for backup notes. Available variables: `{{vmid}}`, `{{guestname}}`, `{{node}}`, `{{cluster}}`. " +
					"Example: `{{guestname}}-{{cluster}}`",
				Optional: true,
			},
			"protected": schema.BoolAttribute{
				Description:         "Mark backups as protected.",
				MarkdownDescription: "Mark backups as protected. Protected backups are not removed by pruning and are not counted toward retention limits.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"repeat_missed": schema.BoolAttribute{
				Description: "Run missed jobs when node becomes available.",
				MarkdownDescription: "Run missed backup jobs when the node becomes available. " +
					"Useful for ensuring backups are not skipped when a node is offline during scheduled time.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"script": schema.StringAttribute{
				Description:         "Hook script path.",
				MarkdownDescription: "Path to a hook script that will be called during backup lifecycle events (before/after backup, on failure, etc.).",
				Optional:            true,
			},
			"stdexcludes": schema.BoolAttribute{
				Description:         "Exclude temporary files and logs.",
				MarkdownDescription: "Exclude standard temporary files and logs from the backup. Recommended for most use cases.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"exclude_path": schema.ListAttribute{
				Description:         "Shell glob patterns to exclude from backup.",
				MarkdownDescription: "Array of shell glob patterns to exclude from backup. Example: `[\"/tmp/**\", \"/var/log/**\"]`",
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"fleecing": schema.SingleNestedAttribute{
				Description: "Backup fleecing configuration.",
				MarkdownDescription: "Backup fleecing configuration. Fleecing caches backup data on fast local storage first, " +
					"reducing I/O impact on guest VMs.",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description:         "Enable backup fleecing.",
						MarkdownDescription: "Enable backup fleecing for this job.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"storage": schema.StringAttribute{
						Description:         "Storage for fleecing cache.",
						MarkdownDescription: "Storage identifier for the fleecing cache. Should be fast local storage.",
						Optional:            true,
					},
				},
			},
			"performance": schema.SingleNestedAttribute{
				Description:         "Performance tuning configuration.",
				MarkdownDescription: "Performance tuning options for backup operations.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"max_workers": schema.Int64Attribute{
						Description:         "Maximum concurrent backup workers.",
						MarkdownDescription: "Maximum number of concurrent backup workers (1-256). Controls parallel backup operations.",
						Optional:            true,
						Validators: []validator.Int64{
							int64validator.Between(1, 256),
						},
					},
					"pbs_entries_max": schema.Int64Attribute{
						Description: "PBS entries maximum.",
						MarkdownDescription: "Maximum number of entries for Proxmox Backup Server. " +
							"Affects memory usage for container backups.",
						Optional: true,
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
						},
					},
				},
			},
			"pbs_change_detection_mode": schema.StringAttribute{
				Description: "PBS change detection mode for containers.",
				MarkdownDescription: "Proxmox Backup Server change detection mode for container backups. " +
					"Options: `data` (hash-based), `metadata` (metadata-based), `legacy` (legacy mode).",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("data", "metadata", "legacy"),
				},
			},
			"lockwait": schema.Int64Attribute{
				Description: "Maximum time to wait for VM lock (minutes).",
				MarkdownDescription: "Maximum time to wait for the global VM lock in minutes. " +
					"Default is 180 minutes.",
				Optional: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"stopwait": schema.Int64Attribute{
				Description:         "Maximum time to wait for VM shutdown (minutes).",
				MarkdownDescription: "Maximum time to wait for guest shutdown in minutes when using `stop` mode. Default is 10 minutes.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"tmpdir": schema.StringAttribute{
				Description:         "Temporary directory for backup data.",
				MarkdownDescription: "Custom temporary directory for storing backup data during the backup process.",
				Optional:            true,
			},
		},
	}
}

func (r *backupJobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan BackupJobModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	reqBody := r.planToAPICreate(ctx, &plan)

	err := r.client.Create(ctx, reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating backup job",
			fmt.Sprintf("Could not create backup job %s: %s", plan.ID.ValueString(), err.Error()),
		)

		return
	}

	jobData, err := r.client.Get(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading backup job after creation",
			fmt.Sprintf("Could not read backup job %s: %s", plan.ID.ValueString(), err.Error()),
		)

		return
	}

	r.apiToModel(jobData, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *backupJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state BackupJobModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jobData, err := r.client.Get(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading backup job",
			fmt.Sprintf("Could not read backup job %s: %s", state.ID.ValueString(), err.Error()),
		)

		return
	}

	r.apiToModel(jobData, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *backupJobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan BackupJobModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	reqBody := r.planToAPIUpdate(ctx, &plan)

	err := r.client.Update(ctx, plan.ID.ValueString(), reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating backup job",
			fmt.Sprintf("Could not update backup job %s: %s", plan.ID.ValueString(), err.Error()),
		)

		return
	}

	jobData, err := r.client.Get(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading backup job after update",
			fmt.Sprintf("Could not read backup job %s: %s", plan.ID.ValueString(), err.Error()),
		)

		return
	}

	r.apiToModel(jobData, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *backupJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state BackupJobModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting backup job",
			fmt.Sprintf("Could not delete backup job %s: %s", state.ID.ValueString(), err.Error()),
		)
	}
}

func (r *backupJobResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *backupJobResource) planToAPICreate(ctx context.Context, plan *BackupJobModel) *backupAPI.CreateRequestBody {
	reqBody := &backupAPI.CreateRequestBody{
		ID:       plan.ID.ValueString(),
		Schedule: plan.Schedule.ValueString(),
		Storage:  plan.Storage.ValueString(),
	}

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		enabled := proxmoxtypes.CustomBool(plan.Enabled.ValueBool())
		reqBody.Enabled = &enabled
	}

	if !plan.Node.IsNull() {
		reqBody.Node = plan.Node.ValueStringPointer()
	}

	if !plan.VMID.IsNull() {
		reqBody.VMID = plan.VMID.ValueStringPointer()
	}

	if !plan.All.IsNull() && !plan.All.IsUnknown() {
		all := proxmoxtypes.CustomBool(plan.All.ValueBool())
		reqBody.All = &all
	}

	if !plan.Mode.IsNull() {
		reqBody.Mode = plan.Mode.ValueStringPointer()
	}

	if !plan.Compress.IsNull() {
		reqBody.Compress = plan.Compress.ValueStringPointer()
	}

	if !plan.StartTime.IsNull() {
		reqBody.StartTime = plan.StartTime.ValueStringPointer()
	}

	if !plan.MaxFiles.IsNull() {
		maxfiles := int(plan.MaxFiles.ValueInt64())
		reqBody.MaxFiles = &maxfiles
	}

	if !plan.MailTo.IsNull() {
		reqBody.MailTo = plan.MailTo.ValueStringPointer()
	}

	if !plan.MailNotification.IsNull() {
		reqBody.MailNotification = plan.MailNotification.ValueStringPointer()
	}

	if !plan.BwLimit.IsNull() {
		bwlimit := int(plan.BwLimit.ValueInt64())
		reqBody.BwLimit = &bwlimit
	}

	if !plan.IONice.IsNull() {
		ionice := int(plan.IONice.ValueInt64())
		reqBody.IONice = &ionice
	}

	if !plan.Pigz.IsNull() {
		pigz := int(plan.Pigz.ValueInt64())
		reqBody.Pigz = &pigz
	}

	if !plan.Zstd.IsNull() {
		zstd := int(plan.Zstd.ValueInt64())
		reqBody.Zstd = &zstd
	}

	if !plan.PruneBackups.IsNull() {
		reqBody.PruneBackups = plan.PruneBackups.ValueStringPointer()
	}

	if !plan.Remove.IsNull() && !plan.Remove.IsUnknown() {
		remove := proxmoxtypes.CustomBool(plan.Remove.ValueBool())
		reqBody.Remove = &remove
	}

	if !plan.NotesTemplate.IsNull() {
		reqBody.NotesTemplate = plan.NotesTemplate.ValueStringPointer()
	}

	if !plan.Protected.IsNull() && !plan.Protected.IsUnknown() {
		protected := proxmoxtypes.CustomBool(plan.Protected.ValueBool())
		reqBody.Protected = &protected
	}

	if !plan.RepeatMissed.IsNull() && !plan.RepeatMissed.IsUnknown() {
		repeatMissed := proxmoxtypes.CustomBool(plan.RepeatMissed.ValueBool())
		reqBody.RepeatMissed = &repeatMissed
	}

	if !plan.Script.IsNull() {
		reqBody.Script = plan.Script.ValueStringPointer()
	}

	if !plan.StdExcludes.IsNull() && !plan.StdExcludes.IsUnknown() {
		stdExcludes := proxmoxtypes.CustomBool(plan.StdExcludes.ValueBool())
		reqBody.StdExcludes = &stdExcludes
	}

	if !plan.ExcludePath.IsNull() && !plan.ExcludePath.IsUnknown() {
		var excludePath []string
		plan.ExcludePath.ElementsAs(ctx, &excludePath, false)

		if len(excludePath) > 0 {
			commaSepList := proxmoxtypes.CustomCommaSeparatedList(excludePath)
			reqBody.ExcludePath = &commaSepList
		}
	}

	if !plan.Pool.IsNull() {
		reqBody.Pool = plan.Pool.ValueStringPointer()
	}

	if !plan.Fleecing.IsNull() && !plan.Fleecing.IsUnknown() {
		var fleecingModel FleecingModel
		plan.Fleecing.As(ctx, &fleecingModel, basetypes.ObjectAsOptions{})

		fleecing := &backupAPI.FleecingConfig{}

		if !fleecingModel.Enabled.IsNull() {
			enabled := proxmoxtypes.CustomBool(fleecingModel.Enabled.ValueBool())
			fleecing.Enabled = &enabled
		}

		if !fleecingModel.Storage.IsNull() {
			fleecing.Storage = fleecingModel.Storage.ValueStringPointer()
		}

		reqBody.Fleecing = fleecing
	}

	if !plan.Performance.IsNull() && !plan.Performance.IsUnknown() {
		var perfModel PerformanceModel
		plan.Performance.As(ctx, &perfModel, basetypes.ObjectAsOptions{})

		perf := &backupAPI.PerformanceConfig{}

		if !perfModel.MaxWorkers.IsNull() {
			maxWorkers := int(perfModel.MaxWorkers.ValueInt64())
			perf.MaxWorkers = &maxWorkers
		}

		if !perfModel.PBSEntriesMax.IsNull() {
			pbsMax := int(perfModel.PBSEntriesMax.ValueInt64())
			perf.PBSEntriesMax = &pbsMax
		}

		reqBody.Performance = perf
	}

	if !plan.PBSChangeDetectionMode.IsNull() {
		reqBody.PBSChangeDetectionMode = plan.PBSChangeDetectionMode.ValueStringPointer()
	}

	if !plan.LockWait.IsNull() {
		lockWait := int(plan.LockWait.ValueInt64())
		reqBody.LockWait = &lockWait
	}

	if !plan.StopWait.IsNull() {
		stopWait := int(plan.StopWait.ValueInt64())
		reqBody.StopWait = &stopWait
	}

	if !plan.TmpDir.IsNull() {
		reqBody.TmpDir = plan.TmpDir.ValueStringPointer()
	}

	return reqBody
}

func (r *backupJobResource) planToAPIUpdate(ctx context.Context, plan *BackupJobModel) *backupAPI.UpdateRequestBody {
	reqBody := &backupAPI.UpdateRequestBody{}

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		enabled := proxmoxtypes.CustomBool(plan.Enabled.ValueBool())
		reqBody.Enabled = &enabled
	}

	if !plan.Schedule.IsNull() {
		reqBody.Schedule = plan.Schedule.ValueStringPointer()
	}

	if !plan.Storage.IsNull() {
		reqBody.Storage = plan.Storage.ValueStringPointer()
	}

	if !plan.Node.IsNull() {
		reqBody.Node = plan.Node.ValueStringPointer()
	}

	if !plan.VMID.IsNull() {
		reqBody.VMID = plan.VMID.ValueStringPointer()
	}

	if !plan.All.IsNull() && !plan.All.IsUnknown() {
		all := proxmoxtypes.CustomBool(plan.All.ValueBool())
		reqBody.All = &all
	}

	if !plan.Mode.IsNull() {
		reqBody.Mode = plan.Mode.ValueStringPointer()
	}

	if !plan.Compress.IsNull() {
		reqBody.Compress = plan.Compress.ValueStringPointer()
	}

	if !plan.StartTime.IsNull() {
		reqBody.StartTime = plan.StartTime.ValueStringPointer()
	}

	if !plan.MaxFiles.IsNull() {
		maxfiles := int(plan.MaxFiles.ValueInt64())
		reqBody.MaxFiles = &maxfiles
	}

	if !plan.MailTo.IsNull() {
		reqBody.MailTo = plan.MailTo.ValueStringPointer()
	}

	if !plan.MailNotification.IsNull() {
		reqBody.MailNotification = plan.MailNotification.ValueStringPointer()
	}

	if !plan.BwLimit.IsNull() {
		bwlimit := int(plan.BwLimit.ValueInt64())
		reqBody.BwLimit = &bwlimit
	}

	if !plan.IONice.IsNull() {
		ionice := int(plan.IONice.ValueInt64())
		reqBody.IONice = &ionice
	}

	if !plan.Pigz.IsNull() {
		pigz := int(plan.Pigz.ValueInt64())
		reqBody.Pigz = &pigz
	}

	if !plan.Zstd.IsNull() {
		zstd := int(plan.Zstd.ValueInt64())
		reqBody.Zstd = &zstd
	}

	if !plan.PruneBackups.IsNull() {
		reqBody.PruneBackups = plan.PruneBackups.ValueStringPointer()
	}

	if !plan.Remove.IsNull() && !plan.Remove.IsUnknown() {
		remove := proxmoxtypes.CustomBool(plan.Remove.ValueBool())
		reqBody.Remove = &remove
	}

	if !plan.NotesTemplate.IsNull() {
		reqBody.NotesTemplate = plan.NotesTemplate.ValueStringPointer()
	}

	if !plan.Protected.IsNull() && !plan.Protected.IsUnknown() {
		protected := proxmoxtypes.CustomBool(plan.Protected.ValueBool())
		reqBody.Protected = &protected
	}

	if !plan.RepeatMissed.IsNull() && !plan.RepeatMissed.IsUnknown() {
		repeatMissed := proxmoxtypes.CustomBool(plan.RepeatMissed.ValueBool())
		reqBody.RepeatMissed = &repeatMissed
	}

	if !plan.Script.IsNull() {
		reqBody.Script = plan.Script.ValueStringPointer()
	}

	if !plan.StdExcludes.IsNull() && !plan.StdExcludes.IsUnknown() {
		stdExcludes := proxmoxtypes.CustomBool(plan.StdExcludes.ValueBool())
		reqBody.StdExcludes = &stdExcludes
	}

	if !plan.ExcludePath.IsNull() && !plan.ExcludePath.IsUnknown() {
		var excludePath []string
		plan.ExcludePath.ElementsAs(ctx, &excludePath, false)

		if len(excludePath) > 0 {
			commaSepList := proxmoxtypes.CustomCommaSeparatedList(excludePath)
			reqBody.ExcludePath = &commaSepList
		}
	}

	if !plan.Pool.IsNull() {
		reqBody.Pool = plan.Pool.ValueStringPointer()
	}

	if !plan.Fleecing.IsNull() && !plan.Fleecing.IsUnknown() {
		var fleecingModel FleecingModel
		plan.Fleecing.As(ctx, &fleecingModel, basetypes.ObjectAsOptions{})

		fleecing := &backupAPI.FleecingConfig{}

		if !fleecingModel.Enabled.IsNull() {
			enabled := proxmoxtypes.CustomBool(fleecingModel.Enabled.ValueBool())
			fleecing.Enabled = &enabled
		}

		if !fleecingModel.Storage.IsNull() {
			fleecing.Storage = fleecingModel.Storage.ValueStringPointer()
		}

		reqBody.Fleecing = fleecing
	}

	if !plan.Performance.IsNull() && !plan.Performance.IsUnknown() {
		var perfModel PerformanceModel
		plan.Performance.As(ctx, &perfModel, basetypes.ObjectAsOptions{})

		perf := &backupAPI.PerformanceConfig{}

		if !perfModel.MaxWorkers.IsNull() {
			maxWorkers := int(perfModel.MaxWorkers.ValueInt64())
			perf.MaxWorkers = &maxWorkers
		}

		if !perfModel.PBSEntriesMax.IsNull() {
			pbsMax := int(perfModel.PBSEntriesMax.ValueInt64())
			perf.PBSEntriesMax = &pbsMax
		}

		reqBody.Performance = perf
	}

	if !plan.PBSChangeDetectionMode.IsNull() {
		reqBody.PBSChangeDetectionMode = plan.PBSChangeDetectionMode.ValueStringPointer()
	}

	if !plan.LockWait.IsNull() {
		lockWait := int(plan.LockWait.ValueInt64())
		reqBody.LockWait = &lockWait
	}

	if !plan.StopWait.IsNull() {
		stopWait := int(plan.StopWait.ValueInt64())
		reqBody.StopWait = &stopWait
	}

	if !plan.TmpDir.IsNull() {
		reqBody.TmpDir = plan.TmpDir.ValueStringPointer()
	}

	return reqBody
}

func (r *backupJobResource) apiToModel(data *backupAPI.GetResponseData, model *BackupJobModel) {
	model.ID = types.StringValue(data.ID)
	model.Schedule = types.StringValue(data.Schedule)
	model.Storage = types.StringValue(data.Storage)

	if data.Enabled != nil {
		model.Enabled = types.BoolValue(bool(*data.Enabled))
	} else {
		model.Enabled = types.BoolValue(true)
	}

	if data.Node != nil {
		model.Node = types.StringPointerValue(data.Node)
	} else {
		model.Node = types.StringNull()
	}

	if data.VMID != nil {
		model.VMID = types.StringPointerValue(data.VMID)
	} else {
		model.VMID = types.StringNull()
	}

	if data.All != nil {
		model.All = types.BoolValue(bool(*data.All))
	} else {
		model.All = types.BoolValue(false)
	}

	if data.Mode != nil {
		model.Mode = types.StringPointerValue(data.Mode)
	} else {
		model.Mode = types.StringValue("snapshot")
	}

	if data.Compress != nil {
		model.Compress = types.StringPointerValue(data.Compress)
	} else {
		model.Compress = types.StringValue("0")
	}

	if data.StartTime != nil {
		model.StartTime = types.StringPointerValue(data.StartTime)
	} else {
		model.StartTime = types.StringNull()
	}

	if data.MaxFiles != nil {
		model.MaxFiles = types.Int64Value(int64(*data.MaxFiles))
	} else {
		model.MaxFiles = types.Int64Null()
	}

	if data.MailTo != nil {
		model.MailTo = types.StringPointerValue(data.MailTo)
	} else {
		model.MailTo = types.StringNull()
	}

	if data.MailNotification != nil {
		model.MailNotification = types.StringPointerValue(data.MailNotification)
	} else {
		model.MailNotification = types.StringValue("always")
	}

	if data.BwLimit != nil {
		model.BwLimit = types.Int64Value(int64(*data.BwLimit))
	} else {
		model.BwLimit = types.Int64Value(0)
	}

	if data.IONice != nil {
		model.IONice = types.Int64Value(int64(*data.IONice))
	} else {
		model.IONice = types.Int64Value(7)
	}

	if data.Pigz != nil {
		model.Pigz = types.Int64Value(int64(*data.Pigz))
	} else {
		model.Pigz = types.Int64Null()
	}

	if data.Zstd != nil {
		model.Zstd = types.Int64Value(int64(*data.Zstd))
	} else {
		model.Zstd = types.Int64Null()
	}

	if data.PruneBackups != nil {
		model.PruneBackups = types.StringPointerValue(data.PruneBackups.Pointer())
	} else {
		model.PruneBackups = types.StringValue("keep-all=1")
	}

	if data.Remove != nil {
		model.Remove = types.BoolValue(bool(*data.Remove))
	} else {
		model.Remove = types.BoolValue(true)
	}

	if data.NotesTemplate != nil {
		model.NotesTemplate = types.StringPointerValue(data.NotesTemplate)
	} else {
		model.NotesTemplate = types.StringNull()
	}

	if data.Protected != nil {
		model.Protected = types.BoolValue(bool(*data.Protected))
	} else {
		model.Protected = types.BoolValue(false)
	}

	if data.RepeatMissed != nil {
		model.RepeatMissed = types.BoolValue(bool(*data.RepeatMissed))
	} else {
		model.RepeatMissed = types.BoolValue(false)
	}

	if data.Script != nil {
		model.Script = types.StringPointerValue(data.Script)
	} else {
		model.Script = types.StringNull()
	}

	if data.StdExcludes != nil {
		model.StdExcludes = types.BoolValue(bool(*data.StdExcludes))
	} else {
		model.StdExcludes = types.BoolValue(true)
	}

	if data.ExcludePath != nil && len(*data.ExcludePath) > 0 {
		excludePathList := make([]attr.Value, len(*data.ExcludePath))
		for i, path := range *data.ExcludePath {
			excludePathList[i] = types.StringValue(path)
		}

		model.ExcludePath = types.ListValueMust(types.StringType, excludePathList)
	} else {
		model.ExcludePath = types.ListNull(types.StringType)
	}

	if data.Pool != nil {
		model.Pool = types.StringPointerValue(data.Pool)
	} else {
		model.Pool = types.StringNull()
	}

	if data.Fleecing != nil {
		fleecingAttrTypes := map[string]attr.Type{
			"enabled": types.BoolType,
			"storage": types.StringType,
		}

		fleecingAttrs := map[string]attr.Value{
			"enabled": types.BoolNull(),
			"storage": types.StringNull(),
		}

		if data.Fleecing.Enabled != nil {
			fleecingAttrs["enabled"] = types.BoolValue(bool(*data.Fleecing.Enabled))
		}

		if data.Fleecing.Storage != nil {
			fleecingAttrs["storage"] = types.StringPointerValue(data.Fleecing.Storage)
		}

		model.Fleecing = types.ObjectValueMust(fleecingAttrTypes, fleecingAttrs)
	} else {
		model.Fleecing = types.ObjectNull(map[string]attr.Type{
			"enabled": types.BoolType,
			"storage": types.StringType,
		})
	}

	if data.Performance != nil {
		perfAttrTypes := map[string]attr.Type{
			"max_workers":     types.Int64Type,
			"pbs_entries_max": types.Int64Type,
		}

		perfAttrs := map[string]attr.Value{
			"max_workers":     types.Int64Null(),
			"pbs_entries_max": types.Int64Null(),
		}

		if data.Performance.MaxWorkers != nil {
			perfAttrs["max_workers"] = types.Int64Value(int64(*data.Performance.MaxWorkers))
		}

		if data.Performance.PBSEntriesMax != nil {
			perfAttrs["pbs_entries_max"] = types.Int64Value(int64(*data.Performance.PBSEntriesMax))
		}

		model.Performance = types.ObjectValueMust(perfAttrTypes, perfAttrs)
	} else {
		model.Performance = types.ObjectNull(map[string]attr.Type{
			"max_workers":     types.Int64Type,
			"pbs_entries_max": types.Int64Type,
		})
	}

	if data.PBSChangeDetectionMode != nil {
		model.PBSChangeDetectionMode = types.StringPointerValue(data.PBSChangeDetectionMode)
	} else {
		model.PBSChangeDetectionMode = types.StringNull()
	}

	if data.LockWait != nil {
		model.LockWait = types.Int64Value(int64(*data.LockWait))
	} else {
		model.LockWait = types.Int64Null()
	}

	if data.StopWait != nil {
		model.StopWait = types.Int64Value(int64(*data.StopWait))
	} else {
		model.StopWait = types.Int64Null()
	}

	if data.TmpDir != nil {
		model.TmpDir = types.StringPointerValue(data.TmpDir)
	} else {
		model.TmpDir = types.StringNull()
	}
}
