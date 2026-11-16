/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package groups

import (
	"encoding/json"
	"testing"

	"github.com/google/go-querystring/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func TestHAGroupCreateRequestBody_QueryEncoding_Bools(t *testing.T) {
	t.Parallel()

	t.Run("nil bools are omitted", func(t *testing.T) {
		t.Parallel()

		values, err := query.Values(HAGroupCreateRequestBody{})
		require.NoError(t, err)

		assert.NotContains(t, values, "nofailback")
		assert.NotContains(t, values, "restricted")
	})

	t.Run("set bools are encoded as 0/1", func(t *testing.T) {
		t.Parallel()

		body := HAGroupCreateRequestBody{}
		body.NoFailback = new(types.CustomBool(true))
		body.Restricted = new(types.CustomBool(false))

		values, err := query.Values(body)
		require.NoError(t, err)

		assert.Equal(t, "1", values.Get("nofailback"))
		assert.Equal(t, "0", values.Get("restricted"))
	})
}

func TestHAGroupGetResponseData_JSONDecoding_Bools(t *testing.T) {
	t.Parallel()

	var withBools HAGroupGetResponseData

	require.NoError(t, json.Unmarshal([]byte(`{"group":"g1","nodes":"pve","nofailback":1,"restricted":0}`), &withBools))
	require.NotNil(t, withBools.NoFailback)
	require.NotNil(t, withBools.Restricted)
	assert.True(t, bool(*withBools.NoFailback))
	assert.False(t, bool(*withBools.Restricted))

	var withoutBools HAGroupGetResponseData

	require.NoError(t, json.Unmarshal([]byte(`{"group":"g1","nodes":"pve"}`), &withoutBools))
	assert.Nil(t, withoutBools.NoFailback)
	assert.Nil(t, withoutBools.Restricted)
}
