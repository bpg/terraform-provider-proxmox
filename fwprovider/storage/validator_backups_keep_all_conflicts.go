/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type backupsKeepAllExcludesOtherKeepSettingsValidator struct{}

func (v backupsKeepAllExcludesOtherKeepSettingsValidator) Description(context.Context) string {
	return "when keep_all is true, other keep_* attributes must not be set"
}

func (v backupsKeepAllExcludesOtherKeepSettingsValidator) MarkdownDescription(context.Context) string {
	return "when `keep_all` is true, other `keep_*` attributes must not be set"
}

func (v backupsKeepAllExcludesOtherKeepSettingsValidator) ValidateObject(
	_ context.Context,
	req validator.ObjectRequest,
	resp *validator.ObjectResponse,
) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	attrs := req.ConfigValue.Attributes()

	keepAllAttr, ok := attrs["keep_all"]
	if !ok {
		return
	}

	keepAll, ok := keepAllAttr.(basetypes.BoolValue)
	if !ok || keepAll.IsNull() || keepAll.IsUnknown() || !keepAll.ValueBool() {
		return
	}

	conflictingNames := []string{
		"keep_last",
		"keep_hourly",
		"keep_daily",
		"keep_weekly",
		"keep_monthly",
		"keep_yearly",
	}

	setNames := make([]string, 0, len(conflictingNames))
	for _, name := range conflictingNames {
		attrValue, exists := attrs[name]
		if !exists {
			continue
		}

		intValue, ok := attrValue.(basetypes.Int64Value)
		if !ok || intValue.IsNull() || intValue.IsUnknown() {
			continue
		}

		setNames = append(setNames, name)
	}

	if len(setNames) == 0 {
		return
	}

	summary := "invalid backup retention settings"
	detail := fmt.Sprintf(
		"when keep_all is true, these attributes must not be set: %s",
		strings.Join(setNames, ", "),
	)

	for _, name := range setNames {
		resp.Diagnostics.AddAttributeError(req.Path.AtName(name), summary, detail)
	}
}
