/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestRealmLDAPModel_toCreateRequest_BindPassword(t *testing.T) {
	t.Parallel()

	t.Run("included when set", func(t *testing.T) {
		t.Parallel()

		m := realmLDAPModel{BindPassword: types.StringValue("secret")}
		req := m.toCreateRequest()

		require.NotNil(t, req.BindPassword)
		require.Equal(t, "secret", *req.BindPassword)
	})

	t.Run("omitted when null", func(t *testing.T) {
		t.Parallel()

		m := realmLDAPModel{BindPassword: types.StringNull()}
		req := m.toCreateRequest()

		require.Nil(t, req.BindPassword)
	})
}

func TestRealmLDAPModel_toUpdateRequest_BindPassword(t *testing.T) {
	t.Parallel()

	t.Run("sent when plan has value and state is null", func(t *testing.T) {
		t.Parallel()

		// WriteOnly: state is always null; plan carries the config value.
		plan := realmLDAPModel{BindPassword: types.StringValue("secret")}
		state := realmLDAPModel{BindPassword: types.StringNull()}

		req := plan.toUpdateRequest(&state)

		require.NotNil(t, req.BindPassword)
		require.Equal(t, "secret", *req.BindPassword)
	})
}
