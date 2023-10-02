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
		minID := 100
		maxID := 2147483647

		var ws []string
		var es []error

		v, ok := i.(int)

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be int", k))
			return ws, es
		}

		if v != -1 {
			if v < minID || v > maxID {
				es = append(es, fmt.Errorf("expected %s to be in the range (%d - %d), got %d", k, minID, maxID, v))
				return ws, es
			}
		}

		return ws, es
	})
}

// BIOS returns a schema validation function for a BIOS type.
func BIOS() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"ovmf",
		"seabios",
	}, false))
}

// ContentType returns a schema validation function for a content type on a storage device.
func ContentType() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"backup",
		"iso",
		"snippets",
		"vztmpl",
	}, false))
}

// CPUType returns a schema validation function for a CPU type.
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

// NetworkDeviceModel is a schema validation function for network device models.
func NetworkDeviceModel() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{"e1000", "rtl8139", "virtio", "vmxnet3"}, false))
}

// QEMUAgentType is a schema validation function for QEMU agent types.
func QEMUAgentType() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{"isa", "virtio"}, false))
}

// KeyboardLayout is a schema validation function for keyboard layouts.
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

// Timeout is a schema validation function for timeouts.
func Timeout() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, k string) ([]string, []error) {
		v, ok := i.(string)

		var ws []string
		var es []error

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return ws, es
		}

		_, err := time.ParseDuration(v)
		if err != nil {
			es = append(es, fmt.Errorf("expected value of %s to be a duration - got: %s", k, v))
			return ws, es
		}

		return ws, es
	})
}

// VGAMemory is a schema validation function for VGA memory sizes.
func VGAMemory() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.IntBetween(4, 512))
}

// VGAType is a schema validation function for VGA device types.
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

// SCSIHardware is a schema validation function for SCSI hardware.
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

// IDEInterface is a schema validation function for IDE interfaces.
func IDEInterface() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"ide0",
		"ide1",
		"ide2",
		"ide3",
	}, false))
}

// CloudInitInterface is a schema validation function that accepts either an IDE interface identifier or an
// empty string, which is used as the default and means "detect which interface should be used automatically".
func CloudInitInterface() schema.SchemaValidateDiagFunc {
	r := regexp.MustCompile(`^ide[0-3]|sata[0-5]|scsi(?:30|[12][0-9]|[0-9])$`)

	return validation.ToDiagFunc(validation.Any(
		validation.StringIsEmpty,
		validation.StringMatch(r, "one of ide0..3|sata0..5|scsi0..30"),
	))
}

// CloudInitType is a schema validation function for cloud-init types.
func CloudInitType() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"configdrive2",
		"nocloud",
	}, false))
}
