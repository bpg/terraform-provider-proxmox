/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package validators

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// AbsoluteFilePathValidator validates that a string is an absolute file path.
func AbsoluteFilePathValidator() validator.String {
	return NewParseValidator(
		func(s string) (string, error) {
			if strings.HasPrefix(s, "/") {
				return s, nil
			}

			return s, fmt.Errorf("%q is not an absolute path", s)
		},
		"must be an absolute file path",
	)
}

// NonEmptyString returns a new validator to ensure a non-empty string.
func NonEmptyString() validator.String {
	return stringvalidator.All(
		stringvalidator.UTF8LengthAtLeast(1),
		stringvalidator.RegexMatches(regexp.MustCompile(`^\S|^$`), "must not start with whitespace"),
		stringvalidator.RegexMatches(regexp.MustCompile(`\S$|^$`), "must not end with whitespace"),
	)
}
