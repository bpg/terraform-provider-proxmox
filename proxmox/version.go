/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package proxmox

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// Version retrieves the version information.
func (c *VirtualEnvironmentClient) Version(ctx context.Context) (*VersionResponseData, error) {
	resBody := &VersionResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, "version", nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("failed to get version information: %w", err)
	}

	if resBody.Data == nil {
		return nil, types.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}
