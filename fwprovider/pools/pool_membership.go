package pools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/pools"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var (
	_ resource.Resource                     = (*poolMembershipResource)(nil)
	_ resource.ResourceWithConfigure        = (*poolMembershipResource)(nil)
	_ resource.ResourceWithImportState      = (*poolMembershipResource)(nil)
	_ resource.ResourceWithConfigValidators = (*poolMembershipResource)(nil)
)

const (
	mkPoolMembershipId        = "id"
	mkPoolMembershipType      = "type"
	mkPoolMembershipPoolId    = "pool_id"
	mkPoolMembershipVmId      = "vm_id"
	mkPoolMembershipStorageId = "storage_id"
)

type poolMembershipResource struct {
	client proxmox.Client
}

func (r *poolMembershipResource) ConfigValidators(context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot(mkPoolMembershipVmId),
			path.MatchRoot(mkPoolMembershipStorageId),
		),
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot(mkPoolMembershipVmId),
			path.MatchRoot(mkPoolMembershipStorageId),
		),
	}
}

func NewPoolMembershipResource() resource.Resource {
	return &poolMembershipResource{}
}

func (r *poolMembershipResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *poolMembershipResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages resource pool memberships for containers, virtual machines and storages",
		Attributes: map[string]schema.Attribute{
			mkPoolMembershipId: attribute.ResourceID(),
			mkPoolMembershipType: schema.StringAttribute{
				Description:         "Resource pool membership type",
				MarkdownDescription: "Resource pool membership type (can be `vm` for VMs and CTs or `storage` for storages)",
				Computed:            true,
			},
			mkPoolMembershipPoolId: schema.StringAttribute{
				Description: "Resource pool id",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			mkPoolMembershipVmId: schema.Int64Attribute{
				Description: "VM or CT id",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			mkPoolMembershipStorageId: schema.StringAttribute{
				Description: "Storage id",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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

	body := &pools.PoolUpdateRequestBody{
		AllowMove: ptr.Ptr(proxmoxtypes.CustomBool(true)),
	}

	if membershipType, err := plan.deduceMembershipType(); err != nil {
		resp.Diagnostics.AddError("Cannot determine pool membership type",
			"Plan does not have enough information to determine pool membership type. This is always an error in the provider.",
		)

		return
	} else {
		plan.Type = types.StringValue(membershipType)
	}

	switch plan.Type.ValueString() {
	case MembershipTypeStorage:
		storageList := (proxmoxtypes.CustomCommaSeparatedList)([]string{plan.StorageID.ValueString()})
		body.Storage = &storageList
	case MembershipTypeVm:
		vmList := (proxmoxtypes.CustomCommaSeparatedList)([]string{strconv.FormatInt(plan.VmID.ValueInt64(), 10)})
		body.VMs = &vmList
	default:
		resp.Diagnostics.AddError("Cannot create pool membership", ErrInvalidMembershipType.Error())
		return
	}

	if err := poolApi.UpdatePool(ctx, poolId, body); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to update resource pool '%s'", poolId),
			err.Error())

		return
	}

	if resourceID, resourceIDErr := plan.generateID(); resourceIDErr != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot create pool membership id for type '%s'", plan.Type.ValueString()),
			resourceIDErr.Error())

		return
	} else {
		plan.ID = resourceID
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *poolMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state poolMembershipModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	poolId := state.PoolId.ValueString()
	membershipType, membershipTypeErr := NewMembershipType(state.Type.ValueString())

	if membershipTypeErr != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Wrong pool membership type '%s' in state", state.Type.ValueString()), membershipTypeErr.Error())
		return
	}

	pool, err := r.client.Pool().GetPool(ctx, poolId)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to get pool '%s'", poolId), err.Error())
		return
	}

	exists := false

	switch membershipType {
	case MembershipTypeStorage:
		exists = checkStorageExists(*pool, state.StorageID.ValueString())
	case MembershipTypeVm:
		exists = checkVmExists(*pool, state.VmID.ValueInt64())
	default:
		resp.Diagnostics.AddError(fmt.Sprintf("Wrong pool membership type '%s' in state", state.Type.ValueString()), ErrInvalidMembershipType.Error())
		return
	}

	if !exists {
		resp.State.RemoveResource(ctx)
	}
}

func checkStorageExists(pool pools.PoolGetResponseData, storageId string) bool {
	for _, member := range pool.Members {
		if member.DatastoreID != nil && *member.DatastoreID == storageId {
			return true
		}
	}

	return false
}

func checkVmExists(pool pools.PoolGetResponseData, vmId int64) bool {
	for _, member := range pool.Members {
		if member.VMID != nil && int64(*member.VMID) == vmId {
			return true
		}
	}

	return false
}

func (r *poolMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state poolMembershipModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	poolId := state.PoolId.ValueString()
	membershipType, membershipTypeErr := NewMembershipType(state.Type.ValueString())

	if membershipTypeErr != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Wrong pool membership type '%s' in state", state.Type.ValueString()), membershipTypeErr.Error())
		return
	}

	body := &pools.PoolUpdateRequestBody{
		Delete: ptr.Ptr(proxmoxtypes.CustomBool(true)),
	}

	switch membershipType {
	case MembershipTypeStorage:
		storageList := (proxmoxtypes.CustomCommaSeparatedList)([]string{state.StorageID.ValueString()})
		body.Storage = &storageList
	case MembershipTypeVm:
		vmList := (proxmoxtypes.CustomCommaSeparatedList)([]string{strconv.FormatInt(state.VmID.ValueInt64(), 10)})
		body.VMs = &vmList
	default:
		resp.Diagnostics.AddError("Cannot create pool membership", ErrInvalidMembershipType.Error())
		return
	}

	if err := r.client.Pool().UpdatePool(ctx, poolId, body); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to update pool '%s", poolId), err.Error())
	}
}

func (r *poolMembershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	model, err := createMembershipModelFromID(req.ID)
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
