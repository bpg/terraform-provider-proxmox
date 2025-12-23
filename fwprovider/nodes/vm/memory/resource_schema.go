/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package memory

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ResourceSchema defines the schema for the memory resource.
//
// Proxmox Memory Ballooning Explained:
//   - size: Total memory available to the VM in MiB (Proxmox API: 'memory')
//   - balloon: Minimum guaranteed memory via balloon device in MiB (Proxmox API: 'balloon')
//   - Setting balloon=0 disables the balloon driver entirely
//   - The range between balloon and size is "balloonable" - can be reclaimed by host
//   - shares: CPU scheduler priority (higher = more CPU time during memory pressure)
//   - hugepages: Use hugepages for VM memory (2, 1024, any)
//   - keep_hugepages: Don't release hugepages when VM shuts down
//
// Example:
//
//	memory = {
//	  size    = 4096  # VM can use up to 4GB
//	  balloon = 2048  # Host guarantees 2GB minimum
//	}
//	# Result: VM gets 2-4GB depending on host memory pressure
func ResourceSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		Description: "Memory configuration. Controls total available RAM and minimum guaranteed RAM via ballooning.",
		MarkdownDescription: "Memory configuration for the VM. Uses Proxmox memory ballooning to allow dynamic memory allocation. " +
			"The `size` sets the total available RAM, while `balloon` sets the guaranteed floor. " +
			"The host can reclaim memory between these values when needed.",
		Optional: true,
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"size": schema.Int64Attribute{
				Description: "Total memory available to the VM in MiB. " +
					"This is the total RAM the VM can use. When ballooning is enabled, memory between `balloon` and `size` can be reclaimed by the host.",
				MarkdownDescription: "Total memory available to the VM in MiB. This is the total RAM the VM can use. " +
					"When ballooning is enabled (balloon > 0), memory between `balloon` and `size` can be reclaimed by the host. " +
					"When ballooning is disabled (balloon = 0), this is the fixed amount of RAM allocated to the VM (defaults to `512` MiB).",
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(512),
				Validators: []validator.Int64{
					int64validator.Between(64, 268435456), // 64 MiB to 256 TiB
				},
			},
			"balloon": schema.Int64Attribute{
				Description: "Minimum guaranteed memory in MiB via balloon device. " +
					"The guaranteed amount of RAM. Set to 0 to disable balloon device.",
				MarkdownDescription: "Minimum guaranteed memory in MiB via balloon device. This is the floor amount of RAM that is always " +
					"guaranteed to the VM. Setting to `0` disables the balloon driver entirely (defaults to `0`). " +
					"\n\n**How it works:** The host can reclaim memory between `balloon` and `size` when under " +
					"memory pressure. The VM is guaranteed to always have at least `balloon` MiB available.",
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.Between(0, 268435456), // 0 to 256 TiB
				},
			},
			"shares": schema.Int64Attribute{
				Description: "CPU scheduler priority for memory ballooning. " +
					"Higher values give the VM more CPU time during memory pressure.",
				MarkdownDescription: "CPU scheduler priority for memory ballooning. This is used by the " +
					"kernel fair scheduler. Higher values mean this VM gets more CPU time during memory ballooning " +
					"operations. The value is relative to other running VMs (defaults to `1000`).",
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(1000),
				Validators: []validator.Int64{
					int64validator.Between(0, 50000),
				},
			},
			"hugepages": schema.StringAttribute{
				Description: "Use hugepages for VM memory. Options: '2' (2 MiB), '1024' (1 GiB), 'any'.",
				MarkdownDescription: "Enable hugepages for VM memory allocation. Hugepages can improve performance " +
					"for memory-intensive workloads by reducing TLB misses. " +
					"\n\n**Options:**" +
					"\n- `2` - Use 2 MiB hugepages" +
					"\n- `1024` - Use 1 GiB hugepages" +
					"\n- `any` - Use any available hugepage size",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("2", "1024", "any"),
				},
			},
			"keep_hugepages": schema.BoolAttribute{
				Description: "Keep hugepages allocated when VM is stopped.",
				MarkdownDescription: "Don't release hugepages when the VM shuts down. By default, hugepages are " +
					"released back to the host when the VM stops. Setting this to `true` keeps them allocated " +
					"for faster VM startup (defaults to `false`).",
				Optional: true,
				Computed: true,
			},
		},
	}
}
