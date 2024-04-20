package vm

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
			"tags": schema.SetAttribute{
				Description: "The tags assigned to the VM.",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(regexp.MustCompile(`(.|\s)*\S(.|\s)*`), "must be a non-empty and non-whitespace string"),
						stringvalidator.LengthAtLeast(1),
					),
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
