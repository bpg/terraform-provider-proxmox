/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resources

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type haResourceTypeListQuery struct {
	ResType *types.HAResourceType `url:"type"`
}

// List retrieves the list of HA resources. If the `resType` argument is `nil`, all resources will be returned;
// otherwise resources will be filtered by the specified type (either `ct` or `vm`).
func (c *Client) List(ctx context.Context, resType *types.HAResourceType) ([]*HAResourceListResponseData, error) {
	options := &haResourceTypeListQuery{resType}
	resBody := &HAResourceListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(""), options, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing HA resources: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].ID.Type < resBody.Data[j].ID.Type ||
			(resBody.Data[i].ID.Type == resBody.Data[j].ID.Type &&
				resBody.Data[i].ID.Name < resBody.Data[j].ID.Name)
	})

	return resBody.Data, nil
}

// Get retrieves the configuration of a single HA resource.
func (c *Client) Get(ctx context.Context, id types.HAResourceID) (*HAResourceGetResponseData, error) {
	resBody := &HAResourceGetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(url.PathEscape(id.String())), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error reading HA resource: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// Create creates a new HA resource.
func (c *Client) Create(ctx context.Context, data *HAResourceCreateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(""), data, nil)
	if err != nil {
		return fmt.Errorf("error creating HA resource: %w", err)
	}

	return nil
}

// Update updates an existing HA resource.
func (c *Client) Update(ctx context.Context, id types.HAResourceID, data *HAResourceUpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(url.PathEscape(id.String())), data, nil)
	if err != nil {
		return fmt.Errorf("error updating HA resource %v: %w", id, err)
	}

	return nil
}

// Delete deletes a HA resource.
func (c *Client) Delete(ctx context.Context, id types.HAResourceID) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(url.PathEscape(id.String())), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting HA resource %v: %w", id, err)
	}

	return nil
}
