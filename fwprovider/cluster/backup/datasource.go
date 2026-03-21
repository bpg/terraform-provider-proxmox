/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package backup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/backup"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &backupJobsDataSource{}
	_ datasource.DataSourceWithConfigure = &backupJobsDataSource{}
)

type backupJobsDataSource struct {
	client *backup.Client
}

// backupJobsDataSourceModel is the top-level model for the backup jobs data source.
type backupJobsDataSourceModel struct {
	ID   types.String               `tfsdk:"id"`
	Jobs []backupJobDatasourceModel `tfsdk:"jobs"`
}

// NewDataSource creates a new backup jobs data source.
func NewDataSource() datasource.DataSource {
	return &backupJobsDataSource{}
}

func (d *backupJobsDataSource) Metadata(
	_ context.Context,
	_ datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "proxmox_backup_jobs"
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
			"Unexpected DataSource Configure Type",
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
		Description: "Retrieves the list of cluster-wide backup jobs.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier for this data source.",
				Computed:    true,
			},
			"jobs": schema.ListNestedAttribute{
				Description: "List of backup jobs.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier of the backup job.",
							Computed:    true,
						},
						"schedule": schema.StringAttribute{
							Description: "Backup schedule in systemd calendar format.",
							Computed:    true,
						},
						"storage": schema.StringAttribute{
							Description: "Target storage for the backup.",
							Computed:    true,
						},
						"enabled": schema.BoolAttribute{
							Description: "Indicates whether the backup job is enabled.",
							Computed:    true,
						},
						"node": schema.StringAttribute{
							Description: "Node on which the backup job runs.",
							Computed:    true,
						},
						"vmid": schema.ListAttribute{
							Description: "List of VM/CT IDs included in the backup job.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"all": schema.BoolAttribute{
							Description: "Indicates whether all VMs and CTs are backed up.",
							Computed:    true,
						},
						"mode": schema.StringAttribute{
							Description: "Backup mode (e.g. snapshot, suspend, stop).",
							Computed:    true,
						},
						"compress": schema.StringAttribute{
							Description: "Compression algorithm used for the backup.",
							Computed:    true,
						},
						"mailto": schema.StringAttribute{
							Description: "Comma-separated list of email addresses for notifications.",
							Computed:    true,
						},
						"mailnotification": schema.StringAttribute{
							Description: "When to send email notifications (always or failure).",
							Computed:    true,
						},
						"notes_template": schema.StringAttribute{
							Description: "Template for backup notes.",
							Computed:    true,
						},
						"pool": schema.StringAttribute{
							Description: "Pool whose members are backed up.",
							Computed:    true,
						},
						"prune_backups": schema.StringAttribute{
							Description: "Prune options in the format `keep-last=N,...`.",
							Computed:    true,
						},
						"protected": schema.BoolAttribute{
							Description: "Indicates whether backups created by this job are protected from pruning.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *backupJobsDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	jobs, err := d.client.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Backup Jobs", err.Error())
		return
	}

	var state backupJobsDataSourceModel

	state.ID = types.StringValue("backup_jobs")
	state.Jobs = make([]backupJobDatasourceModel, len(jobs))

	for i, job := range jobs {
		state.Jobs[i].fromAPI(job)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
