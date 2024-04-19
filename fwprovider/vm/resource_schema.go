package vm

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
			"id": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
				Description: "The unique identifier of the VM in the Proxmox cluster.",
			},

			"description": schema.StringAttribute{
				Description: "The description of the VM.",
				Optional:    true,
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
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
		},
	}
}
