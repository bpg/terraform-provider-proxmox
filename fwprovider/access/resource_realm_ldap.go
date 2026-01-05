/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/attribute"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/config"
	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

var (
	_ resource.Resource                = (*realmLDAPResource)(nil)
	_ resource.ResourceWithConfigure   = (*realmLDAPResource)(nil)
	_ resource.ResourceWithImportState = (*realmLDAPResource)(nil)
)

type realmLDAPResource struct {
	client proxmox.Client
}

// NewRealmLDAPResource creates a new LDAP realm resource.
func NewRealmLDAPResource() resource.Resource {
	return &realmLDAPResource{}
}

func (r *realmLDAPResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_realm_ldap"
}

func (r *realmLDAPResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages an LDAP authentication realm in Proxmox VE.",
		MarkdownDescription: "Manages an LDAP authentication realm in Proxmox VE.\n\n" +
			"LDAP realms allow Proxmox to authenticate users against an LDAP directory service.",
		Attributes: map[string]schema.Attribute{
			"id": attribute.ResourceID("Realm identifier (same as realm)"),
			"realm": schema.StringAttribute{
				Description: "Realm identifier (e.g., 'example.com').",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(32),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[A-Za-z][A-Za-z0-9.\-_]+$`),
						"must be a valid realm identifier",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"server1": schema.StringAttribute{
				Description: "Primary LDAP server hostname or IP address.",
				Required:    true,
			},
			"server2": schema.StringAttribute{
				Description: "Fallback LDAP server hostname or IP address.",
				Optional:    true,
			},
			"base_dn": schema.StringAttribute{
				Description: "LDAP base DN for user searches (e.g., 'ou=users,dc=example,dc=com').",
				Required:    true,
			},
			"bind_dn": schema.StringAttribute{
				Description: "LDAP bind DN for authentication (e.g., 'cn=admin,dc=example,dc=com').",
				Optional:    true,
			},
			"bind_password": schema.StringAttribute{
				Description: "Password for the bind DN. Note: stored in Proxmox but not returned by API.",
				Optional:    true,
				Sensitive:   true,
			},
			"user_attr": schema.StringAttribute{
				Description: "LDAP attribute representing the username.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("uid"),
			},
			"port": schema.Int64Attribute{
				Description: "LDAP server port. Default: 389 (LDAP) or 636 (LDAPS).",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"secure": schema.BoolAttribute{
				Description:        "Use LDAPS (LDAP over SSL/TLS) instead of plain LDAP.",
				Optional:           true,
				Computed:           true,
				Default:            booldefault.StaticBool(false),
				DeprecationMessage: "Deprecated by Proxmox: use mode instead.",
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(path.MatchRoot("mode")),
				},
			},
			"verify": schema.BoolAttribute{
				Description: "Verify LDAP server SSL certificate.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"capath": schema.StringAttribute{
				Description: "Path to CA certificate file for SSL verification.",
				Optional:    true,
			},
			"cert": schema.StringAttribute{
				Description: "Path to client certificate for SSL authentication.",
				Optional:    true,
			},
			"certkey": schema.StringAttribute{
				Description: "Path to client certificate key.",
				Optional:    true,
			},
			"filter": schema.StringAttribute{
				Description: "LDAP filter for user searches.",
				Optional:    true,
			},
			"group_dn": schema.StringAttribute{
				Description: "LDAP base DN for group searches.",
				Optional:    true,
			},
			"group_filter": schema.StringAttribute{
				Description: "LDAP filter for group searches.",
				Optional:    true,
			},
			"group_classes": schema.StringAttribute{
				Description: "LDAP objectClasses for groups (comma-separated).",
				Optional:    true,
			},
			"group_name_attr": schema.StringAttribute{
				Description: "LDAP attribute representing the group name.",
				Optional:    true,
			},
			"mode": schema.StringAttribute{
				Description: "LDAP connection mode (ldap, ldaps, ldap+starttls).",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("ldap", "ldaps", "ldap+starttls"),
					stringvalidator.ConflictsWith(path.MatchRoot("secure")),
				},
			},
			"sslversion": schema.StringAttribute{
				Description: "SSL/TLS version (tlsv1, tlsv1_1, tlsv1_2, tlsv1_3).",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("tlsv1", "tlsv1_1", "tlsv1_2", "tlsv1_3"),
				},
			},
			"user_classes": schema.StringAttribute{
				Description: "LDAP objectClasses for users (comma-separated).",
				Optional:    true,
			},
			"sync_attributes": schema.StringAttribute{
				Description: "Comma-separated list of attributes to sync (e.g., 'email=mail,firstname=givenName').",
				Optional:    true,
			},
			"sync_defaults_options": schema.StringAttribute{
				Description: "Default sync options (e.g., 'scope=users,enable-new=1').",
				Optional:    true,
			},
			"comment": schema.StringAttribute{
				Description: "Description of the realm.",
				Optional:    true,
			},
			"default": schema.BoolAttribute{
				Description: "Use this realm as the default for login.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"case_sensitive": schema.BoolAttribute{
				Description: "Enable case-sensitive username matching.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

func (r *realmLDAPResource) Configure(
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

func (r *realmLDAPResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan realmLDAPModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert Terraform model to API request
	createReq := plan.toCreateRequest()

	err := r.client.Access().CreateRealm(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating LDAP realm", err.Error())
		return
	}

	// Read back the created resource
	plan.ID = plan.Realm

	r.read(ctx, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *realmLDAPResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state realmLDAPModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Verify the realm still exists. If not, remove from state.
	_, err := r.client.Access().GetRealm(ctx, state.Realm.ValueString())
	if err != nil {
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading LDAP Realm",
			fmt.Sprintf("Could not verify existence of realm %q: %v", state.Realm.ValueString(), err),
		)

		return
	}

	r.read(ctx, &state, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *realmLDAPResource) read(
	ctx context.Context,
	model *realmLDAPModel,
	diags *diag.Diagnostics,
) {
	realmData, err := r.client.Access().GetRealm(ctx, model.Realm.ValueString())
	if err != nil {
		diags.AddError("Error reading LDAP realm", err.Error())
		return
	}

	// Preserve the bind password from the plan/state since it's not returned by the API
	bindPassword := model.BindPassword

	model.fromAPIResponse(realmData)

	// Restore the bind password
	model.BindPassword = bindPassword
}

func (r *realmLDAPResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan realmLDAPModel
	var state realmLDAPModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := plan.toUpdateRequest(&state)

	err := r.client.Access().UpdateRealm(ctx, plan.Realm.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating LDAP realm", err.Error())
		return
	}

	r.read(ctx, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *realmLDAPResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state realmLDAPModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Access().DeleteRealm(ctx, state.Realm.ValueString())
	if err != nil {
		// If already deleted, that's fine
		if errors.Is(err, api.ErrResourceDoesNotExist) {
			return
		}
		resp.Diagnostics.AddError("Error deleting LDAP realm", err.Error())
		return
	}
}

func (r *realmLDAPResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("realm"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
