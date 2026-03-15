/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package files

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
)

var (
	_ datasource.DataSource              = &Datasource{}
	_ datasource.DataSourceWithConfigure = &Datasource{}
)

// Datasource is the implementation of the files data source.
type Datasource struct {
	client proxmox.Client
}

// NewDataSource creates a new files data source.
func NewDataSource() datasource.DataSource {
	return &Datasource{}
}

// Metadata defines the name of the data source.
func (d *Datasource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_files"
}

// Configure sets the client for the data source.
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

// Read fetches the list of files from the Proxmox API.
func (d *Datasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model Model

	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}

	storageClient := d.client.Node(model.NodeName.ValueString()).Storage(model.DatastoreID.ValueString())

	var contentType *string

	if !model.ContentType.IsNull() && !model.ContentType.IsUnknown() {
		ct := model.ContentType.ValueString()
		contentType = &ct
	}

	apiFiles, err := storageClient.ListDatastoreFiles(ctx, contentType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Files",
			fmt.Sprintf("Unable to list files from datastore '%s' on node '%s': %s",
				model.DatastoreID.ValueString(), model.NodeName.ValueString(), err.Error()),
		)

		return
	}

	model.Files = make([]File, 0, len(apiFiles))

	for _, apiFile := range apiFiles {
		if apiFile == nil {
			continue
		}

		// Extract filename from volume ID format: "datastore:content/filename"
		var fileName string

		volumeParts := strings.SplitN(apiFile.VolumeID, ":", 2)
		if len(volumeParts) == 2 {
			fileParts := strings.SplitN(volumeParts[1], "/", 2)
			if len(fileParts) == 2 {
				fileName = fileParts[1]
			}
		}

		file := File{
			ID:          types.StringValue(apiFile.VolumeID),
			ContentType: types.StringValue(apiFile.ContentType),
			FileName:    types.StringValue(fileName),
			FileFormat:  types.StringValue(apiFile.FileFormat),
			FileSize:    types.Int64Value(apiFile.FileSize),
		}

		if apiFile.VMID != nil {
			file.VMID = types.Int64Value(int64(*apiFile.VMID))
		} else {
			file.VMID = types.Int64Null()
		}

		model.Files = append(model.Files, file)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
