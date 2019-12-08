/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestDataSourceVirtualEnvironmentGroupsInstantiation tests whether the DataSourceVirtualEnvironmentGroups instance can be instantiated.
func TestDataSourceVirtualEnvironmentGroupsInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentGroups()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentGroups")
	}
}

// TestDataSourceVirtualEnvironmentGroupsSchema tests the dataSourceVirtualEnvironmentGroups schema.
func TestDataSourceVirtualEnvironmentGroupsSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentGroups()

	attributeKeys := []string{
		mkDataSourceVirtualEnvironmentGroupsComments,
		mkDataSourceVirtualEnvironmentGroupsGroupIDs,
	}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in dataSourceVirtualEnvironmentGroups.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in dataSourceVirtualEnvironmentGroups.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
