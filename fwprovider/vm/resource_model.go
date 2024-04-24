package vm

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type vmModel struct {
	Description types.String `tfsdk:"description"`
	Clone       *struct {
		ID      types.Int64 `tfsdk:"id"`
		Retries types.Int64 `tfsdk:"retries"`
	} `tfsdk:"clone"`
	ID       types.Int64  `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	NodeName types.String `tfsdk:"node_name"`
	Tags     types.String `tfsdk:"tags"`
	//Tags     tags.Value     `tfsdk:"tags"`
	Template types.Bool     `tfsdk:"template"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
