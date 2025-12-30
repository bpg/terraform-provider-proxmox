/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fabric_nodes

/*
Fabric used to represent a Fabric in the API.

This part is related to the SDN component: Fabric
Based on docs:
  - https://pve.proxmox.com/pve-docs/chapter-pvesdn.html#pvesdn_config_fabrics
  - https://pve.proxmox.com/pve-docs/api-viewer/index.html#/cluster/sdn/fabrics
*/
type FabricNode struct {
	FabricID    string   `json:"fabric_id"      url:"fabric_id"`
	NodeID      string   `json:"node_id"              url:"node_id"`
	Digest      *string  `json:"digest,omitempty"        url:"digest,omitempty"`
	Protocol    *string  `json:"protocol,omitempty"      url:"protocol,omitempty"`
	Interfaces  []string `json:"interfaces,omitempty"   url:"interfaces,omitempty"`
	IPv4Address *string  `json:"ip,omitempty"    url:"ip,omitempty"`
	IPv6Address *string  `json:"ip6,omitempty"   url:"ip6,omitempty"`
	LockToken   *string  `json:"lock_token,omitempty"   url:"lock_token,omitempty"`
}

type FabricNodeData struct {
	FabricNode
}

type FabricNodeCreate struct {
	FabricNode
}

type FabricNodeUpdate struct {
	FabricNode

	Delete []string `url:"delete,omitempty"`
}
