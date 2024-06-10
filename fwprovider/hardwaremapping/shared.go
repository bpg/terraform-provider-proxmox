/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package hardwaremapping

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

//nolint:gochecknoglobals
var (
	// dataSourceSchemaBaseAttrComment is the base comment attribute for a hardware mapping data source.
	// Note that the Proxmox VE API attribute is named "description", but we map it as a comment since this naming is
	// generally across the Proxmox VE web UI and API documentations. This still follows the [Terraform "best practices"]
	// as it improves the user experience by matching the field name to the naming used in the human-facing interfaces.
	//
	// [Terraform "best practices"]: https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resource-and-attribute-schema-should-closely-match-the-underlying-api
	dataSourceSchemaBaseAttrComment = datasourceschema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			stringvalidator.UTF8LengthAtLeast(1),
			stringvalidator.RegexMatches(regexp.MustCompile(`^\S|^$`), "must not start with whitespace"),
			stringvalidator.RegexMatches(regexp.MustCompile(`\S$|^$`), "must not end with whitespace"),
		},
	}

	// dataSourceSchemaBaseAttrComment is the base comment attribute for a hardware mapping resource.
	// Note that the Proxmox VE API attribute is named "description", but we map it as a comment since this naming is
	// generally across the Proxmox VE web UI and API documentations. This still follows the [Terraform "best practices"]
	// as it improves the user experience by matching the field name to the naming used in the human-facing interfaces.
	//
	// [Terraform "best practices"]: https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resource-and-attribute-schema-should-closely-match-the-underlying-api
	resourceSchemaBaseAttrComment = resourceschema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			stringvalidator.UTF8LengthAtLeast(1),
			stringvalidator.RegexMatches(regexp.MustCompile(`^\S|^$`), "must not start with whitespace"),
			stringvalidator.RegexMatches(regexp.MustCompile(`\S$|^$`), "must not end with whitespace"),
		},
	}
)
