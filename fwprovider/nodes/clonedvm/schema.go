/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package clonedvm

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/cdrom"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/cpu"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/memory"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/rng"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/nodes/vm/vga"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
)

// Schema defines the schema for the cloned VM resource.
func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Clone a VM from a source template/VM and manage only explicitly-defined configuration. " +
			"This resource uses explicit opt-in management: only configuration blocks and devices explicitly " +
			"listed in your Terraform code are managed. Inherited settings from the template are preserved " +
			"unless explicitly overridden or deleted. Removing a configuration from Terraform stops managing " +
			"it but does not delete it from the VM.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
				Description: "The VM identifier in the Proxmox cluster.",
			},
			"node_name": schema.StringAttribute{
				Description: "Target node for the cloned VM.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Description: "Optional VM name override applied after cloning.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "Optional VM description applied after cloning.",
				Optional:    true,
			},
			"tags": schema.SetAttribute{
				CustomType: stringset.Type{
					SetType: types.SetType{
						ElemType: types.StringType,
					},
				},
				Description: "Tags applied after cloning.",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(`(.|\s)*\S(.|\s)*`),
							"must be a non-empty and non-whitespace string",
						),
						stringvalidator.LengthAtLeast(1),
					),
				},
			},
			"cpu":    optInManagedAttribute(cpu.ResourceSchema()),
			"memory": optInManagedAttribute(memory.ResourceSchema()),
			"rng":    optInManagedAttribute(rng.ResourceSchema()),
			"vga":    optInManagedAttribute(vga.ResourceSchema()),
			"cdrom":  optInManagedAttribute(cdrom.ResourceSchema()),
			"stop_on_destroy": schema.BoolAttribute{
				Description: "Stop the VM on destroy (instead of shutdown).",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"purge_on_destroy": schema.BoolAttribute{
				Description: "Purge backup configuration on destroy.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"delete_unreferenced_disks_on_destroy": schema.BoolAttribute{
				Description: "Delete unreferenced disks on destroy. WARNING: When set to true, any disks not " +
					"explicitly managed by Terraform will be deleted on destroy, potentially causing data loss. " +
					"Defaults to false for safety.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"clone":   cloneAttribute(),
			"network": networkAttribute(),
			"disk":    diskAttribute(),
			"delete":  deleteAttribute(),
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
		},
	}
}

func cloneAttribute() schema.Attribute {
	return schema.SingleNestedAttribute{
		Description: "Clone settings. Changes require recreation.",
		Required:    true,
		Attributes: map[string]schema.Attribute{
			"source_vm_id": schema.Int64Attribute{
				Description: "Source VM/template ID to clone from.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"source_node_name": schema.StringAttribute{
				Description: "Source node of the VM/template. Defaults to target node if unset.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"full": schema.BoolAttribute{
				Description: "Perform a full clone (true) or linked clone (false).",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"target_datastore": schema.StringAttribute{
				Description: "Target datastore for cloned disks.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_format": schema.StringAttribute{
				Description: "Target disk format for clone (e.g., raw, qcow2).",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"snapshot_name": schema.StringAttribute{
				Description: "Snapshot name to clone from.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"pool_id": schema.StringAttribute{
				Description: "Pool to assign the cloned VM to.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"retries": schema.Int64Attribute{
				Description: "Number of retries for clone operations.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(3),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
					int64validator.AtMost(10),
				},
			},
			"bandwidth_limit": schema.Int64Attribute{
				Description: "Clone bandwidth limit in MB/s.",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func networkAttribute() schema.Attribute {
	return schema.MapNestedAttribute{
		Description: "Network devices keyed by slot (net0, net1, ...). Only listed keys are managed.",
		Optional:    true,
		Validators: []validator.Map{
			mapvalidator.KeysAre(
				stringvalidator.RegexMatches(regexp.MustCompile(`^net[0-9]+$`), "must be a net interface key (net0, net1, ...)"),
			),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"bridge": schema.StringAttribute{
					Description: "Bridge name.",
					Optional:    true,
				},
				"model": schema.StringAttribute{
					Description: "NIC model (e.g., virtio, e1000).",
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
				"firewall": schema.BoolAttribute{
					Description: "Enable firewall on this interface.",
					Optional:    true,
				},
				"link_down": schema.BoolAttribute{
					Description: "Keep link down.",
					Optional:    true,
				},
				"mac_address": schema.StringAttribute{
					Description: "MAC address (computed if omitted).",
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^(?i:[0-9a-f]{2}(:[0-9a-f]{2}){5})$`),
							"must be a valid MAC address",
						),
					},
				},
				"mtu": schema.Int64Attribute{
					Description: "Interface MTU.",
					Optional:    true,
					Validators: []validator.Int64{
						int64validator.Between(576, 9216),
					},
				},
				"queues": schema.Int64Attribute{
					Description: "Number of multiqueue NIC queues.",
					Optional:    true,
					Validators: []validator.Int64{
						int64validator.Between(0, 16),
					},
				},
				"rate_limit": schema.Float64Attribute{
					Description: "Rate limit (MB/s).",
					Optional:    true,
				},
				"tag": schema.Int64Attribute{
					Description: "VLAN tag.",
					Optional:    true,
					Validators: []validator.Int64{
						int64validator.Between(1, 4094),
					},
				},
				"trunks": schema.SetAttribute{
					Description: "Trunk VLAN IDs.",
					Optional:    true,
					ElementType: types.Int64Type,
					Validators: []validator.Set{
						setvalidator.ValueInt64sAre(
							int64validator.Between(1, 4094),
						),
					},
				},
			},
		},
	}
}

func diskAttribute() schema.Attribute {
	return schema.MapNestedAttribute{
		Description: "Disks keyed by slot (scsi0, virtio0, sata0, ide0, ...). Only listed keys are managed.",
		Optional:    true,
		Validators: []validator.Map{
			mapvalidator.KeysAre(
				stringvalidator.RegexMatches(
					regexp.MustCompile(`^(ide[0-3]|sata[0-5]|scsi[0-9]+|virtio[0-9]+)$`),
					"must be a disk slot like scsi0, virtio0, sata0, ide0",
				),
			),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"file": schema.StringAttribute{
					Description: "Existing volume reference (e.g., local-lvm:vm-100-disk-0).",
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
				"datastore_id": schema.StringAttribute{
					Description: "Target datastore for new disks when file is not provided.",
					Optional:    true,
				},
				"size_gb": schema.Int64Attribute{
					Description: "Disk size (GiB) when creating new disks.",
					Optional:    true,
					Validators: []validator.Int64{
						int64validator.Between(1, 10240),
					},
				},
				"format": schema.StringAttribute{
					Description: "Disk format (raw, qcow2, vmdk).",
					Optional:    true,
				},
				"aio": schema.StringAttribute{
					Description: "AIO mode (io_uring, native, threads).",
					Optional:    true,
				},
				"backup": schema.BoolAttribute{
					Description: "Include disk in backups.",
					Optional:    true,
				},
				"discard": schema.StringAttribute{
					Description: "Discard/trim behavior.",
					Optional:    true,
				},
				"cache": schema.StringAttribute{
					Description: "Cache mode.",
					Optional:    true,
				},
				"iothread": schema.BoolAttribute{
					Description: "Use IO thread.",
					Optional:    true,
				},
				"replicate": schema.BoolAttribute{
					Description: "Consider disk for replication.",
					Optional:    true,
				},
				"serial": schema.StringAttribute{
					Description: "Disk serial number.",
					Optional:    true,
				},
				"ssd": schema.BoolAttribute{
					Description: "Mark disk as SSD.",
					Optional:    true,
				},
				"import_from": schema.StringAttribute{
					Description: "Import source volume/file id.",
					Optional:    true,
				},
				"media": schema.StringAttribute{
					Description: "Disk media (e.g., disk, cdrom).",
					Optional:    true,
				},
			},
		},
	}
}

func deleteAttribute() schema.Attribute {
	return schema.SingleNestedAttribute{
		Description: "Explicit deletions to perform after cloning/updating. Entries persist across applies.",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			"network": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "Network slots to delete (e.g., net1).",
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^net[0-9]+$`),
							"must be a net interface key (net0, net1, ...)",
						),
					),
				},
			},
			"disk": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "Disk slots to delete (e.g., scsi2).",
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^(ide[0-3]|sata[0-5]|scsi[0-9]+|virtio[0-9]+)$`),
							"must be a disk slot like scsi0, virtio0, sata0, ide0",
						),
					),
				},
			},
		},
	}
}

func optInManagedAttribute(attr schema.Attribute) schema.Attribute {
	switch v := attr.(type) {
	case schema.SingleNestedAttribute:
		v.Optional = true
		v.Computed = false
		v.PlanModifiers = nil
		v.Attributes = optInManagedAttributes(v.Attributes)

		return v
	case schema.MapNestedAttribute:
		v.Optional = true
		v.Computed = false
		v.PlanModifiers = nil
		v.NestedObject.Attributes = optInManagedAttributes(v.NestedObject.Attributes)

		return v
	default:
		return attr
	}
}

func optInManagedAttributes(in map[string]schema.Attribute) map[string]schema.Attribute {
	if len(in) == 0 {
		return in
	}

	out := make(map[string]schema.Attribute, len(in))
	for k, v := range in {
		out[k] = optInManagedAttributeAny(v)
	}

	return out
}

func optInManagedAttributeAny(attr schema.Attribute) schema.Attribute {
	switch v := attr.(type) {
	case schema.BoolAttribute:
		v.Optional = true
		v.Computed = false
		v.Default = nil
		v.PlanModifiers = nil

		return v
	case schema.Float64Attribute:
		v.Optional = true
		v.Computed = false
		v.Default = nil
		v.PlanModifiers = nil

		return v
	case schema.Int64Attribute:
		v.Optional = true
		v.Computed = false
		v.Default = nil
		v.PlanModifiers = nil

		return v
	case schema.ListAttribute:
		v.Optional = true
		v.Computed = false
		v.PlanModifiers = nil

		return v
	case schema.MapNestedAttribute:
		v.Optional = true
		v.Computed = false
		v.PlanModifiers = nil
		v.NestedObject.Attributes = optInManagedAttributes(v.NestedObject.Attributes)

		return v
	case schema.SetAttribute:
		v.Optional = true
		v.Computed = false
		v.Default = nil
		v.PlanModifiers = nil

		return v
	case schema.SingleNestedAttribute:
		v.Optional = true
		v.Computed = false
		v.PlanModifiers = nil
		v.Attributes = optInManagedAttributes(v.Attributes)

		return v
	case schema.StringAttribute:
		v.Optional = true
		v.Computed = false
		v.Default = nil
		v.PlanModifiers = nil

		return v
	default:
		return attr
	}
}
