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

// CDROMInterface returns a validator that checks if a string is a valid CD-ROM interface.
func CDROMInterface() validator.String {
	return stringvalidator.RegexMatches(
		regexp.MustCompile(`^ide[0-3]|sata[0-5]|scsi(?:30|[12][0-9]|[0-9])$`),
		"one of `ide[0-3]`, `sata[0-5]`, `scsi[0-30]`",
	)
}
