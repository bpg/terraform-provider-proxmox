package vm

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
)

type vmModel struct {
	Description types.String   `tfsdk:"description"`
	ID          types.Int64    `tfsdk:"id"`
	Name        types.String   `tfsdk:"name"`
	NodeName    types.String   `tfsdk:"node_name"`
	Timeouts    timeouts.Value `tfsdk:"timeouts"`
}

func (m *vmModel) updateFromAPI(config vms.GetResponseData, status vms.GetStatusResponseData) error {
	if status.VMID == nil {
		return errors.New("VM ID is missing in status API response")
	}

	m.ID = types.Int64Value(int64(*status.VMID))

	// Optional fields can be removed from the model, use StringPointerValue to handle removal on nil
	m.Description = types.StringPointerValue(config.Description)
	m.Name = types.StringPointerValue(config.Name)

	return nil
}
