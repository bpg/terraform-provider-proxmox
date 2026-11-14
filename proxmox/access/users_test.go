/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

func TestIsNoSuchUserErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{
			"user gone",
			&api.HTTPError{Code: http.StatusInternalServerError, Message: "no such user ('bot@pve')"},
			true,
		},
		{
			"wrapped",
			fmt.Errorf("wrapped: %w", &api.HTTPError{Code: http.StatusInternalServerError, Message: "no such user ('bot@pve')"}),
			true,
		},
		{
			"permission denied",
			&api.HTTPError{Code: http.StatusForbidden, Message: "Permission check failed"},
			false,
		},
		{"not an HTTP error", errors.New("no such user ('bot@pve')"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, isNoSuchUserErr(tt.err))
		})
	}
}
