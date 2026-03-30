/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cdrom

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
)

const (
	DefaultEnabled = false
	DefaultFileID  = "cdrom"

	MkCDROM          = "cdrom"
	MkCDROMEnabled   = "enabled"
	MkCDROMFileID    = "file_id"
	MkCDROMInterface = "interface"
)

func InterfaceValidator() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringMatch(
		regexp.MustCompile(`^(ide[0-3]|sata[0-5]|scsi([0-9]|1[0-3]))$`),
		"must be one of `ide[0-3]`, `sata[0-5]`, `scsi[0-13]`",
	))
}

func Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		MkCDROM: {
			Type:        schema.TypeList,
			Description: "The CDROM drive",
			Optional:    true,
			DiffSuppressFunc: structure.SuppressIfListsOfMapsAreEqualIgnoringOrderByKey(
				MkCDROMInterface,
			),
			DiffSuppressOnRefresh: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					MkCDROMEnabled: {
						Type:        schema.TypeBool,
						Description: "Whether to enable the CDROM drive",
						Optional:    true,
						Default:     DefaultEnabled,
						Deprecated: "Remove this attribute's configuration as it is no longer used and the attribute will " +
							"be removed in the next version of the provider. Set `file_id` to `none` to leave the CDROM drive empty.",
					},
					MkCDROMFileID: {
						Type:        schema.TypeString,
						Description: "The file id",
						Optional:    true,
						Default:     DefaultFileID,
						ValidateDiagFunc: validation.AnyDiag(
							validation.ToDiagFunc(validation.StringInSlice([]string{"none", "cdrom"}, false)),
							validators.FileID(),
						),
					},
					MkCDROMInterface: {
						Type:             schema.TypeString,
						Description:      "The CDROM interface",
						Required:         true,
						ValidateDiagFunc: InterfaceValidator(),
					},
				},
			},
			MinItems: 0,
		},
	}
}
