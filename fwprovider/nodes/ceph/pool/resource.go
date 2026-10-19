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
	"maps"
	"reflect"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	poolapi "github.com/bpg/terraform-provider-proxmox/proxmox/nodes/ceph/pool"
)

var (
	_ resource.Resource                = &cephPoolResource{}
	_ resource.ResourceWithConfigure   = &cephPoolResource{}
	_ resource.ResourceWithImportState = &cephPoolResource{}
	_ resource.ResourceWithModifyPlan  = &cephPoolResource{}
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
		Description: "Manages a Ceph pool on a Proxmox VE cluster.",
		MarkdownDescription: "Manages a Ceph pool on a Proxmox VE cluster.\n\n" +
			"**Erasure-coded pools.** With `erasure_coding`, Proxmox creates no pool named `name`. It creates an " +
			"erasure-coded data pool `<name>-data` and a replicated metadata pool `<name>-metadata`, as " +
			"`pveceph pool create --erasure-coding` does. As in `pveceph`, `size`, `min_size` and `crush_rule` then " +
			"apply to the **metadata pool**, while `pg_num`, `pg_num_min`, `pg_autoscale_mode`, `application`, " +
			"`target_size` and `target_size_ratio` apply to the **data pool**. The data pool's own size, minimum " +
			"size and CRUSH rule follow from the erasure code profile and are reported in `data_pool`.\n\n" +
			"**Storages.** The provider never asks Proxmox to add a storage for the pool; manage RBD storages on it " +
			"with `proxmox_storage_rbd`. Destroying the pool removes no storage. Destroying an erasure-coded pool " +
			"removes both of its pools, and its erasure code profile when Proxmox created that profile (that is, " +
			"when `erasure_coding.profile` is not set).",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ResourceID("The unique identifier of the Ceph pool, in the form `<node_name>/<name>`."),
			"name": schema.StringAttribute{
				Description: "The pool name. Must be unique within the Ceph cluster. With `erasure_coding`, Proxmox " +
					"creates the pools `<name>-data` and `<name>-metadata` instead, and no pool of this name.",
				Required: true,
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
				Description: "The application using the pool. One of `rbd`, `cephfs`, `rgw`. Defaults to `rbd` server-side. " +
					"With `erasure_coding`, this is the data pool's application (Proxmox also sets it on the metadata " +
					"pool at creation).",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("rbd", "cephfs", "rgw"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"crush_rule": schema.StringAttribute{
				Description: "The CRUSH rule name used for object placement. With `erasure_coding`, it is used solely " +
					"for the metadata pool; the data pool uses the rule Ceph generates from the erasure code profile " +
					"(`data_pool.crush_rule`).",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"erasure_coding": schema.SingleNestedAttribute{
				Description: "Create an erasure-coded pool instead of a replicated one: Proxmox creates an " +
					"erasure-coded data pool `<name>-data` and a replicated metadata pool `<name>-metadata`. While " +
					"this is set, `size`, `min_size` and `crush_rule` describe the metadata pool. Cannot be changed " +
					"after creation, except that after an import the configured value is only recorded.",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"k": schema.Int64Attribute{
						Description: "The number of data chunks. Proxmox requires it even with `profile`, whose own " +
							"value then applies.",
						Required: true,
						Validators: []validator.Int64{
							int64validator.AtLeast(2),
						},
					},
					"m": schema.Int64Attribute{
						Description: "The number of coding chunks. Proxmox requires it even with `profile`, whose own " +
							"value then applies.",
						Required: true,
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
						},
					},
					"failure_domain": schema.StringAttribute{
						Description: "The CRUSH failure domain for the data pool's chunks, e.g. `host` or `osd`. " +
							"Proxmox defaults to `host`.",
						Optional: true,
						Validators: []validator.String{
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("profile")),
						},
					},
					"device_class": schema.StringAttribute{
						Description: "The CRUSH device class to place the data pool on, e.g. `hdd` or `ssd`.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("profile")),
						},
					},
					"profile": schema.StringAttribute{
						Description: "An existing erasure code profile to use instead of the one Proxmox would create " +
							"(`pve_ec_<name>`). Destroying the pool keeps this profile; a profile Proxmox created is " +
							"removed with the pool.",
						Optional: true,
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplaceIf(
						erasureCodingRequiresReplace,
						"Requires replacement if the value changes after creation, unless the pool was imported.",
						"Requires replacement if the value changes after creation, unless the pool was imported.",
					),
				},
			},
			"min_size": schema.Int64Attribute{
				Description: "Minimum number of replicas per object. With `erasure_coding`, this is the metadata " +
					"pool's; the data pool's is `data_pool.min_size`.",
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 7),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"pg_autoscale_mode": schema.StringAttribute{
				Description: "PG autoscaler mode. One of `on`, `off`, `warn`. With `erasure_coding`, this is the data " +
					"pool's (Proxmox also applies it to the metadata pool at creation).",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("on", "off", "warn"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pg_num": schema.Int64Attribute{
				Description: "Number of placement groups. With `erasure_coding`, this is the data pool's (Proxmox " +
					"creates the metadata pool with 32).",
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 32768),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"pg_num_min": schema.Int64Attribute{
				Description: "Minimum number of placement groups (used by the autoscaler). With `erasure_coding`, " +
					"this is the data pool's (Proxmox also applies it to the metadata pool at creation).",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"size": schema.Int64Attribute{
				Description: "Number of replicas per object. With `erasure_coding`, this is the metadata pool's " +
					"replica count; the data pool's size is `k + m` (`data_pool.size`).",
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 7),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"target_size": schema.StringAttribute{
				Description: "Estimated target size for the PG autoscaler. Write-only.",
				MarkdownDescription: "Estimated target size for the PG autoscaler (e.g. `100G`). " +
					"**Write-only:** PVE returns this as a bytes integer (e.g. `1073741824`) while " +
					"the input accepts a unit-suffixed string, so the configured spec is not " +
					"round-tripped from the server. With `erasure_coding`, this is the data pool's (Proxmox also " +
					"applies it to the metadata pool at creation).",
				Optional: true,
			},
			"target_size_ratio": schema.Float64Attribute{
				Description: "Estimated target ratio for the PG autoscaler. Write-only.",
				MarkdownDescription: "Estimated target ratio for the PG autoscaler. " +
					"**Write-only:** kept symmetric with `target_size` (which cannot be round-tripped), " +
					"so the configured value is not refreshed from the server. With `erasure_coding`, this is the " +
					"data pool's (Proxmox also applies it to the metadata pool at creation).",
				Optional: true,
			},
			"pool_id": schema.Int64Attribute{
				Description: "The numeric pool ID assigned by Ceph. Ceph assigns a new ID when the pool is recreated. " +
					"Null with `erasure_coding`, where no pool named `name` exists; see `data_pool.pool_id` and " +
					"`metadata_pool_id`.",
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"data_pool": schema.SingleNestedAttribute{
				Description: "The erasure-coded data pool `<name>-data` of an erasure-coded pool; null for a " +
					"replicated pool. Its size, minimum size and CRUSH rule follow from the erasure code profile.",
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "The data pool's name, `<name>-data`.",
						Computed:    true,
					},
					"pool_id": schema.Int64Attribute{
						Description: "The data pool's numeric ID assigned by Ceph.",
						Computed:    true,
					},
					"size": schema.Int64Attribute{
						Description: "The data pool's size, `k + m`.",
						Computed:    true,
					},
					"min_size": schema.Int64Attribute{
						Description: "The data pool's minimum size (Ceph sets `k + 1` at creation).",
						Computed:    true,
					},
					"crush_rule": schema.StringAttribute{
						Description: "The CRUSH rule Ceph generated for the data pool from the erasure code profile.",
						Computed:    true,
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"metadata_pool_name": schema.StringAttribute{
				Description: "The name of an erasure-coded pool's replicated metadata pool, `<name>-metadata`; null for " +
					"a replicated pool. `size`, `min_size` and `crush_rule` describe this pool.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"metadata_pool_id": schema.Int64Attribute{
				Description: "The numeric ID Ceph assigned to an erasure-coded pool's metadata pool; null for a " +
					"replicated pool.",
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"force_destroy": schema.BoolAttribute{
				Description: "If true, destroy the pool even when in use. Passed as `force=1` on delete. With " +
					"`erasure_coding`, this applies to both pools.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
		},
	}
}

// erasureCodingRequiresReplace replaces the pool when erasure_coding changes, except when a null value is recorded
// for a pool that is already erasure-coded: erasure_coding can't be read back, so it is null after an import.
func erasureCodingRequiresReplace(
	ctx context.Context,
	req planmodifier.ObjectRequest,
	resp *objectplanmodifier.RequiresReplaceIfFuncResponse,
) {
	if !req.StateValue.IsNull() {
		resp.RequiresReplace = true

		return
	}

	var dataPool types.Object

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("data_pool"), &dataPool)...)
	resp.RequiresReplace = dataPool.IsNull()
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

	body, diags := plan.toCreateBody(ctx)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := r.poolClient(plan.NodeName.ValueString())

	result := client.Create(ctx, body)
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
	var plan, state cephPoolModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := r.poolClient(plan.NodeName.ValueString())
	planBodies := plan.toUpdateBodies()
	stateBodies := state.toUpdateBodies()

	for _, name := range slices.Sorted(maps.Keys(planBodies)) {
		planBody, stateBody := planBodies[name], stateBodies[name]

		// erasure_coding and force_destroy exist only in state; don't resend the pool settings for them.
		if reflect.DeepEqual(planBody, stateBody) {
			continue
		}

		// Ceph resets min_size when size changes, and Proxmox compares min_size with the values it read before
		// applying size, so an unchanged min_size would be skipped and lost. Apply size on its own first.
		if planBody.Size != nil && !reflect.DeepEqual(planBody.Size, stateBody.Size) {
			result := client.Update(ctx, name, &poolapi.UpdateRequestBody{Size: planBody.Size})
			if result.AddDiags(&resp.Diagnostics, fmt.Sprintf("Unable to Update Ceph pool %q", name)) {
				return
			}
		}

		result := client.Update(ctx, name, planBody)
		if result.AddDiags(&resp.Diagnostics, fmt.Sprintf("Unable to Update Ceph pool %q", name)) {
			return
		}
	}

	r.readBack(ctx, &plan, &resp.Diagnostics, &resp.State)
}

// ModifyPlan requires an imported erasure-coded pool to be configured as one.
func (r *cephPoolResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	var plan, state cephPoolModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ErasureCoding.IsNull() && state.ErasureCoding.IsNull() && !state.DataPool.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("erasure_coding"),
			fmt.Sprintf("Ceph pool %q is erasure-coded", plan.Name.ValueString()),
			fmt.Sprintf("Proxmox created it as %q and %q. Add an erasure_coding attribute that matches how it was "+
				"created; it can't be read back from the API.",
				plan.Name.ValueString()+ecDataPoolSuffix, plan.Name.ValueString()+ecMetadataPoolSuffix),
		)
	}
}

// Delete destroys the pool, or both pools of an erasure-coded pool, leaving every storage to its own resource.
func (r *cephPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state cephPoolModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params, diags := state.toDeleteParams(ctx)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := r.poolClient(state.NodeName.ValueString())
	name := state.Name.ValueString()

	names := []string{name}
	if state.erasureCoded() {
		// The data pool first: destroying it also removes its CRUSH rule and, if requested, the profile.
		names = []string{name + ecDataPoolSuffix, name + ecMetadataPoolSuffix}
	}

	for _, n := range names {
		// A pool already gone (e.g. after a partly failed delete) needs no task.
		if _, err := client.Get(ctx, n); errors.Is(err, api.ErrResourceDoesNotExist) {
			continue
		}

		result := client.Delete(ctx, n, params)
		if err := result.Err(); err != nil && errors.Is(err, api.ErrResourceDoesNotExist) {
			continue
		}

		if result.AddDiags(&resp.Diagnostics, fmt.Sprintf("Unable to Delete Ceph pool %q", n)) {
			return
		}
	}
}

// ImportState parses the composite import id `node_name/pool_name` and seeds the
// node_name + name attributes and the force_destroy default; the framework then runs Read to populate the rest.
// For an erasure-coded pool, pool_name is the name it was created with, without `-data` or `-metadata`.
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

	// force_destroy can't be read back; record its default so an import plans no update.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("force_destroy"), false)...)
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
	client := r.poolClient(data.NodeName.ValueString())
	name := data.Name.ValueString()

	if !data.erasureCoded() {
		p, found, diags := getPool(ctx, client, name)
		if found {
			data.fromAPI(p)
		}

		// No pool of that name: after an import it may be an erasure-coded pool, which exists only as
		// <name>-data and <name>-metadata.
		if found || diags.HasError() {
			return found, diags
		}
	}

	dataPool, found, diags := getPool(ctx, client, name+ecDataPoolSuffix)
	if !found {
		return false, diags
	}

	metadataPool, found, diags := getPool(ctx, client, name+ecMetadataPoolSuffix)
	if !found {
		return false, diags
	}

	return true, data.fromAPIErasureCoded(ctx, dataPool, metadataPool)
}

// getPool reads one Ceph pool, reporting whether it exists.
func getPool(ctx context.Context, client *poolapi.Client, name string) (*poolapi.StatusResponseData, bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, err := client.Get(ctx, name)
	if err != nil {
		if !errors.Is(err, api.ErrResourceDoesNotExist) {
			diags.AddError(fmt.Sprintf("Unable to Read Ceph pool %q", name), err.Error())
		}

		return nil, false, diags
	}

	return p, true, diags
}
