/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package backup

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

func TestIsNoSuchJobErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{
			"read of a missing job",
			&api.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "Parameter verification failed. (id: No such job 'backup-gone')",
			},
			true,
		},
		{
			"delete of a missing job",
			fmt.Errorf("wrapped: %w", &api.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Parameter verification failed. (id: No such job 'backup-gone')",
			}),
			true,
		},
		{
			"other parameter error",
			&api.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "Parameter verification failed. (schedule: value does not match the regex pattern)",
			},
			false,
		},
		{"not an HTTP error", errors.New("No such job 'backup-gone'"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, isNoSuchJobErr(tt.err))
		})
	}
}
