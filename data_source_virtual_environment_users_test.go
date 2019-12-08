/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestDataSourceVirtualEnvironmentUsersInstantiation tests whether the DataSourceVirtualEnvironmentUsers instance can be instantiated.
func TestDataSourceVirtualEnvironmentUsersInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentUsers()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentUsers")
	}
}

// TestDataSourceVirtualEnvironmentUsersSchema tests the dataSourceVirtualEnvironmentUsers schema.
func TestDataSourceVirtualEnvironmentUsersSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentUsers()

	attributeKeys := []string{
		mkDataSourceVirtualEnvironmentUsersComments,
		mkDataSourceVirtualEnvironmentUsersEmails,
		mkDataSourceVirtualEnvironmentUsersEnabled,
		mkDataSourceVirtualEnvironmentUsersExpirationDates,
		mkDataSourceVirtualEnvironmentUsersFirstNames,
		mkDataSourceVirtualEnvironmentUsersGroups,
		mkDataSourceVirtualEnvironmentUsersKeys,
		mkDataSourceVirtualEnvironmentUsersLastNames,
		mkDataSourceVirtualEnvironmentUsersUserIDs,
	}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in dataSourceVirtualEnvironmentUsers.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in dataSourceVirtualEnvironmentUsers.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
