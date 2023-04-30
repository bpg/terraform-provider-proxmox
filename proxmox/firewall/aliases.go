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

type Alias interface {
	CreateAlias(ctx context.Context, d *AliasCreateRequestBody) error
	DeleteAlias(ctx context.Context, name string) error
	GetAlias(ctx context.Context, name string) (*AliasGetResponseData, error)
	ListAliases(ctx context.Context) ([]*AliasGetResponseData, error)
	UpdateAlias(ctx context.Context, name string, d *AliasUpdateRequestBody) error
}

func (c *Client) aliasesPath() string {
	return c.ExpandPath("firewall/aliases")
}

func (c *Client) aliasPath(name string) string {
	return fmt.Sprintf("%s/%s", c.aliasesPath(), url.PathEscape(name))
}

// CreateAlias create an alias
func (c *Client) CreateAlias(ctx context.Context, d *AliasCreateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.aliasesPath(), d, nil)
	if err != nil {
		return fmt.Errorf("error creating alias: %w", err)
	}
	return nil
}

// DeleteAlias delete an alias
func (c *Client) DeleteAlias(ctx context.Context, name string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.aliasPath(name), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting alias '%s': %w", name, err)
	}
	return nil
}

// GetAlias retrieves an alias
func (c *Client) GetAlias(ctx context.Context, name string) (*AliasGetResponseData, error) {
	resBody := &AliasGetResponseBody{}
	err := c.DoRequest(ctx, http.MethodGet, c.aliasPath(name), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving alias '%s': %w", name, err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// ListAliases retrieves a list of aliases.
func (c *Client) ListAliases(ctx context.Context) ([]*AliasGetResponseData, error) {
	resBody := &AliasListResponseBody{}
	err := c.DoRequest(ctx, http.MethodGet, c.aliasesPath(), nil, resBody)
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
func (c *Client) UpdateAlias(ctx context.Context, name string, d *AliasUpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.aliasPath(name), d, nil)
	if err != nil {
		return fmt.Errorf("error updating alias '%s': %w", name, err)
	}
	return nil
}
