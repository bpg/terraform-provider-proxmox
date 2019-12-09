/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestDataSourceVirtualEnvironmentPoolInstantiation tests whether the DataSourceVirtualEnvironmentPool instance can be instantiated.
func TestDataSourceVirtualEnvironmentPoolInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentPool()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentPool")
	}
}

// TestDataSourceVirtualEnvironmentPoolSchema tests the dataSourceVirtualEnvironmentPool schema.
func TestDataSourceVirtualEnvironmentPoolSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentPool()

	attributeKeys := []string{
		mkDataSourceVirtualEnvironmentPoolComment,
		mkDataSourceVirtualEnvironmentPoolMembers,
	}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in dataSourceVirtualEnvironmentPool.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in dataSourceVirtualEnvironmentPool.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
