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
)

// SDNID returns a validator that checks if a string is a valid SDN ID (Zone, VNet, Subnet, etc).
func SDNID() []validator.String {
	return []validator.String{
		stringvalidator.RegexMatches(
			regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]*[A-Za-z0-9]$`),
			"must be a valid SDN identifier",
		),
		stringvalidator.LengthAtMost(8),
	}
}
