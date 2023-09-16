/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package validators

import (
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// HAResourceIDValidator returns a new HA resource identifier validator.
func HAResourceIDValidator() validator.String {
	return NewParseValidator(types.ParseHAResourceID, "value must be a valid HA resource identifier")
}

// HAResourceStateValidator returns a new HA resource state validator.
func HAResourceStateValidator() validator.String {
	return NewParseValidator(types.ParseHAResourceState, "value must be a valid HA resource state")
}

// HAResourceTypeValidator returns a new HA resource type validator.
func HAResourceTypeValidator() validator.String {
	return NewParseValidator(types.ParseHAResourceType, "value must be a valid HA resource type")
}
