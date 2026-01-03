/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
)

var (
	_ resource.Resource                = (*realmSyncResource)(nil)
	_ resource.ResourceWithConfigure   = (*realmSyncResource)(nil)
	_ resource.ResourceWithImportState = (*realmSyncResource)(nil)
)

type realmSyncResource struct {
	client proxmox.Client
}

// NewRealmSyncResource creates a new realm sync resource.
func NewRealmSyncResource() resource.Resource {
	return &realmSyncResource{}
}

func (r *realmSyncResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_realm_sync"
}

func (r *realmSyncResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Triggers synchronization of an existing authentication realm.",
		MarkdownDescription: "Triggers synchronization of an existing authentication realm using `/access/domains/{realm}/sync`. " +
			"This resource represents the last requested sync configuration; deleting it does not undo the sync.",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ResourceID("Unique sync identifier (same as realm)."),
			"realm": schema.StringAttribute{
				Description: "Name of the realm to synchronize.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"scope": schema.StringAttribute{
				Description: "Sync scope: users, groups, or both.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("users", "groups", "both"),
				},
			},
			"remove_vanished": schema.StringAttribute{
				Description: "How to handle vanished entries (e.g. `acl;properties;entry` or `none`).",
				Optional:    true,
			},
			"enable_new": schema.BoolAttribute{
				Description: "Enable newly synced users.",
				Optional:    true,
			},
			"full": schema.BoolAttribute{
				Description:        "Perform a full sync.",
				Optional:           true,
				DeprecationMessage: "Deprecated by Proxmox: use remove_vanished instead.",
			},
			"purge": schema.BoolAttribute{
				Description:        "Purge removed entries.",
				Optional:           true,
				DeprecationMessage: "Deprecated by Proxmox: use remove_vanished instead.",
			},
			"dry_run": schema.BoolAttribute{
				Description: "Only simulate the sync without applying changes.",
				Optional:    true,
			},
		},
	}
}

func (r *realmSyncResource) Configure(
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

func (r *realmSyncResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan realmSyncModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := plan.toSyncRequest()

	if err := r.client.Access().SyncRealm(ctx, plan.Realm.ValueString(), body); err != nil {
		resp.Diagnostics.AddError("Error syncing realm", err.Error())
		return
	}

	// Use realm name as a stable ID for this sync config.
	plan.ID = plan.Realm

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *realmSyncResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// There is no persistent server-side representation of a sync operation.
	// Keep the current state as-is so Terraform can track the last requested
	// sync configuration.
	var state realmSyncModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if (state.ID.IsNull() || state.ID.ValueString() == "") && !state.Realm.IsNull() && state.Realm.ValueString() != "" {
		state.ID = state.Realm
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *realmSyncResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan realmSyncModel
	var state realmSyncModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := plan.toSyncRequest()

	if err := r.client.Access().SyncRealm(ctx, plan.Realm.ValueString(), body); err != nil {
		resp.Diagnostics.AddError("Error syncing realm", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *realmSyncResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Sync is a one-shot operation; there is nothing to delete remotely.
	// Terraform will simply drop this resource from state.
}

func (r *realmSyncResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Import by realm name.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("realm"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
