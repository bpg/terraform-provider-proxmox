/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package validators

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types/hardwaremapping"
)

const (
	// HardwareMappingDeviceIDValidatorErrMessage is the error message when the validation fails.
	HardwareMappingDeviceIDValidatorErrMessage = `value must be a valid hardware mapping device ID, e.g. "8086:5916"`
)

// HardwareMappingDeviceIDValidator validates a hardware mapping device ID.
func HardwareMappingDeviceIDValidator() validator.String {
	return NewParseValidator(proxmoxtypes.ParseDeviceID, HardwareMappingDeviceIDValidatorErrMessage)
}

// HardwareMappingMapCommentValidator validates that a hardware mapping map comment does not contain
// characters that are used as separators in the Proxmox VE property string format (comma and equals sign).
func HardwareMappingMapCommentValidator() validator.String {
	return stringvalidator.RegexMatches(
		regexp.MustCompile(`^[^,=]*$`),
		`must not contain commas "," or equals signs "=" as these are used as `+
			`separators in the Proxmox VE API property string format`,
	)
}
