/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package mds

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	mdsapi "github.com/bpg/terraform-provider-proxmox/proxmox/nodes/ceph/mds"
)

// cephMDSModel maps the schema data for proxmox_ceph_mds.
type cephMDSModel struct {
	ID       types.String `tfsdk:"id"`
	NodeName types.String `tfsdk:"node_name"`
	Name     types.String `tfsdk:"name"`
}

// fromAPI fills the model from a list entry. node_name follows the host Proxmox reports, so an import naming the
// wrong node is corrected; deleting through the wrong node would leave the daemon's unit and directory behind.
func (m *cephMDSModel) fromAPI(d *mdsapi.ListResponseData) {
	if d.Host != nil && *d.Host != "" {
		m.NodeName = types.StringValue(*d.Host)
	}

	m.Name = types.StringValue(d.Name)
	m.ID = types.StringValue(d.Name)
}
