/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestDataSourceVirtualEnvironmentAccessRolesInstantiation tests whether the DataSourceVirtualEnvironmentAccessRoles instance can be instantiated.
func TestDataSourceVirtualEnvironmentAccessRolesInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentAccessRoles()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentAccessRoles")
	}
}

// TestDataSourceVirtualEnvironmentAccessRolesSchema tests the dataSourceVirtualEnvironmentAccessRoles schema.
func TestDataSourceVirtualEnvironmentAccessRolesSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentAccessRoles()

	attributeKeys := []string{
		mkDataSourceVirtualEnvironmentAccessRolesPrivileges,
		mkDataSourceVirtualEnvironmentAccessRolesRoleIDs,
		mkDataSourceVirtualEnvironmentAccessRolesSpecial,
	}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in dataSourceVirtualEnvironmentAccessRoles.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in dataSourceVirtualEnvironmentAccessRoles.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
