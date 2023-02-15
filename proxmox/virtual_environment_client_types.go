/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"io"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster"
)

const (
	basePathJSONAPI = "api2/json"
	HmDELETE        = "DELETE"
	HmGET           = "GET"
	HmHEAD          = "HEAD"
	HmPOST          = "POST"
	HmPUT           = "PUT"
)

// VirtualEnvironmentClient implements an API client for the Proxmox Virtual Environment API.
type VirtualEnvironmentClient struct {
	Endpoint string
	Insecure bool
	OTP      *string
	Password string
	Username string

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
	Cluster() *cluster.API
}

func (c *VirtualEnvironmentClient) API() API {
	return &api{c}
}

type api struct {
	c *VirtualEnvironmentClient
}

func (a *api) Cluster() *cluster.API {
	return &cluster.API{Client: a.c}
}
