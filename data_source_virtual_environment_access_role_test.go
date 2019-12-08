/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestDataSourceVirtualEnvironmentAccessRoleInstantiation tests whether the DataSourceVirtualEnvironmentAccessRole instance can be instantiated.
func TestDataSourceVirtualEnvironmentAccessRoleInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentAccessRole()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentAccessRole")
	}
}

// TestDataSourceVirtualEnvironmentAccessRoleSchema tests the dataSourceVirtualEnvironmentAccessRole schema.
func TestDataSourceVirtualEnvironmentAccessRoleSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentAccessRole()

	attributeKeys := []string{
		mkDataSourceVirtualEnvironmentAccessRolePrivileges,
	}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in dataSourceVirtualEnvironmentAccessRole.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in dataSourceVirtualEnvironmentAccessRole.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
