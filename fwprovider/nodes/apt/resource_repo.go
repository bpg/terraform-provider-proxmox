/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package apt

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/structure"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	api "github.com/bpg/terraform-provider-proxmox/proxmox/nodes/apt/repositories"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

const (
	// ResourceRepoIDPrefix is the prefix for the resource ID of resourceRepo.
	ResourceRepoIDPrefix = "apt_repository"

	// ResourceRepoActivationStatus is the default activation status for newly created or imported APT repositories.
	// This reflects the same default value used by the Proxmox VE API when the "enabled" parameter is not set.
	ResourceRepoActivationStatus = true
)

// Ensure the resource implements the required interfaces.
var (
	_ resource.Resource                = &resourceRepo{}
	_ resource.ResourceWithConfigure   = &resourceRepo{}
	_ resource.ResourceWithImportState = &resourceRepo{}
)

// resourceRepo contains the APT repository resource's internal data.
type resourceRepo struct {
	// client is the Proxmox VE API client.
	client proxmox.Client
}

// read reads information about an APT repository from the Proxmox VE API.
// Note that the name of the node must be set before this method is called!
func (r *resourceRepo) read(ctx context.Context, rp *modelRepo) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	data, err := r.client.Node(rp.Node.ValueString()).APT().Repositories().Get(ctx)
	if err != nil {
		diags.AddError(fmt.Sprintf("Could not read APT repositories on node %v", rp.Node), err.Error())

		return false, diags
	}

	diags.Append(rp.importFromAPI(ctx, data)...)

	if diags.HasError() {
		return false, diags
	}

	return true, nil
}

// readBack reads information about an APT repository from the Proxmox VE API and then updates the response state
// accordingly.
// Note that the Terraform resource identifier must be set in the state before this method is called!
func (r *resourceRepo) readBack(ctx context.Context, rp *modelRepo, diags *diag.Diagnostics, state *tfsdk.State) {
	found, readDiags := r.read(ctx, rp)

	diags.Append(readDiags...)

	if !found {
		diags.AddError(
			"APT repository resource not found after update",
			"Failed to find the resource when trying to read back the updated APT repository's data.",
		)
	}

	if !diags.HasError() {
		diags.Append(state.Set(ctx, *rp)...)
	}
}

// Configure adds the provider-configured client to the resource.
func (r *resourceRepo) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(proxmox.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected resource configuration type",
			fmt.Sprintf("Expected *proxmox.Client, got: %T", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create modifies the activation state of an existing APT repository, including the addition of standard repositories
// to the repository lists.
// The name of this method might be a bit confusing for this resource, but this is due to the way how the Proxmox VE API
// works for APT repositories.
func (r *resourceRepo) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var rp modelRepo

	resp.Diagnostics.Append(req.Plan.Get(ctx, &rp)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := &api.ModifyRequestBody{
		Enabled: proxmoxtypes.CustomBool(rp.Enabled.ValueBool()),
		Index:   rp.Index.ValueInt64(),
		Path:    rp.FilePath.ValueString(),
	}

	if err := r.client.Node(rp.Node.ValueString()).APT().Repositories().Modify(ctx, body); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not modify APT repository in file %v at index %v on node %v", rp.FilePath, rp.Index, rp.Node),
			err.Error(),
		)

		return
	}

	r.readBack(ctx, &rp, &resp.Diagnostics, &resp.State)
}

// Delete is currently a no-op for APT repositories due to the non-existing capability of the Proxmox VE API of deleting
// a configured APT repository.
// Also see Terraform's "Delete" framework documentation about [recommendations] and [caveats].
//
// [caveats]: https://developer.hashicorp.com/terraform/plugin/framework/resources/delete#caveats
// [recommendations]: https://developer.hashicorp.com/terraform/plugin/framework/resources/delete#recommendations
func (r *resourceRepo) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

// ImportState imports an APT repository from the Proxmox VE API.
func (r *resourceRepo) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	rp := modelRepo{
		Enabled: types.BoolValue(ResourceRepoActivationStatus),
		ID:      types.StringValue(req.ID),
	}

	idFormatErrMsg := "expected import ID as comma-separated list in format " +
		"PROXMOX_VE_NODE_NAME,SOURCE_LIST_FILE_PATH,INDEX (e.g. pve,/etc/apt/sources.list,0)"

	parts := strings.Split(rp.ID.ValueString(), ",")
	if len(parts) != 3 {
		resp.Diagnostics.AddError("Invalid resource ID", fmt.Sprintf("%s, but got %v", idFormatErrMsg, rp.ID))

		return
	}

	rp.Node = types.StringValue(parts[0])

	if !filepath.IsAbs(parts[1]) {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("given source list file path %q is not an absolute path: %s", parts[1], idFormatErrMsg),
		)

		return
	}

	rp.FilePath = types.StringValue(parts[1])

	index, err := strconv.Atoi(parts[2])
	if err != nil {
		resp.Diagnostics.AddError(
			"Parse resource ID",
			fmt.Sprintf("Failed to parse given import ID index parameter %q as number: %s", parts[2], idFormatErrMsg),
		)

		return
	}

	rp.Index = types.Int64Value(int64(index))

	resource.ImportStatePassthroughID(ctx, path.Root(SchemaAttrNameTerraformID), req, resp)
	r.readBack(ctx, &rp, &resp.Diagnostics, &resp.State)
}

// Metadata defines the name of the APT repository resource.
func (r *resourceRepo) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceRepoIDPrefix
}

// Read reads the APT repository.
func (r *resourceRepo) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var rp modelRepo

	resp.Diagnostics.Append(req.State.Get(ctx, &rp)...)

	if resp.Diagnostics.HasError() {
		return
	}

	found, diags := r.read(ctx, &rp)
	resp.Diagnostics.Append(diags...)

	if !resp.Diagnostics.HasError() {
		if found {
			resp.Diagnostics.Append(resp.State.Set(ctx, rp)...)
		} else {
			resp.State.RemoveResource(ctx)
		}
	}
}

// Schema defines the schema for the APT repository.
func (r *resourceRepo) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an APT repository of a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			SchemaAttrNameComment: schema.StringAttribute{
				Computed:    true,
				Description: "The associated comment.",
			},
			SchemaAttrNameComponents: schema.ListAttribute{
				Computed:    true,
				Description: "The list of components.",
				ElementType: types.StringType,
			},
			SchemaAttrNameEnabled: schema.BoolAttribute{
				Computed:    true,
				Default:     booldefault.StaticBool(ResourceRepoActivationStatus),
				Description: "Indicates the activation status.",
				Optional:    true,
			},
			SchemaAttrNameFilePath: schema.StringAttribute{
				Description: "The absolute path of the source list file that contains this repository.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Required: true,
				Validators: []validator.String{
					validators.AbsoluteFilePathValidator(),
					validators.NonEmptyString(),
				},
			},
			SchemaAttrNameFileType: schema.StringAttribute{
				Computed:    true,
				Description: "The format of the defining source list file.",
			},
			SchemaAttrNameIndex: schema.Int64Attribute{
				Description: "The index within the defining source list file.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Required: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
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
			SchemaAttrNamePackageTypes: schema.ListAttribute{
				Computed:    true,
				Description: "The list of package types.",
				ElementType: types.StringType,
			},
			SchemaAttrNameSuites: schema.ListAttribute{
				Computed:    true,
				Description: "The list of package distributions.",
				ElementType: types.StringType,
			},
			SchemaAttrNameTerraformID: structure.IDAttribute("The unique identifier of this APT repository resource."),
			SchemaAttrNameURIs: schema.ListAttribute{
				Computed:    true,
				Description: "The list of repository URIs.",
				ElementType: types.StringType,
			},
		},
	}
}

// Update updates an existing APT repository.
func (r *resourceRepo) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var rpPlan modelRepo

	resp.Diagnostics.Append(req.Plan.Get(ctx, &rpPlan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := &api.ModifyRequestBody{
		Enabled: proxmoxtypes.CustomBool(rpPlan.Enabled.ValueBool()),
		Index:   rpPlan.Index.ValueInt64(),
		Path:    rpPlan.FilePath.ValueString(),
	}

	err := r.client.Node(rpPlan.Node.ValueString()).APT().Repositories().Modify(ctx, body)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Could not modify APT repository in file %v at index %v on node %v",
				rpPlan.FilePath,
				rpPlan.Index,
				rpPlan.Node,
			),
			err.Error(),
		)

		return
	}

	r.readBack(ctx, &rpPlan, &resp.Diagnostics, &resp.State)
}

// NewResourceRepo returns a new resource for managing an APT repository.
// This is a helper function to simplify the provider implementation.
func NewResourceRepo() resource.Resource {
	return &resourceRepo{}
}
