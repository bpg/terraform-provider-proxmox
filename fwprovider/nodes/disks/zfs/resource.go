/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package zfs

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	zfsapi "github.com/bpg/terraform-provider-proxmox/proxmox/nodes/disks/zfs"
)

var (
	_ resource.Resource                   = &zfsPoolResource{}
	_ resource.ResourceWithConfigure      = &zfsPoolResource{}
	_ resource.ResourceWithImportState    = &zfsPoolResource{}
	_ resource.ResourceWithValidateConfig = &zfsPoolResource{}
)

// NewZFSPoolResource creates a new resource for managing ZFS pools.
func NewZFSPoolResource() resource.Resource {
	return &zfsPoolResource{}
}

type zfsPoolResource struct {
	client proxmox.Client
}

// Metadata defines the resource type name.
func (r *zfsPoolResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_disks_zfs"
}

// Schema defines the schema for the resource.
func (r *zfsPoolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ZFS pool (zpool) on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ResourceID("The unique identifier of this resource, in the format `<node_name>/<name>`."),
			"node_name": schema.StringAttribute{
				Description: "The name of the Proxmox node on which to create the ZFS pool.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the ZFS pool (storage identifier).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"devices": schema.ListAttribute{
				Description: "The block devices to use for the ZFS pool (e.g. `[\"/dev/sdb\", \"/dev/sdc\"]`).",
				Required:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					// Only replace when the value actually changes; null state (post-import) is not a change.
					listplanmodifier.RequiresReplaceIf(
						func(_ context.Context, req planmodifier.ListRequest, resp *listplanmodifier.RequiresReplaceIfFuncResponse) {
							resp.RequiresReplace = !req.StateValue.IsNull() && !req.PlanValue.Equal(req.StateValue)
						},
						"Requires replacement if the value changes after initial creation.",
						"Requires replacement if the value changes after initial creation.",
					),
				},
			},
			"raidlevel": schema.StringAttribute{
				Description: "The RAID level for the ZFS pool. One of `single`, `mirror`, `raid10`, `raidz`, `raidz2`, `raidz3`, `draid`, `draid2`, `draid3`.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"single", "mirror", "raid10",
						"raidz", "raidz2", "raidz3",
						"draid", "draid2", "draid3",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIf(
						func(_ context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
							resp.RequiresReplace = !req.StateValue.IsNull() && !req.PlanValue.Equal(req.StateValue)
						},
						"Requires replacement if the value changes after initial creation.",
						"Requires replacement if the value changes after initial creation.",
					),
				},
			},
			"ashift": schema.Int64Attribute{
				Description: "Pool sector size exponent (2^ashift bytes). Defaults to `12` (4 KiB sectors) server-side.",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(9, 16),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplaceIf(
						func(_ context.Context, req planmodifier.Int64Request, resp *int64planmodifier.RequiresReplaceIfFuncResponse) {
							resp.RequiresReplace = !req.StateValue.IsNull() && !req.PlanValue.Equal(req.StateValue)
						},
						"Requires replacement if the value changes after initial creation.",
						"Requires replacement if the value changes after initial creation.",
					),
				},
			},
			"compression": schema.StringAttribute{
				Description: "The compression algorithm for the pool. One of `on`, `off`, `gzip`, `lz4`, `lzjb`, `zle`, `zstd`. Defaults to `on` server-side.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("on", "off", "gzip", "lz4", "lzjb", "zle", "zstd"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIf(
						func(_ context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
							resp.RequiresReplace = !req.StateValue.IsNull() && !req.PlanValue.Equal(req.StateValue)
						},
						"Requires replacement if the value changes after initial creation.",
						"Requires replacement if the value changes after initial creation.",
					),
				},
			},
			"draid_config": schema.SingleNestedAttribute{
				Description: "dRAID configuration. Required when `raidlevel` is `draid`, `draid2`, or `draid3`.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"data": schema.Int64Attribute{
						Description: "Number of data devices per redundancy group.",
						Required:    true,
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
						},
					},
					"spares": schema.Int64Attribute{
						Description: "Number of dRAID distributed spare devices.",
						Required:    true,
						Validators: []validator.Int64{
							int64validator.AtLeast(0),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplaceIf(
						func(_ context.Context, req planmodifier.ObjectRequest, resp *objectplanmodifier.RequiresReplaceIfFuncResponse) {
							resp.RequiresReplace = !req.StateValue.IsNull() && !req.PlanValue.Equal(req.StateValue)
						},
						"Requires replacement if the value changes after initial creation.",
						"Requires replacement if the value changes after initial creation.",
					),
				},
			},
			"add_storage": schema.BoolAttribute{
				Description: "Configure a Proxmox storage entry using this zpool after creation. Applied at create time only.",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplaceIf(
						func(_ context.Context, req planmodifier.BoolRequest, resp *boolplanmodifier.RequiresReplaceIfFuncResponse) {
							resp.RequiresReplace = !req.StateValue.IsNull() && !req.PlanValue.Equal(req.StateValue)
						},
						"Requires replacement if the value changes after initial creation.",
						"Requires replacement if the value changes after initial creation.",
					),
				},
			},
			"cleanup_config": schema.BoolAttribute{
				Description: "On destroy, mark associated Proxmox storage entries as unavailable " +
					"(or remove them if configured for this node only). Defaults to `false`.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"cleanup_disks": schema.BoolAttribute{
				Description: "On destroy, wipe the ZFS member partitions so they can be reused. Defaults to `false`. " +
					"Note: Proxmox wipes the partition contents but leaves the parent disk's GPT partition table intact. " +
					"If you plan to reuse the same device immediately, run `wipefs -a <disk>` and `partprobe <disk>` " +
					"on the node after destroy.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"state": schema.StringAttribute{
				Description: "The current state of the ZFS pool (e.g. `ONLINE`, `DEGRADED`).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"errors": schema.StringAttribute{
				Description: "Error information reported by the ZFS pool.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// ValidateConfig checks attribute combinations before plan/apply.
func (r *zfsPoolResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var cfg zfsPoolModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)

	if resp.Diagnostics.HasError() {
		return
	}

	isDraid := attribute.IsDefined(cfg.RaidLevel) && strings.HasPrefix(cfg.RaidLevel.ValueString(), "draid")

	if isDraid && cfg.DraidConfig == nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("draid_config"),
			"Missing dRAID configuration",
			"The `draid_config` attribute is required when `raidlevel` is `draid`, `draid2`, or `draid3`.",
		)
	}

	if !isDraid && cfg.DraidConfig != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("draid_config"),
			"Invalid attribute combination",
			"The `draid_config` attribute is only valid when `raidlevel` is `draid`, `draid2`, or `draid3`.",
		)
	}
}

// Configure captures the provider-configured API client.
func (r *zfsPoolResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// zfsClient returns the per-node ZFS subclient.
func (r *zfsPoolResource) zfsClient(nodeName string) *zfsapi.Client {
	return r.client.Node(nodeName).Disks().ZFS()
}

// Create provisions a new ZFS pool.
func (r *zfsPoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan zfsPoolModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := plan.toCreateBody(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.zfsClient(plan.NodeName.ValueString())

	result := client.Create(ctx, body)
	if result.AddDiags(&resp.Diagnostics, fmt.Sprintf("Unable to Create ZFS pool %q", plan.Name.ValueString())) {
		return
	}

	plan.ID = types.StringValue(plan.NodeName.ValueString() + "/" + plan.Name.ValueString())

	data, err := client.Get(ctx, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Read ZFS pool %q after creation", plan.Name.ValueString()),
			err.Error(),
		)

		return
	}

	plan.fromAPI(data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the computed attributes from the API.
func (r *zfsPoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state zfsPoolModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := r.zfsClient(state.NodeName.ValueString())

	data, err := client.Get(ctx, state.Name.ValueString())
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Read ZFS pool %q", state.Name.ValueString()),
			err.Error(),
		)

		return
	}

	state.fromAPI(data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place changes to cleanup_config and cleanup_disks, which are
// local-only delete-time flags that require no API call. All other attributes carry
// RequiresReplace, so the Framework will never call Update for them.
func (r *zfsPoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan zfsPoolModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete destroys an existing ZFS pool.
func (r *zfsPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state zfsPoolModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := r.zfsClient(state.NodeName.ValueString())

	result := client.Delete(ctx, state.Name.ValueString(), state.toDeleteParams())
	if result.AddDiags(&resp.Diagnostics, fmt.Sprintf("Unable to Delete ZFS pool %q", state.Name.ValueString())) {
		return
	}
}

// ImportState imports an existing ZFS pool into Terraform state.
// The import ID must be `<node_name>/<pool_name>`.
// Note: write-only attributes (devices, raidlevel, ashift, compression, draid_config)
// are not populated on import and must be added to the configuration manually.
func (r *zfsPoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.SplitN(req.ID, "/", 2)
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format `<node_name>/<pool_name>`. Got: %q", req.ID),
		)

		return
	}

	nodeName := idParts[0]
	poolName := idParts[1]

	client := r.zfsClient(nodeName)

	data, err := client.Get(ctx, poolName)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError(
				"ZFS Pool Not Found",
				fmt.Sprintf("ZFS pool %q on node %q does not exist.", poolName, nodeName),
			)

			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Import ZFS pool %q", poolName),
			err.Error(),
		)

		return
	}

	state := zfsPoolModel{
		NodeName:      types.StringValue(nodeName),
		CleanupConfig: types.BoolValue(false),
		CleanupDisks:  types.BoolValue(false),
		// devices, raidlevel, ashift, compression, draid_config, add_storage are
		// write-only and cannot be reconstructed from the API.
		Devices:   types.ListNull(types.StringType),
		RaidLevel: types.StringNull(),
	}

	state.fromAPI(data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
