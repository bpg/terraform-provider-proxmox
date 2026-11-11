/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package plugins

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

func TestIsPluginNotDefinedErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{
			"plugin not defined",
			&api.HTTPError{Code: http.StatusInternalServerError, Message: "ACME plugin 'my-plugin' not defined"},
			true,
		},
		{
			"other server error",
			&api.HTTPError{Code: http.StatusInternalServerError, Message: "unable to parse plugin config"},
			false,
		},
		{
			"other not defined message",
			&api.HTTPError{Code: http.StatusInternalServerError, Message: "property 'foo' not defined"},
			false,
		},
		{"unrelated error", errors.New("connection refused"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, isPluginNotDefinedErr(tt.err))
		})
	}
}
