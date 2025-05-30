package zones

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// Client is a client for accessing the Proxmox SDN Zones API.
type Client struct {
	api.Client
}

// ExpandPath returns the API path for SDN zones.
func (c *Client) ExpandPath(path string) string {
	return fmt.Sprintf("cluster/sdn/zones/%s", path)
}
