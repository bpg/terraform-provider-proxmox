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

	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
)

// CreateRule creates a firewall rule.
func (c *Client) CreateRule(ctx context.Context, d *firewall.RuleCreateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, "cluster/firewall/rules", d, nil)
	if err != nil {
		return fmt.Errorf("error creating firewall rule: %w", err)
	}
	return nil
}

// CreateGroupRule creates a security group firewall rule.
func (c *Client) CreateGroupRule(ctx context.Context, group string, d *firewall.RuleCreateRequestBody) error {
	err := c.DoRequest(ctx,
		http.MethodPost,
		fmt.Sprintf("cluster/firewall/groups/%s", url.PathEscape(group)),
		d,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error creating security group: %w", err)
	}
	return nil
}

// GetRule retrieves a firewall rule.
func (c *Client) GetRule(ctx context.Context, pos int) (*firewall.RuleGetResponseData, error) {
	resBody := &firewall.RuleGetResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("cluster/firewall/rules/%d", pos),
		nil,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving firewall rule %d: %w", pos, err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// GetGroupRule retrieves a security group firewall rule.
func (c *Client) GetGroupRule(ctx context.Context, group string, pos int) (*firewall.RuleGetResponseData, error) {
	resBody := &firewall.RuleGetResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("cluster/firewall/groups/%s/%d", url.PathEscape(group), pos),
		nil,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving firewall rule %d: %w", pos, err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// ListRules retrieves a list of firewall rules.
func (c *Client) ListRules(ctx context.Context) ([]*firewall.RuleListResponseData, error) {
	resBody := &firewall.RuleListResponseBody{}
	err := c.DoRequest(ctx, http.MethodGet, "cluster/firewall/rules", nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving firewall rules: %w", err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// ListGroupRules retrieves a list of security group firewall rules.
func (c *Client) ListGroupRules(ctx context.Context, group string) ([]*firewall.RuleListResponseData, error) {
	resBody := &firewall.RuleListResponseBody{}
	err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("cluster/firewall/groups/%s", url.PathEscape(group)), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving firewall rules: %w", err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// UpdateRule updates a firewall rule.
func (c *Client) UpdateRule(ctx context.Context, pos int, d *firewall.RuleUpdateRequestBody) error {
	err := c.DoRequest(
		ctx,
		http.MethodPut,
		fmt.Sprintf("cluster/firewall/rules/%d", pos),
		d,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error updating firewall rule %d: %w", pos, err)
	}
	return nil
}

// UpdateGroupRule updates a security group firewall rule.
func (c *Client) UpdateGroupRule(ctx context.Context, group string, pos int, d *firewall.RuleUpdateRequestBody) error {
	err := c.DoRequest(
		ctx,
		http.MethodPut,
		fmt.Sprintf("cluster/firewall/groups/%s/%d", url.PathEscape(group), pos),
		d,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error updating firewall rule %d: %w", pos, err)
	}
	return nil
}

// DeleteRule deletes a firewall rule.
func (c *Client) DeleteRule(ctx context.Context, pos int) error {
	err := c.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("cluster/firewall/rules/%d", pos), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting firewall rule %d: %w", pos, err)
	}
	return nil
}

// DeleteGroupRule deletes a security group firewall rule.
func (c *Client) DeleteGroupRule(ctx context.Context, group string, pos int) error {
	err := c.DoRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("cluster/firewall/groups/%s/%d", url.PathEscape(group), pos),
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error deleting firewall rule %d: %w", pos, err)
	}
	return nil
}
