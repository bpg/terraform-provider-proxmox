/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const basePathJSONAPI = "api2/json"

// VirtualEnvironmentClient implements an API client for the Proxmox Virtual Environment API.
type VirtualEnvironmentClient struct {
	Endpoint string
	Insecure bool
	Password string
	Username string

	authenticationData *VirtualEnvironmentAuthenticationResponseData
	httpClient         *http.Client
}

// NewVirtualEnvironmentClient creates and initializes a VirtualEnvironmentClient instance.
func NewVirtualEnvironmentClient(endpoint, username, password string, insecure bool) (*VirtualEnvironmentClient, error) {
	url, err := url.ParseRequestURI(endpoint)

	if err != nil {
		return nil, errors.New("You must specify a valid endpoint for the Proxmox Virtual Environment API (valid: https://host:port/)")
	}

	if url.Scheme != "https" {
		return nil, errors.New("You must specify a secure endpoint for the Proxmox Virtual Environment API (valid: https://host:port/)")
	}

	if password == "" {
		return nil, errors.New("You must specify a password for the Proxmox Virtual Environment API")
	}

	if username == "" {
		return nil, errors.New("You must specify a username for the Proxmox Virtual Environment API")
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		},
	}

	return &VirtualEnvironmentClient{
		Endpoint:   strings.TrimRight(url.String(), "/"),
		Insecure:   insecure,
		Password:   password,
		Username:   username,
		httpClient: httpClient,
	}, nil
}

// ValidateResponse ensures that a response is valid.
func (c *VirtualEnvironmentClient) ValidateResponse(res *http.Response) error {
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		switch res.StatusCode {
		case 401:
			return fmt.Errorf("Received a HTTP %d response - Please verify that the specified credentials are valid", res.StatusCode)
		case 403:
			return fmt.Errorf("Received a HTTP %d response - Please verify that the user account has the necessary permissions", res.StatusCode)
		case 404:
			return fmt.Errorf("Received a HTTP %d response - Please verify that the endpoint refers to a supported version of the Proxmox Virtual Environment API", res.StatusCode)
		case 500:
			return fmt.Errorf("Received a HTTP %d response - Please verify that Proxmox Virtual Environment is healthy", res.StatusCode)
		default:
			return fmt.Errorf("Received a HTTP %d response", res.StatusCode)
		}
	}

	return nil
}
