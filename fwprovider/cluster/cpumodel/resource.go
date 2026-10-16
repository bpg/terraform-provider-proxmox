/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cpumodel

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/cpumodels"
)

var (
	_ resource.Resource                = &customCPUModelResource{}
	_ resource.ResourceWithConfigure   = &customCPUModelResource{}
	_ resource.ResourceWithImportState = &customCPUModelResource{}
)

type customCPUModelResource struct {
	client *cpumodels.Client
}

// NewResource creates a new custom CPU model resource.
func NewResource() resource.Resource {
	return &customCPUModelResource{}
}

func (r *customCPUModelResource) Metadata(
	_ context.Context,
	_ resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "proxmox_custom_cpu_model"
}

func (r *customCPUModelResource) Configure(
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

	r.client = cfg.Client.Cluster().CustomCPUModels()
}

func (r *customCPUModelResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages a custom CPU model, a cluster-wide CPU type that VMs can use. " +
			"Reference it from a VM's `cpu.type` through the `cputype` attribute.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The name of the custom CPU model, without the `custom-` prefix.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(40),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_-]+$`),
						"must start with a letter and contain only letters, digits, `-` and `_`",
					),
					noCustomPrefix{},
				},
			},
			"cputype": schema.StringAttribute{
				Description: "The CPU type to use in a VM's `cpu.type`: the name with the `custom-` prefix.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"reported_model": schema.StringAttribute{
				Description: "The CPU model and vendor reported to the guest. Must be a model supported by QEMU/KVM, " +
					"e.g. `EPYC-Milan-v3`.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"flags": schema.ListAttribute{
				Description: "CPU flags to enable (`+FLAG`) or disable (`-FLAG`), e.g. `[\"-pcid\", \"+aes\"]`.",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^[+-][a-zA-Z0-9_.-]+$`),
							"must be `+` or `-` followed by a flag name",
						),
					),
				},
			},
			"phys_bits": schema.StringAttribute{
				Description: "The number of physical address bits reported to the guest (8-64), or `host` to use the host CPU's value.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(host|[89]|[1-5][0-9]|6[0-4])$`),
						"must be a number from 8 to 64, or `host`",
					),
				},
			},
			"guest_phys_bits": schema.Int64Attribute{
				Description: "The number of physical address bits available to the guest (32-64).",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(32, 64),
				},
			},
			"hidden": schema.BoolAttribute{
				Description: "Whether to hide that the CPU is a KVM virtual machine from the guest. " +
					"Only affects vCPUs with x86-64 architecture.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"hv_vendor_id": schema.StringAttribute{
				Description: "The Hyper-V vendor ID reported to the guest (1-12 alphanumeric characters).",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9]{1,12}$`),
						"must be 1 to 12 alphanumeric characters",
					),
				},
			},
			"level": schema.Int64Attribute{
				Description: "The maximum input value for the basic CPUID leaves the guest can query. " +
					"Only applies when the vCPU architecture is x86-64.",
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 4294967295),
				},
			},
		},
	}
}

func (r *customCPUModelResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan customCPUModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := plan.toAPICreate(ctx, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.ID.ValueString()

	if err := r.client.Create(ctx, body); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Create Custom CPU Model %q", name), err.Error())
		return
	}

	data, err := r.client.Get(ctx, name)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Read Custom CPU Model %q After Creation", name), err.Error())
		return
	}

	resp.Diagnostics.Append(plan.fromAPI(data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *customCPUModelResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state customCPUModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	name := state.ID.ValueString()

	data, err := r.client.Get(ctx, name)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Read Custom CPU Model %q", name), err.Error())

		return
	}

	resp.Diagnostics.Append(state.fromAPI(data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *customCPUModelResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state customCPUModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	body := plan.toAPIUpdate(ctx, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.ID.ValueString()

	if err := r.client.Update(ctx, name, body); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Update Custom CPU Model %q", name), err.Error())
		return
	}

	data, err := r.client.Get(ctx, name)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Read Custom CPU Model %q After Update", name), err.Error())
		return
	}

	resp.Diagnostics.Append(plan.fromAPI(data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *customCPUModelResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state customCPUModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	name := state.ID.ValueString()

	err := r.client.Delete(ctx, name)
	if err != nil && !errors.Is(err, api.ErrResourceDoesNotExist) {
		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Delete Custom CPU Model %q", name), err.Error())
	}
}

func (r *customCPUModelResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Accept the name with or without the prefix; state always holds it without.
	name := strings.TrimPrefix(req.ID, customPrefix)

	data, err := r.client.Get(ctx, name)
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.Diagnostics.AddError(
				"Custom CPU Model Not Found",
				fmt.Sprintf("Custom CPU model %q was not found", name),
			)

			return
		}

		resp.Diagnostics.AddError(fmt.Sprintf("Unable to Import Custom CPU Model %q", name), err.Error())

		return
	}

	var model customCPUModel

	resp.Diagnostics.Append(model.fromAPI(data)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// noCustomPrefix rejects names that start with `custom-`. PVE strips that prefix on input and always adds it on
// output, so keeping it out of `id` keeps configuration, state and import in one canonical form.
type noCustomPrefix struct{}

func (v noCustomPrefix) Description(_ context.Context) string {
	return fmt.Sprintf("must not start with %q", customPrefix)
}

func (v noCustomPrefix) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v noCustomPrefix) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if strings.HasPrefix(req.ConfigValue.ValueString(), customPrefix) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Custom CPU Model Name",
			fmt.Sprintf("Give the name without the %q prefix; PVE adds it, and `cputype` exposes the prefixed name.", customPrefix),
		)
	}
}
