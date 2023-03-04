/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datasource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestRolesInstantiation tests whether the Roles instance can be instantiated.
func TestRolesInstantiation(t *testing.T) {
	t.Parallel()
	s := Roles()

	if s == nil {
		t.Fatalf("Cannot instantiate Roles")
	}
}

// TestRolesSchema tests the Roles schema.
func TestRolesSchema(t *testing.T) {
	t.Parallel()
	s := Roles()

	test.AssertComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentRolesPrivileges,
		mkDataSourceVirtualEnvironmentRolesRoleIDs,
		mkDataSourceVirtualEnvironmentRolesSpecial,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentRolesPrivileges: schema.TypeList,
		mkDataSourceVirtualEnvironmentRolesRoleIDs:    schema.TypeList,
		mkDataSourceVirtualEnvironmentRolesSpecial:    schema.TypeList,
	})
}
