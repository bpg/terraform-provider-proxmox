/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// VirtualEnvironmentVersionResponseBody contains the body from a version response.
type VirtualEnvironmentVersionResponseBody struct {
	Data *VirtualEnvironmentVersionResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentVersionResponseData contains the data from a version response.
type VirtualEnvironmentVersionResponseData struct {
	Keyboard     string `json:"keyboard"`
	Release      string `json:"release"`
	RepositoryID string `json:"repoid"`
	Version      string `json:"version"`
}

// Version retrieves the version information.
func (c *VirtualEnvironmentClient) Version() (*VirtualEnvironmentVersionResponseData, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/version", c.Endpoint, basePathJSONAPI), new(bytes.Buffer))

	if err != nil {
		return nil, errors.New("Failed to create version information request")
	}

	err = c.AuthenticateRequest(req)

	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, errors.New("Failed to perform version information request")
	}

	err = c.ValidateResponse(res)

	if err != nil {
		return nil, err
	}

	resBody := VirtualEnvironmentVersionResponseBody{}
	err = json.NewDecoder(res.Body).Decode(&resBody)

	if err != nil {
		return nil, errors.New("Failed to decode version information response")
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	return resBody.Data, nil
}
