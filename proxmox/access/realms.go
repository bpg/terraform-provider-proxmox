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

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

func (c *Client) realmsPath() string {
	return c.ExpandPath("domains")
}

func (c *Client) realmPath(realm string) string {
	return fmt.Sprintf("%s/%s", c.realmsPath(), url.PathEscape(realm))
}

// CreateRealm creates an authentication realm.
func (c *Client) CreateRealm(ctx context.Context, d *RealmCreateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.realmsPath(), d, nil)
	if err != nil {
		return fmt.Errorf("error creating realm: %w", err)
	}

	return nil
}

// GetRealm retrieves a realm configuration.
func (c *Client) GetRealm(ctx context.Context, realm string) (*RealmGetResponseData, error) {
	resBody := &RealmGetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.realmPath(realm), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving realm: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// UpdateRealm updates a realm configuration.
func (c *Client) UpdateRealm(ctx context.Context, realm string, d *RealmUpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.realmPath(realm), d, nil)
	if err != nil {
		return fmt.Errorf("error updating realm: %w", err)
	}

	return nil
}

// DeleteRealm deletes a realm.
func (c *Client) DeleteRealm(ctx context.Context, realm string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.realmPath(realm), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting realm: %w", err)
	}

	return nil
}

// ListRealms retrieves all realms.
func (c *Client) ListRealms(ctx context.Context) ([]*RealmListResponseData, error) {
	resBody := &RealmListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.realmsPath(), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing realms: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// SyncRealm triggers user/group synchronization from LDAP.
func (c *Client) SyncRealm(ctx context.Context, realm string, d *RealmSyncRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, fmt.Sprintf("%s/sync", c.realmPath(realm)), d, nil)
	if err != nil {
		return fmt.Errorf("error syncing realm: %w", err)
	}

	return nil
}
