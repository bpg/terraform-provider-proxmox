/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package validator

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// VMID returns a schema validation function for a VM ID.
func VMID() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, k string) ([]string, []error) {
		min := 100
		max := 2147483647

		var ws []string
		var es []error

		v, ok := i.(int)

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be int", k))
			return ws, es
		}

		if v != -1 {
			if v < min || v > max {
				es = append(es, fmt.Errorf("expected %s to be in the range (%d - %d), got %d", k, min, max, v))
				return ws, es
			}
		}

		return ws, es
	})
}

func BIOS() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"ovmf",
		"seabios",
	}, false))
}

func ContentType() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"backup",
		"iso",
		"snippets",
		"vztmpl",
	}, false))
}

func CPUType() schema.SchemaValidateDiagFunc {
	standardTypes := []string{
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
	}

	return validation.ToDiagFunc(validation.Any(
		validation.StringInSlice(standardTypes, false),
		validation.StringMatch(regexp.MustCompile(`^custom-.+$`), "must be a valid custom CPU type"),
	))
}

func NetworkDeviceModel() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{"e1000", "rtl8139", "virtio", "vmxnet3"}, false))
}

func QEMUAgentType() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{"isa", "virtio"}, false))
}

func KeyboardLayout() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"da",
		"de",
		"de-ch",
		"en-gb",
		"en-us",
		"es",
		"fi",
		"fr",
		"fr-be",
		"fr-ca",
		"fr-ch",
		"hu",
		"is",
		"it",
		"ja",
		"lt",
		"mk",
		"nl",
		"no",
		"pl",
		"pt",
		"pt-br",
		"sl",
		"sv",
		"tr",
	}, false))
}

func Timeout() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, k string) (ws []string, es []error) {
		v, ok := i.(string)

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		_, err := time.ParseDuration(v)

		if err != nil {
			es = append(es, fmt.Errorf("expected value of %s to be a duration - got: %s", k, v))
			return
		}

		return
	})
}

func VGAMemory() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.IntBetween(4, 512))
}

func VGAType() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"cirrus",
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
		"vmware",
	}, false))
}

func SCSIHardware() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"lsi",
		"lsi53c810",
		"virtio-scsi-pci",
		"virtio-scsi-single",
		"megasas",
		"pvscsi",
	}, false))
}

func IDEInterface() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"ide0",
		"ide1",
		"ide2",
		"ide3",
	}, false))
}

func CloudInitType() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"configdrive2",
		"nocloud",
	}, false))
}
