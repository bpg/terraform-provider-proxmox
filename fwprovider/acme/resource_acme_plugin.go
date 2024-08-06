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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/acme/plugins"
)

var (
	_ resource.Resource                = &acmePluginResource{}
	_ resource.ResourceWithConfigure   = &acmePluginResource{}
	_ resource.ResourceWithImportState = &acmePluginResource{}
)

// NewACMEPluginResource creates a new resource for managing ACME plugins.
func NewACMEPluginResource() resource.Resource {
	return &acmePluginResource{}
}

// acmePluginResource contains the resource's internal data.
type acmePluginResource struct {
	// The ACME account API client
	client plugins.Client
}

// acmePluginCreateModel maps the schema data for an ACME plugin.
type acmePluginCreateModel struct {
	// API plugin name
	API types.String `tfsdk:"api"`
	// DNS plugin data
	Data types.Map `tfsdk:"data"`
	// A list of settings you want to delete.
	Delete types.String `tfsdk:"delete"`
	// Prevent changes if current configuration file has a different digest.
	// This can be used to prevent concurrent modifications.
	Digest types.String `tfsdk:"digest"`
	// Flag to disable the config
	Disable types.Bool `tfsdk:"disable"`
	// Plugin ID name
	Plugin types.String `tfsdk:"plugin"`
	// List of cluster node names
	Nodes types.String `tfsdk:"nodes"`
	// ACME challenge type (dns, standalone)
	Type types.String `tfsdk:"type"`
	// Extra delay in seconds to wait before requesting validation (0 - 172800)
	ValidationDelay types.Int64 `tfsdk:"validation_delay"`
}

// Metadata defines the name of the resource.
func (r *acmePluginResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_acme_plugin"
}

// Schema defines the schema for the resource.
func (r *acmePluginResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages an ACME plugin in a Proxmox VE cluster.",
		Attributes: map[string]schema.Attribute{
			"api": schema.StringAttribute{
				Description: "API plugin name.",
				Optional:    true,
			},
			"data": schema.MapAttribute{
				Description: "DNS plugin data.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"delete": schema.StringAttribute{
				Description: "A list of settings you want to delete.",
				Optional:    true,
			},
			"digest": schema.StringAttribute{
				Description: "Prevent changes if current configuration file has a different digest. " +
					"This can be used to prevent concurrent modifications.",
				Optional: true,
				Computed: true,
			},
			"disable": schema.BoolAttribute{
				Description: "Flag to disable the config.",
				Optional:    true,
			},
			"nodes": schema.StringAttribute{
				Description: "List of cluster node names.",
				Optional:    true,
			},
			"plugin": schema.StringAttribute{
				Description: "ACME Plugin ID name.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "ACME challenge type (dns, standalone).",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("dns", "standalone"),
				},
			},
			"validation_delay": schema.Int64Attribute{
				Description: "Extra delay in seconds to wait before requesting validation. " +
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

	client, ok := req.ProviderData.(proxmox.Client)
	if ok {
		r.client = *client.Cluster().ACME().Plugins()
	} else {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *proxmox.Client, got: %T",
				req.ProviderData),
		)
	}
}

// Create creates a new ACME plugin on the Proxmox cluster.
func (r *acmePluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan acmePluginCreateModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &plugins.ACMEPluginsCreateRequestBody{}
	createRequest.Plugin = plan.Plugin.ValueString()
	createRequest.Type = plan.Type.ValueString()
	createRequest.API = plan.API.ValueString()
	data := make(plugins.DNSPluginData)

	plan.Data.ElementsAs(ctx, &data, false)

	createRequest.Data = &data
	createRequest.Disable = plan.Disable.ValueBool()
	createRequest.Nodes = plan.Nodes.ValueString()
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
}

// Update modifies an existing ACME plugin on the Proxmox cluster.
func (r *acmePluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete removes an existing ACME plugin from the Proxmox cluster.
func (r *acmePluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

// ImportState retrieves the current state of an existing ACME plugin from the Proxmox cluster.
func (r *acmePluginResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
}
