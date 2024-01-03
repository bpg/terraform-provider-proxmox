/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// GetQueryURLMetadata retrieves the URL filename details for a node.
func (c *Client) GetQueryURLMetadata(
	ctx context.Context,
	d *QueryURLMetadataGetRequestBody,
) (*QueryURLMetadataGetResponseData, error) {
	resBody := &QueryURLMetadataGetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath("query-url-metadata"), d, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Query URL metadata configuration: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}
