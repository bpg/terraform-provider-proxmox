/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package replication

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/replications"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type model struct {
	ID       types.String  `tfsdk:"id"`
	Target   types.String  `tfsdk:"target"`
	Type     types.String  `tfsdk:"type"`
	Comment  types.String  `tfsdk:"comment"`
	Disable  types.Bool    `tfsdk:"disable"`
	Rate     types.Float64 `tfsdk:"rate"`
	Schedule types.String  `tfsdk:"schedule"`
	Source   types.String  `tfsdk:"source"`
	Guest    types.Int64   `tfsdk:"guest"`
	JobNum   types.Int64   `tfsdk:"jobnum"`
}

func (m *model) fromAPI(id string, data *replications.ReplicationData) {
	m.ID = types.StringValue(id)
	m.Target = types.StringValue(data.Target)
	m.Type = types.StringValue(data.Type)
	m.Comment = types.StringPointerValue(data.Comment)

	// API skips returning `disabled` even if created with `disabled=false`
	// so set disabled to false if not retured as true
	if v := data.Disable.PointerBool(); v != nil {
		m.Disable = types.BoolValue(*v)
	} else {
		m.Disable = types.BoolValue(false)
	}

	m.Rate = types.Float64PointerValue(data.Rate)
	m.Schedule = types.StringPointerValue(data.Schedule)
	m.Source = types.StringPointerValue(data.Source)
	m.Guest = types.Int64Value(data.Guest)
	m.JobNum = types.Int64Value(data.JobNum)
}

func (m *model) toAPICreate() *replications.ReplicationCreate {
	data := &replications.ReplicationCreate{}

	data.ID = m.ID.ValueString()
	data.Target = m.Target.ValueString()
	data.Type = m.Type.ValueString()

	if !m.Comment.IsUnknown() {
		data.Comment = m.Comment.ValueStringPointer()
	}

	data.Disable = (*proxmoxtypes.CustomBool)(m.Disable.ValueBoolPointer())

	if !m.Rate.IsUnknown() {
		data.Rate = m.Rate.ValueFloat64Pointer()
	}

	if !m.Schedule.IsUnknown() {
		data.Schedule = m.Schedule.ValueStringPointer()
	}

	return data
}

func (m *model) toAPIUpdate() *replications.ReplicationUpdate {
	data := &replications.ReplicationUpdate{}

	data.ID = m.ID.ValueString()

	if !m.Comment.IsUnknown() {
		data.Comment = m.Comment.ValueStringPointer()
	}

	data.Disable = (*proxmoxtypes.CustomBool)(m.Disable.ValueBoolPointer())

	if !m.Rate.IsUnknown() {
		data.Rate = m.Rate.ValueFloat64Pointer()
	}

	if !m.Schedule.IsUnknown() {
		data.Schedule = m.Schedule.ValueStringPointer()
	}

	return data
}

func (m *model) toAPIDelete() *replications.ReplicationDelete {
	data := &replications.ReplicationDelete{}
	return data
}
