/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cdrom

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

func normalizeFileID(fileID string) string {
	if fileID == "" {
		return DefaultFileID
	}

	return fileID
}

// GetCDROMDeviceObjects converts configured CD-ROM blocks into storage devices keyed by interface.
func GetCDROMDeviceObjects(cdrom []any) vms.CustomStorageDevices {
	cdromMedia := "cdrom"
	deviceObjects := vms.CustomStorageDevices{}

	for _, entry := range cdrom {
		if entry == nil {
			continue
		}

		block := entry.(map[string]any)

		cdromInterface, _ := block[MkCDROMInterface].(string)
		if cdromInterface == "" {
			continue
		}

		cdromFileID, _ := block[MkCDROMFileID].(string)
		deviceObjects[cdromInterface] = &vms.CustomStorageDevice{
			FileVolume: normalizeFileID(cdromFileID),
			Media:      &cdromMedia,
		}
	}

	return deviceObjects
}

// GetCDROMStorageDevices filters VM storage devices down to attached CD-ROM devices only.
func GetCDROMStorageDevices(vmConfig *vms.GetResponseData) vms.CustomStorageDevices {
	cdromDevices := vms.CustomStorageDevices{}

	for iface, dev := range vmConfig.StorageDevices {
		if dev == nil || dev.Media == nil || *dev.Media != "cdrom" {
			continue
		}

		cdromDevices[iface] = dev
	}

	return cdromDevices
}

// OrderedInterfaces returns CD-ROM interfaces in state order first, then any newly discovered interfaces.
func OrderedInterfaces(cdrom []any, deviceObjects vms.CustomStorageDevices) []string {
	interfaces := utils.ListResourcesAttributeValue(cdrom, MkCDROMInterface)
	if len(interfaces) == 0 {
		return orderedKeys(deviceObjects)
	}

	ordered := make([]string, 0, len(deviceObjects))

	for _, iface := range interfaces {
		if _, ok := deviceObjects[iface]; ok {
			ordered = append(ordered, iface)
		}
	}

	for _, iface := range orderedKeys(deviceObjects) {
		if slices.Contains(ordered, iface) {
			continue
		}

		ordered = append(ordered, iface)
	}

	return ordered
}

// MergeCloneDevices adds planned CD-ROM devices into the clone storage device set.
func MergeCloneDevices(planCDROMs vms.CustomStorageDevices, storageDevices vms.CustomStorageDevices) {
	maps.Copy(storageDevices, planCDROMs)
}

// ValidateInterfacesForMachine enforces machine-specific CD-ROM interface restrictions.
func ValidateInterfacesForMachine(machineType string, cdrom []any) error {
	if !isQ35MachineType(machineType) {
		return nil
	}

	for _, entry := range cdrom {
		if entry == nil {
			continue
		}

		block, ok := entry.(map[string]any)
		if !ok {
			continue
		}

		cdromInterface, _ := block[MkCDROMInterface].(string)
		if !strings.HasPrefix(cdromInterface, "ide") {
			continue
		}

		if cdromInterface == "ide0" || cdromInterface == "ide2" {
			continue
		}

		return fmt.Errorf(
			"cdrom interface %q is invalid for q35 machine type: only ide0 and ide2 are supported on the IDE bus",
			cdromInterface,
		)
	}

	return nil
}

// BuildState reconstructs Terraform CD-ROM state from current devices while preserving legacy empty file_id values.
func BuildState(currentCDROM []any, deviceObjects vms.CustomStorageDevices, isClone bool) []any {
	cdromMap := map[string]any{}
	currentCDROMMap := utils.MapResourcesByAttribute(currentCDROM, MkCDROMInterface)

	for iface, dev := range deviceObjects {
		cdromBlock := map[string]any{
			MkCDROMFileID:    dev.FileVolume,
			MkCDROMInterface: iface,
		}

		if currentBlock, ok := currentCDROMMap[iface].(map[string]any); ok {
			preserveLegacyEmptyFileID(currentBlock, cdromBlock)
		}

		cdromMap[iface] = cdromBlock
	}

	if isClone && len(currentCDROM) == 0 {
		return nil
	}

	if len(currentCDROM) > 0 {
		interfaces := utils.ListResourcesAttributeValue(currentCDROM, MkCDROMInterface)
		cdromList := utils.OrderedListFromMapByKeyValues(cdromMap, interfaces)
		cdromList = slices.DeleteFunc(cdromList, func(v any) bool { return v == nil })

		for _, iface := range interfaces {
			delete(cdromMap, iface)
		}

		if len(cdromMap) > 0 {
			cdromList = append(cdromList, utils.OrderedListFromMap(cdromMap)...)
		}

		return cdromList
	}

	return utils.OrderedListFromMap(cdromMap)
}

// Read loads reconstructed CD-ROM state into the Terraform resource data.
func Read(d *schema.ResourceData, deviceObjects vms.CustomStorageDevices, isClone bool) diag.Diagnostics {
	currentCDROM := d.Get(MkCDROM).([]any)

	cdromState := BuildState(currentCDROM, deviceObjects, isClone)
	if cdromState == nil {
		return nil
	}

	return diag.FromErr(d.Set(MkCDROM, cdromState))
}

// Update diffs the prior and planned CD-ROM sets and records the needed storage device changes.
func Update(d *schema.ResourceData, updateBody *vms.UpdateRequestBody, del []string) []string {
	if !d.HasChange(MkCDROM) {
		return del
	}

	old, _ := d.GetChange(MkCDROM)
	oldCDROMs := GetCDROMDeviceObjects(old.([]any))
	newCDROMs := GetCDROMDeviceObjects(d.Get(MkCDROM).([]any))

	return applyDeviceObjectDiff(oldCDROMs, newCDROMs, updateBody, del)
}

func applyDeviceObjectDiff(
	oldCDROMs vms.CustomStorageDevices,
	newCDROMs vms.CustomStorageDevices,
	updateBody *vms.UpdateRequestBody,
	del []string,
) []string {
	toCreate, toUpdate, toDelete := utils.MapDiff(newCDROMs, oldCDROMs)

	for iface := range toDelete {
		del = append(del, iface)
	}

	for iface, dev := range toCreate {
		updateBody.AddCustomStorageDevice(iface, *dev)
	}

	for iface, dev := range toUpdate {
		updateBody.AddCustomStorageDevice(iface, *dev)
	}

	return del
}

func preserveLegacyEmptyFileID(currentBlock map[string]any, nextBlock map[string]any) {
	currentFileID, ok := currentBlock[MkCDROMFileID].(string)
	if !ok || currentFileID != "" {
		return
	}

	nextBlock[MkCDROMFileID] = ""
}

func isQ35MachineType(machineType string) bool {
	return strings.HasPrefix(machineType, "q35") || strings.HasPrefix(machineType, "pc-q35-")
}

func orderedKeys(deviceObjects vms.CustomStorageDevices) []string {
	keys := make([]string, 0, len(deviceObjects))

	for iface := range deviceObjects {
		keys = append(keys, iface)
	}

	sort.Strings(keys)

	return keys
}
