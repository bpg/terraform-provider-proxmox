/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"io"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
	"github.com/bpg/terraform-provider-proxmox/proxmox/container"
	"github.com/bpg/terraform-provider-proxmox/proxmox/vm"
)

const (
	basePathJSONAPI = "api2/json"
)

// VirtualEnvironmentClient implements an API client for the Proxmox Virtual Environment API.
type VirtualEnvironmentClient struct {
	Endpoint string
	Insecure bool
	OTP      *string
	Password string
	Username string
	Agent bool

	authenticationData *VirtualEnvironmentAuthenticationResponseData
	httpClient         *http.Client
}

// VirtualEnvironmentErrorResponseBody contains the body of an error response.
type VirtualEnvironmentErrorResponseBody struct {
	Data   *string
	Errors *map[string]string
}

// VirtualEnvironmentMultiPartData enables multipart uploads in DoRequest.
type VirtualEnvironmentMultiPartData struct {
	Boundary string
	Reader   io.Reader
	Size     *int64
}

type API interface {
	Cluster() *cluster.Client
	VM(nodeName string, vmID int) *vm.Client
	Container(nodeName string, vmID int) *container.Client
}

func (c *VirtualEnvironmentClient) API() API {
	return &client{c}
}

func (c *VirtualEnvironmentClient) ExpandPath(path string) string {
	return path
}

type client struct {
	c *VirtualEnvironmentClient
}

func (c *client) Cluster() *cluster.Client {
	return &cluster.Client{Client: c.c}
}

func (c *client) VM(nodeName string, vmID int) *vm.Client {
	return &vm.Client{Client: c.c, NodeName: nodeName, VMID: vmID}
}

func (c *client) Container(nodeName string, vmID int) *container.Client {
	return &container.Client{Client: c.c, NodeName: nodeName, VMID: vmID}
}
