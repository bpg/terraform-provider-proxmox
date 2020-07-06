/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func getBIOSValidator() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"ovmf",
		"seabios",
	}, false)
}

func getContentTypeValidator() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"backup",
		"iso",
		"snippets",
		"vztmpl",
	}, false)
}

func getCPUFlagsValidator() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (ws []string, es []error) {
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
	}
}

func getCPUTypeValidator() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
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
	}, false)
}

func getFileFormatValidator() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"qcow2",
		"raw",
		"vmdk",
	}, false)
}

func getFileIDValidator() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (ws []string, es []error) {
		v, ok := i.(string)

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if v != "" {
			r := regexp.MustCompile(`^(?i)[a-z0-9\-_]+:([a-z0-9\-_]+/)?.+$`)
			ok := r.MatchString(v)

			if !ok {
				es = append(es, fmt.Errorf("expected %s to be a valid file identifier (datastore-name:iso/some-file.img), got %s", k, v))
				return
			}
		}

		return
	}
}

func getKeyboardLayoutValidator() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
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
	}, false)
}

func diskDigitPrefix(s string) string {
	for i, r := range s {
		if unicode.IsDigit(r) {
			return s[:i]
		}
	}
	return s
}

func getMACAddressValidator() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (ws []string, es []error) {
		v, ok := i.(string)

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if v != "" {
			r := regexp.MustCompile(`^[A-Z0-9]{2}(:[A-Z0-9]{2}){5}$`)
			ok := r.MatchString(v)

			if !ok {
				es = append(es, fmt.Errorf("expected %s to be a valid MAC address (A0:B1:C2:D3:E4:F5), got %s", k, v))
				return
			}
		}

		return
	}
}

func getNetworkDeviceModelValidator() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{"e1000", "rtl8139", "virtio", "vmxnet3"}, false)
}

func getQEMUAgentTypeValidator() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{"isa", "virtio"}, false)
}

func getSchemaBlock(r *schema.Resource, d *schema.ResourceData, m interface{}, k []string, i int, allowDefault bool) (map[string]interface{}, error) {
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
				return resourceBlock, fmt.Errorf("Index out of bounds %d", i)
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

func getTimeoutValidator() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (ws []string, es []error) {
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
	}
}

func getVGAMemoryValidator() schema.SchemaValidateFunc {
	return validation.IntBetween(4, 512)
}

func getVGATypeValidator() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
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
	}, false)
}

func getVLANIDsValidator() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (ws []string, es []error) {
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
	}
}

func getVMIDValidator() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (ws []string, es []error) {
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
	}
}

func getOrderedDiskDeviceList(diskDeviceMap map[string]map[string]proxmox.CustomStorageDevice, diskInterface string) proxmox.CustomStorageDevices {
	diskDevices := diskDeviceMap[diskInterface]

	if diskDevices == nil {
		return nil
	}

	orderedDiskList := make(proxmox.CustomStorageDevices, len(diskDevices))

	for _, value := range diskDevices {
		orderedDiskList = append(orderedDiskList, value)
	}

	return orderedDiskList
}

func getDiskInfo(data *proxmox.VirtualEnvironmentVMGetResponseData) map[string]*proxmox.CustomStorageDevice {
	storageDevices := make(map[string]*proxmox.CustomStorageDevice)
	storageDevices["ide0"] = data.IDEDevice0
	storageDevices["ide1"] = data.IDEDevice1
	storageDevices["ide2"] = data.IDEDevice2

	storageDevices["sata0"] = data.SATADevice0
	storageDevices["sata1"] = data.SATADevice1
	storageDevices["sata2"] = data.SATADevice2
	storageDevices["sata3"] = data.SATADevice3
	storageDevices["sata4"] = data.SATADevice4
	storageDevices["sata5"] = data.SATADevice5

	storageDevices["scsi0"] = data.SCSIDevice0
	storageDevices["scsi1"] = data.SCSIDevice1
	storageDevices["scsi2"] = data.SCSIDevice2
	storageDevices["scsi3"] = data.SCSIDevice3
	storageDevices["scsi4"] = data.SCSIDevice4
	storageDevices["scsi5"] = data.SCSIDevice5
	storageDevices["scsi6"] = data.SCSIDevice6
	storageDevices["scsi7"] = data.SCSIDevice7
	storageDevices["scsi8"] = data.SCSIDevice8
	storageDevices["scsi9"] = data.SCSIDevice9
	storageDevices["scsi10"] = data.SCSIDevice10
	storageDevices["scsi11"] = data.SCSIDevice11
	storageDevices["scsi12"] = data.SCSIDevice12
	storageDevices["scsi13"] = data.SCSIDevice13

	storageDevices["virtio0"] = data.VirtualIODevice0
	storageDevices["virtio1"] = data.VirtualIODevice1
	storageDevices["virtio2"] = data.VirtualIODevice2
	storageDevices["virtio3"] = data.VirtualIODevice3
	storageDevices["virtio4"] = data.VirtualIODevice4
	storageDevices["virtio5"] = data.VirtualIODevice5
	storageDevices["virtio6"] = data.VirtualIODevice6
	storageDevices["virtio7"] = data.VirtualIODevice7
	storageDevices["virtio8"] = data.VirtualIODevice8
	storageDevices["virtio9"] = data.VirtualIODevice9
	storageDevices["virtio10"] = data.VirtualIODevice10
	storageDevices["virtio11"] = data.VirtualIODevice11
	storageDevices["virtio12"] = data.VirtualIODevice12
	storageDevices["virtio13"] = data.VirtualIODevice13
	storageDevices["virtio14"] = data.VirtualIODevice14
	storageDevices["virtio15"] = data.VirtualIODevice15

	return storageDevices
}

func parseDiskSize(size *string) (int, error) {
	var diskSize int
	var err error
	if size != nil {
		if strings.HasSuffix(*size, "T") {
			diskSize, err = strconv.Atoi(strings.TrimSuffix(*size, "T"))

			if err != nil {
				return -1, err
			}

			diskSize = int(math.Ceil(float64(diskSize) * 1024))
		} else if strings.HasSuffix(*size, "G") {
			diskSize, err = strconv.Atoi(strings.TrimSuffix(*size, "G"))

			if err != nil {
				return -1, err
			}
		} else if strings.HasSuffix(*size, "M") {
			diskSize, err = strconv.Atoi(strings.TrimSuffix(*size, "M"))

			if err != nil {
				return -1, err
			}

			diskSize = int(math.Ceil(float64(diskSize) / 1024))
		} else {
			return -1, fmt.Errorf("Cannot parse storage size \"%s\"", *size)
		}
	}
	return diskSize, err
}

func testComputedAttributes(t *testing.T, s *schema.Resource, keys []string) {
	for _, v := range keys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in Schema: Missing definition for \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in Schema: Attribute \"%s\" is not computed", v)
		}
	}
}

func testNestedSchemaExistence(t *testing.T, s *schema.Resource, key string) *schema.Resource {
	schema, ok := s.Schema[key].Elem.(*schema.Resource)

	if !ok {
		t.Fatalf("Error in Schema: Missing nested schema for \"%s\"", key)

		return nil
	}

	return schema
}

func testOptionalArguments(t *testing.T, s *schema.Resource, keys []string) {
	for _, v := range keys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in Schema: Missing definition for \"%s\"", v)
		}

		if s.Schema[v].Optional != true {
			t.Fatalf("Error in Schema: Argument \"%s\" is not optional", v)
		}
	}
}

func testRequiredArguments(t *testing.T, s *schema.Resource, keys []string) {
	for _, v := range keys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in Schema: Missing definition for \"%s\"", v)
		}

		if s.Schema[v].Required != true {
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
