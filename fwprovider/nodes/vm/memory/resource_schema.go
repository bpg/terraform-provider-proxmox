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
// This implementation uses clearer naming (maximum/minimum) compared to the legacy
// SDK VM resource which uses (dedicated/floating). See GitHub discussion #2198.
//
// Proxmox Memory Ballooning Explained:
//   - maximum: The max RAM available to the VM (Proxmox API: 'memory')
//   - minimum: The guaranteed minimum RAM (Proxmox API: 'balloon')
//   - Setting minimum=0 disables the balloon driver entirely
//   - The range between minimum and maximum is "balloonable" - can be reclaimed by host
//   - shares: CPU scheduler priority (higher = more CPU time during memory pressure)
//   - hugepages: Use hugepages for VM memory (2, 1024, any)
//   - keep_hugepages: Don't release hugepages when VM shuts down
//
// Example:
//
//	memory = {
//	  maximum = 4096  # VM can use up to 4GB
//	  minimum = 2048  # Host guarantees 2GB minimum
//	}
//	# Result: VM gets 2-4GB depending on host memory pressure
//
// Legacy SDK Mapping (for migration):
//   - maximum = dedicated (Proxmox: memory)
//   - minimum = floating (Proxmox: balloon)
//   - shares = shared (Proxmox: shares)
func ResourceSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		Description: "Memory configuration. Controls maximum available RAM and minimum guaranteed RAM via ballooning.",
		MarkdownDescription: "Memory configuration for the VM. Uses Proxmox memory ballooning to allow dynamic memory allocation. " +
			"The `maximum` sets the upper limit, while `minimum` sets the guaranteed floor. " +
			"The host can reclaim memory between these values when needed. " +
			"\n\n**Note:** This uses clearer naming (`maximum`/`minimum`) compared to the legacy `vm` resource " +
			"which uses `dedicated`/`floating`. See the [migration guide](/docs/guides/migration-vm2-clone.md#memory-terminology) " +
			"for mapping details.",
		Optional: true,
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"maximum": schema.Int64Attribute{
				Description: "Maximum available memory in MiB (Proxmox API: 'memory'). " +
					"This is the upper limit of RAM the VM can use when balloon device is enabled.",
				MarkdownDescription: "Maximum available memory in MiB. This is the upper limit of RAM the VM can use " +
					"when the balloon device is enabled (defaults to `512` MiB). " +
					"\n\n**Proxmox API:** `memory` parameter " +
					"\n\n**Legacy SDK:** `dedicated` parameter",
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(512),
				Validators: []validator.Int64{
					int64validator.Between(64, 268435456), // 64 MiB to 256 TiB
				},
			},
			"minimum": schema.Int64Attribute{
				Description: "Minimum guaranteed memory in MiB (Proxmox API: 'balloon'). " +
					"The guaranteed amount of RAM. Set to 0 to disable balloon device.",
				MarkdownDescription: "Minimum guaranteed memory in MiB. This is the floor amount of RAM that is always " +
					"guaranteed to the VM. Setting to `0` disables the balloon driver entirely (defaults to `0`). " +
					"\n\n**How it works:** The host can reclaim memory between `minimum` and `maximum` when under " +
					"memory pressure. The VM is guaranteed to always have at least `minimum` MiB available. " +
					"\n\n**Proxmox API:** `balloon` parameter " +
					"\n\n**Legacy SDK:** `floating` parameter",
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.Between(0, 268435456), // 0 to 256 TiB
				},
			},
			"shares": schema.Int64Attribute{
				Description: "CPU scheduler priority for memory ballooning (Proxmox API: 'shares'). " +
					"Higher values give the VM more CPU time during memory pressure.",
				MarkdownDescription: "CPU scheduler priority for memory ballooning. This is used by the " +
					"kernel fair scheduler. Higher values mean this VM gets more CPU time during memory ballooning " +
					"operations. The value is relative to other running VMs (defaults to `1000`). " +
					"\n\n**Proxmox API:** `shares` parameter " +
					"\n\n**Legacy SDK:** `shared` parameter",
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
					"\n- `any` - Use any available hugepage size" +
					"\n\n**Proxmox API:** `hugepages` parameter",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("2", "1024", "any"),
				},
			},
			"keep_hugepages": schema.BoolAttribute{
				Description: "Keep hugepages allocated when VM is stopped (Proxmox API: 'keephugepages').",
				MarkdownDescription: "Don't release hugepages when the VM shuts down. By default, hugepages are " +
					"released back to the host when the VM stops. Setting this to `true` keeps them allocated " +
					"for faster VM startup (defaults to `false`). " +
					"\n\n**Proxmox API:** `keephugepages` parameter",
				Optional: true,
				Computed: true,
			},
		},
	}
}
