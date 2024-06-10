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
		regexp.MustCompile(`^(?i)[a-z\d\-_.]+:([a-z\d\-_]+/)?.+$`),
		"must be in the format `<datastore name>:<content type>/<file name>`",
	)
}
