/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vnets

import "github.com/bpg/terraform-provider-proxmox/proxmox/types"

/*
VNETS

This part is related to the SDN component : VNETS
Based on docs :
https://pve.proxmox.com/pve-docs/chapter-pvesdn.html#pvesdn_config_vnet
https://pve.proxmox.com/pve-docs/api-viewer/index.html#/cluster/sdn/vnets
*/
type Vnet struct {
	Alias        *string           `json:"alias,omitempty"         url:"alias,omitempty"`
	IsolatePorts *types.CustomBool `json:"isolate-ports,omitempty" url:"isolate-ports,omitempty,int"`
	Tag          *int64            `json:"tag,omitempty"           url:"tag,omitempty"`
	Type         *string           `json:"type,omitempty"          url:"type,omitempty"`
	VlanAware    *types.CustomBool `json:"vlanaware,omitempty"     url:"vlanaware,omitempty,int"`
	Zone         *string           `json:"zone,omitempty"          url:"zone,omitempty"`
}

type VnetData struct {
	Vnet

	Pending *Vnet `json:"pending,omitempty" url:"pending,omitempty"`
}

type VnetCreate struct {
	Vnet

	ID string `json:"vnet" url:"vnet"`
}

type VnetUpdate struct {
	Vnet

	Delete []string `url:"delete,omitempty"`
}

type vnetResponse struct {
	Data *VnetData `json:"data"`
}

type vnetsResponse struct {
	Data *[]VnetData `json:"data"`
}
