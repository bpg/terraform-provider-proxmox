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
)

type SecurityGroup interface {
	CreateGroup(ctx context.Context, d *GroupCreateRequestBody) error
	ListGroups(ctx context.Context) ([]*GroupListResponseData, error)
	UpdateGroup(ctx context.Context, d *GroupUpdateRequestBody) error
	DeleteGroup(ctx context.Context, group string) error
}

// GroupCreateRequestBody contains the data for a security group create request.
type GroupCreateRequestBody struct {
	Group   string  `json:"group"             url:"group"`
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	Digest  *string `json:"digest,omitempty"  url:"digest,omitempty"`
}

// GroupListResponseData contains the data from a group list response.
type GroupListResponseData struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	Group   string  `json:"group"             url:"group"`
	Digest  string  `json:"digest"            url:"digest"`
}

// GroupListResponseBody contains the data from a group get response.
type GroupListResponseBody struct {
	Data []*GroupListResponseData `json:"data,omitempty"`
}

// GroupUpdateRequestBody contains the data for a group update request.
type GroupUpdateRequestBody struct {
	Group string `json:"group"             url:"group"`

	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	ReName  *string `json:"rename,omitempty"  url:"rename,omitempty"`
	Digest  *string `json:"digest,omitempty"  url:"digest,omitempty"`
}

func (c *Client) securityGroupsPath() string {
	return "cluster/firewall/groups"
}

// CreateGroup create new security group.
func (c *Client) CreateGroup(ctx context.Context, d *GroupCreateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.securityGroupsPath(), d, nil)
	if err != nil {
		return fmt.Errorf("error creating security group: %w", err)
	}
	return nil
}

// ListGroups retrieve list of security groups.
func (c *Client) ListGroups(ctx context.Context) ([]*GroupListResponseData, error) {
	resBody := &GroupListResponseBody{}
	err := c.DoRequest(ctx, http.MethodGet, c.securityGroupsPath(), nil, resBody)
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
func (c *Client) UpdateGroup(ctx context.Context, d *GroupUpdateRequestBody) error {
	err := c.DoRequest(
		ctx,
		http.MethodPost,
		c.securityGroupsPath(),
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
		fmt.Sprintf("%s/%s", c.securityGroupsPath(), url.PathEscape(group)),
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error deleting security group '%s': %w", group, err)
	}
	return nil
}
