/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package mapping

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

func TestIsMappingNotFoundErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{
			"mapping not found",
			&api.HTTPError{Code: http.StatusInternalServerError, Message: "mapping 'gpu' not found"},
			true,
		},
		{
			"wrapped mapping not found",
			fmt.Errorf("wrapped: %w", &api.HTTPError{Code: http.StatusInternalServerError, Message: "mapping 'gpu' not found"}),
			true,
		},
		{
			"another mapping not found",
			&api.HTTPError{Code: http.StatusInternalServerError, Message: "mapping 'gpu-2' not found"},
			false,
		},
		{
			"other server error",
			&api.HTTPError{Code: http.StatusInternalServerError, Message: "unable to parse mapping config"},
			false,
		},
		{"not an HTTP error", errors.New("mapping 'gpu' not found"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, isMappingNotFoundErr(tt.err, "gpu"))
		})
	}
}
