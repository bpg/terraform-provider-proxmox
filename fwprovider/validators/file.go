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

// FileID returns a validator that checks if a string is a valid file ID.
func FileID() validator.String {
	return stringvalidator.RegexMatches(
		// PVE cannot parse ',' or '=' in a drive's volume, and cannot delete a volume with whitespace
		regexp.MustCompile(`^(?i)[a-z\d\-_.]+:([a-z\d\-_]+/)?[^\s,=]+$`),
		"must be in the format `<datastore name>:<content type>/<file name>` without whitespace, `,` or `=`",
	)
}

// DiskSerial returns validators for a disk serial number, limited to the characters PVE accepts.
func DiskSerial() []validator.String {
	return []validator.String{
		stringvalidator.LengthAtMost(20),
		stringvalidator.RegexMatches(
			regexp.MustCompile(`^[-%A-Za-z0-9_.!~*'()]*$`),
			"must contain only letters, digits and `-%_.!~*'()`",
		),
	}
}

// ImportFrom returns a validator for a disk import source: a file ID or an absolute path, as PVE accepts both.
func ImportFrom() validator.String {
	return stringvalidator.RegexMatches(
		regexp.MustCompile(`^(?i)([a-z\d\-_.]+:([a-z\d\-_]+/)?|/)[^\s,=]+$`),
		"must be a file ID (`<datastore name>:<content type>/<file name>`) or an absolute path, "+
			"without whitespace, `,` or `=`",
	)
}
