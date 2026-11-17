/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package rules

import (
	"encoding/json"
	"testing"

	"github.com/google/go-querystring/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

func TestHARuleCreateRequestBody_QueryEncoding_Bools(t *testing.T) {
	t.Parallel()

	t.Run("nil bools are omitted", func(t *testing.T) {
		t.Parallel()

		values, err := query.Values(HARuleCreateRequestBody{})
		require.NoError(t, err)

		assert.NotContains(t, values, "disable")
		assert.NotContains(t, values, "strict")
	})

	t.Run("set bools are encoded as 0/1", func(t *testing.T) {
		t.Parallel()

		body := HARuleCreateRequestBody{}
		body.Disable = new(types.CustomBool(false))
		body.Strict = new(types.CustomBool(true))

		values, err := query.Values(body)
		require.NoError(t, err)

		assert.Equal(t, "0", values.Get("disable"))
		assert.Equal(t, "1", values.Get("strict"))
	})
}

func TestHARuleGetResponseData_JSONDecoding_Disable(t *testing.T) {
	t.Parallel()

	var withDisable HARuleGetResponseData

	require.NoError(t, json.Unmarshal([]byte(`{"rule":"r1","disable":1}`), &withDisable))
	require.NotNil(t, withDisable.Disable)
	assert.True(t, bool(*withDisable.Disable))

	var withoutDisable HARuleGetResponseData

	require.NoError(t, json.Unmarshal([]byte(`{"rule":"r1"}`), &withoutDisable))
	assert.Nil(t, withoutDisable.Disable)
}
