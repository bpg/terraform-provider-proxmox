package vm

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

type vmModel struct {
	// ID is the identifier used by Terraform.
	// Other fields are sorted alphabetically.
	ID types.Int64 `tfsdk:"id"`
	// Timeouts are the timeouts for the resource, defined in terraform-plugin-framework-timeouts
	Timeouts timeouts.Value `tfsdk:"timeouts"`

	Description types.String `tfsdk:"description"`
	Name        types.String `tfsdk:"name"`
	NodeName    types.String `tfsdk:"node_name"`
	// VMID        types.Int64  `tfsdk:"vm_id"`
}

func (m *vmModel) updateFromAPI(config vms.GetResponseData, status vms.GetStatusResponseData) error {
	if status.VMID == nil {
		return errors.New("VM ID is missing in status API response")
	}

	// m.VMID = types.Int64Value(int64(*status.VMID))
	m.ID = types.Int64Value(int64(*status.VMID))

	// Optional fields can be removed from the model, use StringPointerValue to handle removal on nil

	m.Description = types.StringPointerValue(config.Description)
	m.Name = types.StringPointerValue(config.Name)

	return nil
}
