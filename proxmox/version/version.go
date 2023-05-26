/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package version

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// Version retrieves the version information.
func (c *Client) Version(ctx context.Context) (*ResponseData, error) {
	resBody := &ResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, "version", nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("failed to get version information: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}
