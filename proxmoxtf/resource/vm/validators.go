/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// VMIDValidator returns a schema validation function for a VM ID.
func VMIDValidator() schema.SchemaValidateDiagFunc {
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

// BIOSValidator returns a schema validation function for a BIOSValidator type.
func BIOSValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"ovmf",
		"seabios",
	}, false))
}

// CPUArchitectureValidator returns a schema validation function for a CPU architecture.
func CPUArchitectureValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"aarch64",
		"x86_64",
	}, false))
}

// CPUTypeValidator returns a schema validation function for a CPU type.
func CPUTypeValidator() schema.SchemaValidateDiagFunc {
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

// CPUAffinityValidator returns a schema validation function for a CPU affinity.
func CPUAffinityValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(
		validation.StringMatch(regexp.MustCompile(`^\d+[\d-,]*$`), "must contain numbers but also number ranges"),
	)
}

// QEMUAgentTypeValidator is a schema validation function for QEMU agent types.
func QEMUAgentTypeValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{"isa", "virtio"}, false))
}

// KeyboardLayoutValidator is a schema validation function for keyboard layouts.
func KeyboardLayoutValidator() schema.SchemaValidateDiagFunc {
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

// MachineTypeValidator is a schema validation function for machine types.
func MachineTypeValidator() schema.SchemaValidateDiagFunc {
	//nolint:lll
	r := regexp.MustCompile(`^$|^(pc|pc(-i440fx)?-\d+(\.\d+)+(\+pve\d+)?(\.pxe)?|q35|pc-q35-\d+(\.\d+)+(\+pve\d+)?(\.pxe)?|virt(?:-\d+(\.\d+)+)?(\+pve\d+)?)$`)

	return validation.ToDiagFunc(validation.StringMatch(r, "must be a valid machine type"))
}

// TimeoutValidator is a schema validation function for timeouts.
func TimeoutValidator() schema.SchemaValidateDiagFunc {
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

// VGAMemoryValidator is a schema validation function for VGA memory sizes.
func VGAMemoryValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.IntBetween(4, 512))
}

// VGATypeValidator is a schema validation function for VGA device types.
func VGATypeValidator() schema.SchemaValidateDiagFunc {
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

// SCSIHardwareValidator is a schema validation function for SCSI hardware.
func SCSIHardwareValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"lsi",
		"lsi53c810",
		"virtio-scsi-pci",
		"virtio-scsi-single",
		"megasas",
		"pvscsi",
	}, false))
}

// IDEInterfaceValidator is a schema validation function for IDE interfaces.
func IDEInterfaceValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"ide0",
		"ide1",
		"ide2",
		"ide3",
	}, false))
}

// CloudInitInterfaceValidator is a schema validation function that accepts either an IDE interface identifier or an
// empty string, which is used as the default and means "detect which interface should be used automatically".
func CloudInitInterfaceValidator() schema.SchemaValidateDiagFunc {
	r := regexp.MustCompile(`^ide[0-3]|sata[0-5]|scsi(?:30|[12][0-9]|[0-9])$`)

	return validation.ToDiagFunc(validation.Any(
		validation.StringIsEmpty,
		validation.StringMatch(r, "one of ide0..3|sata0..5|scsi0..30"),
	))
}

// CloudInitTypeValidator is a schema validation function for cloud-init types.
func CloudInitTypeValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"configdrive2",
		"nocloud",
	}, false))
}

// AudioDeviceValidator is a schema validation function for audio devices.
func AudioDeviceValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"AC97",
		"ich9-intel-hda",
		"intel-hda",
	}, false))
}

// AudioDriverValidator is a schema validation function for audio drivers.
func AudioDriverValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"spice",
	}, false))
}

// OperatingSystemTypeValidator is a schema validation function for operating system types.
func OperatingSystemTypeValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"l24",
		"l26",
		"other",
		"solaris",
		"w2k",
		"w2k3",
		"w2k8",
		"win7",
		"win8",
		"win10",
		"win11",
		"wvista",
		"wxp",
	}, false))
}

// SerialDeviceValidator is a schema validation function for serial devices.
func SerialDeviceValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, k string) ([]string, []error) {
		v, ok := i.(string)

		var es []error

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return nil, es
		}

		if !strings.HasPrefix(v, "/dev/") && v != "socket" {
			es = append(es, fmt.Errorf("expected %s to be '/dev/*' or 'socket'", k))
			return nil, es
		}

		return nil, es
	})
}

// RangeSemicolonValidator is a proxmox list validation function for ranges with semicolon.
func RangeSemicolonValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(
		validation.StringMatch(
			regexp.MustCompile(`^\d+(?:-\d+)?(?:;\d+(?:-\d+)?)*`),
			"must contain numbers but also number ranges",
		),
	)
}
