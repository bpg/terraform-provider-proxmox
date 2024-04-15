/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ha

import (
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// resourceIDValidator returns a new HA resource identifier validator.
func resourceIDValidator() validator.String {
	return validators.NewParseValidator(types.ParseHAResourceID, "value must be a valid HA resource identifier")
}

// resourceStateValidator returns a new HA resource state validator.
func resourceStateValidator() validator.String {
	return validators.NewParseValidator(types.ParseHAResourceState, "value must be a valid HA resource state")
}

// resourceTypeValidator returns a new HA resource type validator.
func resourceTypeValidator() validator.String {
	return validators.NewParseValidator(types.ParseHAResourceType, "value must be a valid HA resource type")
}
