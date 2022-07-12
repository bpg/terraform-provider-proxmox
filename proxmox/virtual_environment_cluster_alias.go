/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sort"
)

// CreateAlias create an alias
func (c *VirtualEnvironmentClient) CreateAlias(ctx context.Context, d *VirtualEnvironmentClusterAliasCreateRequestBody) error {
	return c.DoRequest(ctx, hmPOST, "cluster/firewall/aliases", d, nil)
}

// DeleteAlias delete an alias
func (c *VirtualEnvironmentClient) DeleteAlias(ctx context.Context, id string) error {
	return c.DoRequest(ctx, hmDELETE, fmt.Sprintf("cluster/firewall/aliases/%s", url.PathEscape(id)), nil, nil)
}

// GetAlias retrieves an alias
func (c *VirtualEnvironmentClient) GetAlias(ctx context.Context, id string) (*VirtualEnvironmentClusterAliasGetResponseData, error) {
	resBody := &VirtualEnvironmentClusterAliasGetResponseBody{}
	err := c.DoRequest(ctx, hmGET, fmt.Sprintf("cluster/firewall/aliases/%s", url.PathEscape(id)), nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// ListAliases retrieves a list of aliases.
func (c *VirtualEnvironmentClient) ListAliases(ctx context.Context) ([]*VirtualEnvironmentClusterAliasGetResponseData, error) {
	resBody := &VirtualEnvironmentClusterAliasListResponseBody{}
	err := c.DoRequest(ctx, hmGET, "cluster/firewall/aliases", nil, resBody)

	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Name < resBody.Data[j].Name
	})

	return resBody.Data, nil
}

// UpdateAlias updates an alias.
func (c *VirtualEnvironmentClient) UpdateAlias(ctx context.Context, id string, d *VirtualEnvironmentClusterAliasUpdateRequestBody) error {
	return c.DoRequest(ctx, hmPUT, fmt.Sprintf("cluster/firewall/aliases/%s", url.PathEscape(id)), d, nil)
}
