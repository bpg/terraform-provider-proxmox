package vnets

/*
--------------------------------- VNETS ---------------------------------

This part is related to the SDN component : VNETS
Based on docs :
https://pve.proxmox.com/pve-docs/chapter-pvesdn.html#pvesdn_config_vnet
https://pve.proxmox.com/pve-docs/api-viewer/index.html#/cluster/sdn/vnets

Notes:

 1. IsolatePorts is a boolean in the docs but needs to be passed as 0 or 1
    and is therefore defined as int.

 2. Type field can be 'vnet' but other values are unknown

 3. Tag cannot be set on Vnets created in simple Zones, might actually be
    only usable on vlan or vxlan zones as it sets the vlan or vxlan id.

 4. Currently in the API there are Delete and Digest options which are not available
    in the UI so the choice was made to remove them temporary, waiting for a fix.

-------------------------------------------------------------------------
*/
type VnetData struct {
	ID           string  `json:"vnet,omitempty"              url:"vnet,omitempty"`
	Zone         *string `json:"zone,omitempty"              url:"zone,omitempty"`
	Alias        *string `json:"alias,omitempty"             url:"alias,omitempty"`
	IsolatePorts *int64  `json:"isolate-ports,omitempty"     url:"isolate-ports,omitempty"`
	Tag          *int64  `json:"tag,omitempty"               url:"tag,omitempty"`
	Type         *string `json:"type,omitempty"              url:"type,omitempty"`
	VlanAware    *int64  `json:"vlanaware,omitempty"         url:"vlanaware,omitempty"`
	// DeleteSettings *string `json:"delete,omitempty"            url:"delete,omitempty"`
	// Digest         *string `json:"digest,omitempty"            url:"digest,omitempty"`
}

type VnetRequestData struct {
	VnetData
	Delete []string `url:"delete,omitempty"`
}

type VnetResponseBody struct {
	Data *VnetData `json:"data"`
}

type VnetsResponseBody struct {
	Data *[]VnetData `json:"data"`
}
