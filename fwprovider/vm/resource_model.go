package vm

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/tags"
)

type vmModel struct {
	Description types.String   `tfsdk:"description"`
	ID          types.Int64    `tfsdk:"id"`
	Name        types.String   `tfsdk:"name"`
	NodeName    types.String   `tfsdk:"node_name"`
	Tags        tags.Value     `tfsdk:"tags"`
	Timeouts    timeouts.Value `tfsdk:"timeouts"`
}
