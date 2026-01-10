/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package storage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// DownloadOCIImageByReference downloads the OCI image by its reference.
func (c *Client) DownloadOCIImageByReference(
	ctx context.Context,
	d *OCIRegistryPullRequestBody,
) error {
	resBody := &OCIRegistryPullResponseBody{}

	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("oci-registry-pull"), d, resBody)
	if err != nil {
		return fmt.Errorf("error download OCI image by reference: %w", err)
	}

	if resBody.TaskID == nil {
		return api.ErrNoDataObjectInResponse
	}

	taskErr := c.Tasks().WaitForTask(ctx, *resBody.TaskID)
	if taskErr != nil {
		err = fmt.Errorf(
			"error download OCI image to datastore %s: failed waiting for OCI image download: %w",
			c.StorageName,
			taskErr,
		)

		deleteErr := c.Tasks().DeleteTask(context.WithoutCancel(ctx), *resBody.TaskID)
		if deleteErr != nil {
			return fmt.Errorf("%w \n %w", err, deleteErr)
		}

		return err
	}

	return nil
}
