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

	"github.com/bpg/terraform-provider-proxmox/fwprovider/structure"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

var (
	_ resource.Resource              = &downloadFileResource{}
	_ resource.ResourceWithConfigure = &downloadFileResource{}
	// _      resource.ResourceWithModifyPlan = &downloadFileResource{}
	httpRe = regexp.MustCompile(`https?://.*`)
)

func RequiresReplace() planmodifier.Int64 {
	return RequiresReplaceModifier{}
}

type RequiresReplaceModifier struct{}

func (r RequiresReplaceModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	// Do not replace on resource creation.
	if req.State.Raw.IsNull() {
		return
	}

	// // Do not replace on resource destroy.
	if req.Plan.Raw.IsNull() {
		return
	}
	var plan, state downloadFileModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	resp.Diagnostics.AddWarning(
		"xxxxxxx", plan.Size.String(),
	)
	resp.Diagnostics.AddWarning(
		"accccccccc", state.Size.String(),
	)

	xxx, _ := req.Private.GetKey(ctx, "ProviderData")
	resp.Diagnostics.AddWarning(
		"123123", string(xxx),
	)

	// // Do not replace if the plan and state values are equal.
	// if req.PlanValue.Equal(req.StateValue) {
	// 	return
	// }

	resp.RequiresReplace = true
	resp.PlanValue = types.Int64Value(1231)
}

func (r RequiresReplaceModifier) Description(ctx context.Context) string {
	return "aaa"
}

func (r RequiresReplaceModifier) MarkdownDescription(ctx context.Context) string {
	return "ccc"
}

type downloadFileModel struct {
	ID                types.String `tfsdk:"id"`
	Content           types.String `tfsdk:"content_type"`
	FileName          types.String `tfsdk:"file_name"`
	Storage           types.String `tfsdk:"datastore_id"`
	Node              types.String `tfsdk:"node_name"`
	Size              types.Int64  `tfsdk:"size"`
	URL               types.String `tfsdk:"url"`
	Checksum          types.String `tfsdk:"checksum"`
	Compression       types.String `tfsdk:"compression"`
	UploadTimeout     types.Int64  `tfsdk:"upload_timeout"`
	ChecksumAlgorithm types.String `tfsdk:"checksum_algorithm"`
	Path              types.String `tfsdk:"path"`
	Verify            types.Bool   `tfsdk:"verify"`
	Overwrite         types.Bool   `tfsdk:"overwrite"`
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
		Description: "Manages files upload using PVE download-url API. ",
		MarkdownDescription: "Manages files upload using PVE download-url API. " +
			"It can be fully compatibile and faster replacement for image files created using " +
			"`proxmox_virtual_environment_file`. Supports `iso` and `vztmpl` content types. " +
			"For some file extenstions, like `.qcow2` of debian images, you must manually " +
			"enter `file_name` with one of PVE supported extensions like `.img`.",
		Attributes: map[string]schema.Attribute{
			"id": structure.IDAttribute(),
			"content_type": schema.StringAttribute{
				Description: "The file content type. Must be `iso` | `vztmpl`.",
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
				Description: "The file name. If not provided, it is calculated" +
					"using `url`.",
				Computed: true,
				Required: false,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"path": schema.StringAttribute{
				Description: "The file path on host.",
				Computed:    true,
				Required:    false,
				Optional:    false,
				PlanModifiers: []planmodifier.String{
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
				Default:     nil,
				// PlanModifiers: []planmodifier.Int64{
				// 	// int64planmodifier.UseStateForUnknown(),
				// 	// RequiresReplace(),
				// 	// int64planmodifier.RequiresReplace(),
				// },
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
					stringvalidator.RegexMatches(httpRe, "Must match http url regex"),
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
			"compression": schema.StringAttribute{
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
		if strings.Contains(err.Error(), "refusing to override existing file") {
			resp.Diagnostics.AddError(
				"File already exists in a datastore, it was created outside of Terraform "+
					"or is managed by another resource.",
				fmt.Sprintf(
					"File already exists in a datastore: `%s`, "+
						"error: %s",
					plan.FileName.ValueString(),
					err.Error(),
				),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error creating Download File interface",
				fmt.Sprintf(
					"Could not DownloadFileByURL: `%s`, "+
						"unexpected error: %s",
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
			"Error fetching metadata from download url, "+
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

	fileData, err := storageClient.GetDatastoreFile(
		ctx,
		model.ID.ValueString(),
		model.Node.ValueString(),
	)
	if err != nil {
		return fmt.Errorf("file does not exists in datastore: %w", err)
	}

	pathSplit := strings.Split(*fileData.Path, "/")
	filename := pathSplit[len(pathSplit)-1]

	model.FileName = types.StringValue(filename)
	model.Size = types.Int64Value(*fileData.FileSize)
	model.Path = types.StringValue(*fileData.Path)

	return nil
}

func (r *downloadFileResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var state, plan downloadFileModel
	diags := req.State.Get(ctx, &state)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
	}

	diags = req.Plan.Get(ctx, &plan)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	plan.Size = types.Int64Value(21300)
	plan.Path = types.StringValue("asddsadsa")
	// // req.State.Set(ctx, state)
	// // force resource recreation
	resp.Diagnostics.AddWarning(
		"asdddddddddddddddddddddd",
		fmt.Sprintf("%d %d",
			plan.Size.ValueInt64(),
			state.Size.ValueInt64(),
		),
	)

	resp.Diagnostics.Append(req.Plan.Set(ctx, plan)...)
	// resp.Diagnostics.Append(req.Plan.SetAttribute(ctx, path.Root("size"), 2137)...)
	// resp.Diagnostics.Append(req.State.SetAttribute(ctx, path.Root("size"), 2137)...)
	// resp.RequiresReplace = resp.RequiresReplace.Append(path.Root("size"))
	// size := path.Root("size")
	resp.RequiresReplace.Append(path.Root("size"))
	// resp.RequiresReplace(path.Paths{path.Root("field?")})
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

	err := r.read(ctx, &state)
	if err != nil {
		if strings.Contains(err.Error(), "failed to authenticate") {
			resp.Diagnostics.AddError("Failed to authenticate", err.Error())

			return
		}
		resp.Diagnostics.AddWarning(
			"The file does not exist in datastore and must be replaced.",
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
			if *urlMetadata.Size != state.Size.ValueInt64() {
				resp.Diagnostics.AddWarning(
					"File size in datastore does not match size from url.",
					fmt.Sprintf("File current size %d does not match target size from url: %d",
						state.Size.ValueInt64(),
						*urlMetadata.Size),
				)
			}
			state.Size = types.Int64Value(*urlMetadata.Size)
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
	if err != nil {
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
