/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package clone

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// ResourceSchema defines the schema for the clone block on a VM resource.
// All fields are RequiresReplace: any change to clone config destroys and recreates the VM.
func ResourceSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		Description: "Clone an existing VM as the base for this VM. All fields cause replacement when changed.",
		Optional:    true,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplace(),
		},
		Attributes: map[string]schema.Attribute{
			"vm_id": schema.Int64Attribute{
				Description: "The VM ID of the source VM to clone.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"node_name": schema.StringAttribute{
				Description: "The node where the source VM resides. Defaults to the target node if not specified.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"full": schema.BoolAttribute{
				Description: "Whether to create a full clone rather than a linked clone. Defaults to `true`.",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"datastore_id": schema.StringAttribute{
				Description: "The storage location for the clone's disks.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"retries": schema.Int64Attribute{
				Description: "The number of retries if the clone operation fails due to transient errors.",
				Optional:    true,
			},
		},
	}
}
