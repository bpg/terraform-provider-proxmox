/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodestorage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// DownloadFileByURL downloads the file using URL.
func (c *Client) DownloadFileByURL(
	ctx context.Context,
	d *DownloadURLPostRequestBody,
	uploadTimeout int,
) (*DownloadURLResponseBody, error) {
	resBody := &DownloadURLResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath("download-url"), d, resBody)
	if err != nil {
		return nil, fmt.Errorf("error download file by URL: %w", err)
	}

	if resBody.TaskID == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	err = c.Tasks().WaitForTask(ctx, *resBody.TaskID, uploadTimeout, 5)
	if err != nil {
		return nil, fmt.Errorf(
			"error download file to datastore %s: failed waiting for url download - %w",
			c.StorageName,
			err,
		)
	}

	return resBody, nil
}
