/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cpumodels

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// List retrieves all custom CPU models, sorted by name.
func (c *Client) List(ctx context.Context) ([]*GetResponseData, error) {
	resBody := &ListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(""), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing custom CPU models: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].CPUType < resBody.Data[j].CPUType
	})

	return resBody.Data, nil
}

// Get retrieves a single custom CPU model. The name may be given with or without the `custom-` prefix.
func (c *Client) Get(ctx context.Context, name string) (*GetResponseData, error) {
	resBody := &GetResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(url.PathEscape(name)), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error reading custom CPU model %q: %w", name, translateNotFound(err))
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// Create creates a custom CPU model.
func (c *Client) Create(ctx context.Context, data *CreateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath(""), data, nil)
	if err != nil {
		return fmt.Errorf("error creating custom CPU model %q: %w", data.CPUType, err)
	}

	return nil
}

// Update updates an existing custom CPU model.
func (c *Client) Update(ctx context.Context, name string, data *UpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(url.PathEscape(name)), data, nil)
	if err != nil {
		return fmt.Errorf("error updating custom CPU model %q: %w", name, translateNotFound(err))
	}

	return nil
}

// Delete deletes a custom CPU model.
func (c *Client) Delete(ctx context.Context, name string) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(url.PathEscape(name)), nil, nil)
	if err != nil {
		return fmt.Errorf("error deleting custom CPU model %q: %w", name, translateNotFound(err))
	}

	return nil
}

// translateNotFound maps PVE's HTTP 500 "Custom cputype 'X' not found" to api.ErrResourceDoesNotExist,
// which the shared client does not recognize on its own.
func translateNotFound(err error) error {
	if errors.Is(err, api.ErrResourceDoesNotExist) {
		return err
	}

	msg := err.Error()
	if strings.Contains(msg, "Custom cputype") && strings.Contains(msg, "not found") {
		return errors.Join(api.ErrResourceDoesNotExist, err)
	}

	return err
}
