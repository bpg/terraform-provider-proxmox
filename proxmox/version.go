/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package proxmox

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// Version retrieves the version information.
func (c *VirtualEnvironmentClient) Version(
	ctx context.Context,
) (*VersionResponseData, error) {
	resBody := &VersionResponseBody{}
	err := c.DoRequest(ctx, http.MethodGet, "version", nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("failed to get version information: %w", err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}
