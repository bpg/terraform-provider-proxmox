/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"errors"
)

// VirtualEnvironmentAccessGroupListResponseBody contains the body from an access group list response.
type VirtualEnvironmentAccessGroupListResponseBody struct {
	Data []*VirtualEnvironmentAccessGroupListResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentAccessGroupListResponseData contains the data from an access group list response.
type VirtualEnvironmentAccessGroupListResponseData struct {
	Comment string `json:"comment"`
	ID      string `json:"groupid"`
}

// ListAccessGroups retrieves a list of access groups.
func (c *VirtualEnvironmentClient) ListAccessGroups() ([]*VirtualEnvironmentAccessGroupListResponseData, error) {
	resBody := &VirtualEnvironmentAccessGroupListResponseBody{}
	err := c.DoRequest(hmGET, "access/groups", nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("The server did not include a data object in the response")
	}

	return resBody.Data, nil
}
