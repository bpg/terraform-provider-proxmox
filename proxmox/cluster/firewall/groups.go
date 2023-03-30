/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
)

// CreateGroup create new security group.
func (c *Client) CreateGroup(ctx context.Context, d *firewall.GroupCreateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, "cluster/firewall/groups", d, nil)
	if err != nil {
		return fmt.Errorf("error creating security group: %w", err)
	}
	return nil
}

// GetGroupRules retrieve positions of defined security group rules.
func (c *Client) GetGroupRules(ctx context.Context, group string) ([]*firewall.GroupGetResponseData, error) {
	resBody := &firewall.GroupGetResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("cluster/firewall/groups/%s", url.PathEscape(group)),
		nil,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving security group '%s': %w", group, err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// ListGroups retrieve list of security groups.
func (c *Client) ListGroups(ctx context.Context) ([]*firewall.GroupListResponseData, error) {
	resBody := &firewall.GroupListResponseBody{}
	err := c.DoRequest(ctx, http.MethodGet, "cluster/firewall/groups", nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving security groups: %w", err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Group < resBody.Data[j].Group
	})

	return resBody.Data, nil
}

// UpdateGroup update security group.
func (c *Client) UpdateGroup(ctx context.Context, d *firewall.GroupUpdateRequestBody) error {
	err := c.DoRequest(
		ctx,
		http.MethodPost,
		"cluster/firewall/groups",
		d,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error updating security group: %w", err)
	}
	return nil
}

// DeleteGroup delete security group.
func (c *Client) DeleteGroup(ctx context.Context, group string) error {
	err := c.DoRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("cluster/firewall/groups/%s", url.PathEscape(group)),
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error deleting security group '%s': %w", group, err)
	}
	return nil
}
