/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"errors"
	"fmt"
	"net/url"
)

// GetHosts retrieves the Hosts configuration for a node.
func (c *VirtualEnvironmentClient) GetHosts(nodeName string) (*VirtualEnvironmentHostsGetResponseData, error) {
	resBody := &VirtualEnvironmentHostsGetResponseBody{}
	err := c.DoRequest(hmGET, fmt.Sprintf("nodes/%s/hosts", url.PathEscape(nodeName)), nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// UpdateHosts updates the Hosts configuration for a node.
func (c *VirtualEnvironmentClient) UpdateHosts(nodeName string, d *VirtualEnvironmentHostsUpdateRequestBody) error {
	return c.DoRequest(hmPOST, fmt.Sprintf("nodes/%s/hosts", url.PathEscape(nodeName)), d, nil)
}
