/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabrics

/*
Fabric used to represent a Fabric in the API.

This part is related to the SDN component: Fabric
Based on docs:
  - https://pve.proxmox.com/pve-docs/chapter-pvesdn.html#pvesdn_config_fabrics
  - https://pve.proxmox.com/pve-docs/api-viewer/index.html#/cluster/sdn/fabrics
*/
type Fabric struct {
	ID         string  `json:"fabric"              url:"fabric"`
	Digest     *string `json:"digest,omitempty"        url:"digest,omitempty"`
	Protocol   *string `json:"protocol,omitempty"      url:"protocol,omitempty"`
	IPv6Prefix *string `json:"ip6_prefix,omitempty"   url:"ip6_prefix,omitempty"`
	IPv4Prefix *string `json:"ip_prefix,omitempty"    url:"ip_prefix,omitempty"`
	LockToken  *string `json:"lock_token,omitempty"   url:"lock_token,omitempty"`

	// OSPF
	Area *string `json:"area,omitempty"         url:"area,omitempty"`

	// OpenFabric
	CsnpInterval  *int64 `json:"csnp_interval,omitempty" url:"csnp_interval,omitempty"`
	HelloInterval *int64 `json:"hello_interval,omitempty" url:"hello_interval,omitempty"`
}

type FabricData struct {
	Fabric

	Pending *Fabric `json:"pending,omitempty" url:"pending,omitempty"`
}

type FabricCreate struct {
	Fabric
}

type FabricUpdate struct {
	Fabric

	Delete []string `url:"delete,omitempty"`
}

type fabricResponse struct {
	Data *FabricData `json:"data"`
}

type fabricsResponse struct {
	Data *[]FabricData `json:"data"`
}
