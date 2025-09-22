/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
)

var (
	_ datasource.DataSource              = &fileDataSource{}
	_ datasource.DataSourceWithConfigure = &fileDataSource{}
)

type fileDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	NodeName    types.String `tfsdk:"node_name"`
	DatastoreID types.String `tfsdk:"datastore_id"`
	ContentType types.String `tfsdk:"content_type"`
	FileName    types.String `tfsdk:"file_name"`
	FileSize    types.Int64  `tfsdk:"file_size"`
	FileFormat  types.String `tfsdk:"file_format"`
	VMID        types.Int64  `tfsdk:"vmid"`
}

type fileDataSource struct {
	client proxmox.Client
}

func NewFileDataSource() datasource.DataSource {
	return &fileDataSource{}
}

func (d *fileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (d *fileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about an existing file in a Proxmox Virtual Environment node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the file (volume ID).",
				Computed:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the node.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"datastore_id": schema.StringAttribute{
				Description: "The identifier of the datastore.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"content_type": schema.StringAttribute{
				Description: "The content type of the file.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"backup",
						"iso",
						"vztmpl",
						"rootdir",
						"images",
						"snippets",
					),
				},
			},
			"file_name": schema.StringAttribute{
				Description: "The name of the file.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"file_size": schema.Int64Attribute{
				Description: "The size of the file in bytes.",
				Computed:    true,
			},
			"file_format": schema.StringAttribute{
				Description: "The format of the file.",
				Computed:    true,
			},
			"vmid": schema.Int64Attribute{
				Description: "The VM ID associated with the file (if applicable).",
				Computed:    true,
			},
		},
	}
}

func (d *fileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *fileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data fileDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	storageClient := d.client.Node(data.NodeName.ValueString()).Storage(data.DatastoreID.ValueString())

	files, err := storageClient.ListDatastoreFiles(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Files",
			fmt.Sprintf("Unable to list files from datastore %s on node %s: %s",
				data.DatastoreID.ValueString(), data.NodeName.ValueString(), err.Error()),
		)

		return
	}

	var foundFile *storage.DatastoreFileListResponseData
	targetContentType := data.ContentType.ValueString()
	targetFileName := data.FileName.ValueString()

	for _, file := range files {
		if file == nil {
			continue
		}

		if file.ContentType != targetContentType {
			continue
		}

		// Extract filename from volume ID format: datastore:content/filename
		volumeParts := strings.SplitN(file.VolumeID, ":", 2)
		if len(volumeParts) < 2 {
			continue
		}

		contentAndFile := volumeParts[1]

		fileParts := strings.SplitN(contentAndFile, "/", 2)
		if len(fileParts) < 2 {
			continue
		}

		fileName := fileParts[1]
		if fileName == targetFileName {
			foundFile = file
			break
		}
	}

	if foundFile == nil {
		resp.Diagnostics.AddError(
			"File Not Found",
			fmt.Sprintf("File '%s' with content type '%s' not found in datastore '%s' on node '%s'",
				targetFileName, targetContentType, data.DatastoreID.ValueString(), data.NodeName.ValueString()),
		)

		return
	}

	data.ID = types.StringValue(foundFile.VolumeID)
	data.FileSize = types.Int64Value(foundFile.FileSize)
	data.FileFormat = types.StringValue(foundFile.FileFormat)

	if foundFile.VMID != nil {
		data.VMID = types.Int64Value(int64(*foundFile.VMID))
	} else {
		data.VMID = types.Int64Null()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
