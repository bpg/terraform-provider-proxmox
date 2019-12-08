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
	"net/url"
)

// VirtualEnvironmentAuthenticationResponseBody contains the body from an authentication response.
type VirtualEnvironmentAuthenticationResponseBody struct {
	Data *VirtualEnvironmentAuthenticationResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentAuthenticationResponseCapabilities contains the supported capabilities for a session.
type VirtualEnvironmentAuthenticationResponseCapabilities struct {
	Access          *CustomPrivileges `json:"access,omitempty"`
	Datacenter      *CustomPrivileges `json:"dc,omitempty"`
	Nodes           *CustomPrivileges `json:"nodes,omitempty"`
	Storage         *CustomPrivileges `json:"storage,omitempty"`
	VirtualMachines *CustomPrivileges `json:"vms,omitempty"`
}

// VirtualEnvironmentAuthenticationResponseData contains the data from an authentication response.
type VirtualEnvironmentAuthenticationResponseData struct {
	ClusterName         string                                                `json:"clustername,omitempty"`
	CSRFPreventionToken string                                                `json:"CSRFPreventionToken"`
	Capabilities        *VirtualEnvironmentAuthenticationResponseCapabilities `json:"cap,omitempty"`
	Ticket              string                                                `json:"ticket"`
	Username            string                                                `json:"username"`
}

// Authenticate authenticates against the specified endpoint.
func (c *VirtualEnvironmentClient) Authenticate(reset bool) error {
	if c.authenticationData != nil && !reset {
		return nil
	}

	body := bytes.NewBufferString(fmt.Sprintf("username=%s&password=%s", url.QueryEscape(c.Username), url.QueryEscape(c.Password)))
	req, err := http.NewRequest(hmPOST, fmt.Sprintf("%s/%s/access/ticket", c.Endpoint, basePathJSONAPI), body)

	if err != nil {
		return errors.New("Failed to create authentication request")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.httpClient.Do(req)

	if err != nil {
		return errors.New("Failed to retrieve authentication response")
	}

	err = c.ValidateResponseCode(res)

	if err != nil {
		return err
	}

	resBody := VirtualEnvironmentAuthenticationResponseBody{}
	err = json.NewDecoder(res.Body).Decode(&resBody)

	if err != nil {
		return errors.New("Failed to decode authentication response")
	}

	if resBody.Data == nil {
		return errors.New("The server did not include a data object in the authentication response")
	}

	if resBody.Data.CSRFPreventionToken == "" {
		return errors.New("The server did not include a CSRF prevention token in the authentication response")
	}

	if resBody.Data.Ticket == "" {
		return errors.New("The server did not include a ticket in the authentication response")
	}

	if resBody.Data.Username == "" {
		return errors.New("The server did not include the username in the authentication response")
	}

	c.authenticationData = resBody.Data

	return nil
}

// AuthenticateRequest adds authentication data to a new request.
func (c *VirtualEnvironmentClient) AuthenticateRequest(req *http.Request) error {
	err := c.Authenticate(false)

	if err != nil {
		return err
	}

	req.AddCookie(&http.Cookie{
		Name:  "PVEAuthCookie",
		Value: c.authenticationData.Ticket,
	})

	if req.Method != "GET" {
		req.Header.Add("CSRFPreventionToken", c.authenticationData.CSRFPreventionToken)
	}

	return nil
}
