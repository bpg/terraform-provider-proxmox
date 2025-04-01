/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

func (c *Client) RealmsPath() string {
	return c.ExpandPath("domains/")
}

func (c *Client) RealmPath(id string) string {
	return fmt.Sprintf("%s/%s", c.RealmsPath(), url.PathEscape(id))
}

// CreateRealm creates an access Realm.
func (c *Client) CreateRealm(ctx context.Context, d *RealmCreateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.RealmsPath(), d, nil)
	if err != nil {
		return fmt.Errorf("error creating Realm: %w", err)
	}

	return nil
}

// DeleteRealm deletes an access Realm.
func (c *Client) DeleteRealm(ctx context.Context, id string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.RealmPath(id), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting Realm: %w", err)
	}

	return nil
}

// GetRealm retrieves an access Realm.
func (c *Client) GetRealm(ctx context.Context, id string) (*RealmGetResponseBody, error) {
	resBody := &RealmGetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.RealmPath(id), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error getting Realm: %w", err)
	}

	if resBody == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody, nil
}

// ListRealms retrieves a list of access Realms.
func (c *Client) ListRealms(ctx context.Context) ([]*RealmListResponseData, error) {
	resBody := &RealmListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.RealmsPath(), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing Realms: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Realm < resBody.Data[j].Realm
	})

	return resBody.Data, nil
}

// UpdateRealm updates an access Realm.
func (c *Client) UpdateRealm(ctx context.Context, id string, d *RealmUpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.RealmPath(id), d, nil)
	if err != nil {
		return fmt.Errorf("error updating Realm: %w", err)
	}

	return nil
}
