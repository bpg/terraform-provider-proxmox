/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package replication

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
)

var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithConfigure   = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

type Resource struct {
	client *cluster.Client
}

func NewResource() resource.Resource {
	return &Resource{}
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_replication"
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

	r.client = cfg.Client.Cluster()
}

func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Proxmox VE Replication.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Replication Job ID. The ID is composed of a Guest ID and a job number, separated by a hyphen, i.e. '<GUEST>-<JOBNUM>'.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[0-9]+-[0-9]+$`),
						"id must be <GUEST>-<JOBNUM>",
					),
				},
			},
			"target": schema.StringAttribute{
				Description: "Target node.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "Section type.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("local"),
				},
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Description: "Description.",
			},
			"disable": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Flag to disable/deactivate this replication.",
			},
			"rate": schema.Float64Attribute{
				Optional:    true,
				Description: "Rate limit in mbps (megabytes per second) as floating point number.",
			},
			"schedule": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "Storage replication schedule. The format is a subset of `systemd` calendar events. Defaults to */15",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"source": schema.StringAttribute{
				Computed:    true,
				Description: "For internal use, to detect if the guest was stolen.",
			},
			"guest": schema.Int64Attribute{
				Computed:    true,
				Description: "Guest ID.",
			},
			"jobnum": schema.Int64Attribute{
				Computed:    true,
				Description: "Unique, sequential ID assigned to each job.",
			},
		},
	}
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	repl := plan.toAPICreate()

	err := r.client.Replication(plan.ID.ValueString()).CreateReplication(ctx, repl)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Replication", err.Error())
		return
	}

	// Read the created Replication to get the actual state including pending
	data, err := r.client.Replication(plan.ID.ValueString()).GetReplication(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Replication After Creation", err.Error())
		return
	}

	readModel := &model{}
	readModel.fromAPI(plan.ID.ValueString(), data)

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.Replication(state.ID.ValueString()).GetReplication(ctx)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Unable to Read Replication", err.Error())

		return
	}

	readModel := &model{}
	readModel.fromAPI(state.ID.ValueString(), data)

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan model

	var state model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	repl := plan.toAPIUpdate()

	var toDelete []string

	attribute.CheckDelete(plan.Comment, state.Comment, &toDelete, "comment")
	attribute.CheckDelete(plan.Disable, state.Disable, &toDelete, "disable")
	attribute.CheckDelete(plan.Rate, state.Rate, &toDelete, "rate")
	attribute.CheckDelete(plan.Schedule, state.Schedule, &toDelete, "schedule")

	repl.Delete = toDelete

	err := r.client.Replication(plan.ID.ValueString()).UpdateReplication(ctx, repl)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Update Replication", err.Error())
		return
	}

	// Read the updated Replication to get the actual state including pending
	data, err := r.client.Replication(plan.ID.ValueString()).GetReplication(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Replication After Update", err.Error())
		return
	}

	readModel := &model{}
	readModel.fromAPI(plan.ID.ValueString(), data)

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	repl := state.toAPIDelete()

	err := r.client.Replication(state.ID.ValueString()).DeleteReplication(ctx, repl)
	if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
		resp.Diagnostics.AddError("Unable to Delete Replication", err.Error())
	}
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	data, err := r.client.Replication(req.ID).GetReplication(ctx)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError("Replication Not Found", fmt.Sprintf("Replication with ID '%s' was not found", req.ID))
			return
		}

		resp.Diagnostics.AddError("Unable to Import Replication", err.Error())

		return
	}

	readModel := &model{}
	readModel.fromAPI(req.ID, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}
