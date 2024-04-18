/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package nodes

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

const (
	networkReloadTimeout = 5 * time.Second
)

// reloadLock is used to prevent concurrent network reloads.
// global variable by design.
//
//nolint:gochecknoglobals
var reloadLock sync.Mutex

// ListNetworkInterfaces retrieves a list of network interfaces for a specific nodes.
func (c *Client) ListNetworkInterfaces(ctx context.Context) ([]*NetworkInterfaceListResponseData, error) {
	resBody := &NetworkInterfaceListResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath("network"), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces for node \"%s\": %w", c.NodeName, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	sort.Slice(resBody.Data, func(i, j int) bool {
		return resBody.Data[i].Priority < resBody.Data[j].Priority
	})

	return resBody.Data, nil
}

// CreateNetworkInterface creates a network interface for a specific node.
func (c *Client) CreateNetworkInterface(ctx context.Context, d *NetworkInterfaceCreateUpdateRequestBody) error {
	err := c.DoRequest(ctx, http.MethodPost, c.ExpandPath("network"), d, nil)
	if err != nil {
		return fmt.Errorf(
			"failed to create network interface \"%s\" for node \"%s\": %w",
			d.Iface, c.NodeName, err,
		)
	}

	return nil
}

// ReloadNetworkConfiguration reloads the network configuration for a specific node.
func (c *Client) ReloadNetworkConfiguration(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, networkReloadTimeout)
	defer cancel()

	reloadLock.Lock()
	defer reloadLock.Unlock()

	resBody := &ReloadNetworkResponseBody{}

	err := retry.Do(
		func() error {
			err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath("network"), nil, resBody)
			if err != nil {
				return err //nolint:wrapcheck
			}

			if resBody.Data == nil {
				return api.ErrNoDataObjectInResponse
			}

			return c.Tasks().WaitForTask(ctx, *resBody.Data)
		},
		retry.Context(ctx),
		retry.Delay(1*time.Second),
		retry.Attempts(3),
		retry.RetryIf(func(err error) bool {
			return strings.Contains(err.Error(), "exit code 89")
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to reload network configuration for node \"%s\": %w", c.NodeName, err)
	}

	return nil
}

// RevertNetworkConfiguration reverts the network configuration changes for a specific node.
func (c *Client) RevertNetworkConfiguration(ctx context.Context) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath("network"), nil, nil)
	if err != nil {
		return fmt.Errorf("failed to revert network configuration for node \"%s\": %w", c.NodeName, err)
	}

	return nil
}

// UpdateNetworkInterface updates a network interface for a specific node.
func (c *Client) UpdateNetworkInterface(
	ctx context.Context,
	iface string,
	d *NetworkInterfaceCreateUpdateRequestBody,
) error {
	err := c.DoRequest(
		ctx,
		http.MethodPut,
		c.ExpandPath(fmt.Sprintf("network/%s", url.PathEscape(iface))),
		d,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to update network interface \"%s\" for node \"%s\": %w",
			d.Iface, c.NodeName, err,
		)
	}

	return nil
}

// DeleteNetworkInterface deletes a network interface configuration for a specific node.
func (c *Client) DeleteNetworkInterface(ctx context.Context, iface string) error {
	err := c.DoRequest(
		ctx,
		http.MethodDelete,
		c.ExpandPath(fmt.Sprintf("network/%s", url.PathEscape(iface))),
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to delete network interface \"%s\" for node \"%s\": %w",
			iface, c.NodeName, err,
		)
	}

	return nil
}
