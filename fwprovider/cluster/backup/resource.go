/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package backup

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/backup"
)

var (
	_ resource.Resource                     = &backupJobResource{}
	_ resource.ResourceWithConfigure        = &backupJobResource{}
	_ resource.ResourceWithImportState      = &backupJobResource{}
	_ resource.ResourceWithConfigValidators = &backupJobResource{}
)

type backupJobResource struct {
	client *backup.Client
}

// NewResource creates a new backup job resource.
func NewResource() resource.Resource {
	return &backupJobResource{}
}

func (r *backupJobResource) Metadata(
	_ context.Context,
	_ resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "proxmox_backup_job"
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

func (r *backupJobResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("all"),
			path.MatchRoot("vmid"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("all"),
			path.MatchRoot("pool"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("vmid"),
			path.MatchRoot("pool"),
		),
	}
}

func (r *backupJobResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE cluster backup job.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The identifier of the backup job.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"schedule": schema.StringAttribute{
				Description: "Backup schedule in systemd calendar event format.",
				Required:    true,
			},
			"storage": schema.StringAttribute{
				Description: "The storage identifier for the backup.",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the backup job is enabled.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"node": schema.StringAttribute{
				Description: "The cluster node name to limit the backup job to.",
				Optional:    true,
			},
			"vmid": schema.ListAttribute{
				Description: "A list of guest VM/CT IDs to include in the backup job.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"all": schema.BoolAttribute{
				Description: "Whether to back up all known guests on the node.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"mode": schema.StringAttribute{
				Description: "The backup mode (snapshot, suspend, or stop).",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("snapshot", "suspend", "stop"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"compress": schema.StringAttribute{
				Description: "The compression algorithm (0, 1, gzip, lzo, or zstd).",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("0", "1", "gzip", "lzo", "zstd"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"starttime": schema.StringAttribute{
				Description: "The scheduled start time (HH:MM).",
				Optional:    true,
			},
			"maxfiles": schema.Int64Attribute{
				Description: "Deprecated: use prune_backups instead. Maximum number of backup files per guest.",
				Optional:    true,
			},
			"mailto": schema.ListAttribute{
				Description: "A list of email addresses to send notifications to.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"mailnotification": schema.StringAttribute{
				Description: "Email notification setting (always or failure).",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("always", "failure"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bwlimit": schema.Int64Attribute{
				Description: "I/O bandwidth limit in KiB/s.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"ionice": schema.Int64Attribute{
				Description: "I/O priority (0-8).",
				Optional:    true,
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 8),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"pigz": schema.Int64Attribute{
				Description: "Number of pigz threads (0 disables, 1 uses single-threaded gzip).",
				Optional:    true,
			},
			"zstd": schema.Int64Attribute{
				Description: "Number of zstd threads (0 uses half of available cores).",
				Optional:    true,
			},
			"prune_backups": schema.MapAttribute{
				Description: "Retention options as a map of keep policies " +
					"(e.g. keep-last = \"3\", keep-weekly = \"2\").",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"remove": schema.BoolAttribute{
				Description: "Whether to remove old backups if there are more than maxfiles.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"notes_template": schema.StringAttribute{
				Description: "Template for notes attached to the backup.",
				Optional:    true,
			},
			"protected": schema.BoolAttribute{
				Description: "Whether the backup should be marked as protected.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"repeat_missed": schema.BoolAttribute{
				Description: "Whether to repeat missed backup jobs as soon as possible.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"script": schema.StringAttribute{
				Description: "Path to a script to execute before/after the backup job.",
				Optional:    true,
			},
			"stdexcludes": schema.BoolAttribute{
				Description: "Whether to exclude common temporary files from the backup.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"exclude_path": schema.ListAttribute{
				Description: "A list of paths to exclude from the backup.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"pool": schema.StringAttribute{
				Description: "Limit backup to guests in the specified pool.",
				Optional:    true,
			},
			"fleecing": schema.SingleNestedAttribute{
				Description: "Fleecing configuration for the backup job.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: "Whether fleecing is enabled.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"storage": schema.StringAttribute{
						Description: "The storage identifier for fleecing.",
						Optional:    true,
					},
				},
			},
			"performance": schema.SingleNestedAttribute{
				Description: "Performance-related settings for the backup job.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"max_workers": schema.Int64Attribute{
						Description: "Maximum number of workers for parallel backup.",
						Optional:    true,
					},
					"pbs_entries_max": schema.Int64Attribute{
						Description: "Maximum number of entries for PBS catalog.",
						Optional:    true,
					},
				},
			},
			"pbs_change_detection_mode": schema.StringAttribute{
				Description: "PBS change detection mode (legacy, data, or metadata).",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("legacy", "data", "metadata"),
				},
			},
			"lockwait": schema.Int64Attribute{
				Description: "Maximum wait time in minutes for the global lock.",
				Optional:    true,
			},
			"stopwait": schema.Int64Attribute{
				Description: "Maximum wait time in minutes for a guest to stop.",
				Optional:    true,
			},
			"tmpdir": schema.StringAttribute{
				Description: "Path to the temporary directory for the backup job.",
				Optional:    true,
			},
		},
	}
}

func (r *backupJobResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan backupJobModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createBody := plan.toCreateAPI(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Create(ctx, createBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Backup Job", err.Error())
		return
	}

	// Read back to get server-assigned defaults.
	data, err := r.client.Get(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Backup Job After Creation", err.Error())
		return
	}

	resp.Diagnostics.Append(plan.fromAPI(ctx, data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *backupJobResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state backupJobModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.Get(ctx, state.ID.ValueString())
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Unable to Read Backup Job", err.Error())

		return
	}

	resp.Diagnostics.Append(state.fromAPI(ctx, data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *backupJobResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state backupJobModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateBody := plan.toUpdateAPI(ctx, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Update(ctx, plan.ID.ValueString(), updateBody)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Update Backup Job", err.Error())
		return
	}

	// Read back to get server state.
	data, err := r.client.Get(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Backup Job After Update", err.Error())
		return
	}

	resp.Diagnostics.Append(plan.fromAPI(ctx, data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *backupJobResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state backupJobModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, state.ID.ValueString())
	if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
		resp.Diagnostics.AddError("Unable to Delete Backup Job", err.Error())
	}
}

func (r *backupJobResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	data, err := r.client.Get(ctx, req.ID)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError("Backup Job Not Found",
				fmt.Sprintf("Backup job with ID '%s' was not found", req.ID))

			return
		}

		resp.Diagnostics.AddError("Unable to Import Backup Job", err.Error())

		return
	}

	var model backupJobModel

	resp.Diagnostics.Append(model.fromAPI(ctx, data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
