/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vm

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/cdrom"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/cpu"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/rng"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/vga"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
)

// Schema defines the schema for the resource.
func (r *Resource) Schema(
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
			"cdrom": cdrom.ResourceSchema(),
			"cpu":   cpu.ResourceSchema(),
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
			"rng": rng.ResourceSchema(),
			"stop_on_destroy": schema.BoolAttribute{
				Description:         "Set to true to stop (rather than shutdown) the VM on destroy.",
				MarkdownDescription: "Set to true to stop (rather than shutdown) the VM on destroy (defaults to `false`).",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"tags": stringset.ResourceAttribute("The tags assigned to the VM.", ""),
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
			"vga": vga.ResourceSchema(),
		},
	}
}
