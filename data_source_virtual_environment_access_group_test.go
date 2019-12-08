/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestDataSourceVirtualEnvironmentGroupInstantiation tests whether the DataSourceVirtualEnvironmentGroup instance can be instantiated.
func TestDataSourceVirtualEnvironmentGroupInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentGroup()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentGroup")
	}
}

// TestDataSourceVirtualEnvironmentGroupSchema tests the dataSourceVirtualEnvironmentGroup schema.
func TestDataSourceVirtualEnvironmentGroupSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentGroup()

	attributeKeys := []string{
		mkDataSourceVirtualEnvironmentGroupComment,
		mkDataSourceVirtualEnvironmentGroupMembers,
	}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in dataSourceVirtualEnvironmentGroup.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in dataSourceVirtualEnvironmentGroup.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
