/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package fs

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

func TestIsRetryableCreateError(t *testing.T) {
	t.Parallel()

	// The client wraps HTTP errors as Create does.
	httpErr := func(code int, message string) error {
		return fmt.Errorf("error creating CephFS %q: %w", "test", &api.HTTPError{Code: code, Message: message})
	}

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"existing filesystem", httpErr(http.StatusInternalServerError, "ceph fs 'test' already exists"), false},
		{
			"existing pools",
			httpErr(http.StatusInternalServerError, "ceph pools 'test_data' and/or 'test_metadata' already exist"),
			false,
		},
		{"no standby MDS yet", httpErr(http.StatusInternalServerError, "no standby Metadata Server (MDS) found!"), true},
		{"bad gateway", httpErr(http.StatusBadGateway, "Bad Gateway"), true},
		{"timeout", errors.New("got timeout"), true},
		{"permission denied", httpErr(http.StatusForbidden, "Permission check failed"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, isRetryableCreateError(tt.err))
		})
	}
}
