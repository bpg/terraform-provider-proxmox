/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestDataSourceVirtualEnvironmentRolesInstantiation tests whether the DataSourceVirtualEnvironmentRoles instance can be instantiated.
func TestDataSourceVirtualEnvironmentRolesInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentRoles()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentRoles")
	}
}

// TestDataSourceVirtualEnvironmentRolesSchema tests the dataSourceVirtualEnvironmentRoles schema.
func TestDataSourceVirtualEnvironmentRolesSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentRoles()

	attributeKeys := []string{
		mkDataSourceVirtualEnvironmentRolesPrivileges,
		mkDataSourceVirtualEnvironmentRolesRoleIDs,
		mkDataSourceVirtualEnvironmentRolesSpecial,
	}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in dataSourceVirtualEnvironmentRoles.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in dataSourceVirtualEnvironmentRoles.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
