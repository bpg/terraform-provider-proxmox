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

// CreateAlias create an alias
func (c *Client) CreateAlias(ctx context.Context, d *firewall.AliasCreateRequestBody) error {
	err := c.DoRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf("nodes/%s/qemu/%d/firewall/aliases", url.PathEscape(c.NodeName), c.VMID),
		d,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error creating alias: %w", err)
	}
	return nil
}

// DeleteAlias delete an alias
func (c *Client) DeleteAlias(ctx context.Context, name string) error {
	err := c.DoRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("nodes/%s/qemu/%d/firewall/aliases/%s",
			url.PathEscape(c.NodeName), c.VMID, url.PathEscape(name)),
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error deleting alias '%s': %w", name, err)
	}
	return nil
}

// GetAlias retrieves an alias
func (c *Client) GetAlias(ctx context.Context, name string) (*firewall.AliasGetResponseData, error) {
	resBody := &firewall.AliasGetResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("nodes/%s/qemu/%d/firewall/aliases/%s",
			url.PathEscape(c.NodeName), c.VMID, url.PathEscape(name)),
		nil,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving alias '%s': %w", name, err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// ListAliases retrieves a list of aliases.
func (c *Client) ListAliases(ctx context.Context) ([]*firewall.AliasGetResponseData, error) {
	resBody := &firewall.AliasListResponseBody{}
	err := c.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("nodes/%s/qemu/%d/firewall/aliases", c.NodeName, c.VMID),
		nil,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving aliases: %w", err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Name < resBody.Data[j].Name
	})

	return resBody.Data, nil
}

// UpdateAlias updates an alias.
func (c *Client) UpdateAlias(ctx context.Context, name string, d *firewall.AliasUpdateRequestBody) error {
	err := c.DoRequest(
		ctx,
		http.MethodPut,
		fmt.Sprintf("nodes/%s/qemu/%d/firewall/aliases/%s",
			url.PathEscape(c.NodeName), c.VMID, url.PathEscape(name)),
		d,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error updating alias '%s': %w", name, err)
	}
	return nil
}
