package vga

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Model represents the VGA model.
type Model struct {
	Clipboard types.String `tfsdk:"clipboard"`
	Type      types.String `tfsdk:"type"`
	Memory    types.Int64  `tfsdk:"memory"`
}

func attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"clipboard": types.StringType,
		"type":      types.StringType,
		"memory":    types.Int64Type,
	}
}
