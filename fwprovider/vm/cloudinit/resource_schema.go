/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cloudinit

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	customtypes "github.com/bpg/terraform-provider-proxmox/fwprovider/types"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
)

// ResourceSchema defines the schema for the CPU resource.
func ResourceSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		Description: "The cloud-init configuration.",
		Optional:    true,
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"datastore_id": schema.StringAttribute{
				Description: "The identifier for the datastore to create the cloud-init disk in (defaults to `local-lvm`)",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("local-lvm"),
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"interface": schema.StringAttribute{
				Description: "The hardware interface to connect the cloud-init image to.",
				MarkdownDescription: "The hardware interface to connect the cloud-init image to. " +
					"Must be one of `ideN`, `sataN`, `scsiN`, where N is the index of the interface. " +
					"Will be detected if the setting is missing but a cloud-init image is present, " +
					"otherwise defaults to `ide2`. Note that `q35` machine type only supports " +
					"`ide0` and `ide2` of IDE interfaces.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("ide2"),
				Validators: []validator.String{
					validators.CDROMInterface(),
				},
			},
			"dns": schema.SingleNestedAttribute{
				Description: "The DNS configuration.",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"domain": schema.StringAttribute{
						Description: "The domain name to use for the VM.",
						Optional:    true,
						Computed:    true,
					},
					"servers": schema.ListAttribute{
						Description: "The list of DNS servers to use.",
						ElementType: customtypes.IPAddrType{},
						Optional:    true,
						Computed:    true,
					},
				},
			},
		},
	}
}
