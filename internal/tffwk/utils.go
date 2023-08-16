/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tffwk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// IDAttribute generates an attribute definition suitable for the always-present `id` attribute.
func IDAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

// NewParseValidator creates a validator which uses a parsing function to validate a string. The function is expected
// to return a value of type `T` and an error. If the error is non-nil, the validator will fail. The `description`
// argument should contain a description of the validator's effect.
func NewParseValidator[T any](parseFunction func(string) (T, error), description string) validator.String {
	return &parseValidator[T]{
		parseFunction: parseFunction,
		description:   description,
	}
}

// parseValidator is a validator which uses a parsing function to validate a string.
type parseValidator[T any] struct {
	parseFunction func(string) (T, error)
	description   string
}

func (val *parseValidator[T]) Description(_ context.Context) string {
	return val.description
}

func (val *parseValidator[T]) MarkdownDescription(_ context.Context) string {
	return val.description
}

func (val *parseValidator[T]) ValidateString(
	ctx context.Context,
	request validator.StringRequest,
	response *validator.StringResponse,
) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue

	_, err := val.parseFunction(value.ValueString())
	if err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			val.Description(ctx),
			value.String(),
		))
	}
}
