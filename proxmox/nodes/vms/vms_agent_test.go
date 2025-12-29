/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package vms

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

func TestIsAgentNotReadyError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		err     error
		want    bool
		comment string
	}{
		{
			name:    "HTTP 400 BadRequest",
			err:     &api.HTTPError{Code: http.StatusBadRequest, Message: "any message"},
			want:    true,
			comment: "HTTP 400 should be treated as agent not ready",
		},
		{
			name:    "HTTP 500 with 'QEMU guest agent is not running'",
			err:     &api.HTTPError{Code: http.StatusInternalServerError, Message: "QEMU guest agent is not running"},
			want:    true,
			comment: "HTTP 500 with agent not running message should be retried",
		},
		{
			name:    "HTTP 500 with 'qemu guest agent not available'",
			err:     &api.HTTPError{Code: http.StatusInternalServerError, Message: "qemu guest agent not available"},
			want:    true,
			comment: "HTTP 500 with agent not available message should be retried",
		},
		{
			name:    "HTTP 500 with 'QEMU guest agent is not ready'",
			err:     &api.HTTPError{Code: http.StatusInternalServerError, Message: "QEMU guest agent is not ready"},
			want:    true,
			comment: "HTTP 500 with agent not ready message should be retried",
		},
		{
			name:    "HTTP 500 with case-insensitive matching",
			err:     &api.HTTPError{Code: http.StatusInternalServerError, Message: "Qemu Guest Agent Is Not Running"},
			want:    true,
			comment: "Case-insensitive matching should work",
		},
		{
			name:    "HTTP 500 with different error message",
			err:     &api.HTTPError{Code: http.StatusInternalServerError, Message: "Internal Server Error"},
			want:    false,
			comment: "HTTP 500 without agent-related message should not be retried",
		},
		{
			name:    "HTTP 500 with unrelated message",
			err:     &api.HTTPError{Code: http.StatusInternalServerError, Message: "Something went wrong"},
			want:    false,
			comment: "HTTP 500 with unrelated message should not be retried",
		},
		{
			name:    "HTTP 403 Forbidden",
			err:     &api.HTTPError{Code: http.StatusForbidden, Message: "Forbidden"},
			want:    false,
			comment: "HTTP 403 should not be retried",
		},
		{
			name:    "HTTP 404 Not Found",
			err:     &api.HTTPError{Code: http.StatusNotFound, Message: "Not Found"},
			want:    false,
			comment: "HTTP 404 should not be retried",
		},
		{
			name:    "non-HTTP error",
			err:     assert.AnError,
			want:    false,
			comment: "Non-HTTP errors should not be treated as agent not ready",
		},
		{
			name:    "nil error",
			err:     nil,
			want:    false,
			comment: "Nil error should not be treated as agent not ready",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := isAgentNotReadyError(tt.err)
			assert.Equal(t, tt.want, got, tt.comment)
		})
	}
}
