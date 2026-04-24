/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cpu

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceSchema defines the schema for the CPU resource.
//
// Scope matches the PVE web UI's "Processors" dialog: cores, sockets, vcpus, type+flags,
// limit, units, affinity, arch, numa.
//
// All attributes are `Optional` only. PVE surfaces only the keys the user explicitly wrote to
// the config — there is no implicit read-back to reconcile, so an attribute absent from the
// plan is sent as `delete=<key>` on the wire and the prior value is dropped cleanly.
func ResourceSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		Description: "The CPU configuration.",
		Optional:    true,
		Attributes: map[string]schema.Attribute{
			"affinity": schema.StringAttribute{
				Description: "List of host cores used to execute guest processes, for example: '0,5,8-11'",
				MarkdownDescription: "The CPU cores that are used to run the VM’s vCPU. The value is a list of CPU IDs, " +
					"separated by commas. The CPU IDs are zero-based.  For example, `0,1,2,3` " +
					"(which also can be shortened to `0-3`) means that the VM’s vCPUs are run on the first " +
					"four CPU cores. Setting `affinity` is only allowed for `root@pam` authenticated user.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\d+[\d-,]*$`),
						"must contain numbers or number ranges separated by ','"),
				},
			},
			"architecture": schema.StringAttribute{
				Description: "The CPU architecture.",
				MarkdownDescription: "The CPU architecture `<aarch64 | x86_64>` (defaults to the host). " +
					"Setting `architecture` is only allowed for `root@pam` authenticated user.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("aarch64", "x86_64"),
				},
			},
			"cores": schema.Int64Attribute{
				Description:         "The number of CPU cores per socket.",
				MarkdownDescription: "The number of CPU cores per socket (PVE defaults to `1` when unset).",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 1024),
				},
			},
			"flags": schema.SetAttribute{
				Description: "Set of additional CPU flags.",
				MarkdownDescription: "Set of additional CPU flags. " +
					"Use `+FLAG` to enable, `-FLAG` to disable a flag. Custom CPU models can specify any flag " +
					"supported by QEMU/KVM, VM-specific flags must be from the following set for security reasons: " +
					"`pcid`, `spec-ctrl`, `ibpb`, `ssbd`, `virt-ssbd`, `amd-ssbd`, `amd-no-ssb`, `pdpe1gb`, " +
					"`md-clear`, `hv-tlbflush`, `hv-evmcs`, `aes`.",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("type")),
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(`(.|\s)*\S(.|\s)*`),
							"must be a non-empty and non-whitespace string",
						),
						stringvalidator.LengthAtLeast(1),
					),
				},
			},
			"limit": schema.Float64Attribute{
				Description:         "Limit of CPU usage.",
				MarkdownDescription: "Limit of CPU usage. `0` means no limit (PVE default).",
				Optional:            true,
				Validators: []validator.Float64{
					float64validator.Between(0, 128),
				},
			},
			"numa": schema.BoolAttribute{
				Description:         "Enable NUMA.",
				MarkdownDescription: "Enable NUMA topology emulation. Matches the PVE Processors → **Enable NUMA** checkbox.",
				Optional:            true,
			},
			"sockets": schema.Int64Attribute{
				Description:         "The number of CPU sockets.",
				MarkdownDescription: "The number of CPU sockets (PVE defaults to `1` when unset).",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 16),
				},
			},
			"type": schema.StringAttribute{
				Description: "Emulated CPU type.",
				MarkdownDescription: "Emulated CPU type, " +
					"it's recommended to use `x86-64-v2-AES` or higher. See " +
					"[the PVE admin guide](https://pve.proxmox.com/pve-docs/pve-admin-guide.html#qm_virtual_machines_settings) " +
					"for the full list of supported types.",
				Optional: true,
				// No OneOf validator: the PVE-supported CPU type list evolves across releases (new
				// microarchitectures, `-vN` spin-offs); PVE validates server-side. See ADR-004
				// §Enum Validator Rule.
			},
			"units": schema.Int64Attribute{
				Description: "CPU weight for a VM.",
				MarkdownDescription: "CPU weight for a VM. Argument is used in the kernel fair scheduler. " +
					"The larger the number is, the more CPU time this VM gets. " +
					"Number is relative to weights of all the other running VMs. " +
					"On cgroup v2 `0` is a valid value meaning disable CPU share weighting.",
				Optional: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"vcpus": schema.Int64Attribute{
				Description: "Number of active vCPUs (for CPU hotplug).",
				MarkdownDescription: "Number of vCPUs started with the VM, bounded by `cores * sockets`. " +
					"Matches the PVE Processors → **VCPUs** field. Leave unset to start with `cores * sockets` vCPUs. " +
					"Requires PVE hotplug feature enabled to change at runtime.",
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 1024),
				},
			},
		},
	}
}
