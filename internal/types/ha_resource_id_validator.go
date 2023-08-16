/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package types

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type haResourceIDValidator struct{}

var _ validator.String = &haResourceIDValidator{}

// HAResourceIDValidator checks that the String held in the attribute is a valid Proxmox HA resource identifier.
func HAResourceIDValidator() validator.String {
	return &haResourceIDValidator{}
}

func (v *haResourceIDValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v *haResourceIDValidator) MarkdownDescription(_ context.Context) string {
	return "value must be a valid Proxmox HA resource identifier"
}

func (v *haResourceIDValidator) ValidateString(
	ctx context.Context,
	request validator.StringRequest,
	response *validator.StringResponse,
) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue

	_, err := ParseHAResourceID(value.ValueString())
	if err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value.String(),
		))
	}
}
