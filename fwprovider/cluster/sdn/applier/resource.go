/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package applier

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/applier"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &Resource{}
	_ resource.ResourceWithConfigure = &Resource{}
)

type model struct {
	// Opaque ID set timestamp at creation time.
	ID types.String `tfsdk:"id"`
}

type Resource struct {
	client *applier.Client
}

func NewResource() resource.Resource {
	return &Resource{}
}

func (r *Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdn_applier"
}

func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Applies pending Proxmox SDN configuration (cluster-wide).",
		MarkdownDescription: "**EXPERIMENTAL** Triggers Proxmox's SDN **Apply** (equivalent to `PUT /cluster/sdn`)." +
			"Intended to be used with `replace_triggered_by` so it runs after SDN objects change.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Opaque identifier set to the Unix timestamp (milliseconds) when the apply was executed.",
			},
		},
	}
}

func (r *Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = cfg.Client.Cluster().SDNApplier()
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if err := r.client.ApplyConfig(ctx); err != nil {
		resp.Diagnostics.AddError("Unable to Apply SDN Configuration", err.Error())
		return
	}

	state := &model{
		ID: types.StringValue(strconv.FormatInt(time.Now().UTC().UnixMilli(), 10)),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Nothing to refresh
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// We expect replacements only. But if someone does in-place Update,
	// we just re-run apply for safety and bump the ID timestamp.
	if err := r.client.ApplyConfig(ctx); err != nil {
		resp.Diagnostics.AddError("Unable to Re-Apply SDN Configuration", err.Error())
		return
	}

	var plan model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(strconv.FormatInt(time.Now().UTC().UnixMilli(), 10))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if err := r.client.ApplyConfig(ctx); err != nil {
		resp.Diagnostics.AddError("Unable to Re-Apply SDN Configuration", err.Error())
		return
	}
}
