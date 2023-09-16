/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package validators

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// LanguageValidator returns a new validator for language codes.
func LanguageValidator() validator.String {
	return stringvalidator.OneOf([]string{
		`ca`, `da`, `de`, `en`, `es`, `eu`, `fa`, `fr`, `he`, `it`, `ja`, `nb`,
		`nn`, `pl`, `pt_BR`, `ru`, `sl`, `sv`, `tr`, `zh_CN`, `zh_TW`,
	}...)
}

// KeyboardLayoutValidator returns a new validator for keyboard layouts.
func KeyboardLayoutValidator() validator.String {
	return stringvalidator.OneOf([]string{
		`de`, `de-ch`, `da`, `en-gb`, `en-us`, `es`, `fi`, `fr`, `fr-be`, `fr-ca`, `fr-ch`,
		`hu`, `is`, `it`, `ja`, `lt`, `mk`, `nl`, `no`, `pl`, `pt`, `pt-br`, `sv`, `sl`, `tr`,
	}...)
}
