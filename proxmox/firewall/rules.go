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
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type Rule interface {
	GetStableID() string
	CreateRule(ctx context.Context, d *RuleCreateRequestBody) error
	GetRule(ctx context.Context, pos int) (*RuleGetResponseData, error)
	ListRules(ctx context.Context) ([]*RuleListResponseData, error)
	UpdateRule(ctx context.Context, pos int, d *RuleUpdateRequestBody) error
	DeleteRule(ctx context.Context, pos int) error
}

// RuleCreateRequestBody contains the data for a firewall rule create request.
type RuleCreateRequestBody struct {
	BaseRule

	Action string `json:"action" url:"action"`
	Type   string `json:"type"   url:"type"`

	Group *string `json:"group,omitempty" url:"group,omitempty"`
}

// RuleGetResponseBody contains the body from a firewall rule get response.
type RuleGetResponseBody struct {
	Data *RuleGetResponseData `json:"data,omitempty"`
}

// RuleGetResponseData contains the data from a firewall rule get response.
type RuleGetResponseData struct {
	BaseRule

	// NOTE: This is `int` in the PVE API docs, but it's actually a string in the response.
	Pos    string `json:"pos"     url:"pos"`
	Action string `json:"action"  url:"action"`
	Type   string `json:"type"    url:"type"`
}

// RuleListResponseBody contains the data from a firewall rule get response.
type RuleListResponseBody struct {
	Data []*RuleListResponseData `json:"data,omitempty"`
}

// RuleListResponseData contains the data from a firewall rule get response.
type RuleListResponseData struct {
	Pos int `json:"pos" url:"pos"`
}

// RuleUpdateRequestBody contains the data for a firewall rule update request.
type RuleUpdateRequestBody struct {
	BaseRule

	Pos    *int    `json:"pos,omitempty"    url:"pos,omitempty"`
	Action *string `json:"action,omitempty" url:"action,omitempty"`
	Type   *string `json:"type,omitempty"   url:"type,omitempty"`

	Group *string `json:"group,omitempty"   url:"group,omitempty"`
}

type BaseRule struct {
	Comment  *string           `json:"comment,omitempty"   url:"comment,omitempty"`
	Dest     *string           `json:"dest,omitempty"      url:"dest,omitempty"`
	Digest   *string           `json:"digest,omitempty"    url:"digest,omitempty"`
	DPort    *string           `json:"dport,omitempty"     url:"dport,omitempty"`
	Enable   *types.CustomBool `json:"enable,omitempty"    url:"enable,omitempty,int"`
	ICMPType *string           `json:"icmp-type,omitempty" url:"icmp-type,omitempty"`
	IFace    *string           `json:"iface,omitempty"     url:"iface,omitempty"`
	Log      *string           `json:"log,omitempty"       url:"log,omitempty"`
	Macro    *string           `json:"macro,omitempty"     url:"macro,omitempty"`
	Proto    *string           `json:"proto,omitempty"     url:"proto,omitempty"`
	Source   *string           `json:"source,omitempty"    url:"source,omitempty"`
	SPort    *string           `json:"sport,omitempty"     url:"sport,omitempty"`
}

// "cluster/firewall/groups/%s" -> "cluster/firewall/rules"

func (c *Client) rulesPath() string {
	return c.AdjustPath("firewall/rules")
}

func (c *Client) GetStableID() string {
	return "rule-" + strconv.Itoa(schema.HashString(c.rulesPath()))
}

// CreateRule creates a firewall rule.
func (c *Client) CreateRule(ctx context.Context, d *RuleCreateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.rulesPath(), d, nil)
	if err != nil {
		return fmt.Errorf("error creating firewall rule: %w", err)
	}
	return nil
}

// GetRule retrieves a firewall rule.
func (c *Client) GetRule(ctx context.Context, pos int) (*RuleGetResponseData, error) {
	resBody := &RuleGetResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/%d", c.rulesPath(), pos),
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
func (c *Client) ListRules(ctx context.Context) ([]*RuleListResponseData, error) {
	resBody := &RuleListResponseBody{}
	err := c.DoRequest(ctx, http.MethodGet, c.rulesPath(), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving firewall rules: %w", err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// UpdateRule updates a firewall rule.
func (c *Client) UpdateRule(ctx context.Context, pos int, d *RuleUpdateRequestBody) error {
	err := c.DoRequest(
		ctx,
		http.MethodPut,
		fmt.Sprintf("%s/%d", c.rulesPath(), pos),
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
	err := c.DoRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/%d", c.rulesPath(), pos),
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error deleting firewall rule %d: %w", pos, err)
	}
	return nil
}
