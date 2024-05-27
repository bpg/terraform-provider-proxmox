/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fwprovider

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

var (
	_         resource.Resource              = &downloadFileResource{}
	_         resource.ResourceWithConfigure = &downloadFileResource{}
	httpRegex                                = regexp.MustCompile(`https?://.*`)
)

func sizeRequiresReplace() planmodifier.Int64 {
	return sizeRequiresReplaceModifier{}
}

type sizeRequiresReplaceModifier struct{}

func (r sizeRequiresReplaceModifier) PlanModifyInt64(
	ctx context.Context,
	req planmodifier.Int64Request,
	resp *planmodifier.Int64Response,
) {
	// Do not replace on resource creation.
	if req.State.Raw.IsNull() {
		return
	}

	// Do not replace on resource destroy.
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state downloadFileModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	originalStateSizeBytes, diags := req.Private.GetKey(ctx, "original_state_size")

	resp.Diagnostics.Append(diags...)

	if originalStateSizeBytes != nil {
		originalStateSize, err := strconv.ParseInt(string(originalStateSizeBytes), 10, 64)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unexpected error when reading originalStateSize from Private",
				fmt.Sprintf(
					"Unexpected error in ParseInt: %s",
					err.Error(),
				),
			)

			return
		}

		if state.Size.ValueInt64() != originalStateSize {
			resp.RequiresReplace = true
			resp.PlanValue = types.Int64Value(originalStateSize)

			resp.Diagnostics.AddWarning(
				"The file size in datastore has changed.",
				fmt.Sprintf(
					"Previous size %d does not match size from datastore: %d",
					originalStateSize,
					state.Size.ValueInt64(),
				),
			)

			return
		}
	}

	urlSizeBytes, diags := req.Private.GetKey(ctx, "url_size")

	resp.Diagnostics.Append(diags...)

	if (urlSizeBytes != nil) && (plan.URL.ValueString() == state.URL.ValueString()) {
		urlSize, err := strconv.ParseInt(string(urlSizeBytes), 10, 64)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unexpected error when reading urlSize from Private",
				fmt.Sprintf(
					"Unexpected error in ParseInt: %s",
					err.Error(),
				),
			)

			return
		}

		if state.Size.ValueInt64() != urlSize {
			resp.RequiresReplace = true
			resp.PlanValue = types.Int64Value(urlSize)

			resp.Diagnostics.AddWarning(
				"The file size from url has changed.",
				fmt.Sprintf(
					"Size from url %d does not match size from datastore: %d",
					urlSize,
					state.Size.ValueInt64(),
				),
			)

			return
		}
	}
}

func (r sizeRequiresReplaceModifier) Description(_ context.Context) string {
	return "Triggers resource force replacement if `size` in state does not match remote value."
}

func (r sizeRequiresReplaceModifier) MarkdownDescription(_ context.Context) string {
	return "Triggers resource force replacement if `size` in state does not match remote value."
}

type downloadFileModel struct {
	ID                     types.String `tfsdk:"id"`
	Content                types.String `tfsdk:"content_type"`
	FileName               types.String `tfsdk:"file_name"`
	Storage                types.String `tfsdk:"datastore_id"`
	Node                   types.String `tfsdk:"node_name"`
	Size                   types.Int64  `tfsdk:"size"`
	URL                    types.String `tfsdk:"url"`
	Checksum               types.String `tfsdk:"checksum"`
	DecompressionAlgorithm types.String `tfsdk:"decompression_algorithm"`
	UploadTimeout          types.Int64  `tfsdk:"upload_timeout"`
	ChecksumAlgorithm      types.String `tfsdk:"checksum_algorithm"`
	Verify                 types.Bool   `tfsdk:"verify"`
	Overwrite              types.Bool   `tfsdk:"overwrite"`
	OverwriteUnmanaged     types.Bool   `tfsdk:"overwrite_unmanaged"`
}

// NewDownloadFileResource manages files downloaded using Proxmox API.
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
		Description: "Manages files upload using PVE download-url API. ",
		MarkdownDescription: "Manages files upload using PVE download-url API. " +
			"It can be fully compatible and faster replacement for image files created using " +
			"`proxmox_virtual_environment_file`. Supports images for VMs (ISO images) and LXC (CT Templates).",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ID(),
			"content_type": schema.StringAttribute{
				Description: "The file content type. Must be `iso` for VM images or `vztmpl` for LXC images.",
				Required:    true,
				Validators: []validator.String{stringvalidator.OneOf([]string{
					"iso",
					"vztmpl",
				}...)},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"file_name": schema.StringAttribute{
				Description: "The file name. If not provided, it is calculated " +
					"using `url`. PVE will raise 'wrong file extension' error for some popular " +
					"extensions file `.raw` or `.qcow2`. Workaround is to use e.g. `.img` instead.",
				Computed: true,
				Required: false,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"datastore_id": schema.StringAttribute{
				Description: "The identifier for the target datastore.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"node_name": schema.StringAttribute{
				Description: "The node name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"size": schema.Int64Attribute{
				Description: "The file size.",
				Optional:    false,
				Required:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
					sizeRequiresReplace(),
				},
			},
			"upload_timeout": schema.Int64Attribute{
				Description: "The file download timeout seconds. Default is 600 (10min).",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(600),
			},
			"url": schema.StringAttribute{
				Description: "The URL to download the file from. Format `https?://.*`.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(httpRegex, "Must match http url regex"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"checksum": schema.StringAttribute{
				Description: "The expected checksum of the file.",
				Optional:    true,
				Default:     nil,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("checksum_algorithm")),
				},
			},
			"decompression_algorithm": schema.StringAttribute{
				Description: "Decompress the downloaded file using the " +
					"specified compression algorithm. Must be one of `gz` | `lzo` | `zst`.",
				Optional: true,
				Default:  nil,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"gz",
						"lzo",
						"zst",
					}...),
				},
			},
			"checksum_algorithm": schema.StringAttribute{
				Description: "The algorithm to calculate the checksum of the file. " +
					"Must be `md5` | `sha1` | `sha224` | `sha256` | `sha384` | `sha512`.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"md5",
						"sha1",
						"sha224",
						"sha256",
						"sha384",
						"sha512",
					}...),
					stringvalidator.AlsoRequires(path.MatchRoot("checksum")),
				},
				Default: nil,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"verify": schema.BoolAttribute{
				Description: "By default `true`. If `false`, no SSL/TLS certificates will be verified.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"overwrite": schema.BoolAttribute{
				Description: "If `true` and size of uploaded file is different, " +
					"than size from `url` Content-Length header, file will be downloaded again. " +
					"If `false`, there will be no checks.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"overwrite_unmanaged": schema.BoolAttribute{
				Description: "If `true` and a file with the same name already exists in the datastore, " +
					"it will be deleted and the new file will be downloaded. If `false` and the file already exists, " +
					"an error will be returned.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
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

	timeout := time.Duration(plan.UploadTimeout.ValueInt64()) * time.Second

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	fileMetadata, err := r.getURLMetadata(
		ctx,
		&plan,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error initiating file download",
			"Could not get file metadata, unexpected error: "+err.Error(),
		)

		return
	}

	if plan.FileName.IsUnknown() {
		plan.FileName = types.StringValue(*fileMetadata.Filename)
	}

	nodesClient := r.client.Node(plan.Node.ValueString())
	verify := proxmoxtypes.CustomBool(plan.Verify.ValueBool())

	downloadFileReq := storage.DownloadURLPostRequestBody{
		Node:              plan.Node.ValueStringPointer(),
		Storage:           plan.Storage.ValueStringPointer(),
		Content:           plan.Content.ValueStringPointer(),
		Checksum:          plan.Checksum.ValueStringPointer(),
		ChecksumAlgorithm: plan.ChecksumAlgorithm.ValueStringPointer(),
		Compression:       plan.DecompressionAlgorithm.ValueStringPointer(),
		FileName:          plan.FileName.ValueStringPointer(),
		URL:               plan.URL.ValueStringPointer(),
		Verify:            &verify,
	}

	storageClient := nodesClient.Storage(plan.Storage.ValueString())

	err = storageClient.DownloadFileByURL(ctx, &downloadFileReq)

	if isErrFileAlreadyExists(err) && plan.OverwriteUnmanaged.ValueBool() {
		fileID := plan.Content.ValueString() + "/" + plan.FileName.ValueString()

		err = storageClient.DeleteDatastoreFile(ctx, fileID)
		if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError("Error deleting file from datastore",
				fmt.Sprintf("Could not delete file '%s', unexpected error: %s", fileID, err.Error()),
			)
		}

		err = storageClient.DownloadFileByURL(ctx, &downloadFileReq)
	}

	if err != nil {
		if isErrFileAlreadyExists(err) {
			resp.Diagnostics.AddError(
				"File already exists in the datastore, it was created outside of Terraform "+
					"or is managed by another resource.",
				fmt.Sprintf("File already exists in the datastore: '%s', error: %s",
					plan.FileName.ValueString(), err.Error(),
				),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error downloading file from url",
				fmt.Sprintf("Could not download file '%s', unexpected error: %s",
					plan.FileName.ValueString(),
					err.Error(),
				),
			)
		}

		return
	}

	plan.ID = types.StringValue(plan.Storage.ValueString() + ":" +
		plan.Content.ValueString() + "/" + plan.FileName.ValueString())

	err = r.read(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading file from datastore", err.Error(),
		)
	}

	resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *downloadFileResource) getURLMetadata(
	ctx context.Context,
	model *downloadFileModel,
) (*nodes.QueryURLMetadataGetResponseData, error) {
	nodesClient := r.client.Node(model.Node.ValueString())
	verify := proxmoxtypes.CustomBool(model.Verify.ValueBool())

	queryURLMetadataReq := nodes.QueryURLMetadataGetRequestBody{
		URL:    model.URL.ValueStringPointer(),
		Verify: &verify,
	}

	fileMetadata, err := nodesClient.GetQueryURLMetadata(
		ctx,
		&queryURLMetadataReq,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error fetching metadata from download url, "+
				"unexpected error in GetQueryURLMetadata: %w",
			err,
		)
	}

	return fileMetadata, nil
}

func (r *downloadFileResource) read(
	ctx context.Context,
	model *downloadFileModel,
) error {
	nodesClient := r.client.Node(model.Node.ValueString())
	storageClient := nodesClient.Storage(model.Storage.ValueString())

	datastoresFiles, err := storageClient.ListDatastoreFiles(ctx)
	if err != nil {
		return fmt.Errorf("unexpected error when listing datastore files: %w", err)
	}

	for _, file := range datastoresFiles {
		if file != nil {
			if file.VolumeID != model.ID.ValueString() {
				continue
			}

			model.Size = types.Int64Value(file.FileSize)

			return nil
		}
	}

	return fmt.Errorf("file does not exists in datastore")
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

	setOriginalValue := []byte(strconv.FormatInt(state.Size.ValueInt64(), 10))
	resp.Private.SetKey(ctx, "original_state_size", setOriginalValue)

	err := r.read(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "failed to authenticate") {
			resp.Diagnostics.AddError("Failed to authenticate", err.Error())

			return
		}

		resp.Diagnostics.AddWarning(
			"The file does not exist in datastore and resource must be recreated.",
			err.Error(),
		)
		resp.State.RemoveResource(ctx)

		return
	}

	if state.Overwrite.ValueBool() {
		// with overwrite, use url to get proper target size
		urlMetadata, err := r.getURLMetadata(
			ctx,
			&state,
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not get file metadata from url.",
				err.Error(),
			)

			return
		}

		if urlMetadata.Size != nil {
			setValue := []byte(strconv.FormatInt(*urlMetadata.Size, 10))
			resp.Private.SetKey(ctx, "url_size", setValue)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update file resource.
func (r *downloadFileResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state downloadFileModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	err := r.read(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading file from datastore", err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete removes file resource.
func (r *downloadFileResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state downloadFileModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	nodesClient := r.client.Node(state.Node.ValueString())
	storageClient := nodesClient.Storage(state.Storage.ValueString())

	err := storageClient.DeleteDatastoreFile(
		ctx,
		state.ID.ValueString(),
	)
	if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
		if strings.Contains(err.Error(), "unable to parse") {
			resp.Diagnostics.AddWarning(
				"Datastore file does not exists",
				fmt.Sprintf(
					"Could not delete datastore file '%s', it does not exist or has been deleted outside of Terraform.",
					state.ID.ValueString(),
				),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error deleting datastore file",
				fmt.Sprintf("Could not delete datastore file '%s', unexpected error: %s",
					state.ID.ValueString(), err.Error()),
			)
		}
	}
}

func isErrFileAlreadyExists(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), "refusing to override existing file")
}
