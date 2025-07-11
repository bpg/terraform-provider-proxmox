package subnets

import (
	"fmt"
)

/*
SUBNETS

This part is related to the SDN component : SubNets
Based on docs :
https://pve.proxmox.com/pve-docs/chapter-pvesdn.html#pvesdn_config_subnet
https://pve.proxmox.com/pve-docs/api-viewer/index.html#/cluster/sdn/vnets/{vnet}/subnets

Notes:
 1. The Type is once again defined as an enum type in the API docs but isn't referenced
    anywhere. Therefore no way to check what are allowed types. 'subnet' works
 2. Currently in the API there are Delete and Digest options which are not available
    in the UI so the choice was made to remove them temporary, waiting for a fix.
 3. It is also not really in the terraform spirit to update elements like this.
*/
type SubnetData struct {
	ID            string        `json:"subnet,omitempty"          url:"subnet,omitempty"`
	Type          *string       `json:"type,omitempty"            url:"type,omitempty"`
	Vnet          *string       `json:"vnet,omitempty"            url:"vnet,omitempty"`
	DHCPDNSServer *string       `json:"dhcp-dns-server,omitempty" url:"dhcp-dns-server,omitempty"`
	DHCPRange     DHCPRangeList `json:"dhcp-range,omitempty"      url:"dhcp-range,omitempty"`
	DNSZonePrefix *string       `json:"dnszoneprefix,omitempty"   url:"dnszoneprefix,omitempty"`
	Gateway       *string       `json:"gateway,omitempty"         url:"gateway,omitempty"`
	SNAT          *int64        `json:"snat,omitempty"            url:"snat,omitempty"`
}

type SubnetRequestData struct {
	EncodedSubnetData
	Delete []string `url:"delete,omitempty"`
}

type SubnetResponseBody struct {
	Data *SubnetData `json:"data"`
}

type SubnetsResponseBody struct {
	Data *[]SubnetData `json:"data"`
}

type DHCPRangeList []DHCPRangeEntry

type DHCPRangeEntry struct {
	StartAddress string `json:"start-address"`
	EndAddress   string `json:"end-address"`
}

/*
This structure had to be defined and added after realizing a weird behavior in Proxmox's API.
When creating or updating Subnets, the dhcpRange needs to be passed as string array.
But when reading (GET), it arrives as an array of JSON structures.
*/
type EncodedSubnetData struct {
	ID            string   `url:"subnet,omitempty"`
	Type          *string  `url:"type,omitempty"`
	Vnet          *string  `url:"vnet,omitempty"`
	DHCPDNSServer *string  `url:"dhcp-dns-server,omitempty"`
	DHCPRange     []string `url:"dhcp-range,omitempty"`
	DNSZonePrefix *string  `url:"dnszoneprefix,omitempty"`
	Gateway       *string  `url:"gateway,omitempty"`
	SNAT          *int64   `url:"snat,omitempty"`
}

func (s *SubnetData) ToEncoded() *EncodedSubnetData {
	encodedRanges := make([]string, 0, len(s.DHCPRange))
	for _, r := range s.DHCPRange {
		encodedRanges = append(encodedRanges, fmt.Sprintf("start-address=%s,end-address=%s", r.StartAddress, r.EndAddress))
	}

	return &EncodedSubnetData{
		ID:            s.ID,
		Type:          s.Type,
		Vnet:          s.Vnet,
		DHCPDNSServer: s.DHCPDNSServer,
		DHCPRange:     encodedRanges,
		DNSZonePrefix: s.DNSZonePrefix,
		Gateway:       s.Gateway,
		SNAT:          s.SNAT,
	}
}
