/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package replications

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/retry"
)

// GetReplication retrieves a single Replication by ID.
func (c *Client) GetReplication(ctx context.Context) (*ReplicationData, error) {
	resBody := &replicationResponse{}

	err := c.DoRequest(ctx, http.MethodGet, c.ExpandPath(), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error reading Replication %s: %w", c.ID, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// GetReplications lists all Replications.
func (c *Client) GetReplications(ctx context.Context) ([]ReplicationData, error) {
	resBody := &replicationsResponse{}

	err := c.DoRequest(ctx, http.MethodGet, c.basePath(), nil, resBody)
	if err != nil {
		return nil, fmt.Errorf("error listing Replications: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return *resBody.Data, nil
}

// CreateReplication creates a new Replication.
func (c *Client) CreateReplication(ctx context.Context, data *ReplicationCreate) error {
	err := c.DoRequest(ctx, http.MethodPost, c.basePath(), data, nil)
	if err != nil {
		return fmt.Errorf("error creating Replication: %w", err)
	}

	return nil
}

// UpdateReplication Updates an existing Replication.
func (c *Client) UpdateReplication(ctx context.Context, data *ReplicationUpdate) error {
	err := c.DoRequest(ctx, http.MethodPut, c.ExpandPath(), data, nil)
	if err != nil {
		return fmt.Errorf("error updating Replication: %w", err)
	}

	return nil
}

// DeleteReplication deletes a Replication and wait for it to complete.
func (c *Client) DeleteReplication(ctx context.Context, data *ReplicationDelete) error {
	err := c.DeleteReplicationAsync(ctx, data)
	if err != nil {
		return err
	}

	stillExists := errors.New("replication still exists")

	op := retry.NewPollOperation("replication delete",
		retry.WithRetryIf(func(err error) bool {
			return errors.Is(err, stillExists)
		}),
	)

	err = op.DoPoll(ctx, func() error {
		repls, err := c.GetReplications(ctx)
		if err != nil {
			return err
		}

		for _, r := range repls {
			if r.ID == c.ID {
				return stillExists
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error waiting for deleting replication: %w", err)
	}

	return nil
}

// DeleteReplicationAsync deletes a Replication but not wait for it to complete.
func (c *Client) DeleteReplicationAsync(ctx context.Context, data *ReplicationDelete) error {
	err := c.DoRequest(ctx, http.MethodDelete, c.ExpandPath(), data, nil)
	if err != nil {
		return fmt.Errorf("error deleting replication: %w", err)
	}

	return nil
}
