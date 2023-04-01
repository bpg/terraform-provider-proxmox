/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

/**
* Reference: https://pve.proxmox.com/pve-docs/api-viewer/#/cluster/firewall/ipset
 */

package firewall

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type IPSet interface {
	CreateIPSet(ctx context.Context, d *IPSetCreateRequestBody) error
	AddCIDRToIPSet(ctx context.Context, id string, d IPSetGetResponseData) error
	UpdateIPSet(ctx context.Context, d *IPSetUpdateRequestBody) error
	DeleteIPSet(ctx context.Context, id string) error
	DeleteIPSetContent(ctx context.Context, id string, cidr string) error
	GetIPSetContent(ctx context.Context, id string) ([]*IPSetGetResponseData, error)
	ListIPSets(ctx context.Context) ([]*IPSetListResponseData, error)
}

// IPSetListResponseBody contains the data from an IPSet get response.
type IPSetListResponseBody struct {
	Data []*IPSetListResponseData `json:"data,omitempty"`
}

// IPSetCreateRequestBody contains the data for an IPSet create request
type IPSetCreateRequestBody struct {
	Comment string `json:"comment,omitempty" url:"comment,omitempty"`
	Name    string `json:"name"              url:"name"`
}

// IPSetGetResponseBody contains the body from an IPSet get response.
type IPSetGetResponseBody struct {
	Data []*IPSetGetResponseData `json:"data,omitempty"`
}

// IPSetGetResponseData contains the data from an IPSet get response.
type IPSetGetResponseData struct {
	CIDR    string            `json:"cidr"              url:"cidr"`
	NoMatch *types.CustomBool `json:"nomatch,omitempty" url:"nomatch,omitempty,int"`
	Comment *string           `json:"comment,omitempty" url:"comment,omitempty"`
}

// IPSetUpdateRequestBody contains the data for an IPSet update request.
type IPSetUpdateRequestBody struct {
	ReName  string  `json:"rename,omitempty"  url:"rename,omitempty"`
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	Name    string  `json:"name"              url:"name"`
}

// IPSetListResponseData contains list of IPSets from
type IPSetListResponseData struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	Name    string  `json:"name"              url:"name"`
}

// IPSetContent is an array of IPSetGetResponseData.
type IPSetContent []IPSetGetResponseData

func (c *Client) ipsetPath() string {
	return c.AdjustPath("firewall/ipset")
}

// CreateIPSet create an IPSet
func (c *Client) CreateIPSet(ctx context.Context, d *IPSetCreateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ipsetPath(), d, nil)
	if err != nil {
		return fmt.Errorf("error creating IPSet: %w", err)
	}
	return nil
}

// AddCIDRToIPSet adds IP or Network to IPSet
func (c *Client) AddCIDRToIPSet(ctx context.Context, id string, d IPSetGetResponseData) error {
	err := c.DoRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/%s", c.ipsetPath(), url.PathEscape(id)),
		&d,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error adding CIDR to IPSet: %w", err)
	}
	return nil
}

// UpdateIPSet updates an IPSet.
func (c *Client) UpdateIPSet(ctx context.Context, d *IPSetUpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ipsetPath(), d, nil)
	if err != nil {
		return fmt.Errorf("error updating IPSet: %w", err)
	}
	return nil
}

// DeleteIPSet delete an IPSet
func (c *Client) DeleteIPSet(ctx context.Context, id string) error {
	err := c.DoRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/%s", c.ipsetPath(), url.PathEscape(id)),
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error deleting IPSet %s: %w", id, err)
	}
	return nil
}

// DeleteIPSetContent remove IP or Network from IPSet.
func (c *Client) DeleteIPSetContent(ctx context.Context, id string, cidr string) error {
	err := c.DoRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/%s/%s", c.ipsetPath(), url.PathEscape(id), url.PathEscape(cidr)),
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error deleting IPSet content %s: %w", id, err)
	}
	return nil
}

// GetIPSetContent retrieve a list of IPSet content
func (c *Client) GetIPSetContent(ctx context.Context, id string) ([]*IPSetGetResponseData, error) {
	resBody := &IPSetGetResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/%s", c.ipsetPath(), url.PathEscape(id)),
		nil,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting IPSet content: %w", err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// ListIPSets retrieves list of IPSets.
func (c *Client) ListIPSets(ctx context.Context) ([]*IPSetListResponseData, error) {
	resBody := &IPSetListResponseBody{}
	err := c.DoRequest(ctx, http.MethodGet, c.ipsetPath(), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error getting IPSet list: %w", err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Name < resBody.Data[j].Name
	})

	return resBody.Data, nil
}
