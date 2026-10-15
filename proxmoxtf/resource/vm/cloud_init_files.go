/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"github.com/hashicorp/go-cty/cty"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// vmCloudInitFilesExplicitlyCleared reports whether any of the cloud-init file ids in the first
// initialization block is set to an empty string in the raw config. These attributes are
// optional and computed, so on its own an omitted id does not clear anything; only an explicit empty
// string does.
func vmCloudInitFilesExplicitlyCleared(config cty.Value) bool {
	if config.IsNull() || !config.IsKnown() || !config.Type().IsObjectType() || !config.Type().HasAttribute(mkInitialization) {
		return false
	}

	blocks := config.GetAttr(mkInitialization)
	if !blocks.IsKnown() || blocks.IsNull() || blocks.LengthInt() == 0 {
		return false
	}

	block := blocks.AsValueSlice()[0]
	if !block.IsKnown() || block.IsNull() || !block.Type().IsObjectType() {
		return false
	}

	for _, key := range []string{
		mkInitializationUserDataFileID,
		mkInitializationVendorDataFileID,
		mkInitializationNetworkDataFileID,
		mkInitializationMetaDataFileID,
	} {
		if !block.Type().HasAttribute(key) {
			continue
		}

		v, _ := block.GetAttr(key).Unmark()
		if v.IsKnown() && !v.IsNull() && v.Type() == cty.String && v.AsString() == "" {
			return true
		}
	}

	return false
}

func vmCloudInitFilesHaveVolumes(files *vms.CustomCloudInitFiles) bool {
	return files != nil &&
		(files.UserVolume != nil || files.VendorVolume != nil || files.NetworkVolume != nil || files.MetaVolume != nil)
}

// vmCloneCicustomDeletes appends "cicustom" to the delete list for a cloned VM when the template has
// custom cloud-init files, the config explicitly clears at least one of them, and no file id is
// configured with a value.
//
// cicustom is a single PVE key, and the provider sends it as a whole. So an explicit clear of any id
// clears the whole inherited cicustom collection when no ids remain configured, including ids that were
// omitted from the config. This matches the existing behavior when only some ids are set: the new
// cicustom value replaces the template's one, and omitted inherited ids are dropped.
func vmCloneCicustomDeletes(
	del []string,
	templateFiles *vms.CustomCloudInitFiles,
	wantedFiles *vms.CustomCloudInitFiles,
	explicitlyCleared bool,
) []string {
	if !explicitlyCleared || !vmCloudInitFilesHaveVolumes(templateFiles) || vmCloudInitFilesHaveVolumes(wantedFiles) {
		return del
	}

	return append(del, "cicustom")
}
