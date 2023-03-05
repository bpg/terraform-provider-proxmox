/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package firewall

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
)

// CreateAlias create an alias
func (a *API) CreateAlias(ctx context.Context, d *AliasCreateRequestBody) error {
	err := a.DoRequest(ctx, http.MethodPost, "cluster/firewall/aliases", d, nil)
	if err != nil {
		return fmt.Errorf("error creating alias: %w", err)
	}
	return nil
}

// DeleteAlias delete an alias
func (a *API) DeleteAlias(ctx context.Context, id string) error {
	err := a.DoRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("cluster/firewall/aliases/%s", url.PathEscape(id)),
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error deleting alias %s: %w", id, err)
	}
	return nil
}

// GetAlias retrieves an alias
func (a *API) GetAlias(ctx context.Context, id string) (*AliasGetResponseData, error) {
	resBody := &AliasGetResponseBody{}
	err := a.DoRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("cluster/firewall/aliases/%s", url.PathEscape(id)),
		nil,
		resBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving alias %s: %w", id, err)
	}

	if resBody.Data == nil {
		return nil, errors.New("the server did not include a data object in the response")
	}

	return resBody.Data, nil
}

// ListAliases retrieves a list of aliases.
func (a *API) ListAliases(ctx context.Context) ([]*AliasGetResponseData, error) {
	resBody := &AliasListResponseBody{}
	err := a.DoRequest(ctx, http.MethodGet, "cluster/firewall/aliases", nil, resBody)
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
func (a *API) UpdateAlias(ctx context.Context, id string, d *AliasUpdateRequestBody) error {
	err := a.DoRequest(
		ctx,
		http.MethodPut,
		fmt.Sprintf("cluster/firewall/aliases/%s", url.PathEscape(id)),
		d,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error updating alias %s: %w", id, err)
	}
	return nil
}
