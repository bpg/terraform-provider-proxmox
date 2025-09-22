/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package subnets

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

/*
Subnet used to represent a Subnet in the API.

This part is related to the SDN component: Subnet
Based on docs:
  - https://pve.proxmox.com/pve-docs/chapter-pvesdn.html#pvesdn_config_subnet
  - https://pve.proxmox.com/pve-docs/api-viewer/index.html#/cluster/sdn/vnets/{vnet}/subnets
*/
type Subnet struct {
	ID   string  `json:"subnet,omitempty" url:"subnet,omitempty"`
	VNet *string `json:"vnet,omitempty"   url:"vnet,omitempty"`

	DHCPDNSServer *string           `json:"dhcp-dns-server,omitempty" url:"dhcp-dns-server,omitempty"`
	DHCPRange     DHCPRange         `json:"dhcp-range,omitempty"      url:"dhcp-range,omitempty"`
	DNSZonePrefix *string           `json:"dnszoneprefix,omitempty"   url:"dnszoneprefix,omitempty"`
	Gateway       *string           `json:"gateway,omitempty"         url:"gateway,omitempty"`
	SNAT          *types.CustomBool `json:"snat,omitempty"            url:"snat,omitempty,int"`
	Type          *string           `json:"type,omitempty"            url:"type,omitempty"`
}

type SubnetData struct {
	Subnet

	Pending *Subnet `json:"pending,omitempty" url:"pending,omitempty"`
}

type SubnetCreate = Subnet

type SubnetUpdate struct {
	Subnet

	Delete []string `url:"delete,omitempty"`
}

type subnetResponse struct {
	Data *SubnetData `json:"data"`
}

type subnetsResponse struct {
	Data *[]SubnetData `json:"data"`
}

type DHCPRange []DHCPRangeEntry

type DHCPRangeEntry struct {
	StartAddress string `json:"start-address"`
	EndAddress   string `json:"end-address"`
}

// EncodeValues converts a DHCPRange struct to a URL value.
func (r DHCPRange) EncodeValues(key string, v *url.Values) error {
	if r == nil {
		return nil
	}

	encodedRanges := make([]string, 0, len(r))
	for _, entry := range r {
		encodedRanges = append(encodedRanges, fmt.Sprintf("start-address=%s,end-address=%s", entry.StartAddress, entry.EndAddress))
	}

	v.Add(key, strings.Join(encodedRanges, ","))

	return nil
}
