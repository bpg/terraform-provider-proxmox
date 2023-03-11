/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
)

func getBIOSValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"ovmf",
		"seabios",
	}, false))
}

func getContentTypeValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"backup",
		"iso",
		"snippets",
		"vztmpl",
	}, false))
}

//nolint:unused
func getCPUFlagsValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, k string) (ws []string, es []error) {
		list, ok := i.([]interface{})

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be []interface{}", k))
			return
		}

		validator := validation.StringInSlice([]string{
			"+aes",
			"-aes",
			"+amd-no-ssb",
			"-amd-no-ssb",
			"+amd-ssbd",
			"-amd-ssbd",
			"+hv-evmcs",
			"-hv-evmcs",
			"+hv-tlbflush",
			"-hv-tlbflush",
			"+ibpb",
			"-ibpb",
			"+md-clear",
			"-md-clear",
			"+pcid",
			"-pcid",
			"+pdpe1gb",
			"-pdpe1gb",
			"+spec-ctrl",
			"-spec-ctrl",
			"+ssbd",
			"-ssbd",
			"+virt-ssbd",
			"-virt-ssbd",
		}, false)

		for li, lv := range list {
			v, ok := lv.(string)

			if !ok {
				es = append(es, fmt.Errorf("expected type of %s[%d] to be string", k, li))
				return
			}

			warns, errs := validator(v, k)

			ws = append(ws, warns...)
			es = append(es, errs...)

			if len(es) > 0 {
				return
			}
		}

		return
	})
}

func getCPUTypeValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, k string) (warnings []string, errors []error) {
		ignoreCase := false
		r, _ := regexp.Compile("^custom-.+")
		valid := []string{
			"486",
			"Broadwell",
			"Broadwell-IBRS",
			"Broadwell-noTSX",
			"Broadwell-noTSX-IBRS",
			"Cascadelake-Server",
			"Conroe",
			"EPYC",
			"EPYC-IBPB",
			"Haswell",
			"Haswell-IBRS",
			"Haswell-noTSX",
			"Haswell-noTSX-IBRS",
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
			"Skylake-Client",
			"Skylake-Client-IBRS",
			"Skylake-Server",
			"Skylake-Server-IBRS",
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
		}

		v, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
			return warnings, errors
		}

		for _, str := range valid {
			if v == str || (ignoreCase && strings.EqualFold(v, str)) {
				return warnings, errors
			}
		}

		if r.MatchString(v) {
			return warnings, errors
		}

		errors = append(errors, fmt.Errorf("expected %s to be one of %v or a custom cpu model, got %s", k, valid, v))
		return warnings, errors
	})
}

func getFileFormatValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"qcow2",
		"raw",
		"vmdk",
	}, false))
}

func getFileIDValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, k string) (ws []string, es []error) {
		v, ok := i.(string)

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if v != "" {
			r := regexp.MustCompile(`^(?i)[a-z\d\-_]+:([a-z\d\-_]+/)?.+$`)
			ok := r.MatchString(v)

			if !ok {
				es = append(es, fmt.Errorf("expected %s to be a valid file identifier (datastore-name:iso/some-file.img), got %s", k, v))
				return
			}
		}

		return
	})
}

func getKeyboardLayoutValidator() schema.SchemaValidateDiagFunc {
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

func diskDigitPrefix(s string) string {
	for i, r := range s {
		if unicode.IsDigit(r) {
			return s[:i]
		}
	}
	return s
}

func getMACAddressValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, k string) (ws []string, es []error) {
		v, ok := i.(string)

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if v != "" {
			r := regexp.MustCompile(`^[A-Z\d]{2}(:[A-Z\d]{2}){5}$`)
			ok := r.MatchString(v)

			if !ok {
				es = append(es, fmt.Errorf("expected %s to be a valid MAC address (A0:B1:C2:D3:E4:F5), got %s", k, v))
				return
			}
		}

		return
	})
}

func getNetworkDeviceModelValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{"e1000", "rtl8139", "virtio", "vmxnet3"}, false))
}

func getQEMUAgentTypeValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{"isa", "virtio"}, false))
}

func getSchemaBlock(r *schema.Resource, d *schema.ResourceData, k []string, i int, allowDefault bool) (map[string]interface{}, error) {
	var resourceBlock map[string]interface{}
	var resourceData interface{}
	var resourceSchema *schema.Schema

	for ki, kv := range k {
		if ki == 0 {
			resourceData = d.Get(kv)
			resourceSchema = r.Schema[kv]
		} else {
			mapValues := resourceData.([]interface{})

			if len(mapValues) <= i {
				return resourceBlock, fmt.Errorf("index out of bounds %d", i)
			}

			mapValue := mapValues[i].(map[string]interface{})

			resourceData = mapValue[kv]
			resourceSchema = resourceSchema.Elem.(*schema.Resource).Schema[kv]
		}
	}

	list := resourceData.([]interface{})

	if len(list) == 0 {
		if allowDefault {
			listDefault, err := resourceSchema.DefaultValue()

			if err != nil {
				return nil, err
			}

			list = listDefault.([]interface{})
		}
	}

	if len(list) > i {
		resourceBlock = list[i].(map[string]interface{})
	}

	return resourceBlock, nil
}

func getTimeoutValidator() schema.SchemaValidateDiagFunc {
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

func getVGAMemoryValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.IntBetween(4, 512))
}

func getVGATypeValidator() schema.SchemaValidateDiagFunc {
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

//nolint:unused
func getVLANIDsValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, k string) (ws []string, es []error) {
		min := 1
		max := 4094

		list, ok := i.([]interface{})

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be []interface{}", k))
			return
		}

		for li, lv := range list {
			v, ok := lv.(int)

			if !ok {
				es = append(es, fmt.Errorf("expected type of %s[%d] to be int", k, li))
				return
			}

			if v != -1 {
				if v < min || v > max {
					es = append(es, fmt.Errorf("expected %s[%d] to be in the range (%d - %d), got %d", k, li, min, max, v))
					return
				}
			}
		}

		return
	})
}

func getVMIDValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, k string) (ws []string, es []error) {
		min := 100
		max := 2147483647

		v, ok := i.(int)

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be int", k))
			return
		}

		if v != -1 {
			if v < min || v > max {
				es = append(es, fmt.Errorf("expected %s to be in the range (%d - %d), got %d", k, min, max, v))
				return
			}
		}

		return
	})
}

func getDiskInfo(vm *proxmox.VirtualEnvironmentVMGetResponseData, d *schema.ResourceData) map[string]*proxmox.CustomStorageDevice {
	currentDisk := d.Get(mkResourceVirtualEnvironmentVMDisk)

	currentDiskList := currentDisk.([]interface{})
	currentDiskMap := map[string]map[string]interface{}{}

	for _, v := range currentDiskList {
		diskMap := v.(map[string]interface{})
		diskInterface := diskMap[mkResourceVirtualEnvironmentVMDiskInterface].(string)

		currentDiskMap[diskInterface] = diskMap
	}

	storageDevices := map[string]*proxmox.CustomStorageDevice{}

	storageDevices["ide0"] = vm.IDEDevice0
	storageDevices["ide1"] = vm.IDEDevice1
	storageDevices["ide2"] = vm.IDEDevice2

	storageDevices["sata0"] = vm.SATADevice0
	storageDevices["sata1"] = vm.SATADevice1
	storageDevices["sata2"] = vm.SATADevice2
	storageDevices["sata3"] = vm.SATADevice3
	storageDevices["sata4"] = vm.SATADevice4
	storageDevices["sata5"] = vm.SATADevice5

	storageDevices["scsi0"] = vm.SCSIDevice0
	storageDevices["scsi1"] = vm.SCSIDevice1
	storageDevices["scsi2"] = vm.SCSIDevice2
	storageDevices["scsi3"] = vm.SCSIDevice3
	storageDevices["scsi4"] = vm.SCSIDevice4
	storageDevices["scsi5"] = vm.SCSIDevice5
	storageDevices["scsi6"] = vm.SCSIDevice6
	storageDevices["scsi7"] = vm.SCSIDevice7
	storageDevices["scsi8"] = vm.SCSIDevice8
	storageDevices["scsi9"] = vm.SCSIDevice9
	storageDevices["scsi10"] = vm.SCSIDevice10
	storageDevices["scsi11"] = vm.SCSIDevice11
	storageDevices["scsi12"] = vm.SCSIDevice12
	storageDevices["scsi13"] = vm.SCSIDevice13

	storageDevices["virtio0"] = vm.VirtualIODevice0
	storageDevices["virtio1"] = vm.VirtualIODevice1
	storageDevices["virtio2"] = vm.VirtualIODevice2
	storageDevices["virtio3"] = vm.VirtualIODevice3
	storageDevices["virtio4"] = vm.VirtualIODevice4
	storageDevices["virtio5"] = vm.VirtualIODevice5
	storageDevices["virtio6"] = vm.VirtualIODevice6
	storageDevices["virtio7"] = vm.VirtualIODevice7
	storageDevices["virtio8"] = vm.VirtualIODevice8
	storageDevices["virtio9"] = vm.VirtualIODevice9
	storageDevices["virtio10"] = vm.VirtualIODevice10
	storageDevices["virtio11"] = vm.VirtualIODevice11
	storageDevices["virtio12"] = vm.VirtualIODevice12
	storageDevices["virtio13"] = vm.VirtualIODevice13
	storageDevices["virtio14"] = vm.VirtualIODevice14
	storageDevices["virtio15"] = vm.VirtualIODevice15

	for k, v := range storageDevices {
		if v != nil {
			if currentDiskMap[k] != nil {
				if currentDiskMap[k][mkResourceVirtualEnvironmentVMDiskFileID] != nil {
					fileID := currentDiskMap[k][mkResourceVirtualEnvironmentVMDiskFileID].(string)
					v.FileID = &fileID
				}
			}

			v.Interface = &k
		}
	}

	return storageDevices
}

// getDiskDatastores returns a list of the used datastores in a VM
func getDiskDatastores(vm *proxmox.VirtualEnvironmentVMGetResponseData, d *schema.ResourceData) []string {
	storageDevices := getDiskInfo(vm, d)
	datastoresSet := map[string]int{}

	for _, diskInfo := range storageDevices {
		// Ignore empty storage devices and storage devices (like ide) which may not have any media mounted
		if diskInfo == nil || diskInfo.FileVolume == "none" {
			continue
		}
		fileIDParts := strings.Split(diskInfo.FileVolume, ":")
		datastoresSet[fileIDParts[0]] = 1
	}

	datastores := []string{}
	for datastore := range datastoresSet {
		datastores = append(datastores, datastore)
	}

	return datastores
}

func getPCIInfo(vm *proxmox.VirtualEnvironmentVMGetResponseData, d *schema.ResourceData) map[string]*proxmox.CustomPCIDevice {
	pciDevices := map[string]*proxmox.CustomPCIDevice{}

	pciDevices["hostpci0"] = vm.PCIDevice0
	pciDevices["hostpci1"] = vm.PCIDevice1
	pciDevices["hostpci2"] = vm.PCIDevice2
	pciDevices["hostpci3"] = vm.PCIDevice3

	return pciDevices
}

func getCloudInitTypeValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"configdrive2",
		"nocloud",
	}, false))
}

func testComputedAttributes(t *testing.T, s *schema.Resource, keys []string) {
	for _, v := range keys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in Schema: Missing definition for \"%s\"", v)
		}

		if !s.Schema[v].Computed {
			t.Fatalf("Error in Schema: Attribute \"%s\" is not computed", v)
		}
	}
}

func testNestedSchemaExistence(t *testing.T, s *schema.Resource, key string) *schema.Resource {
	sh, ok := s.Schema[key].Elem.(*schema.Resource)

	if !ok {
		t.Fatalf("Error in Schema: Missing nested schema for \"%s\"", key)

		return nil
	}

	return sh
}

func testOptionalArguments(t *testing.T, s *schema.Resource, keys []string) {
	for _, v := range keys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in Schema: Missing definition for \"%s\"", v)
		}

		if !s.Schema[v].Optional {
			t.Fatalf("Error in Schema: Argument \"%s\" is not optional", v)
		}
	}
}

func testRequiredArguments(t *testing.T, s *schema.Resource, keys []string) {
	for _, v := range keys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in Schema: Missing definition for \"%s\"", v)
		}

		if !s.Schema[v].Required {
			t.Fatalf("Error in Schema: Argument \"%s\" is not required", v)
		}
	}
}

func testValueTypes(t *testing.T, s *schema.Resource, f map[string]schema.ValueType) {
	for fn, ft := range f {
		if s.Schema[fn] == nil {
			t.Fatalf("Error in Schema: Missing definition for \"%s\"", fn)
		}

		if s.Schema[fn].Type != ft {
			t.Fatalf("Error in Schema: Argument or attribute \"%s\" is not of type \"%v\"", fn, ft)
		}
	}
}

type ErrorDiags diag.Diagnostics

func (diags ErrorDiags) Errors() []error {
	var es []error
	for i := range diags {
		if diags[i].Severity == diag.Error {
			s := fmt.Sprintf("Error: %s", diags[i].Summary)
			if diags[i].Detail != "" {
				s = fmt.Sprintf("%s: %s", s, diags[i].Detail)
			}
			es = append(es, errors.New(s))
		}
	}
	return es
}

func (diags ErrorDiags) Error() string {
	return multierror.ListFormatFunc(diags.Errors())
}
