/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tasks

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/retry"
)

func TestDoTask_ResultFromLastAttempt(t *testing.T) {
	t.Parallel()

	lockErr := errors.New("can't lock file '/var/lock/pve-manager/pve-storage-local-lvm' - got timeout")
	errGone := errors.New("resource does not exist")

	tests := []struct {
		name      string
		retryDone bool
		wantErr   error
	}{
		{name: "retry finds the work already done", retryDone: true},
		{name: "retry dispatch fails", wantErr: errGone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			op := retry.NewTaskOperation("test",
				retry.WithBaseDelay(time.Millisecond),
				retry.WithRetryIf(func(error) bool { return false }),
				retry.WithAlreadyDoneCheck(func(err error) bool { return tt.retryDone && errors.Is(err, errGone) }),
			)

			dispatches := 0
			taskID := "UPID:test"

			result := doTask(context.Background(), op,
				func() (*string, error) {
					dispatches++
					if dispatches == 1 {
						return &taskID, nil
					}

					return nil, errGone
				},
				func(context.Context, string) TaskResult { return TaskFailed(lockErr) },
			)

			require.Equal(t, 2, dispatches)

			if tt.wantErr == nil {
				require.NoError(t, result.Err())
				return
			}

			require.ErrorIs(t, result.Err(), tt.wantErr)
		})
	}
}
