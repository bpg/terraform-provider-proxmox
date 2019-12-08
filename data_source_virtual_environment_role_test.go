/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestDataSourceVirtualEnvironmentRoleInstantiation tests whether the DataSourceVirtualEnvironmentRole instance can be instantiated.
func TestDataSourceVirtualEnvironmentRoleInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentRole()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentRole")
	}
}

// TestDataSourceVirtualEnvironmentRoleSchema tests the dataSourceVirtualEnvironmentRole schema.
func TestDataSourceVirtualEnvironmentRoleSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentRole()

	attributeKeys := []string{
		mkDataSourceVirtualEnvironmentRolePrivileges,
	}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in dataSourceVirtualEnvironmentRole.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in dataSourceVirtualEnvironmentRole.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
