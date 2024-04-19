/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package validators

import (
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
