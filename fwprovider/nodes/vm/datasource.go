/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &Datasource{}
	_ datasource.DataSourceWithConfigure = &Datasource{}
)

// Datasource is the implementation of VM datasource.
type Datasource struct {
	client proxmox.Client
}

// NewDataSource creates a new VM datasource.
func NewDataSource() datasource.DataSource {
	return &Datasource{}
}

// Metadata defines the name of the resource.
func (d *Datasource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_vm2"
}

// Configure sets the client for the resource.
func (d *Datasource) Configure(
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

	d.client = cfg.Client
}

func (d *Datasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model DatasourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := model.Timeouts.Read(ctx, defaultReadTimeout)
	resp.Diagnostics.Append(diags...)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	exists := readForDatasource(ctx, d.client, &model, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	if !exists {
		resp.Diagnostics.AddError(
			"VM Not Found",
			fmt.Sprintf("VM with ID %d was not found on node '%s'", model.ID.ValueInt64(), model.NodeName.ValueString()),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
