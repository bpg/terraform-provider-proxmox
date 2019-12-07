/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestDataSourceVirtualEnvironmentAccessGroupInstantiation tests whether the DataSourceVirtualEnvironmentAccessGroup instance can be instantiated.
func TestDataSourceVirtualEnvironmentAccessGroupInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentAccessGroup()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentAccessGroup")
	}
}

// TestDataSourceVirtualEnvironmentAccessGroupSchema tests the dataSourceVirtualEnvironmentAccessGroup schema.
func TestDataSourceVirtualEnvironmentAccessGroupSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentAccessGroup()

	attributeKeys := []string{
		mkDataSourceVirtualEnvironmentAccessGroupComment,
		mkDataSourceVirtualEnvironmentAccessGroupMembers,
	}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in dataSourceVirtualEnvironmentAccessGroup.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in dataSourceVirtualEnvironmentAccessGroup.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
