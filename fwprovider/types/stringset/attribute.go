/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package stringset

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceAttribute returns a resource schema attribute for string set.
func ResourceAttribute(desc, markdownDesc string) schema.SetAttribute {
	return schema.SetAttribute{
		CustomType: Type{
			SetType: types.SetType{
				ElemType: types.StringType,
			},
		},
		Description:         desc,
		MarkdownDescription: markdownDesc,
		Optional:            true,
		Computed:            true,
		ElementType:         types.StringType,
		Validators: []validator.Set{
			// NOTE: we allow empty list to remove all previously set values
			setvalidator.ValueStringsAre(
				stringvalidator.RegexMatches(
					regexp.MustCompile(`(.|\s)*\S(.|\s)*`),
					"must be a non-empty and non-whitespace string",
				),
				stringvalidator.LengthAtLeast(1),
			),
		},
	}
}

// DataSourceAttribute returns a data source schema attribute for string set.
func DataSourceAttribute(desc, markdownDesc string, optional bool) schema.SetAttribute {
	attribute := schema.SetAttribute{
		CustomType: Type{
			SetType: types.SetType{
				ElemType: types.StringType,
			},
		},
		Description:         desc,
		MarkdownDescription: markdownDesc,
		ElementType:         types.StringType,
	}

	if optional {
		attribute.Optional = true
	} else {
		attribute.Required = true
	}

	return attribute
}
