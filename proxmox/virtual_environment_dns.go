/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// GetDNS retrieves the DNS configuration for a node.
func (c *VirtualEnvironmentClient) GetDNS(
	ctx context.Context,
	nodeName string,
) (*VirtualEnvironmentDNSGetResponseData, error) {
	resBody := &VirtualEnvironmentDNSGetResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("nodes/%s/dns", url.PathEscape(nodeName)),
		nil,
		resBody,
	)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// UpdateDNS updates the DNS configuration for a node.
func (c *VirtualEnvironmentClient) UpdateDNS(
	ctx context.Context,
	nodeName string,
	d *VirtualEnvironmentDNSUpdateRequestBody,
) error {
	return c.DoRequest(ctx, http.MethodPut, fmt.Sprintf("nodes/%s/dns", url.PathEscape(nodeName)), d, nil)
}
