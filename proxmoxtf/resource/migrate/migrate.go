/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

// Package migrate provides HA-aware guest migration helpers shared by the VM and container resources.
package migrate

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	pollInterval = 2 * time.Second
	maxAttempts  = 150 // 5 minutes
)

// LocateFunc reports which node currently hosts the guest.
type LocateFunc func(ctx context.Context, vmID int) (*string, error)

// ReadyFunc reports whether the guest has settled on the given node — unlocked, and for containers
// no longer mid-relocate.
type ReadyFunc func(ctx context.Context, node string, vmID int) (bool, error)

// WaitForResourceOnNode polls until the guest is on the target node and unlocked.
func WaitForResourceOnNode(
	ctx context.Context,
	vmID int,
	targetNode string,
	locate LocateFunc,
	ready ReadyFunc,
) error {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	var lastNode *string

	onTargetNode := false

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("guest %d did not migrate to node %s: %w", vmID, targetNode, ctx.Err())
		case <-ticker.C:
			currentNode, err := locate(ctx, vmID)
			if err != nil {
				tflog.Debug(ctx, "failed to get guest location, retrying...", map[string]any{
					"vm_id":   vmID,
					"attempt": attempt,
					"error":   err.Error(),
				})

				continue
			}

			lastNode = currentNode

			if currentNode == nil || *currentNode != targetNode {
				continue
			}

			onTargetNode = true

			settled, err := ready(ctx, targetNode, vmID)
			if err != nil {
				tflog.Debug(ctx, "failed to get guest status, retrying...", map[string]any{
					"vm_id":   vmID,
					"attempt": attempt,
					"error":   err.Error(),
				})

				continue
			}

			if !settled {
				continue
			}

			return nil
		}
	}

	if onTargetNode {
		return fmt.Errorf("guest %d on node %s but did not settle before timeout", vmID, targetNode)
	}

	lastNodeStr := "<unknown>"
	if lastNode != nil {
		lastNodeStr = *lastNode
	}

	return fmt.Errorf("guest %d did not migrate to node %s within timeout (last seen on %s)", vmID, targetNode, lastNodeStr)
}
