/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestDataSourceVirtualEnvironmentAccessGroupsInstantiation tests whether the DataSourceVirtualEnvironmentAccessGroups instance can be instantiated.
func TestDataSourceVirtualEnvironmentAccessGroupsInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentAccessGroups()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentAccessGroups")
	}
}

// TestDataSourceVirtualEnvironmentAccessGroupsSchema tests the dataSourceVirtualEnvironmentAccessGroups schema.
func TestDataSourceVirtualEnvironmentAccessGroupsSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentAccessGroups()

	attributeKeys := []string{
		mkDataSourceVirtualEnvironmentAccessGroupsComments,
		mkDataSourceVirtualEnvironmentAccessGroupsIDs,
	}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in dataSourceVirtualEnvironmentAccessGroups.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in dataSourceVirtualEnvironmentAccessGroups.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
