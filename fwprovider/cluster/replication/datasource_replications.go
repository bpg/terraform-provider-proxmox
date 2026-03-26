/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package replication

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/replications"
)

func replicationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":         types.StringType,
		"target":     types.StringType,
		"type":       types.StringType,
		"comment":    types.StringType,
		"disable":    types.BoolType,
		"rate":       types.Float64Type,
		"remove_job": types.StringType,
		"schedule":   types.StringType,
		"source":     types.StringType,
		"guest":      types.Int64Type,
		"jobnum":     types.Int64Type,
	}
}

// Ensure the implementation satisfies the required interfaces.
var (
	_ datasource.DataSource              = &replicationsDataSource{}
	_ datasource.DataSourceWithConfigure = &replicationsDataSource{}
)

// replicationsDataSource is the data source implementation for Replications.
type replicationsDataSource struct {
	client *replications.Client
}

// replicationsDataSourceModel represents the data source model for listing Replications.
type replicationsDataSourceModel struct {
	Replications types.List `tfsdk:"replications"`
}

// replicationDataModel represents individual Replication data in the list.
type replicationDataModel struct {
	ID        types.String  `tfsdk:"id"`
	Target    types.String  `tfsdk:"target"`
	Type      types.String  `tfsdk:"type"`
	Comment   types.String  `tfsdk:"comment"`
	Disable   types.Bool    `tfsdk:"disable"`
	Rate      types.Float64 `tfsdk:"rate"`
	RemoveJob types.String  `tfsdk:"remove_job"`
	Schedule  types.String  `tfsdk:"schedule"`
	Source    types.String  `tfsdk:"source"`
	Guest     types.Int64   `tfsdk:"guest"`
	JobNum    types.Int64   `tfsdk:"jobnum"`
}

// Configure adds the provider-configured client to the data source.
func (d *replicationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = &replications.Client{Client: cfg.Client.Cluster()}
}

// Metadata returns the data source type name.
func (d *replicationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replications"
}

// Schema defines the schema for the data source.
func (d *replicationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		DeprecationMessage:  migration.DeprecationMessage("proxmox_replications"),
		Description:         "Retrieves information about all Replications in Proxmox.",
		MarkdownDescription: "Retrieves information about all Replications in Proxmox.",
		Attributes: map[string]schema.Attribute{
			"replications": schema.ListAttribute{
				Description: "List of Replications.",
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: replicationAttrTypes(),
				},
			},
		},
	}
}

// Read fetches all Replications from the Proxmox VE API.
func (d *replicationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data replicationsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	replList, err := d.client.GetReplications(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Replications",
			err.Error(),
		)

		return
	}

	// Convert Replications to list elements
	replElements := make([]attr.Value, len(replList))
	for i, repl := range replList {
		replData := replicationDataModel{
			ID:        types.StringValue(repl.ID),
			Target:    types.StringValue(repl.Target),
			Type:      types.StringValue(repl.Type),
			Comment:   types.StringPointerValue(repl.Comment),
			Disable:   types.BoolPointerValue(repl.Disable.PointerBool()),
			Rate:      types.Float64PointerValue(repl.Rate),
			RemoveJob: types.StringPointerValue(repl.RemoveJob),
			Schedule:  types.StringPointerValue(repl.Schedule),
			Source:    types.StringPointerValue(repl.Source),
			Guest:     types.Int64Value(repl.Guest),
			JobNum:    types.Int64Value(repl.JobNum),
		}

		objValue, objDiag := types.ObjectValueFrom(ctx, replicationAttrTypes(), replData)
		resp.Diagnostics.Append(objDiag...)

		if resp.Diagnostics.HasError() {
			return
		}

		replElements[i] = objValue
	}

	listValue, listDiag := types.ListValue(types.ObjectType{
		AttrTypes: replicationAttrTypes(),
	}, replElements)
	resp.Diagnostics.Append(listDiag...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Replications = listValue
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// NewReplicationsDataSource returns a new data source for Replications.
func NewReplicationsDataSource() datasource.DataSource {
	return &replicationsDataSource{}
}

// Short-name alias for the replications data source (ADR-007).

var (
	_ datasource.DataSource              = &replicationsDataSourceShort{}
	_ datasource.DataSourceWithConfigure = &replicationsDataSourceShort{}
)

type replicationsDataSourceShort struct {
	replicationsDataSource
}

// NewReplicationsShortDataSource creates a short-name alias for the replications data source.
func NewReplicationsShortDataSource() datasource.DataSource {
	return &replicationsDataSourceShort{}
}

func (d *replicationsDataSourceShort) Metadata(
	_ context.Context,
	_ datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "proxmox_replications"
}

func (d *replicationsDataSourceShort) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	d.replicationsDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}
