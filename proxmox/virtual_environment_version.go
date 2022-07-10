/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"errors"
)

// Version retrieves the version information.
func (c *VirtualEnvironmentClient) Version() (*VirtualEnvironmentVersionResponseData, error) {
	resBody := &VirtualEnvironmentVersionResponseBody{}
	err := c.DoRequest(hmGET, "version", nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}
