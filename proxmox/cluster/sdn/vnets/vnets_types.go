/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vnets

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

/*
VNet used to represent a VNet in the API.

This part is related to the SDN component : VNETS
Based on docs :
https://pve.proxmox.com/pve-docs/chapter-pvesdn.html#pvesdn_config_vnet
https://pve.proxmox.com/pve-docs/api-viewer/index.html#/cluster/sdn/vnets
*/
type VNet struct {
	Alias        *string           `json:"alias,omitempty"         url:"alias,omitempty"`
	IsolatePorts *types.CustomBool `json:"isolate-ports,omitempty" url:"isolate-ports,omitempty,int"`
	Tag          *int64            `json:"tag,omitempty"           url:"tag,omitempty"`
	Type         *string           `json:"type,omitempty"          url:"type,omitempty"`
	VlanAware    *types.CustomBool `json:"vlanaware,omitempty"     url:"vlanaware,omitempty,int"`
	Zone         *string           `json:"zone,omitempty"          url:"zone,omitempty"`
}

type VNetData struct {
	VNet

	Pending *VNet `json:"pending,omitempty" url:"pending,omitempty"`
}

type VNetCreate struct {
	VNet

	ID string `json:"vnet" url:"vnet"`
}

type VNetUpdate struct {
	VNet

	Delete []string `url:"delete,omitempty"`
}

type vnetResponse struct {
	Data *VNetData `json:"data"`
}

type vnetsResponse struct {
	Data *[]VNetData `json:"data"`
}
