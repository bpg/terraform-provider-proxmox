/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package replication

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
)

var _ datasource.DataSource = &DataSource{}

var _ datasource.DataSourceWithConfigure = &DataSource{}

type DataSource struct {
	client *cluster.Client
}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_replication"
}

func (d *DataSource) Configure(
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
			"Unexpected Provider Data",
			fmt.Sprintf("Expected config.DataSource, got: %T", req.ProviderData),
		)

		return
	}

	d.client = cfg.Client.Cluster()
}

func (d *DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about an existing Replication.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Replication Job ID. The ID is composed of a Guest ID and a job number, separated by a hyphen, i.e. '<GUEST>-<JOBNUM>'.",
			},
			"target": schema.StringAttribute{
				Computed:    true,
				Description: "Target node.",
			},
			"jobnum": schema.Int64Attribute{
				Computed:    true,
				Description: "Unique, sequential ID assigned to each job.",
			},
			"guest": schema.Int64Attribute{
				Computed:    true,
				Description: "Guest ID.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "Section type.",
			},
			"comment": schema.StringAttribute{
				Computed:    true,
				Description: "Description.",
			},
			"disable": schema.BoolAttribute{
				Computed:    true,
				Description: "Flag to disable/deactivate this replication.",
			},
			"rate": schema.Float64Attribute{
				Computed:    true,
				Description: "Rate limit in mbps (megabytes per second) as floating point number.",
			},
			"schedule": schema.StringAttribute{
				Computed:    true,
				Description: "Storage replication schedule. The format is a subset of `systemd` calendar events. Defaults to */15",
			},
			"source": schema.StringAttribute{
				Computed:    true,
				Description: "For internal use, to detect if the guest was stolen.",
			},
		},
	}
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var readModel model

	resp.Diagnostics.Append(req.Config.Get(ctx, &readModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	repl, err := d.client.Replication(readModel.ID.ValueString()).GetReplication(ctx)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError("Replication Not Found", fmt.Sprintf("Replication with ID '%s' was not found", readModel.ID.ValueString()))
			return
		}

		resp.Diagnostics.AddError("Unable to Read Replication", err.Error())

		return
	}

	state := model{}
	state.fromAPI(readModel.ID.ValueString(), repl)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
