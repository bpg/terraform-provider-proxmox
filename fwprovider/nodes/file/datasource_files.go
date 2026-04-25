/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package file

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
)

var (
	_ datasource.DataSource              = &listDatasource{}
	_ datasource.DataSourceWithConfigure = &listDatasource{}
)

// listModel is the data model for the files list data source.
type listModel struct {
	NodeName      types.String    `tfsdk:"node_name"`
	DatastoreID   types.String    `tfsdk:"datastore_id"`
	ContentType   types.String    `tfsdk:"content_type"`
	FileNameRegex types.String    `tfsdk:"file_name_regex"`
	Files         []listFileEntry `tfsdk:"files"`
}

// listFileEntry represents a single file in the list data source output.
type listFileEntry struct {
	ID          types.String `tfsdk:"id"`
	ContentType types.String `tfsdk:"content_type"`
	FileName    types.String `tfsdk:"file_name"`
	FileFormat  types.String `tfsdk:"file_format"`
	FileSize    types.Int64  `tfsdk:"file_size"`
	VMID        types.Int64  `tfsdk:"vmid"`
}

// listDatasource is the implementation of the files list data source.
type listDatasource struct {
	client proxmox.Client
}

// NewListDataSource creates a new files list data source.
func NewListDataSource() datasource.DataSource {
	return &listDatasource{}
}

// Metadata defines the name of the data source.
func (d *listDatasource) Metadata(
	_ context.Context,
	_ datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "proxmox_files"
}

// Schema defines the schema for the files list data source.
func (d *listDatasource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a list of files available in a datastore on a specific Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
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
				Description: "The content type to filter by. When set, only files of this type " +
					"are returned. Valid values are `backup`, `images`, `import`, `iso`, " +
					"`rootdir`, `snippets`, `vztmpl`.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(storage.ValidContentTypes()...),
				},
			},
			"file_name_regex": schema.StringAttribute{
				Description: "A regular expression to filter files by name. When set, only files " +
					"whose name matches the expression are returned.",
				Optional: true,
				Validators: []validator.String{
					validators.IsValidRegularExpression(),
				},
			},
			"files": schema.ListNestedAttribute{
				Description: "The list of files in the datastore.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the file (volume ID), " +
								"e.g. `local:iso/ubuntu.iso`.",
							Computed: true,
						},
						"content_type": schema.StringAttribute{
							Description: "The content type of the file.",
							Computed:    true,
						},
						"file_name": schema.StringAttribute{
							Description: "The name of the file.",
							Computed:    true,
						},
						"file_format": schema.StringAttribute{
							Description: "The format of the file.",
							Computed:    true,
						},
						"file_size": schema.Int64Attribute{
							Description: "The size of the file in bytes.",
							Computed:    true,
						},
						"vmid": schema.Int64Attribute{
							Description: "The VM ID associated with the file, if applicable.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure sets the client for the data source.
func (d *listDatasource) Configure(
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
func (d *listDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model listModel

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

	var nameRegex *regexp.Regexp

	if !model.FileNameRegex.IsNull() && !model.FileNameRegex.IsUnknown() {
		nameRegex, err = regexp.Compile(model.FileNameRegex.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Files",
				fmt.Sprintf("Invalid file_name_regex: %s", err.Error()),
			)

			return
		}
	}

	model.Files = make([]listFileEntry, 0, len(apiFiles))

	for _, apiFile := range apiFiles {
		if apiFile == nil {
			continue
		}

		// Extract filename from volume ID format: "datastore:content/filename"
		var fileName string

		if _, afterColon, found := strings.Cut(apiFile.VolumeID, ":"); found {
			if _, afterSlash, foundSlash := strings.Cut(afterColon, "/"); foundSlash {
				fileName = afterSlash
			}
		}

		if nameRegex != nil && !nameRegex.MatchString(fileName) {
			continue
		}

		file := listFileEntry{
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
