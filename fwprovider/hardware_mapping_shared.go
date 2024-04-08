package fwprovider

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

//nolint:gochecknoglobals
var (
	// hardwareMappingDataSourceSchemaWithBaseAttrComment is the base comment attribute for a hardware mapping data
	// source.
	// Note that the Proxmox VE API attribute is named "description", but we map it as a comment since this naming is
	// generally across the Proxmox VE
	// web UI and API documentations. This still follows the Terraform "best practices" as it improves the user experience
	// by matching the field name to
	// the naming used in the human-facing interfaces.
	// References:
	//   1. https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resource-and-attribute-schema-should-closely-match-the-underlying-api
	hardwareMappingDataSourceSchemaWithBaseAttrComment = datasourceschema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			stringvalidator.UTF8LengthAtLeast(1),
			stringvalidator.RegexMatches(regexp.MustCompile(`^\S|^$`), "must not start with whitespace"),
			stringvalidator.RegexMatches(regexp.MustCompile(`\S$|^$`), "must not end with whitespace"),
		},
	}

	// hardwareMappingDataSourceSchemaWithBaseAttrComment is the base comment attribute for a hardware mapping resource.
	// Note that the Proxmox VE API attribute is named "description", but we map it as a comment since this naming is
	// generally across the Proxmox VE
	// web UI and API documentations. This still follows the Terraform "best practices" as it improves the user experience
	// by matching the field name to
	// the naming used in the human-facing interfaces.
	// References:
	//   1. https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles#resource-and-attribute-schema-should-closely-match-the-underlying-api
	hardwareMappingResourceSchemaWithBaseAttrComment = resourceschema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			stringvalidator.UTF8LengthAtLeast(1),
			stringvalidator.RegexMatches(regexp.MustCompile(`^\S|^$`), "must not start with whitespace"),
			stringvalidator.RegexMatches(regexp.MustCompile(`\S$|^$`), "must not end with whitespace"),
		},
	}
)
