/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func (c *Client) aclPath() string {
	return c.ExpandPath("acl")
}

// GetACL retrieves the access control list.
func (c *Client) GetACL(ctx context.Context) ([]*ACLGetResponseData, error) {
	resBody := &ACLGetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.aclPath(), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("failed to get access control list: %w", err)
	}

	if resBody.Data == nil {
		return nil, types.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Path < resBody.Data[j].Path
	})

	return resBody.Data, nil
}

// UpdateACL updates the access control list.
func (c *Client) UpdateACL(ctx context.Context, d *ACLUpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.aclPath(), d, nil)
	if err != nil {
		return fmt.Errorf("failed to update access control list: %w", err)
	}

	return nil
}
