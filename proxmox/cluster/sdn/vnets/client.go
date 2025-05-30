package vnets

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// Client is a client for accessing the Proxmox SDN VNETs API.
type Client struct {
	api.Client
}

// ExpandPath returns the API path for SDN VNETS.
func (c *Client) ExpandPath(path string) string {
	return fmt.Sprintf("cluster/sdn/vnets/%s", path)
}

func (c *Client) ParentPath(parentId string) string {
	return fmt.Sprintf("cluster/sdn/zones/%s", parentId)
}
