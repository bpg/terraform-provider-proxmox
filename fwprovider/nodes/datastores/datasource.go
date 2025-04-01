/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datastores

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &Datasource{}
	_ datasource.DataSourceWithConfigure = &Datasource{}
)

// Datasource is the implementation of datastores datasource.
type Datasource struct {
	client proxmox.Client
}

// NewDataSource creates a new datastores datasource.
func NewDataSource() datasource.DataSource {
	return &Datasource{}
}

// Metadata defines the name of the resource.
func (d *Datasource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_datastores"
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
	var model Model

	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}

	storageAPI := d.client.Node(model.NodeName.ValueString()).Storage("")

	r := storage.DatastoreListRequestBody{}
	if model.Filters != nil {
		r.ContentTypes = model.Filters.ContentTypes.ValueList(ctx, &resp.Diagnostics)
		r.ID = model.Filters.ID.ValueStringPointer()
		r.Target = model.Filters.Target.ValueStringPointer()
	}

	dsList, err := storageAPI.ListDatastores(ctx, &r)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read datastores",
			err.Error(),
		)

		return
	}

	model.Datastores = make([]Datastore, 0, len(dsList))

	for _, ds := range dsList {
		datastore := Datastore{}

		if ds.ContentTypes != nil {
			datastore.ContentTypes = stringset.NewValueList(*ds.ContentTypes, &resp.Diagnostics)
		}

		datastore.Active = types.BoolPointerValue(ds.Active.PointerBool())
		datastore.Enabled = types.BoolPointerValue(ds.Enabled.PointerBool())
		datastore.ID = types.StringValue(ds.ID)
		datastore.NodeName = types.StringValue(model.NodeName.ValueString())
		datastore.Shared = types.BoolPointerValue(ds.Shared.PointerBool())
		datastore.SpaceAvailable = types.Int64PointerValue(ds.SpaceAvailable.PointerInt64())
		datastore.SpaceTotal = types.Int64PointerValue(ds.SpaceTotal.PointerInt64())
		datastore.SpaceUsed = types.Int64PointerValue(ds.SpaceUsed.PointerInt64())
		datastore.SpaceUsedFraction = types.Float64PointerValue(ds.SpaceUsedPercentage.PointerFloat64())
		datastore.Type = types.StringValue(ds.Type)

		model.Datastores = append(model.Datastores, datastore)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
