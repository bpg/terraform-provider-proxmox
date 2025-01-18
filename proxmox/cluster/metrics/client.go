package metrics

import (
	"fmt"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// Client is an interface for accessing the Proxmox metrics management API.
type Client struct {
	api.Client
}

// ExpandPath expands a relative path to the Proxmox metrics server management API path.
func (c *Client) ExpandPath(path string) string {
	return fmt.Sprintf("cluster/metrics/server/%s", path)
}
