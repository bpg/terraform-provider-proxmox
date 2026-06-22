/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package acme

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/migration"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/acme/plugins"
)

var (
	_ resource.Resource                     = &acmePluginResource{}
	_ resource.ResourceWithConfigure        = &acmePluginResource{}
	_ resource.ResourceWithImportState      = &acmePluginResource{}
	_ resource.ResourceWithConfigValidators = &acmePluginResource{}
)

// NewACMEPluginResource creates a new resource for managing ACME plugins.
func NewACMEPluginResource() resource.Resource {
	return &acmePluginResource{}
}

// acmePluginResource contains the resource's internal data.
type acmePluginResource struct {
	// The ACME plugin API client
	client *plugins.Client
}

// Metadata defines the name of the resource.
func (r *acmePluginResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_acme_dns_plugin"
}

// Schema defines the schema for the resource.
func (r *acmePluginResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description:        "Manages an ACME plugin in a Proxmox VE cluster.",
		DeprecationMessage: migration.DeprecationMessage("proxmox_acme_dns_plugin"),
		Attributes: map[string]schema.Attribute{
			"api": schema.StringAttribute{
				Description: "API plugin name.",
				Required:    true,
			},
			"data": schema.MapAttribute{
				Description: "DNS plugin data.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"data_wo": schema.MapAttribute{
				Description: "DNS plugin data (write-only).",
				MarkdownDescription: "DNS plugin data, supplied as a [write-only argument](https://developer.hashicorp.com/terraform/language/resources/ephemeral/write-only) " +
					"so credentials are never stored in Terraform state. Requires Terraform 1.11+. Mutually exclusive with `data`. " +
					"Pair with `data_wo_version` to push rotated values.",
				Optional:    true,
				WriteOnly:   true,
				ElementType: types.StringType,
			},
			"data_wo_version": schema.Int64Attribute{
				Description: "Version counter for data_wo.",
				MarkdownDescription: "Version counter for `data_wo`. Because write-only values are not stored in state, Terraform cannot " +
					"detect when `data_wo` changes; increment this value to signal a rotation and force the new `data_wo` to be sent.",
				Optional: true,
				Validators: []validator.Int64{
					int64validator.AlsoRequires(path.MatchRoot("data_wo")),
				},
			},
			"digest": schema.StringAttribute{
				Description: "SHA1 digest of the current configuration.",
				MarkdownDescription: "SHA1 digest of the current configuration. " +
					"Prevent changes if current configuration file has a different digest. " +
					"This can be used to prevent concurrent modifications.",
				Optional: true,
				Computed: true,
			},
			"disable": schema.BoolAttribute{
				Description: "Flag to disable the config.",
				Optional:    true,
			},
			"plugin": schema.StringAttribute{
				Description: "ACME plugin ID name.",
				Required:    true,
			},
			"validation_delay": schema.Int64Attribute{
				Description: "Extra delay in seconds to wait before requesting validation.",
				MarkdownDescription: "Extra delay in seconds to wait before requesting validation. " +
					"Allows to cope with a long TTL of DNS records (0 - 172800).",
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(30),
				Validators: []validator.Int64{
					int64validator.Between(0, 172800),
				},
			},
		},
	}
}

// Configure accesses the provider-configured Proxmox API client on behalf of the resource.
func (r *acmePluginResource) Configure(
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

	r.client = cfg.Client.Cluster().ACME().Plugins()
}

// ConfigValidators enforces mutual exclusion between data and its write-only sibling.
func (r *acmePluginResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("data"),
			path.MatchRoot("data_wo"),
		),
	}
}

// Create creates a new ACME plugin on the Proxmox cluster.
func (r *acmePluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan acmePluginCreateModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Write-only data_wo is only present in config, never in plan.
	var dataWO types.Map

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("data_wo"), &dataWO)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &plugins.ACMEPluginsCreateRequestBody{}
	createRequest.Plugin = plan.Plugin.ValueString()
	createRequest.Type = "dns"
	createRequest.API = plan.API.ValueString()
	data := make(plugins.DNSPluginData)

	if !dataWO.IsNull() {
		dataWO.ElementsAs(ctx, &data, false)
	} else {
		plan.Data.ElementsAs(ctx, &data, false)
	}

	createRequest.Data = &data
	createRequest.Disable = plan.Disable.ValueBool()
	createRequest.ValidationDelay = plan.ValidationDelay.ValueInt64()

	err := r.client.Create(ctx, createRequest)
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unable to create ACME plugin '%s'", createRequest.Plugin),
				err.Error(),
			)

			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("ACME plugin '%s' already exists", createRequest.Plugin),
			err.Error(),
		)
	}

	plugin, err := r.client.Get(ctx, plan.Plugin.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read ACME plugin",
			err.Error(),
		)

		return
	}

	plan.Digest = types.StringValue(plugin.Digest)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read retrieves the current state of the ACME plugin from the Proxmox cluster.
func (r *acmePluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state acmePluginCreateModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Plugin.ValueString()

	plugin, err := r.client.Get(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read ACME plugin",
			err.Error(),
		)

		return
	}

	// A freshly imported resource has only `plugin` set; `digest` is always populated
	// by Create/Update, so a null digest means "imported, not yet reconciled" and we
	// must hydrate the `data` attribute from the API for import to work.
	imported := state.Digest.IsNull()

	state.API = types.StringValue(plugin.API)
	state.Digest = types.StringValue(plugin.Digest)
	state.ValidationDelay = types.Int64Value(plugin.ValidationDelay)

	// Mirror the API data back into the `data` attribute only when it is the source of
	// truth (the `data` path, or import). When credentials were supplied via the
	// write-only `data_wo`, `data` is null in state and must stay null; otherwise GET
	// (which returns the stored data) would create a perpetual diff against null config.
	if imported || !state.Data.IsNull() {
		mapValue, diags := types.MapValueFrom(ctx, types.StringType, plugin.Data)
		resp.Diagnostics.Append(diags...)

		state.Data = mapValue
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update modifies an existing ACME plugin on the Proxmox cluster.
func (r *acmePluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state acmePluginCreateModel

	toDelete := make([]string, 0)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Write-only data_wo is only present in config, never in plan/state.
	var dataWO types.Map

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("data_wo"), &dataWO)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &plugins.ACMEPluginsUpdateRequestBody{}
	updateRequest.API = plan.API.ValueString()

	data := make(plugins.DNSPluginData)

	switch {
	case !dataWO.IsNull():
		dataWO.ElementsAs(ctx, &data, false)
		updateRequest.Data = &data
	case !plan.Data.IsNull():
		plan.Data.ElementsAs(ctx, &data, false)
		updateRequest.Data = &data
	case !state.Data.IsNull():
		toDelete = append(toDelete, "data")
	}

	updateRequest.Digest = plan.Digest.ValueString()

	if plan.Disable.IsNull() && !state.Disable.IsNull() || !plan.Disable.ValueBool() {
		toDelete = append(toDelete, "disable")
	} else {
		updateRequest.Disable = plan.Disable.ValueBool()
	}

	if plan.ValidationDelay.IsNull() && !state.ValidationDelay.IsNull() {
		toDelete = append(toDelete, "validation_delay")
	} else {
		updateRequest.ValidationDelay = plan.ValidationDelay.ValueInt64()
	}

	if len(toDelete) > 0 {
		updateRequest.Delete = toDelete
	}

	err := r.client.Update(ctx, plan.Plugin.ValueString(), updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to update ACME plugin '%s'", plan.Plugin.ValueString()),
			err.Error(),
		)

		return
	}

	plugin, err := r.client.Get(ctx, plan.Plugin.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read ACME plugin",
			err.Error(),
		)

		return
	}

	plan.Digest = types.StringValue(plugin.Digest)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete removes an existing ACME plugin from the Proxmox cluster.
func (r *acmePluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state acmePluginCreateModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, state.Plugin.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to delete ACME plugin '%s'", state.Plugin.ValueString()),
			err.Error(),
		)
	}
}

// ImportState retrieves the current state of an existing ACME plugin from the Proxmox cluster.
func (r *acmePluginResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("plugin"), req, resp)
}

// Short-name alias for proxmox_acme_dns_plugin (ADR-007).

var (
	_ resource.Resource                = &acmePluginShort{}
	_ resource.ResourceWithConfigure   = &acmePluginShort{}
	_ resource.ResourceWithImportState = &acmePluginShort{}
	_ resource.ResourceWithMoveState   = &acmePluginShort{}
)

type acmePluginShort struct{ acmePluginResource }

// NewACMEPluginShortResource creates the short-name version of the ACME DNS plugin resource.
func NewACMEPluginShortResource() resource.Resource {
	return &acmePluginShort{}
}

func (r *acmePluginShort) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "proxmox_acme_dns_plugin"
}

func (r *acmePluginShort) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.acmePluginResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = ""
}

func (r *acmePluginShort) MoveState(ctx context.Context) []resource.StateMover {
	var schemaResp resource.SchemaResponse

	r.acmePluginResource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)

	return []resource.StateMover{
		migration.PrefixMoveState("proxmox_virtual_environment_acme_dns_plugin", &schemaResp.Schema),
	}
}
