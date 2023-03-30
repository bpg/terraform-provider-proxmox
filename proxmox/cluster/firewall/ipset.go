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

// CreateIPSet create an IPSet
func (c *Client) CreateIPSet(ctx context.Context, d *firewall.IPSetCreateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, "cluster/firewall/ipset", d, nil)
	if err != nil {
		return fmt.Errorf("error creating IPSet: %w", err)
	}
	return nil
}

// AddCIDRToIPSet adds IP or Network to IPSet
func (c *Client) AddCIDRToIPSet(ctx context.Context, id string, d firewall.IPSetGetResponseData) error {
	err := c.DoRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("cluster/firewall/ipset/%s/", url.PathEscape(id)),
		&d,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error adding CIDR to IPSet: %w", err)
	}
	return nil
}

// UpdateIPSet updates an IPSet.
func (c *Client) UpdateIPSet(ctx context.Context, d *firewall.IPSetUpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, "cluster/firewall/ipset/", d, nil)
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
		fmt.Sprintf("cluster/firewall/ipset/%s", url.PathEscape(id)),
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
		fmt.Sprintf("cluster/firewall/ipset/%s/%s", url.PathEscape(id), url.PathEscape(cidr)),
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error deleting IPSet content %s: %w", id, err)
	}
	return nil
}

// GetIPSetContent retrieve a list of IPSet content
func (c *Client) GetIPSetContent(ctx context.Context, id string) ([]*firewall.IPSetGetResponseData, error) {
	resBody := &firewall.IPSetGetResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("cluster/firewall/ipset/%s", url.PathEscape(id)),
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
func (c *Client) ListIPSets(ctx context.Context) ([]*firewall.IPSetListResponseData, error) {
	resBody := &firewall.IPSetListResponseBody{}
	err := c.DoRequest(ctx, http.MethodGet, "cluster/firewall/ipset", nil, resBody)
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
