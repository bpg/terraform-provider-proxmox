/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	basePathJSONAPI = "api2/json"
	hmDELETE        = "DELETE"
	hmGET           = "GET"
	hmHEAD          = "HEAD"
	hmPOST          = "POST"
	hmPUT           = "PUT"
)

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

// DoRequest performs a HTTP request against a JSON API endpoint.
func (c *VirtualEnvironmentClient) DoRequest(method, path string, requestBody interface{}, responseBody interface{}) error {
	log.Printf("[DEBUG] Performing HTTP %s request (path: %s)", method, path)

	modifiedPath := path
	urlEncodedRequestBody := new(bytes.Buffer)

	if requestBody != nil {
		v, err := query.Values(requestBody)

		if err != nil {
			return fmt.Errorf("Failed to encode HTTP %s request (path: %s) - Reason: %s", method, path, err.Error())
		}

		encodedValues := v.Encode()

		if method == hmGET || method == hmHEAD {
			if !strings.Contains(modifiedPath, "?") {
				modifiedPath = fmt.Sprintf("%s?%s", modifiedPath, encodedValues)
			} else {
				modifiedPath = fmt.Sprintf("%s&%s", modifiedPath, encodedValues)
			}
		} else {
			urlEncodedRequestBody = bytes.NewBufferString(encodedValues)
		}

		log.Printf("[DEBUG] Added request body to HTTP %s request (path: %s) - Body: %s", method, path, encodedValues)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s/%s", c.Endpoint, basePathJSONAPI, modifiedPath), urlEncodedRequestBody)

	if err != nil {
		return fmt.Errorf("Failed to create HTTP %s request (path: %s) - Reason: %s", method, path, err.Error())
	}

	req.Header.Add("Accept", "application/json")

	if req.Method != hmGET && req.Method != hmHEAD {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	err = c.AuthenticateRequest(req)

	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)

	if err != nil {
		return fmt.Errorf("Failed to perform HTTP %s request (path: %s) - Reason: %s", method, path, err.Error())
	}

	err = c.ValidateResponseCode(res)

	if err != nil {
		return err
	}

	if responseBody != nil {
		err = json.NewDecoder(res.Body).Decode(responseBody)

		if err != nil {
			return fmt.Errorf("Failed to decode HTTP %s response (path: %s) - Reason: %s", method, path, err.Error())
		}
	}

	return nil
}

// ValidateResponseCode ensures that a response is valid.
func (c *VirtualEnvironmentClient) ValidateResponseCode(res *http.Response) error {
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		switch res.StatusCode {
		case 400:
			return fmt.Errorf("Received a HTTP %d response - This is most likely caused by a bug in the code, so please create a new issue on https://github.com/danitso/terraform-provider-proxmox/issues", res.StatusCode)
		case 401:
			return fmt.Errorf("Received a HTTP %d response - Please verify that the specified credentials are valid", res.StatusCode)
		case 403:
			return fmt.Errorf("Received a HTTP %d response - Please verify that the user account has the necessary permissions", res.StatusCode)
		case 404:
			return fmt.Errorf("Received a HTTP %d response - Please verify that the endpoint refers to a supported version of the Proxmox Virtual Environment API", res.StatusCode)
		case 500, 501, 502, 503:
			return fmt.Errorf("Received a HTTP %d response - Please verify that the Proxmox Virtual Environment API is healthy (status: %s)", res.StatusCode, res.Status)
		default:
			return fmt.Errorf("Received a HTTP %d response", res.StatusCode)
		}
	}

	return nil
}
