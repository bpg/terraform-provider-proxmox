package vm

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/types/stringset"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/vm/cpu"
)

// Model represents the VM model.
//
// Note: for computed fields / blocks we have to use an Object type (or an alias),
// or a custom type in order to hold an unknown value.
type Model struct {
	Description types.String `tfsdk:"description"`
	CPU         cpu.Value    `tfsdk:"cpu"`
	Clone       *struct {
		ID      types.Int64 `tfsdk:"id"`
		Retries types.Int64 `tfsdk:"retries"`
	} `tfsdk:"clone"`
	ID       types.Int64     `tfsdk:"id"`
	Name     types.String    `tfsdk:"name"`
	NodeName types.String    `tfsdk:"node_name"`
	Tags     stringset.Value `tfsdk:"tags"`
	Template types.Bool      `tfsdk:"template"`
	Timeouts timeouts.Value  `tfsdk:"timeouts"`
}
