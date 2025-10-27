package pools

import (
	"context"
	"fmt"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/pools"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"strconv"
)

var (
	_ resource.Resource                = (*poolMembershipResource)(nil)
	_ resource.ResourceWithConfigure   = (*poolMembershipResource)(nil)
	_ resource.ResourceWithImportState = (*poolMembershipResource)(nil)
	//_ resource.ResourceWithConfigValidators = (*poolMembershipResource)(nil)
)

type poolMembershipResource struct {
	client proxmox.Client
}

func NewPoolMembershipResource() resource.Resource {
	return &poolMembershipResource{}
}

func (r *poolMembershipResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *poolMembershipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"pool_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vm_id": schema.Int64Attribute{
				Required: true, // consider changing to Optional if storage membership is supported
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *poolMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan poolMembershipModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	poolApi := r.client.Pool()

	poolId := plan.PoolId.ValueString()
	vmId := plan.VmID.ValueInt64()

	vmList := (types.CustomCommaSeparatedList)([]string{strconv.FormatInt(vmId, 10)})

	trueValue := types.CustomBool(true)
	body := &pools.PoolUpdateRequestBody{
		VMs:       &vmList,
		AllowMove: &trueValue,
	}

	if err := poolApi.UpdatePool(ctx, poolId, body); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to update resource pool '%s'", poolId),
			err.Error())
		return
	}
	plan.ID = plan.generateID()
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *poolMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var state poolMembershipModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	poolId := state.PoolId.ValueString()
	vmId := state.VmID.ValueInt64()

	pool, err := r.client.Pool().GetPool(ctx, poolId)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to get pool '%s'", poolId), err.Error())
		return
	}

	for _, member := range pool.Members {
		if member.VMID != nil && int64(*member.VMID) == vmId {
			return
		}
	}

	// Membership not found
	resp.State.RemoveResource(ctx)
}

func (r *poolMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state poolMembershipModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	poolId := state.PoolId.ValueString()
	vmId := state.VmID.ValueInt64()

	vmList := (types.CustomCommaSeparatedList)([]string{strconv.FormatInt(vmId, 10)})

	trueValue := types.CustomBool(true)
	body := &pools.PoolUpdateRequestBody{
		VMs:    &vmList,
		Delete: &trueValue,
	}

	if err := r.client.Pool().UpdatePool(ctx, poolId, body); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to update pool '%s", poolId), err.Error())
	}
}

func (r *poolMembershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	model, err := parseMembershipModelFromID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Unable to import pool membership", fmt.Sprintf("failed to parse ID: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *poolMembershipResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pool_membership"
}

func (r *poolMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan poolMembershipModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
