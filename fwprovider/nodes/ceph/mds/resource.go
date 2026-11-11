/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package mds

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	mdsapi "github.com/bpg/terraform-provider-proxmox/proxmox/nodes/ceph/mds"
)

var (
	_ resource.Resource                = &cephMDSResource{}
	_ resource.ResourceWithConfigure   = &cephMDSResource{}
	_ resource.ResourceWithImportState = &cephMDSResource{}
)

// NewCephMDSResource creates a new resource for managing Ceph metadata servers.
func NewCephMDSResource() resource.Resource {
	return &cephMDSResource{}
}

// cephMDSResource holds the provider-wide API client. The MDS subclient is resolved per call from node_name.
type cephMDSResource struct {
	client proxmox.Client
}

// Metadata defines the resource type name.
func (r *cephMDSResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_ceph_mds"
}

// Schema defines the schema for the resource.
func (r *cephMDSResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Ceph metadata server (MDS) on a Proxmox VE node.",
		MarkdownDescription: "Manages a Ceph metadata server (MDS) on a Proxmox VE node. CephFS needs a running MDS, and " +
			"Proxmox VE only creates a CephFS while one is in standby.",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ResourceID("The unique identifier of the MDS: its ID within the Ceph cluster, the same as `name`."),
			"node_name": schema.StringAttribute{
				Description: "The node that runs the MDS.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The MDS ID, unique within the Ceph cluster. It must start with a letter and contain only " +
					"letters, digits and hyphens, not ending with a hyphen.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
					// Proxmox's service ID pattern, which createmds further restricts to a leading letter.
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z]([a-zA-Z0-9-]*[a-zA-Z0-9])?$`),
						"must start with a letter and contain only letters, digits and hyphens, not ending with a hyphen",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Configure captures the provider-configured API client.
func (r *cephMDSResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// mdsClient returns the per-node Ceph MDS subclient for the given node.
func (r *cephMDSResource) mdsClient(nodeName string) *mdsapi.Client {
	return r.client.Node(nodeName).Ceph().MDS()
}

// Create creates the MDS. It returns once the Proxmox task ends, without waiting for a particular Ceph state.
func (r *cephMDSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan cephMDSModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := r.mdsClient(plan.NodeName.ValueString())

	result := client.Create(ctx, plan.Name.ValueString())
	if result.AddDiags(&resp.Diagnostics, fmt.Sprintf("Unable to Create Ceph MDS %q", plan.Name.ValueString())) {
		return
	}

	r.readBack(ctx, &plan, "After Creation", &resp.Diagnostics, &resp.State)
}

// Read refreshes the MDS from the cluster's MDS list.
func (r *cephMDSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cephMDSModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	found, diags := r.read(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !found {
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update only reads the MDS back: every configurable attribute requires replacement, and Proxmox has no update call.
func (r *cephMDSResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan cephMDSModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	r.readBack(ctx, &plan, "After Update", &resp.Diagnostics, &resp.State)
}

// Delete destroys the MDS through the node that runs it.
func (r *cephMDSResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state cephMDSModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := r.mdsClient(state.NodeName.ValueString())

	result := client.Delete(ctx, state.Name.ValueString())
	if err := result.Err(); err != nil && errors.Is(err, api.ErrResourceDoesNotExist) {
		return
	}

	result.AddDiags(&resp.Diagnostics, fmt.Sprintf("Unable to Delete Ceph MDS %q", state.Name.ValueString()))
}

// ImportState parses the import ID `node_name/name`; Read then fills in the rest.
func (r *cephMDSResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	nodeName, name, found := strings.Cut(req.ID, "/")
	if !found || nodeName == "" || name == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier in format 'node_name/name', got %q.", req.ID),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("node_name"), nodeName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)
}

// readBack reads the MDS after a create or update and saves it to state. The MDS must be listed by then.
func (r *cephMDSResource) readBack(
	ctx context.Context,
	data *cephMDSModel,
	after string,
	respDiags *diag.Diagnostics,
	respState *tfsdk.State,
) {
	found, diags := r.read(ctx, data)
	respDiags.Append(diags...)

	if respDiags.HasError() {
		return
	}

	if !found {
		respDiags.AddError(
			fmt.Sprintf("Unable to Read Ceph MDS %q %s", data.Name.ValueString(), after),
			"The task succeeded, but the cluster does not list the MDS.",
		)

		return
	}

	respDiags.Append(respState.Set(ctx, data)...)
}

// read looks the MDS up in the cluster's list and merges it into data. It returns false when the MDS isn't listed.
func (r *cephMDSResource) read(ctx context.Context, data *cephMDSModel) (bool, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	d, err := r.mdsClient(data.NodeName.ValueString()).Get(ctx, data.Name.ValueString())
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			return false, diags
		}

		diags.AddError(fmt.Sprintf("Unable to Read Ceph MDS %q", data.Name.ValueString()), err.Error())

		return false, diags
	}

	data.fromAPI(d)

	return true, diags
}
