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
			// NOTE: we allow empty list to remove all previously set tags
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
