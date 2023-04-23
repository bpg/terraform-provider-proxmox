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
	"sort"
)

// GetACL retrieves the access control list.
func (c *VirtualEnvironmentClient) GetACL(ctx context.Context) ([]*ACLGetResponseData, error) {
	resBody := &ACLGetResponseBody{}
	err := c.DoRequest(ctx, http.MethodGet, "access/acl", nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("failed to get access control list: %w", err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Path < resBody.Data[j].Path
	})

	return resBody.Data, nil
}

// UpdateACL updates the access control list.
func (c *VirtualEnvironmentClient) UpdateACL(ctx context.Context, d *ACLUpdateRequestBody) error {
	return c.DoRequest(ctx, http.MethodPut, "access/acl", d, nil)
}
