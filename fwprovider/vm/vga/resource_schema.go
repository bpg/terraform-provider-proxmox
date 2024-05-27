package vga

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ResourceSchema defines the schema for the CPU resource.
func ResourceSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		CustomType: basetypes.ObjectType{
			AttrTypes: attributeTypes(),
		},
		Description: "The VGA configuration.",
		MarkdownDescription: "Configure the VGA Hardware. If you want to use high resolution modes (>= 1280x1024x16) " +
			"you may need to increase the vga memory option. Since QEMU 2.9 the default VGA display type is `std` " +
			"for all OS types besides some Windows versions (XP and older) which use `cirrus`. The `qxl` option " +
			"enables the SPICE display server. For win* OS you can select how many independent displays you want, " +
			"Linux guests can add displays themself. You can also run without any graphic card, using a serial device " +
			"as terminal. See the [Proxmox documentation](https://pve.proxmox.com/pve-docs/pve-admin-guide.html#" +
			"qm_virtual_machines_settings) section 10.2.8 for more information and available configuration parameters.",
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
		Attributes: map[string]schema.Attribute{
			"clipboard": schema.StringAttribute{
				Description: "Enable a specific clipboard.",
				MarkdownDescription: "Enable a specific clipboard. If not set, depending on the display type the SPICE " +
					"one will be added. Currently only `vnc` is available. Migration with VNC clipboard is not " +
					"supported by Proxmox.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("vnc"),
				},
			},
			"type": schema.StringAttribute{
				Description:         "The VGA type.",
				MarkdownDescription: "The VGA type (defaults to `std`).",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"cirrus",
						"none",
						"qxl",
						"qxl2",
						"qxl3",
						"qxl4",
						"serial0",
						"serial1",
						"serial2",
						"serial3",
						"std",
						"virtio",
						"virtio-gl",
						"vmware",
					),
				},
			},
			"memory": schema.Int64Attribute{
				Description:         "The VGA memory in megabytes (4-512 MB)",
				MarkdownDescription: "The VGA memory in megabytes (4-512 MB). Has no effect with serial display. ",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(4, 512),
				},
			},
		},
	}
}
