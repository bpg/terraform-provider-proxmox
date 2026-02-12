/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package backup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	backupAPI "github.com/bpg/terraform-provider-proxmox/proxmox/cluster/backup"
)

var (
	_ datasource.DataSource              = &backupJobsDataSource{}
	_ datasource.DataSourceWithConfigure = &backupJobsDataSource{}
)

type backupJobsDataSource struct {
	client *backupAPI.Client
}

// BackupJobsDataSourceModel maps the data source schema data.
type BackupJobsDataSourceModel struct {
	Jobs types.List `tfsdk:"jobs"`
}

// BackupJobDataModel represents a single backup job in the list.
type BackupJobDataModel struct {
	ID               types.String `tfsdk:"id"`
	Enabled          types.Bool   `tfsdk:"enabled"`
	Schedule         types.String `tfsdk:"schedule"`
	Storage          types.String `tfsdk:"storage"`
	Node             types.String `tfsdk:"node"`
	VMID             types.String `tfsdk:"vmid"`
	All              types.Bool   `tfsdk:"all"`
	Mode             types.String `tfsdk:"mode"`
	Compress         types.String `tfsdk:"compress"`
	MailTo           types.String `tfsdk:"mailto"`
	MailNotification types.String `tfsdk:"mailnotification"`
	PruneBackups     types.String `tfsdk:"prune_backups"`
	NotesTemplate    types.String `tfsdk:"notes_template"`
	Protected        types.Bool   `tfsdk:"protected"`
	Pool             types.String `tfsdk:"pool"`
}

// NewBackupJobsDataSource creates a new data source.
func NewBackupJobsDataSource() datasource.DataSource {
	return &backupJobsDataSource{}
}

func (d *backupJobsDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_backup_jobs"
}

func (d *backupJobsDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.DataSource)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected config.DataSource, got: %T", req.ProviderData),
		)

		return
	}

	d.client = cfg.Client.Cluster().Backup()
}

func (d *backupJobsDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a list of Proxmox VE cluster backup jobs.",
		MarkdownDescription: "Retrieves information about all backup jobs configured in the Proxmox VE cluster.",
		Attributes: map[string]schema.Attribute{
			"jobs": schema.ListNestedAttribute{
				Description:         "List of backup jobs",
				MarkdownDescription: "List of all backup jobs configured in the cluster.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description:         "The unique identifier of the backup job.",
							MarkdownDescription: "The unique identifier of the backup job.",
							Computed:            true,
						},
						"enabled": schema.BoolAttribute{
							Description:         "Whether the backup job is enabled.",
							MarkdownDescription: "Whether the backup job is enabled.",
							Computed:            true,
						},
						"schedule": schema.StringAttribute{
							Description:         "The backup schedule in systemd calendar event format.",
							MarkdownDescription: "The backup schedule in systemd calendar event format.",
							Computed:            true,
						},
						"storage": schema.StringAttribute{
							Description:         "The storage where backups are stored.",
							MarkdownDescription: "The storage identifier where backups are stored.",
							Computed:            true,
						},
						"node": schema.StringAttribute{
							Description:         "The node where the backup job runs.",
							MarkdownDescription: "The node name where the backup job runs. Empty if it runs on all nodes.",
							Computed:            true,
						},
						"vmid": schema.StringAttribute{
							Description:         "Comma-separated list of VM/container IDs to backup.",
							MarkdownDescription: "Comma-separated list of VM/container IDs to backup.",
							Computed:            true,
						},
						"all": schema.BoolAttribute{
							Description:         "Whether to backup all VMs and containers.",
							MarkdownDescription: "Whether to backup all VMs and containers on the node.",
							Computed:            true,
						},
						"mode": schema.StringAttribute{
							Description:         "The backup mode.",
							MarkdownDescription: "The backup mode (snapshot, suspend, or stop).",
							Computed:            true,
						},
						"compress": schema.StringAttribute{
							Description:         "The compression algorithm.",
							MarkdownDescription: "The compression algorithm used for backups.",
							Computed:            true,
						},
						"mailto": schema.StringAttribute{
							Description:         "Email recipients for notifications.",
							MarkdownDescription: "Email recipients for backup notifications.",
							Computed:            true,
						},
						"mailnotification": schema.StringAttribute{
							Description:         "When to send email notifications.",
							MarkdownDescription: "When to send email notifications (always or failure).",
							Computed:            true,
						},
						"prune_backups": schema.StringAttribute{
							Description:         "Retention policy for backups.",
							MarkdownDescription: "Retention policy for backups.",
							Computed:            true,
						},
						"notes_template": schema.StringAttribute{
							Description:         "Template for backup notes.",
							MarkdownDescription: "Template string for backup notes with variable substitution.",
							Computed:            true,
						},
						"protected": schema.BoolAttribute{
							Description:         "Whether backups are marked as protected.",
							MarkdownDescription: "Whether backups are marked as protected from pruning.",
							Computed:            true,
						},
						"pool": schema.StringAttribute{
							Description:         "The pool to backup.",
							MarkdownDescription: "The resource pool to backup.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *backupJobsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state BackupJobsDataSourceModel

	jobsData, err := d.client.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading backup jobs",
			fmt.Sprintf("Could not read backup jobs: %s", err.Error()),
		)

		return
	}

	jobModels := make([]BackupJobDataModel, len(jobsData))

	for i, job := range jobsData {
		jobModels[i] = BackupJobDataModel{
			ID:               types.StringValue(job.ID),
			Schedule:         types.StringValue(job.Schedule),
			Storage:          types.StringValue(job.Storage),
			Enabled:          customBoolPtrOr(job.Enabled, true),
			Node:             stringPtrOr(job.Node, types.StringValue("")),
			VMID:             stringPtrOr(job.VMID, types.StringValue("")),
			All:              customBoolPtrOr(job.All, false),
			Mode:             stringPtrOr(job.Mode, types.StringValue("")),
			Compress:         stringPtrOr(job.Compress, types.StringValue("")),
			MailTo:           stringPtrOr(job.MailTo, types.StringValue("")),
			MailNotification: stringPtrOr(job.MailNotification, types.StringValue("")),
			NotesTemplate:    stringPtrOr(job.NotesTemplate, types.StringValue("")),
			Protected:        customBoolPtrOr(job.Protected, false),
			Pool:             stringPtrOr(job.Pool, types.StringValue("")),
		}

		if job.PruneBackups != nil {
			jobModels[i].PruneBackups = types.StringPointerValue(job.PruneBackups.Pointer())
		} else {
			jobModels[i].PruneBackups = types.StringValue("")
		}
	}

	jobAttrTypes := map[string]attr.Type{
		"id":               types.StringType,
		"enabled":          types.BoolType,
		"schedule":         types.StringType,
		"storage":          types.StringType,
		"node":             types.StringType,
		"vmid":             types.StringType,
		"all":              types.BoolType,
		"mode":             types.StringType,
		"compress":         types.StringType,
		"mailto":           types.StringType,
		"mailnotification": types.StringType,
		"prune_backups":    types.StringType,
		"notes_template":   types.StringType,
		"protected":        types.BoolType,
		"pool":             types.StringType,
	}

	jobsList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: jobAttrTypes}, jobModels)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	state.Jobs = jobsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
