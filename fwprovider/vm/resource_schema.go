package vm

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/tags"
)

// Schema defines the schema for the resource.
func (r *vmResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "This is an experimental implementation of a Proxmox VM resource using Plugin Framework.",
		MarkdownDescription: "This is an experimental implementation of a Proxmox VM resource using Plugin Framework." +
			"<br><br>It is a Proof of Concept, highly experimental and **will** change in future. " +
			"It does not support all features of the Proxmox API for VMs and **MUST NOT** be used in production.",
		Attributes: map[string]schema.Attribute{
			"clone": schema.SingleNestedAttribute{
				Description: "The cloning configuration.",
				Optional:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "The ID of the VM to clone.",
						Required:    true,
					},
					"retries": schema.Int64Attribute{
						Description: "The number of retries to perform when cloning the VM (default: 3).",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(3),
					},
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the VM.",
				Optional:    true,
			},
			"id": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
				Description: "The unique identifier of the VM in the Proxmox cluster.",
			},
			"name": schema.StringAttribute{
				Description:         "The name of the VM.",
				MarkdownDescription: "The name of the VM. Doesn't have to be unique.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^([a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9])$`),
						"must be a valid DNS name",
					),
				},
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the node where the VM is provisioned.",
				Required:    true,
			},
			"tags": tags.ResourceAttribute(),
			"template": schema.BoolAttribute{
				Description: "Set to true to create a VM template.",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
		},
	}
}

//
//// useStateForUnknownModifier implements the plan modifier.
//type forceNullModifier struct{}
//
//// Description returns a human-readable description of the plan modifier.
//func (m forceNullModifier) Description(_ context.Context) string {
//	return "Forces null value."
//}
//
//// MarkdownDescription returns a markdown description of the plan modifier.
//func (m forceNullModifier) MarkdownDescription(_ context.Context) string {
//	return "Forces null value."
//}
//
//// PlanModifyBool implements the plan modification logic.
//func (m forceNullModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
//	//if req.StateValue.IsNull() {
//	//	return
//	//}
//	//
//	//if !req.PlanValue.IsUnknown() {
//	//	return
//	//}
//
//	// forceNullIfUnconfiguredComputed
//
//	path, diagnostics := req.Plan.Schema.AttributeAtPath(ctx, req.Path)
//	path.IsComputed()
//	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
//		resp.PlanValue = types.StringNull()
//	}
//}
