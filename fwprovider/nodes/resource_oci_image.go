/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
)

var (
	_              resource.Resource              = &ociImageResource{}
	_              resource.ResourceWithConfigure = &ociImageResource{}
	referenceRegex                                = regexp.MustCompile(
		`^(?:(?:[a-zA-Z\d]|[a-zA-Z\d][a-zA-Z\d-]*[a-zA-Z\d])(?:\.(?:[a-zA-Z\d]|[a-zA-Z\d][a-zA-Z\d-]*[a-zA-Z\d]))*` +
			`(?::\d+)?/)?[a-z\d]+(?:(?:[._]|__|[-]*)[a-z\d]+)*(?:/[a-z\d]+(?:(?:[._]|__|[-]*)[a-z\d]+)*)*:\w[\w.-]{0,127}$`,
	)
)

const ociSizeRequiresReplaceDescription = "Triggers resource force replacement if `size` in state does not match remote value."

type ociSizeRequiresReplaceModifier struct{}

func (r ociSizeRequiresReplaceModifier) PlanModifyInt64(
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

	var plan, state ociImageModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	originalStateSizeBytes, diags := req.Private.GetKey(ctx, "original_state_size")

	resp.Diagnostics.Append(diags...)

	if originalStateSizeBytes != nil {
		originalStateSize, err := strconv.ParseInt(string(originalStateSizeBytes), 10, 64)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to convert original state OCI image size to int64",
				"Unexpected error in parsing string to int64, key original_state_size. "+
					"Please retry the operation or report this issue to the provider developers.\n\n"+
					"Error: "+err.Error(),
			)

			return
		}

		if state.Size.ValueInt64() != originalStateSize && plan.Overwrite.ValueBool() {
			resp.RequiresReplace = true
			resp.PlanValue = types.Int64Value(originalStateSize)

			resp.Diagnostics.AddWarning(
				"The OCI image size in datastore has changed outside of terraform.",
				fmt.Sprintf(
					"Previous size: %d saved in state does not match current size from datastore: %d. "+
						"You can disable this behaviour by using overwrite=false",
					originalStateSize,
					state.Size.ValueInt64(),
				),
			)

			return
		}
	}
}

func (r ociSizeRequiresReplaceModifier) Description(_ context.Context) string {
	return ociSizeRequiresReplaceDescription
}

func (r ociSizeRequiresReplaceModifier) MarkdownDescription(_ context.Context) string {
	return ociSizeRequiresReplaceDescription
}

type ociImageModel struct {
	ID                 types.String `tfsdk:"id"`
	FileName           types.String `tfsdk:"file_name"`
	Storage            types.String `tfsdk:"datastore_id"`
	Node               types.String `tfsdk:"node_name"`
	Size               types.Int64  `tfsdk:"size"`
	Reference          types.String `tfsdk:"reference"`
	UploadTimeout      types.Int64  `tfsdk:"upload_timeout"`
	Overwrite          types.Bool   `tfsdk:"overwrite"`
	OverwriteUnmanaged types.Bool   `tfsdk:"overwrite_unmanaged"`
}

// NewOCIImageResource manages OCI images downloaded using Proxmox API.
func NewOCIImageResource() resource.Resource {
	return &ociImageResource{}
}

type ociImageResource struct {
	client proxmox.Client
}

func (r *ociImageResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_oci_image"
}

// Schema defines the schema for the resource.
func (r *ociImageResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages OCI images pulled from OCI registries using PVE oci-registry-pull API. ",
		MarkdownDescription: "Manages OCI images pulled from OCI registries using PVE oci-registry-pull API. " +
			"Pulls OCI container images and stores them as tar files in Proxmox VE datastores.",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ResourceID(),
			"file_name": schema.StringAttribute{
				Description: "The file name for the pulled OCI image. If not provided, " +
					"it will be generated automatically. The file will be stored as a .tar file.",
				Computed: true,
				Required: false,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`\.tar$`),
						"file name must end with .tar",
					),
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
				Description: "The image size in PVE.",
				Optional:    false,
				Required:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					ociSizeRequiresReplaceModifier{},
				},
			},
			"upload_timeout": schema.Int64Attribute{
				Description: "The OCI image pull timeout in seconds. Default is 600 (10min).",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(600),
			},
			"reference": schema.StringAttribute{
				Description: "The reference to the OCI image.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(referenceRegex, "must match OCI image reference regex `"+referenceRegex.String()+"`"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"overwrite": schema.BoolAttribute{
				Description: "By default `true`. If `true` and the OCI image size has changed in the datastore, " +
					"it will be replaced. If `false`, there will be no check.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"overwrite_unmanaged": schema.BoolAttribute{
				Description: "If `true` and an OCI image with the same name already exists in the datastore, " +
					"it will be deleted and the new image will be pulled. If `false` and the image already exists, " +
					"an error will be returned.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
		},
	}
}

func (r *ociImageResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	cfg, ok := req.ProviderData.(config.Resource)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *proxmox.Client, got: %T", req.ProviderData),
		)

		return
	}

	r.client = cfg.Client
}

func (r *ociImageResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan ociImageModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	timeout := time.Duration(plan.UploadTimeout.ValueInt64()) * time.Second

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// If filename is not provided, generate one from reference
	if plan.FileName.IsUnknown() || plan.FileName.IsNull() {
		// Generate filename from reference (e.g., "docker.io/library/ubuntu:latest" -> "ubuntu_latest.tar")
		refParts := strings.Split(plan.Reference.ValueString(), "/")
		lastPart := refParts[len(refParts)-1]
		filename := strings.ReplaceAll(lastPart, ":", "_")

		plan.FileName = types.StringValue(filename + ".tar")
	}

	nodesClient := r.client.Node(plan.Node.ValueString())

	// Proxmox API expects filename without .tar extension
	filenameWithoutTar := strings.TrimSuffix(plan.FileName.ValueString(), ".tar")

	ociPullReq := storage.OCIRegistryPullRequestBody{
		Storage:   plan.Storage.ValueStringPointer(),
		FileName:  &filenameWithoutTar,
		Reference: plan.Reference.ValueStringPointer(),
	}

	storageClient := nodesClient.Storage(plan.Storage.ValueString())

	err := storageClient.DownloadOCIImageByReference(ctx, &ociPullReq)
	if isErrFileAlreadyExists(err) && plan.OverwriteUnmanaged.ValueBool() {
		// OCI images are stored as vztmpl content type in Proxmox
		fileID := storage.ContentTypeVZTmpl + "/" + plan.FileName.ValueString()

		err = storageClient.DeleteDatastoreFile(ctx, fileID)
		if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError("Error deleting OCI image from datastore",
				fmt.Sprintf("Could not delete OCI image '%s', unexpected error: %s", fileID, err.Error()),
			)
		}

		err = storageClient.DownloadOCIImageByReference(ctx, &ociPullReq)
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
				"Error pulling OCI image from reference",
				fmt.Sprintf("Could not pull OCI image '%s', unexpected error: %s",
					plan.FileName.ValueString(), err.Error()),
			)
		}

		return
	}

	// OCI images are stored as vztmpl content type in Proxmox
	plan.ID = types.StringValue(plan.Storage.ValueString() + ":" + storage.ContentTypeVZTmpl + "/" + plan.FileName.ValueString())

	err = r.read(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading file from datastore", err.Error(),
		)
	}

	resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ociImageResource) read(
	ctx context.Context,
	model *ociImageModel,
) error {
	nodesClient := r.client.Node(model.Node.ValueString())
	storageClient := nodesClient.Storage(model.Storage.ValueString())

	// OCI images are stored as `vztmpl` content type in Proxmox
	contentType := storage.ContentTypeVZTmpl

	datastoresFiles, err := storageClient.ListDatastoreFiles(ctx, &contentType)
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

	return fmt.Errorf("OCI image does not exist in datastore")
}

// Read reads file from datastore.
func (r *ociImageResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state ociImageModel

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
			"The OCI image does not exist in datastore and resource must be recreated.",
			err.Error(),
		)
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update file resource.
func (r *ociImageResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state ociImageModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	err := r.read(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error when reading OCI Image from datastore", err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

//nolint:dupl // delete path mirrors download file resource implementation
func (r *ociImageResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state ociImageModel

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
				"Datastore OCI image does not exists",
				fmt.Sprintf(
					"Could not delete datastore OCI image '%s', it does not exist or has been deleted outside of Terraform.",
					state.ID.ValueString(),
				),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error deleting datastore OCI image",
				fmt.Sprintf("Could not delete datastore OCI image '%s', unexpected error: %s",
					state.ID.ValueString(), err.Error()),
			)
		}
	}
}
