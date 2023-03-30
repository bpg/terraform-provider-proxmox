package cluster

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
	fw "github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type Client struct {
	types.Client
}

func (c *Client) Firewall() fw.API {
	return &firewall.Client{Client: c}
}
