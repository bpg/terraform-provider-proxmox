/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package pool

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	poolapi "github.com/bpg/terraform-provider-proxmox/proxmox/nodes/ceph/pool"
)

var (
	_ resource.Resource                = &cephPoolResource{}
	_ resource.ResourceWithConfigure   = &cephPoolResource{}
	_ resource.ResourceWithImportState = &cephPoolResource{}
)

// NewCephPoolResource creates a new resource for managing Ceph pools.
func NewCephPoolResource() resource.Resource {
	return &cephPoolResource{}
}

// cephPoolResource holds the provider-wide API client. The pool subclient is resolved
// per-call from the model's node_name attribute because that lives on the resource,
// not the provider.
type cephPoolResource struct {
	client proxmox.Client
}

// Metadata defines the resource type name.
func (r *cephPoolResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_ceph_pool"
}

// Schema defines the schema for the resource.
func (r *cephPoolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages a Ceph pool on a Proxmox VE cluster.",
		MarkdownDescription: "Manages a Ceph pool on a Proxmox VE cluster.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The pool name. Must be unique within the Ceph cluster.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"node_name": schema.StringAttribute{
				Description: "The cluster node used to dispatch the API call. Any node running Ceph is acceptable; the pool itself is cluster-wide.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"application": schema.StringAttribute{
				Description: "The application using the pool. One of `rbd`, `cephfs`, `rgw`.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("rbd"),
				Validators: []validator.String{
					stringvalidator.OneOf("rbd", "cephfs", "rgw"),
				},
			},
			"crush_rule": schema.StringAttribute{
				Description: "The CRUSH rule name used for object placement.",
				Optional:    true,
				Computed:    true,
			},
			"erasure_coding": schema.StringAttribute{
				Description: "Create an erasure coded pool. Specified as `k+m[,profile=name]` (e.g. `4+2`). Cannot be changed after creation.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"min_size": schema.Int64Attribute{
				Description: "Minimum number of replicas per object.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 7),
				},
			},
			"pg_autoscale_mode": schema.StringAttribute{
				Description: "PG autoscaler mode. One of `on`, `off`, `warn`.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("on", "off", "warn"),
				},
			},
			"pg_num": schema.Int64Attribute{
				Description: "Number of placement groups.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 32768),
				},
			},
			"pg_num_min": schema.Int64Attribute{
				Description: "Minimum number of placement groups (used by the autoscaler).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"size": schema.Int64Attribute{
				Description: "Number of replicas per object.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 7),
				},
			},
			"target_size": schema.StringAttribute{
				Description: "Estimated target size for the PG autoscaler (e.g. `100G`). Write-only: " +
					"the PVE list endpoint returns this in bytes, so the configured spec is not round-tripped.",
				Optional: true,
			},
			"target_size_ratio": schema.Float64Attribute{
				Description:   "Estimated target ratio for the PG autoscaler.",
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.Float64{
					// API omits the field entirely when unset; preserve user intent across plan cycles.
				},
			},
			"add_storages": schema.BoolAttribute{
				Description:   "Configure VM and CT storage entries using the new pool. Applied at create time only.",
				Optional:      true,
				PlanModifiers: []planmodifier.Bool{
					// Create-only side effect.
				},
			},
			"force_destroy": schema.BoolAttribute{
				Description: "If true, destroy the pool even when in use. Passed as `force=1` on delete.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"remove_storages": schema.BoolAttribute{
				Description: "If true, remove all pveceph-managed storages configured for this pool on destroy.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"remove_ecprofile": schema.BoolAttribute{
				Description: "If true, remove the erasure code profile on destroy. Defaults to true. Only relevant for EC pools.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

// Configure captures the provider-configured API client.
func (r *cephPoolResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// poolClient returns the per-node Ceph pool subclient for the given node.
func (r *cephPoolResource) poolClient(nodeName string) *poolapi.Client {
	return r.client.Node(nodeName).Ceph().Pool()
}

// Create provisions a new Ceph pool.
func (r *cephPoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan cephPoolModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := r.poolClient(plan.NodeName.ValueString())

	result := client.Create(ctx, plan.toCreateBody())
	if result.AddDiags(&resp.Diagnostics, fmt.Sprintf("Unable to Create Ceph pool %q", plan.Name.ValueString())) {
		return
	}

	r.readBack(ctx, &plan, &resp.Diagnostics, &resp.State)
}

// Read fetches the current pool state from the cluster.
func (r *cephPoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cephPoolModel

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

// Update applies in-place changes to an existing pool.
func (r *cephPoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan cephPoolModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := r.poolClient(plan.NodeName.ValueString())

	result := client.Update(ctx, plan.Name.ValueString(), plan.toUpdateBody())
	if result.AddDiags(&resp.Diagnostics, fmt.Sprintf("Unable to Update Ceph pool %q", plan.Name.ValueString())) {
		return
	}

	r.readBack(ctx, &plan, &resp.Diagnostics, &resp.State)
}

// Delete destroys the pool, passing through the delete-only flags from state.
func (r *cephPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state cephPoolModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := r.poolClient(state.NodeName.ValueString())

	result := client.Delete(ctx, state.Name.ValueString(), state.toDeleteParams())
	if err := result.Err(); err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			return
		}

		result.AddDiags(&resp.Diagnostics, fmt.Sprintf("Unable to Delete Ceph pool %q", state.Name.ValueString()))
	}
}

// ImportState parses the composite import id `node_name/pool_name` and seeds the
// node_name + name attributes; the framework then runs Read to populate the rest.
func (r *cephPoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			"Expected import identifier in format 'node_name/pool_name'.",
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("node_name"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), parts[1])...)
}

// readBack performs Read and persists the result; surfaces an error if the pool is missing
// after a Create/Update (which would indicate eventual-consistency or a phantom failure).
func (r *cephPoolResource) readBack(
	ctx context.Context,
	data *cephPoolModel,
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
			fmt.Sprintf("Ceph pool %q not found after create/update", data.Name.ValueString()),
			"Failed to find the Ceph pool when reading it back after a create or update operation.",
		)

		return
	}

	respDiags.Append(respState.Set(ctx, data)...)
}

// read fetches the pool's current settings and merges them into data. Returns false
// when the pool no longer exists.
func (r *cephPoolResource) read(ctx context.Context, data *cephPoolModel) (bool, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	client := r.poolClient(data.NodeName.ValueString())

	p, err := client.Get(ctx, data.Name.ValueString())
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			return false, diags
		}

		diags.AddError(
			fmt.Sprintf("Unable to Read Ceph pool %q", data.Name.ValueString()),
			err.Error(),
		)

		return false, diags
	}

	data.fromAPI(p)

	return true, diags
}
