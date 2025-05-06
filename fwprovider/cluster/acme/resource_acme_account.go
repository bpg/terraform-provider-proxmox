/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package acme

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/acme/account"
)

var (
	_ resource.Resource                = &acmeAccountResource{}
	_ resource.ResourceWithConfigure   = &acmeAccountResource{}
	_ resource.ResourceWithImportState = &acmeAccountResource{}
)

// NewACMEAccountResource creates a new resource for managing ACME accounts.
func NewACMEAccountResource() resource.Resource {
	return &acmeAccountResource{}
}

// acmeAccountResource contains the resource's internal data.
type acmeAccountResource struct {
	// The ACME account API client
	client *account.Client
}

// acmeAccountModel maps the schema data for the ACME account resource.
type acmeAccountModel struct {
	// Contact email addresses.
	Contact types.String `tfsdk:"contact"`
	// CreatedAt timestamp of the account creation.
	CreatedAt types.String `tfsdk:"created_at"`
	// URL of ACME CA directory endpoint.
	Directory types.String `tfsdk:"directory"`
	// HMAC key for External Account Binding.
	EABHMACKey types.String `tfsdk:"eab_hmac_key"`
	// Key Identifier for External Account Binding.
	EABKID types.String `tfsdk:"eab_kid"`
	// Location of the ACME account.
	Location types.String `tfsdk:"location"`
	// ACME account config file name.
	Name types.String `tfsdk:"name"`
	// URL of CA TermsOfService - setting this indicates agreement.
	TOS types.String `tfsdk:"tos"`
}

// Metadata defines the name of the resource.
func (r *acmeAccountResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_acme_account"
}

// Schema defines the schema for the resource.
func (r *acmeAccountResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages an ACME account in a Proxmox VE cluster.",
		MarkdownDescription: "Manages an ACME account in a Proxmox VE cluster.\n\n" +
			"~> This resource requires `root@pam` authentication.",
		Attributes: map[string]schema.Attribute{
			"contact": schema.StringAttribute{
				Description: "The contact email addresses.",
				Required:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The timestamp of the ACME account creation.",
				Computed:    true,
			},
			"directory": schema.StringAttribute{
				Description: "The URL of the ACME CA directory endpoint.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^https?://.*$`),
						"must be a valid URL",
					),
				},
				Optional: true,
			},
			"eab_hmac_key": schema.StringAttribute{
				Description: "The HMAC key for External Account Binding.",
				Optional:    true,
			},
			"eab_kid": schema.StringAttribute{
				Description: "The Key Identifier for External Account Binding.",
				Optional:    true,
			},
			"location": schema.StringAttribute{
				Description: "The location of the ACME account.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The ACME account config file name.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("default"),
			},
			"tos": schema.StringAttribute{
				Description: "The URL of CA TermsOfService - setting this indicates agreement.",
				Optional:    true,
			},
		},
	}
}

// Configure accesses the provider-configured Proxmox API client on behalf of the resource.
func (r *acmeAccountResource) Configure(
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

	r.client = cfg.Client.Cluster().ACME().Account()
}

// Create creates a new ACME account on the Proxmox cluster.
func (r *acmeAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan acmeAccountModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &account.ACMEAccountCreateRequestBody{}
	createRequest.Contact = plan.Contact.ValueString()
	createRequest.Directory = plan.Directory.ValueString()
	createRequest.EABHMACKey = plan.EABHMACKey.ValueString()
	createRequest.EABKID = plan.EABKID.ValueString()
	createRequest.Name = plan.Name.ValueString()
	createRequest.TOS = plan.TOS.ValueString()

	err := r.client.Create(ctx, createRequest)
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unable to create ACME account '%s'", plan.Name),
				err.Error(),
			)

			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("ACME account '%s' already exists", plan.Name),
			err.Error(),
		)
	}

	r.readBack(ctx, &plan, &resp.Diagnostics, &resp.State)
}

// Read retrieves the current state of the ACME account from the Proxmox cluster.
func (r *acmeAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state acmeAccountModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	found, diags := r.read(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if !resp.Diagnostics.HasError() {
		if found {
			resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
		} else {
			resp.State.RemoveResource(ctx)
		}
	}
}

// Update modifies an existing ACME account on the Proxmox cluster.
func (r *acmeAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan acmeAccountModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := &account.ACMEAccountUpdateRequestBody{}
	updateRequest.Contact = plan.Contact.ValueString()

	err := r.client.Update(ctx, plan.Name.ValueString(), updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to update ACME account '%s'", plan.Name),
			err.Error(),
		)

		return
	}

	r.readBack(ctx, &plan, &resp.Diagnostics, &resp.State)
}

// Delete removes an existing ACME account from the Proxmox cluster.
func (r *acmeAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state acmeAccountModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to delete ACME account '%s'", state.Name),
			err.Error(),
		)

		return
	}
}

// ImportState retrieves the current state of an existing ACME account from the Proxmox cluster.
func (r *acmeAccountResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func (r *acmeAccountResource) readBack(
	ctx context.Context,
	data *acmeAccountModel,
	respDiags *diag.Diagnostics,
	respState *tfsdk.State,
) {
	found, diags := r.read(ctx, data)

	respDiags.Append(diags...)

	if !found {
		respDiags.AddError(
			fmt.Sprintf("ACME account '%s' not found after update", data.Name),
			"Failed to find ACME account when trying to read back the updated ACME account's data.",
		)
	}

	if !respDiags.HasError() {
		respDiags.Append(respState.Set(ctx, data)...)
	}
}

func (r *acmeAccountResource) read(ctx context.Context, data *acmeAccountModel) (bool, diag.Diagnostics) {
	name := data.Name.ValueString()

	acc, err := r.client.Get(ctx, name)
	if err != nil {
		diags := diag.Diagnostics{}

		if !strings.Contains(err.Error(), "does not exist") {
			diags.AddError(
				fmt.Sprintf("Unable to read ACME account '%s'", name),
				err.Error(),
			)
		}

		return false, diags
	}

	var contact string
	if len(acc.Account.Contact) > 0 {
		contact = strings.Replace(acc.Account.Contact[0], "mailto:", "", 1)
	}

	data.Directory = types.StringValue(acc.Directory)
	data.TOS = types.StringValue(acc.TOS)
	data.Location = types.StringValue(acc.Location)
	data.Contact = types.StringValue(contact)
	data.CreatedAt = types.StringValue(acc.Account.CreatedAt)

	return true, nil
}
