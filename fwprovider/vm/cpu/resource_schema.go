/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cpu

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceSchema defines the schema for the CPU resource.
func ResourceSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		Description: "The CPU configuration.",
		Optional:    true,
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"affinity": schema.StringAttribute{
				Description: "List of host cores used to execute guest processes, for example: '0,5,8-11'",
				MarkdownDescription: "The CPU cores that are used to run the VM’s vCPU. The value is a list of CPU IDs, " +
					"separated by commas. The CPU IDs are zero-based.  For example, `0,1,2,3` " +
					"(which also can be shortened to `0-3`) means that the VM’s vCPUs are run on the first " +
					"four CPU cores. Setting `affinity` is only allowed for `root@pam` authenticated user.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\d+[\d-,]*$`),
						"must contain numbers or number ranges separated by ','"),
				},
			},
			"architecture": schema.StringAttribute{
				Description: "The CPU architecture.",
				MarkdownDescription: "The CPU architecture `<aarch64 | x86_64>` (defaults to the host). " +
					"Setting `affinity` is only allowed for `root@pam` authenticated user.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("aarch64", "x86_64"),
				},
			},
			"cores": schema.Int64Attribute{
				Description:         "The number of CPU cores per socket.",
				MarkdownDescription: "The number of CPU cores per socket (defaults to `1`).",
				Optional:            true,
				Computed:            true,
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
				Computed:    true,
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
			"hotplugged": schema.Int64Attribute{
				Description:         "The number of hotplugged vCPUs.",
				MarkdownDescription: "The number of hotplugged vCPUs (defaults to `0`).",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 1024),
				},
			},
			"limit": schema.Int64Attribute{
				Description:         "Limit of CPU usage.",
				MarkdownDescription: "Limit of CPU usage (defaults to `0` which means no limit).",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 128),
				},
			},
			"numa": schema.BoolAttribute{
				Description:         "Enable NUMA.",
				MarkdownDescription: "Enable NUMA (defaults to `false`).",
				Optional:            true,
				Computed:            true,
			},
			"sockets": schema.Int64Attribute{
				Description:         "The number of CPU sockets.",
				MarkdownDescription: "The number of CPU sockets (defaults to `1`).",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 16),
				},
			},
			"type": schema.StringAttribute{
				Description: "Emulated CPU type.",
				MarkdownDescription: "Emulated CPU type, " +
					"it's recommended to use `x86-64-v2-AES` or higher (defaults to `kvm64`). " +
					"See https://pve.proxmox.com/pve-docs/pve-admin-guide.html#qm_virtual_machines_settings " +
					"for more information.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"486",
						"Broadwell",
						"Broadwell-IBRS",
						"Broadwell-noTSX",
						"Broadwell-noTSX-IBRS",
						"Cascadelake-Server",
						"Cascadelake-Server-noTSX",
						"Cascadelake-Server-v2",
						"Cascadelake-Server-v4",
						"Cascadelake-Server-v5",
						"Conroe",
						"Cooperlake",
						"Cooperlake-v2",
						"EPYC",
						"EPYC-IBPB",
						"EPYC-Milan",
						"EPYC-Rome",
						"EPYC-Rome-v2",
						"EPYC-v3",
						"Haswell",
						"Haswell-IBRS",
						"Haswell-noTSX",
						"Haswell-noTSX-IBRS",
						"Icelake-Client",
						"Icelake-Client-noTSX",
						"Icelake-Server",
						"Icelake-Server-noTSX",
						"Icelake-Server-v3",
						"Icelake-Server-v4",
						"Icelake-Server-v5",
						"Icelake-Server-v6",
						"IvyBridge",
						"IvyBridge-IBRS",
						"KnightsMill",
						"Nehalem",
						"Nehalem-IBRS",
						"Opteron_G1",
						"Opteron_G2",
						"Opteron_G3",
						"Opteron_G4",
						"Opteron_G5",
						"Penryn",
						"SandyBridge",
						"SandyBridge-IBRS",
						"SapphireRapids",
						"Skylake-Client",
						"Skylake-Client-IBRS",
						"Skylake-Client-noTSX-IBRS",
						"Skylake-Client-v4",
						"Skylake-Server",
						"Skylake-Server-IBRS",
						"Skylake-Server-noTSX-IBRS",
						"Skylake-Server-v4",
						"Skylake-Server-v5",
						"Westmere",
						"Westmere-IBRS",
						"athlon",
						"core2duo",
						"coreduo",
						"host",
						"kvm32",
						"kvm64",
						"max",
						"pentium",
						"pentium2",
						"pentium3",
						"phenom",
						"qemu32",
						"qemu64",
						"x86-64-v2",
						"x86-64-v2-AES",
						"x86-64-v3",
						"x86-64-v4",
					),
				},
			},
			"units": schema.Int64Attribute{
				Description: "CPU weight for a VM.",
				MarkdownDescription: "CPU weight for a VM. Argument is used in the kernel fair scheduler. " +
					"The larger the number is, the more CPU time this VM gets. " +
					"Number is relative to weights of all the other running VMs.",
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.Between(2, 262144),
				},
			},
		},
	}
}
