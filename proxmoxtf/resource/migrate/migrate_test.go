/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package migrate_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/migrate"
)

func TestWaitForResourceOnNodeReturnsWhenUnlockedOnTarget(t *testing.T) {
	t.Parallel()

	calls := 0

	locate := func(_ context.Context, _ int) (*string, error) {
		calls++
		if calls < 3 {
			return new("pve1"), nil
		}

		return new("pve2"), nil
	}

	ready := func(_ context.Context, _ string, _ int) (bool, error) {
		return true, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := migrate.WaitForResourceOnNode(ctx, 101, "pve2", locate, ready)
	require.NoError(t, err)
	require.Equal(t, 3, calls)
}

func TestWaitForResourceOnNodeKeepsWaitingUntilReady(t *testing.T) {
	t.Parallel()

	readyChecks := 0

	locate := func(_ context.Context, _ int) (*string, error) {
		return new("pve2"), nil
	}

	ready := func(_ context.Context, _ string, _ int) (bool, error) {
		readyChecks++
		return readyChecks >= 2, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := migrate.WaitForResourceOnNode(ctx, 101, "pve2", locate, ready)
	require.NoError(t, err)
	require.Equal(t, 2, readyChecks)
}

func TestWaitForResourceOnNodeFailsOnContextCancel(t *testing.T) {
	t.Parallel()

	locate := func(_ context.Context, _ int) (*string, error) {
		return new("pve1"), nil
	}

	ready := func(_ context.Context, _ string, _ int) (bool, error) {
		return true, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := migrate.WaitForResourceOnNode(ctx, 101, "pve2", locate, ready)
	require.Error(t, err)
	require.Contains(t, err.Error(), "did not migrate")
}
