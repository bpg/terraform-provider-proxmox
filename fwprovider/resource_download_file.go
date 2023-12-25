/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fwprovider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/structure"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/nodestorage"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

var (
	_      resource.Resource              = &downloadFileResource{}
	_      resource.ResourceWithConfigure = &downloadFileResource{}
	httpRe                                = regexp.MustCompile(`https?://.*`)
)

type downloadFileModel struct {
	ID                types.String `tfsdk:"id"`
	Content           types.String `tfsdk:"content"`
	FileName          types.String `tfsdk:"filename"`
	Storage           types.String `tfsdk:"datastore_id"`
	Node              types.String `tfsdk:"node_name"`
	Size              types.Int64  `tfsdk:"size"`
	URL               types.String `tfsdk:"download_url"`
	Checksum          types.String `tfsdk:"checksum"`
	Compression       types.String `tfsdk:"compression"`
	UploadTimeout     types.Int64  `tfsdk:"upload_timeout"`
	ChecksumAlgorithm types.String `tfsdk:"checksum_algorithm"`
	Verify            types.Bool   `tfsdk:"verify"`
}

// NewDownloadFileResource manages files downloaded using proxmomx API.
func NewDownloadFileResource() resource.Resource {
	return &downloadFileResource{}
}

type downloadFileResource struct {
	client proxmox.Client
}

func (r *downloadFileResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_download_file"
}

// Schema defines the schema for the resource.
func (r *downloadFileResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages files downloaded directly using proxmomx API. " +
			"Supports only `iso` and `vztmpl` content types.",
		Attributes: map[string]schema.Attribute{
			"id": structure.IDAttribute(),
			"content": schema.StringAttribute{
				MarkdownDescription: "The file content type. Must be `iso` | `vztmpl`.",
				Required:            true,
				Validators: []validator.String{stringvalidator.OneOf([]string{
					"iso",
					"vztmpl",
				}...)},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"filename": schema.StringAttribute{
				MarkdownDescription: "The file name.",
				Computed:            true,
				Required:            false,
				Optional:            false,
			},
			"datastore_id": schema.StringAttribute{
				MarkdownDescription: "The identifier for the target datastore.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"node_name": schema.StringAttribute{
				MarkdownDescription: "The node name.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"size": schema.Int64Attribute{
				MarkdownDescription: "The file size.",
				Optional:            false,
				Required:            false,
				Computed:            true,
			},
			"upload_timeout": schema.Int64Attribute{
				MarkdownDescription: "The file download timeout seconds. Default is 600 (10min).",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(600),
			},
			"download_url": schema.StringAttribute{
				MarkdownDescription: "The URL to download the file from. Format `https?://.*`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(httpRe, "Must match http url regex"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"checksum": schema.StringAttribute{
				MarkdownDescription: "The expected checksum of the file.",
				Optional:            true,
				Computed:            true,
				Default:             nil,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"compression": schema.StringAttribute{
				MarkdownDescription: "Decompress the downloaded file using the " +
					"specified compression algorithm.",
				Optional: true,
				Computed: true,
				Default:  nil,
			},
			"checksum_algorithm": schema.StringAttribute{
				MarkdownDescription: "The algorithm to calculate the checksum of the file. " +
					"Must be `md5` | `sha1` | `sha224` | `sha256` | `sha384` | `sha512`.",
				Optional: true,
				Validators: []validator.String{stringvalidator.OneOf([]string{
					"md5",
					"sha1",
					"sha224",
					"sha256",
					"sha384",
					"sha512",
				}...)},
				Computed: true,
				Default:  nil,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"verify": schema.BoolAttribute{
				MarkdownDescription: "By default `true`. If `false`, no SSL/TLS certificates will be verified.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
		},
	}
}

func (r *downloadFileResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(proxmox.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *proxmox.Client, got: %T", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *downloadFileResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan downloadFileModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	nodesClient := r.client.Node(plan.Node.ValueString())
	verify := proxmoxtypes.CustomBool(*plan.Verify.ValueBoolPointer())

	queryURLMetadataReq := nodes.QueryURLMetadataGetRequestBody{
		URL:    plan.URL.ValueStringPointer(),
		Verify: &verify,
	}

	fileMetadata, err := nodesClient.GetQueryURLMetadata(
		ctx,
		&queryURLMetadataReq,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Download File interface",
			"Could not get GetQueryURLMetadata, unexpected error: "+err.Error(),
		)

		return
	}

	plan.FileName = types.StringValue(*fileMetadata.Filename)
	plan.Size = types.Int64Value(*fileMetadata.Size)

	downloadFileReq := nodestorage.DownloadURLPostRequestBody{
		Node:              plan.Node.ValueStringPointer(),
		Storage:           plan.Storage.ValueStringPointer(),
		Content:           plan.Content.ValueStringPointer(),
		Checksum:          plan.Checksum.ValueStringPointer(),
		ChecksumAlgorithm: plan.ChecksumAlgorithm.ValueStringPointer(),
		Compression:       plan.Compression.ValueStringPointer(),
		FileName:          plan.FileName.ValueStringPointer(),
		URL:               plan.URL.ValueStringPointer(),
		Verify:            &verify,
	}

	storageClient := nodesClient.Storage(plan.Storage.ValueString())
	err = storageClient.DownloadFileByURL(
		ctx,
		&downloadFileReq,
		plan.UploadTimeout.ValueInt64(),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Download File interface",
			"Could not  GetQueryURLMetadata, unexpected error: "+err.Error(),
		)

		return
	}

	plan.ID = types.StringValue(plan.Storage.ValueString() + ":" +
		plan.Content.ValueString() + "/" + plan.FileName.ValueString())
}

// Read reads file from datastore.
func (r *downloadFileResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state downloadFileModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	nodesClient := r.client.Node(state.Node.ValueString())
	storageClient := nodesClient.Storage(state.Storage.ValueString())

	fileData, err := storageClient.GetDatastoreFile(
		ctx,
		state.ID.ValueString(),
		state.Node.ValueString(),
	)
	if err != nil {
		diags.AddError("Could not get file from datastore", err.Error())
	}

	if resp.Diagnostics.HasError() {
		resp.State.RemoveResource(ctx)
	}

	state.Content = types.StringValue(*fileData.FileFormat)
	state.Size = types.Int64Value(*fileData.FileSize)
	state.ID = types.StringValue(*fileData.Path)
}

// Update force-reacreate resource.
func (r *downloadFileResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var data, state downloadFileModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		resp.State.RemoveResource(ctx)
	}
}

// Delete removes file resource.
func (r *downloadFileResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data downloadFileModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	nodesClient := r.client.Node(data.Node.ValueString())
	storageClient := nodesClient.Storage(data.Storage.ValueString())

	err := storageClient.DeleteDatastoreFile(
		ctx,
		data.ID.ValueString(),
	)
	if err != nil {
		if strings.Contains(err.Error(), "no such") {
			resp.Diagnostics.AddWarning(
				"Datastore file does not exists",
				fmt.Sprintf(
					"Could not delete Datastore file '%s', it does not exist or has been deleted outside of Terraform.",
					data.ID.ValueString(),
				),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error deleting Datastore file",
				fmt.Sprintf("Could not delete Datastore file '%s', unexpected error: %s",
					data.ID.ValueString(), err.Error()),
			)
		}
	}
}
