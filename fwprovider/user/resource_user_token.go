package user

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/access"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

var (
	_ resource.Resource                = &userTokenResource{}
	_ resource.ResourceWithConfigure   = &userTokenResource{}
	_ resource.ResourceWithImportState = &userTokenResource{}
)

type userTokenResource struct {
	client proxmox.Client
}

type userTokenModel struct {
	Comment        types.String `tfsdk:"comment"`
	ExpirationDate types.String `tfsdk:"expiration_date"`
	ID             types.String `tfsdk:"id"`
	PrivSeparation types.Bool   `tfsdk:"privileges_separation"`
	UserID         types.String `tfsdk:"user_id"`
	Value          types.String `tfsdk:"value"`
}

// NewUserTokenResource creates a new user token resource.
func NewUserTokenResource() resource.Resource {
	return &userTokenResource{}
}

func (r *userTokenResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "User API tokens.",
		Attributes: map[string]schema.Attribute{
			"comment": schema.StringAttribute{
				Description: "Comment for the token.",
				Optional:    true,
			},
			"expiration_date": schema.StringAttribute{
				Description: "Expiration date for the token.",
				Optional:    true,
				// TODO: add validator
			},
			"id": schema.StringAttribute{
				Description: "User-specific token identifier.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`[A-Za-z][A-Za-z0-9.\-_]+`), "must be a valid token identifier"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"privileges_separation": schema.BoolAttribute{
				Description: "Restrict API token privileges with separate ACLs (default)",
				MarkdownDescription: "Restrict API token privileges with separate ACLs (default), " +
					"Restrict API token privileges with separate ACLs (default), " +
					"or give full privileges of corresponding user.",
				Optional: true,
			},
			"user_id": schema.StringAttribute{
				Description: "User identifier.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value": schema.StringAttribute{
				Description: "API token value used for authentication.",
				Optional:    false,
				Computed:    true,
			},
		},
	}
}

func (r *userTokenResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(proxmox.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *proxmox.Client, got: %T", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *userTokenResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_user_token"
}

func (r *userTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan userTokenModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	body := access.UserTokenCreateRequestBody{
		Comment:        plan.Comment.ValueStringPointer(),
		ExpirationDate: nil,
		PrivSeparate:   proxmoxtypes.CustomBoolPtr(plan.PrivSeparation.ValueBoolPointer()),
	}

	value, err := r.client.Access().CreateUserToken(ctx, plan.UserID.ValueString(), plan.ID.ValueString(), &body)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user token", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	plan.Value = types.StringValue(value)
	resp.State.Set(ctx, plan)
}

func (r *userTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

func (r *userTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (r *userTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (r *userTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

}
