/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package apt

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types/nodes/apt"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	api "github.com/bpg/terraform-provider-proxmox/proxmox/nodes/apt/repositories"
)

const (
	// ResourceStandardRepoIDPrefix is the prefix for the resource ID of standardRepositoryResource.
	ResourceStandardRepoIDPrefix = "apt_standard_repository"
)

// Ensure the resource implements the required interfaces.
var (
	_ resource.Resource                = &standardRepositoryResource{}
	_ resource.ResourceWithConfigure   = &standardRepositoryResource{}
	_ resource.ResourceWithImportState = &standardRepositoryResource{}
)

// standardRepositoryResource contains the APT standard repository resource's internal data.
type standardRepositoryResource struct {
	// client is the Proxmox VE API client.
	client proxmox.Client
}

// read reads information about an APT standard repository from the Proxmox VE API.
// Note that the name of the node must be set before this method is called!
func (r *standardRepositoryResource) read(ctx context.Context, srp *modelStandardRepo) (bool, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	data, err := r.client.Node(srp.Node.ValueString()).APT().Repositories().Get(ctx)
	if err != nil {
		diags.AddError("Could not read APT repositories", err.Error())

		return false, diags
	}

	for _, stdRepo := range data.StandardRepos {
		// Check if the APT standard repository is configured…
		if stdRepo.Handle == srp.Handle.ValueString() && stdRepo.Status == nil {
			// …handle the situation gracefully if not to signal that the repository has been removed outside of Terraform and
			// must be added back again.
			return false, diags
		}
	}

	srp.importFromAPI(ctx, data)

	return true, nil
}

// readBack reads information about an APT standard repository from the Proxmox VE API and then updates the response
// state accordingly.
func (r *standardRepositoryResource) readBack(
	ctx context.Context,
	srp *modelStandardRepo,
	diags *diag.Diagnostics,
	state *tfsdk.State,
) {
	found, readDiags := r.read(ctx, srp)

	diags.Append(readDiags...)

	if !found {
		diags.AddError(
			"APT standard repository resource not found after update",
			"Failed to find the resource when trying to read back the updated APT standard repository's data.",
		)
	}

	if !diags.HasError() {
		diags.Append(state.Set(ctx, *srp)...)
	}
}

// Configure adds the provider-configured client to the resource.
func (r *standardRepositoryResource) Configure(
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
			fmt.Sprintf("Expected config.Resource, got: %T", req.ProviderData),
		)

		return
	}

	r.client = cfg.Client
}

// Create adds an APT standard repository to the repository source lists.
// The name of this method might be a bit confusing for this resource, but this is due to the way how the Proxmox VE API
// works for APT standard repositories.
func (r *standardRepositoryResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var srp modelStandardRepo

	resp.Diagnostics.Append(req.Plan.Get(ctx, &srp)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := &api.AddRequestBody{
		Handle: srp.Handle.ValueString(),
		Node:   srp.Node.ValueString(),
	}

	if err := r.client.Node(srp.Node.ValueString()).APT().Repositories().Add(ctx, body); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not add APT standard repository with handle %v on node %v", srp.Handle, srp.Node),
			err.Error(),
		)
	}

	r.readBack(ctx, &srp, &resp.Diagnostics, &resp.State)
}

// Delete is currently a no-op for APT standard repositories due to the non-existing capability of the Proxmox VE API
// of deleting a configured APT standard repository.
// Also see Terraform's "Delete" framework documentation about [recommendations] and [caveats].
//
// [caveats]: https://developer.hashicorp.com/terraform/plugin/framework/resources/delete#caveats
// [recommendations]: https://developer.hashicorp.com/terraform/plugin/framework/resources/delete#recommendations
func (r *standardRepositoryResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

// ImportState imports an APT standard repository from the Proxmox VE API.
func (r *standardRepositoryResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	srp := modelStandardRepo{
		ID: types.StringValue(req.ID),
	}

	idFormatErrMsg := "expected import ID as comma-separated list in format " +
		"PROXMOX_VE_NODE_NAME,STANDARD_REPOSITORY_HANDLE (e.g. pve,no-subscription)"

	parts := strings.Split(srp.ID.ValueString(), ",")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid resource ID", fmt.Sprintf("%s, but got %v", idFormatErrMsg, srp.ID))

		return
	}

	srp.Node = types.StringValue(parts[0])
	srp.Handle = customtypes.StandardRepoHandleValue{StringValue: types.StringValue(parts[1])}

	resource.ImportStatePassthroughID(ctx, path.Root(SchemaAttrNameTerraformID), req, resp)
	r.readBack(ctx, &srp, &resp.Diagnostics, &resp.State)
}

// Metadata defines the name of the APT standard repository resource.
func (r *standardRepositoryResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceStandardRepoIDPrefix
}

// Read reads the APT standard repository.
func (r *standardRepositoryResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var srp modelStandardRepo

	resp.Diagnostics.Append(req.State.Get(ctx, &srp)...)

	if resp.Diagnostics.HasError() {
		return
	}

	found, diags := r.read(ctx, &srp)
	resp.Diagnostics.Append(diags...)

	if !resp.Diagnostics.HasError() {
		if found {
			resp.Diagnostics.Append(resp.State.Set(ctx, srp)...)
		} else {
			resp.State.RemoveResource(ctx)
		}
	}
}

// Schema defines the schema for the APT standard repository.
func (r *standardRepositoryResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages an APT standard repository of a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			SchemaAttrNameStandardDescription: schema.StringAttribute{
				Computed:    true,
				Description: "The description of the APT standard repository.",
			},
			SchemaAttrNameFilePath: schema.StringAttribute{
				Computed:    true,
				Description: "The absolute path of the source list file that contains this standard repository.",
			},
			SchemaAttrNameStandardHandle: schema.StringAttribute{
				CustomType:  customtypes.StandardRepoHandleType{},
				Description: "The handle of the APT standard repository.",
				MarkdownDescription: "The handle of the APT standard repository. Must be `ceph-quincy-enterprise` | " +
					"`ceph-quincy-no-subscription` | `ceph-quincy-test` | `ceph-reef-enterprise` | `ceph-reef-no-subscription` " +
					"| `ceph-reef-test` | `ceph-squid-enterprise` | `ceph-squid-no-subscription` | `ceph-squid-test` " +
					"| `enterprise` | `no-subscription` | `test`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Required: true,
				Validators: []validator.String{
					validators.NonEmptyString(),
				},
			},
			SchemaAttrNameIndex: schema.Int64Attribute{
				Computed:    true,
				Description: "The index within the defining source list file.",
			},
			SchemaAttrNameStandardName: schema.StringAttribute{
				Computed:    true,
				Description: "The name of the APT standard repository.",
			},
			SchemaAttrNameNode: schema.StringAttribute{
				Description: "The name of the target Proxmox VE node.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Required: true,
				Validators: []validator.String{
					validators.NonEmptyString(),
				},
			},
			SchemaAttrNameStandardStatus: schema.Int64Attribute{
				Computed:    true,
				Description: "Indicates the activation status.",
			},
			SchemaAttrNameTerraformID: attribute.ResourceID(
				"The unique identifier of this APT standard repository resource.",
			),
		},
	}
}

// Update is currently a no-op for APT repositories due to the non-existing capability of the Proxmox VE API of updating
// a configured APT standard repository.
// Also see Terraform's "Delete" framework documentation about [recommendations] and [caveats].
//
// [caveats]: https://developer.hashicorp.com/terraform/plugin/framework/resources/delete#caveats
// [recommendations]: https://developer.hashicorp.com/terraform/plugin/framework/resources/delete#recommendations
func (r *standardRepositoryResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

// NewStandardRepositoryResource returns a new resource for managing an APT standard repository.
// This is a helper function to simplify the provider implementation.
func NewStandardRepositoryResource() resource.Resource {
	return &standardRepositoryResource{}
}
