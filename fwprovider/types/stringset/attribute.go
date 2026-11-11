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

	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
)

type ResourceAttributeOption func(*schema.SetAttribute)

// WithRequired makes the attribute required and rejects an empty set.
func WithRequired() ResourceAttributeOption {
	return func(attribute *schema.SetAttribute) {
		attribute.Required = true
		attribute.Optional = false
		attribute.Computed = false
		attribute.Validators = append(attribute.Validators, setvalidator.SizeAtLeast(1))
	}
}

// WithOptionalNonComputed makes the attribute optional without Computed, so removing it from the configuration plans
// a null value. An empty set is rejected: it would be sent as nothing and could not be told apart from a removal.
func WithOptionalNonComputed() ResourceAttributeOption {
	return func(attribute *schema.SetAttribute) {
		attribute.Optional = true
		attribute.Required = false
		attribute.Computed = false
		attribute.Validators = append(attribute.Validators, setvalidator.SizeAtLeast(1))
	}
}

// WithIPAddresses requires every element to be a plain IPv4 or IPv6 address, matching PVE's "ip-list" format.
func WithIPAddresses() ResourceAttributeOption {
	return func(attribute *schema.SetAttribute) {
		attribute.Validators = append(attribute.Validators, setvalidator.ValueStringsAre(validators.IPAddress()))
	}
}

func WithOptional() ResourceAttributeOption {
	return func(attribute *schema.SetAttribute) {
		attribute.Optional = true
		attribute.Required = false
		attribute.Computed = true
	}
}

// ResourceAttribute returns a resource schema attribute for string set.
func ResourceAttribute(desc, markdownDesc string, options ...ResourceAttributeOption) schema.SetAttribute {
	attribute := schema.SetAttribute{
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
			setvalidator.ValueStringsAre(
				stringvalidator.RegexMatches(
					regexp.MustCompile(`(.|\s)*\S(.|\s)*`),
					"must be a non-empty and non-whitespace string",
				),
				stringvalidator.LengthAtLeast(1),
			),
		},
	}

	for _, option := range options {
		option(&attribute)
	}

	return attribute
}

// DataSourceAttribute returns a computed-only data source schema attribute for string set.
// Use this for read-only output attributes in datasources.
func DataSourceAttribute(desc, markdownDesc string) schema.SetAttribute {
	return schema.SetAttribute{
		CustomType: Type{
			SetType: types.SetType{
				ElemType: types.StringType,
			},
		},
		Description:         desc,
		MarkdownDescription: markdownDesc,
		ElementType:         types.StringType,
		Computed:            true,
	}
}

// DataSourceFilterAttribute returns an optional data source schema attribute for string set.
// Use this for input filter attributes in datasources.
func DataSourceFilterAttribute(desc, markdownDesc string) schema.SetAttribute {
	return schema.SetAttribute{
		CustomType: Type{
			SetType: types.SetType{
				ElemType: types.StringType,
			},
		},
		Description:         desc,
		MarkdownDescription: markdownDesc,
		ElementType:         types.StringType,
		Optional:            true,
	}
}
