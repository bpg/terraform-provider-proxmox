/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
	"github.com/google/uuid"
)

// Rule is an interface for the Proxmox firewall rule API.
type Rule interface {
	GetRulesID() string
	CreateRule(ctx context.Context, d *RuleCreateRequestBody) error
	GetRule(ctx context.Context, pos int) (*RuleGetResponseData, error)
	ListRules(ctx context.Context) ([]*RuleListResponseData, error)
	UpdateRule(ctx context.Context, pos int, d *RuleUpdateRequestBody) error
	DeleteRule(ctx context.Context, pos int) error
}

// BaseRule is the base struct for firewall rules.
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

func (c *Client) rulesPath() string {
	return c.ExpandPath("firewall/rules")
}

func (c *Client) rulePath(pos int) string {
	return fmt.Sprintf("%s/%d", c.rulesPath(), pos)
}

// GetRulesID returns the ID of the rules object.
func (c *Client) GetRulesID() string {
	// Creates unique id in every being called by using unix timestamp
	timestamp := time.Now().UnixNano() / int64(time.Microsecond)
	id := uuid.New().ID()
	finalID := timestamp + int64(id)

	return fmt.Sprintf("sg-%d", finalID)
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

	err := c.DoRequest(ctx, http.MethodGet, c.rulePath(pos), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving firewall rule %d: %w", pos, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
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
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// UpdateRule updates a firewall rule.
func (c *Client) UpdateRule(ctx context.Context, pos int, d *RuleUpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.rulePath(pos), d, nil)
	if err != nil {
		return fmt.Errorf("error updating firewall rule %d: %w", pos, err)
	}

	return nil
}

// DeleteRule deletes a firewall rule.
func (c *Client) DeleteRule(ctx context.Context, pos int) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.rulePath(pos), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting firewall rule %d: %w", pos, err)
	}

	return nil
}
