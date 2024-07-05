/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cdrom

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

// Model represents the CD-ROM model.
type Model struct {
	FileID types.String `tfsdk:"file_id"`
}

func attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"file_id": types.StringType,
	}
}

func (m *Model) exportToCustomStorageDevice() vms.CustomStorageDevice {
	return vms.CustomStorageDevice{
		FileVolume: m.FileID.ValueString(),
		Media:      ptr.Ptr("cdrom"),
	}
}

func (m *Model) importFromCustomStorageDevice(d vms.CustomStorageDevice) {
	m.FileID = types.StringValue(d.FileVolume)
}
