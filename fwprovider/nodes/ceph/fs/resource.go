/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fs

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	fsapi "github.com/bpg/terraform-provider-proxmox/proxmox/nodes/ceph/fs"
	poolapi "github.com/bpg/terraform-provider-proxmox/proxmox/nodes/ceph/pool"
)

var (
	_ resource.Resource                = &cephFSResource{}
	_ resource.ResourceWithConfigure   = &cephFSResource{}
	_ resource.ResourceWithImportState = &cephFSResource{}
)

// destroyGuideURL is the admin guide's section on removing a CephFS by hand.
const destroyGuideURL = "https://pve.proxmox.com/pve-docs/chapter-pveceph.html#_destroy_cephfs"

// NewCephFSResource creates a new resource for managing a CephFS.
func NewCephFSResource() resource.Resource {
	return &cephFSResource{}
}

// cephFSResource holds the provider-wide API client. The subclients are resolved per call from node_name.
type cephFSResource struct {
	client proxmox.Client
}

// Metadata defines the resource type name.
func (r *cephFSResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_ceph_fs"
}

// Schema defines the schema for the resource.
func (r *cephFSResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a CephFS on the Proxmox VE cluster's own Ceph cluster.",
		MarkdownDescription: "Manages a CephFS on the Proxmox VE cluster's own Ceph cluster. Proxmox VE creates the " +
			"filesystem with a new metadata pool and data pool, and only while a metadata server (MDS) is in standby. " +
			"Mount it with a `proxmox_storage_cephfs`.",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ResourceID("The unique identifier of the filesystem: its name."),
			"node_name": schema.StringAttribute{
				Description: "The node used to dispatch the API calls. Any node running Ceph works; the filesystem is " +
					"cluster-wide.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The filesystem's name. Proxmox VE names its pools `<name>_metadata` and `<name>_data`.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[^:/\s]+$`), "must not contain ':', '/' or whitespace"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"pg_num": schema.Int64Attribute{
				Description: "The number of placement groups for the data pool at creation; the metadata pool gets a " +
					"quarter of it, and at least 8. Proxmox VE defaults to 128. Only used at creation, since the PG " +
					"autoscaler may change it later; changing it replaces the filesystem.",
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(8, 32768),
				},
				PlanModifiers: []planmodifier.Int64{
					// An import has no pg_num to record, so a null prior value isn't a change.
					int64planmodifier.RequiresReplaceIf(
						func(_ context.Context, req planmodifier.Int64Request, resp *int64planmodifier.RequiresReplaceIfFuncResponse) {
							resp.RequiresReplace = !req.StateValue.IsNull() && !req.PlanValue.IsNull() &&
								!req.PlanValue.Equal(req.StateValue)
						},
						"Changing pg_num replaces the filesystem.",
						"Changing `pg_num` replaces the filesystem.",
					),
				},
			},
			"metadata_pool": poolAttribute("The filesystem's metadata pool, `<name>_metadata`."),
			"data_pool":     poolAttribute("The filesystem's data pool, `<name>_data`."),
		},
	}
}

// poolAttribute describes one of the filesystem's pools. Proxmox VE creates them under the cluster's default CRUSH
// rule and size; any setting given here is applied right after creation and kept in place afterwards.
func poolAttribute(description string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Description: description + " Its replication settings default to the cluster's defaults.",
		MarkdownDescription: description + " Proxmox VE creates it under the cluster's default CRUSH rule and size; " +
			"settings given here are applied right after creation and can be changed in place.",
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The pool's name.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pool_id": schema.Int64Attribute{
				Description: "Ceph's numeric ID of the pool. Null with pve-manager releases before 9.1.18, which don't " +
					"report it.",
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"crush_rule": schema.StringAttribute{
				Description: "The name of the CRUSH rule that places the pool's data.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"size": schema.Int64Attribute{
				Description: "The number of replicas per object.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 7),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"min_size": schema.Int64Attribute{
				Description: "The minimum number of replicas for the pool to accept writes.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 7),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure captures the provider-configured API client.
func (r *cephFSResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cephFSResource) fsClient(nodeName string) *fsapi.Client {
	return r.client.Node(nodeName).Ceph().FS()
}

func (r *cephFSResource) poolClient(nodeName string) *poolapi.Client {
	return r.client.Node(nodeName).Ceph().Pool()
}

// Create creates the filesystem, then applies the planned pool settings.
func (r *cephFSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan cephFSModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()

	result := r.fsClient(plan.NodeName.ValueString()).Create(ctx, name, &fsapi.CreateRequestBody{
		PGNum: attribute.Int64PtrFromValue(plan.PGNum),
	})
	if result.AddDiags(&resp.Diagnostics, fmt.Sprintf("Unable to Create CephFS %q", name)) {
		return
	}

	// The filesystem exists from here on, so its state is saved even if a pool setting fails.
	r.applyAndReadBack(ctx, &plan, "After Creation", &resp.Diagnostics)

	if plan.ID.IsNull() || plan.ID.IsUnknown() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the filesystem and its pools.
func (r *cephFSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cephFSModel

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

// Update applies changed pool settings; every other attribute requires replacement.
func (r *cephFSResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan cephFSModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	r.applyAndReadBack(ctx, &plan, "After Update", &resp.Diagnostics)

	if plan.ID.IsNull() || plan.ID.IsUnknown() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete doesn't remove the filesystem: Proxmox VE documents that as a manual procedure. It fails while the
// filesystem exists, and succeeds once it has been removed by hand.
func (r *cephFSResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state cephFSModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	name := state.Name.ValueString()

	_, err := r.fsClient(state.NodeName.ValueString()).Get(ctx, name)
	if errors.Is(err, api.ErrResourceDoesNotExist) {
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Delete CephFS %q", name), err.Error())

		return
	}

	resp.Diagnostics.AddError(
		fmt.Sprintf("Unable to Delete CephFS %q", name),
		fmt.Sprintf("Proxmox VE documents removing a CephFS as a manual procedure (%s), so the provider doesn't "+
			"remove it. Remove it by hand, then run the destroy again:\n"+
			"1. Disconnect its clients, and disable and remove its CephFS storages.\n"+
			"2. Unmount those storages on every node: umount /mnt/pve/<storage>\n"+
			"3. Make sure no MDS serves the filesystem: stop its MDS daemons, standbys first, or run "+
			"`ceph fs fail %s`.\n"+
			"4. Destroy it with its pools: pveceph fs destroy %s --remove-pools 1",
			destroyGuideURL, name, name),
	)
}

// ImportState parses the import ID `node_name/name`; Read then fills in the rest. pg_num stays unset.
func (r *cephFSResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

// applyAndReadBack applies the planned pool settings and reads the filesystem back into data. It leaves data.ID
// unset if the filesystem can't be read.
func (r *cephFSResource) applyAndReadBack(ctx context.Context, data *cephFSModel, after string, diags *diag.Diagnostics) {
	name := data.Name.ValueString()
	nodeName := data.NodeName.ValueString()

	fs, err := r.fsClient(nodeName).Get(ctx, name)
	if err != nil {
		diags.AddError(fmt.Sprintf("Unable to Read CephFS %q %s", name, after), err.Error())

		return
	}

	dataPool, _ := fs.FirstDataPool()

	for _, p := range []struct {
		name    string
		planned types.Object
	}{
		{fs.MetadataPool, data.MetadataPool},
		{dataPool, data.DataPool},
	} {
		settings, d := poolSettings(ctx, p.planned)
		diags.Append(d...)

		if err := r.applyPoolSettings(ctx, nodeName, p.name, settings); err != nil {
			diags.AddError(fmt.Sprintf("Unable to Update Ceph pool %q of CephFS %q", p.name, name), err.Error())
		}
	}

	found, d := r.read(ctx, data)
	diags.Append(d...)

	if !found && !d.HasError() {
		diags.AddError(
			fmt.Sprintf("Unable to Read CephFS %q %s", name, after),
			"The task succeeded, but the cluster does not list the filesystem.",
		)
	}
}

// applyPoolSettings changes a pool's replication settings where they differ from the wanted ones. Size goes first and
// on its own: Ceph resets min_size whenever size changes, and Proxmox VE drops a min_size sent in the same request
// (Proxmox bug 8193).
func (r *cephFSResource) applyPoolSettings(ctx context.Context, nodeName, pool string, want *poolapi.UpdateRequestBody) error {
	client := r.poolClient(nodeName)

	current, err := client.Get(ctx, pool)
	if err != nil {
		return err
	}

	if want.Size != nil && *want.Size != current.Size {
		if err := client.Update(ctx, pool, &poolapi.UpdateRequestBody{Size: want.Size}).Err(); err != nil {
			return err
		}

		if current, err = client.Get(ctx, pool); err != nil {
			return err
		}
	}

	rest := &poolapi.UpdateRequestBody{}

	if want.CrushRule != nil && *want.CrushRule != current.CrushRule {
		rest.CrushRule = want.CrushRule
	}

	if want.MinSize != nil && *want.MinSize != current.MinSize {
		rest.MinSize = want.MinSize
	}

	if rest.CrushRule == nil && rest.MinSize == nil {
		return nil
	}

	return client.Update(ctx, pool, rest).Err()
}

// read looks the filesystem and its pools up and merges them into data. It returns false when the filesystem isn't
// listed. pg_num is kept as it is, since it only matters at creation.
func (r *cephFSResource) read(ctx context.Context, data *cephFSModel) (bool, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	name := data.Name.ValueString()
	nodeName := data.NodeName.ValueString()

	fs, err := r.fsClient(nodeName).Get(ctx, name)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			return false, diags
		}

		diags.AddError(fmt.Sprintf("Unable to Read CephFS %q", name), err.Error())

		return false, diags
	}

	metadata, err := r.poolClient(nodeName).Get(ctx, fs.MetadataPool)
	if err != nil {
		diags.AddError(fmt.Sprintf("Unable to Read Ceph pool %q of CephFS %q", fs.MetadataPool, name), err.Error())

		return false, diags
	}

	dataPool, dataPoolID := fs.FirstDataPool()

	dataStatus, err := r.poolClient(nodeName).Get(ctx, dataPool)
	if err != nil {
		diags.AddError(fmt.Sprintf("Unable to Read Ceph pool %q of CephFS %q", dataPool, name), err.Error())

		return false, diags
	}

	data.ID = data.Name
	data.MetadataPool = poolObject(fs.MetadataPool, fs.MetadataPoolID, metadata)
	data.DataPool = poolObject(dataPool, dataPoolID, dataStatus)

	return true, diags
}
