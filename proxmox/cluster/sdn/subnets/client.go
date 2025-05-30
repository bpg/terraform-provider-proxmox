package subnets

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// Client is a client for accessing the Proxmox SDN VNETs API.
type Client struct {
	api.Client
}

// ExpandPath returns the API path for SDN VNETS.
func (c *Client) ExpandPath(vnet_id string, path string) string {
	return fmt.Sprintf("cluster/sdn/vnets/%s/subnets/%s", vnet_id, path)
}
