/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestDataSourceVirtualEnvironmentUserInstantiation tests whether the DataSourceVirtualEnvironmentUser instance can be instantiated.
func TestDataSourceVirtualEnvironmentUserInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentUser()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentUser")
	}
}

// TestDataSourceVirtualEnvironmentUserSchema tests the dataSourceVirtualEnvironmentUser schema.
func TestDataSourceVirtualEnvironmentUserSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentUser()

	attributeKeys := []string{
		mkDataSourceVirtualEnvironmentUserComment,
		mkDataSourceVirtualEnvironmentUserEmail,
		mkDataSourceVirtualEnvironmentUserEnabled,
		mkDataSourceVirtualEnvironmentUserExpirationDate,
		mkDataSourceVirtualEnvironmentUserFirstName,
		mkDataSourceVirtualEnvironmentUserGroups,
		mkDataSourceVirtualEnvironmentUserKeys,
		mkDataSourceVirtualEnvironmentUserLastName,
	}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in dataSourceVirtualEnvironmentUser.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in dataSourceVirtualEnvironmentUser.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
