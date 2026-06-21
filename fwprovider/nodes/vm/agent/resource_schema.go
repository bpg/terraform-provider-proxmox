/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package agent

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ResourceSchema defines the schema for the QEMU guest agent block.
func ResourceSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		Description: "QEMU guest agent configuration.",
		MarkdownDescription: "Configure the QEMU guest agent. The agent enables the hypervisor to communicate " +
			"with the guest OS for operations like graceful shutdown, IP address retrieval, and file system freeze.",
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				Description: "Whether the QEMU guest agent is enabled.",
				Optional:    true,
			},
			"trim": schema.BoolAttribute{
				Description: "Whether to run fstrim after cloning or moving a disk.",
				Optional:    true,
			},
			"type": schema.StringAttribute{
				Description: "Guest agent channel type.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("isa", "virtio"),
				},
			},
		},
	}
}
