/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/bpg/terraform-provider-proxmox/internal/types"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
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

func getFileSizeValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(func(i interface{}, k string) ([]string, []error) {
		v, ok := i.(string)
		var es []error

		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return nil, es
		}

		if v != "" {
			_, err := types.ParseDiskSize(v)
			if err != nil {
				es = append(es, fmt.Errorf("expected %s to be a valid file size (100, 1M, 1G), got %s", k, v))
				return nil, es
			}
		}

		return []string{}, es
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

func getSCSIHardwareValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"lsi",
		"lsi53c810",
		"virtio-scsi-pci",
		"virtio-scsi-single",
		"megasas",
		"pvscsi",
	}, false))
}

func getIDEInterfaceValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"ide0",
		"ide1",
		"ide2",
		"ide3",
	}, false))
}

// suppressIfListsAreEqualIgnoringOrder is a customdiff.SuppressionFunc that suppresses
// changes to a list if the old and new lists are equal, ignoring the order of the
// elements.
// It will be called for each list item, so it is not super efficient. It is
// recommended to use it only for small lists.
// Ref: https://github.com/hashicorp/terraform-plugin-sdk/issues/477
func suppressIfListsAreEqualIgnoringOrder(key, _, _ string, d *schema.ResourceData) bool {
	// the key is a path to the list item, not the list itself, e.g. "tags.0"
	lastDotIndex := strings.LastIndex(key, ".")
	if lastDotIndex != -1 {
		key = key[:lastDotIndex]
	}

	oldData, newData := d.GetChange(key)
	if oldData == nil || newData == nil {
		return false
	}

	oldArray := oldData.([]interface{})
	newArray := newData.([]interface{})

	if len(oldArray) != len(newArray) {
		return false
	}

	oldEvents := make([]string, len(oldArray))
	newEvents := make([]string, len(newArray))

	for i, oldEvt := range oldArray {
		oldEvents[i] = fmt.Sprint(oldEvt)
	}

	for j, newEvt := range newArray {
		newEvents[j] = fmt.Sprint(newEvt)
	}

	sort.Strings(oldEvents)
	sort.Strings(newEvents)

	return reflect.DeepEqual(oldEvents, newEvents)
}

func getDiskInfo(resp *vms.GetResponseData, d *schema.ResourceData) map[string]*vms.CustomStorageDevice {
	currentDisk := d.Get(mkResourceVirtualEnvironmentVMDisk)

	currentDiskList := currentDisk.([]interface{})
	currentDiskMap := map[string]map[string]interface{}{}

	for _, v := range currentDiskList {
		diskMap := v.(map[string]interface{})
		diskInterface := diskMap[mkResourceVirtualEnvironmentVMDiskInterface].(string)

		currentDiskMap[diskInterface] = diskMap
	}

	storageDevices := map[string]*vms.CustomStorageDevice{}

	storageDevices["ide0"] = resp.IDEDevice0
	storageDevices["ide1"] = resp.IDEDevice1
	storageDevices["ide2"] = resp.IDEDevice2
	storageDevices["ide3"] = resp.IDEDevice3

	storageDevices["sata0"] = resp.SATADevice0
	storageDevices["sata1"] = resp.SATADevice1
	storageDevices["sata2"] = resp.SATADevice2
	storageDevices["sata3"] = resp.SATADevice3
	storageDevices["sata4"] = resp.SATADevice4
	storageDevices["sata5"] = resp.SATADevice5

	storageDevices["scsi0"] = resp.SCSIDevice0
	storageDevices["scsi1"] = resp.SCSIDevice1
	storageDevices["scsi2"] = resp.SCSIDevice2
	storageDevices["scsi3"] = resp.SCSIDevice3
	storageDevices["scsi4"] = resp.SCSIDevice4
	storageDevices["scsi5"] = resp.SCSIDevice5
	storageDevices["scsi6"] = resp.SCSIDevice6
	storageDevices["scsi7"] = resp.SCSIDevice7
	storageDevices["scsi8"] = resp.SCSIDevice8
	storageDevices["scsi9"] = resp.SCSIDevice9
	storageDevices["scsi10"] = resp.SCSIDevice10
	storageDevices["scsi11"] = resp.SCSIDevice11
	storageDevices["scsi12"] = resp.SCSIDevice12
	storageDevices["scsi13"] = resp.SCSIDevice13

	storageDevices["virtio0"] = resp.VirtualIODevice0
	storageDevices["virtio1"] = resp.VirtualIODevice1
	storageDevices["virtio2"] = resp.VirtualIODevice2
	storageDevices["virtio3"] = resp.VirtualIODevice3
	storageDevices["virtio4"] = resp.VirtualIODevice4
	storageDevices["virtio5"] = resp.VirtualIODevice5
	storageDevices["virtio6"] = resp.VirtualIODevice6
	storageDevices["virtio7"] = resp.VirtualIODevice7
	storageDevices["virtio8"] = resp.VirtualIODevice8
	storageDevices["virtio9"] = resp.VirtualIODevice9
	storageDevices["virtio10"] = resp.VirtualIODevice10
	storageDevices["virtio11"] = resp.VirtualIODevice11
	storageDevices["virtio12"] = resp.VirtualIODevice12
	storageDevices["virtio13"] = resp.VirtualIODevice13
	storageDevices["virtio14"] = resp.VirtualIODevice14
	storageDevices["virtio15"] = resp.VirtualIODevice15

	for k, v := range storageDevices {
		if v != nil {
			if currentDiskMap[k] != nil {
				if currentDiskMap[k][mkResourceVirtualEnvironmentVMDiskFileID] != nil {
					fileID := currentDiskMap[k][mkResourceVirtualEnvironmentVMDiskFileID].(string)
					v.FileID = &fileID
				}
			}
			// defensive copy of the loop variable
			iface := k
			v.Interface = &iface
		}
	}

	return storageDevices
}

// getDiskDatastores returns a list of the used datastores in a VM
func getDiskDatastores(vm *vms.GetResponseData, d *schema.ResourceData) []string {
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

	if vm.EFIDisk != nil {
		fileIDParts := strings.Split(vm.EFIDisk.FileVolume, ":")
		datastoresSet[fileIDParts[0]] = 1
	}

	datastores := []string{}
	for datastore := range datastoresSet {
		datastores = append(datastores, datastore)
	}

	return datastores
}

func getPCIInfo(resp *vms.GetResponseData, _ *schema.ResourceData) map[string]*vms.CustomPCIDevice {
	pciDevices := map[string]*vms.CustomPCIDevice{}

	pciDevices["hostpci0"] = resp.PCIDevice0
	pciDevices["hostpci1"] = resp.PCIDevice1
	pciDevices["hostpci2"] = resp.PCIDevice2
	pciDevices["hostpci3"] = resp.PCIDevice3

	return pciDevices
}

func getCloudInitTypeValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"configdrive2",
		"nocloud",
	}, false))
}

func parseImportIDWithNodeName(id string) (string, string, error) {
	nodeName, id, found := strings.Cut(id, "/")

	if !found {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected node/id", id)
	}

	return nodeName, id, nil
}
