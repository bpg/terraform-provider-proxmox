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

// CreateGroup creates an access group.
func (c *VirtualEnvironmentClient) CreateGroup(
	ctx context.Context,
	d *VirtualEnvironmentGroupCreateRequestBody,
) error {
	return c.DoRequest(ctx, HmPOST, "access/groups", d, nil)
}

// DeleteGroup deletes an access group.
func (c *VirtualEnvironmentClient) DeleteGroup(ctx context.Context, id string) error {
	return c.DoRequest(ctx, HmDELETE, fmt.Sprintf("access/groups/%s", url.PathEscape(id)), nil, nil)
}

// GetGroup retrieves an access group.
func (c *VirtualEnvironmentClient) GetGroup(
	ctx context.Context,
	id string,
) (*VirtualEnvironmentGroupGetResponseData, error) {
	resBody := &VirtualEnvironmentGroupGetResponseBody{}
	err := c.DoRequest(
		ctx,
		HmGET,
		fmt.Sprintf("access/groups/%s", url.PathEscape(id)),
		nil,
		resBody,
	)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	sort.Strings(resBody.Data.Members)

	return resBody.Data, nil
}

// ListGroups retrieves a list of access groups.
func (c *VirtualEnvironmentClient) ListGroups(
	ctx context.Context,
) ([]*VirtualEnvironmentGroupListResponseData, error) {
	resBody := &VirtualEnvironmentGroupListResponseBody{}
	err := c.DoRequest(ctx, HmGET, "access/groups", nil, resBody)
	if err != nil {
		return nil, err
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].ID < resBody.Data[j].ID
	})

	return resBody.Data, nil
}

// UpdateGroup updates an access group.
func (c *VirtualEnvironmentClient) UpdateGroup(
	ctx context.Context,
	id string,
	d *VirtualEnvironmentGroupUpdateRequestBody,
) error {
	return c.DoRequest(ctx, HmPUT, fmt.Sprintf("access/groups/%s", url.PathEscape(id)), d, nil)
}
