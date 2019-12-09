/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestDataSourceVirtualEnvironmentPoolsInstantiation tests whether the DataSourceVirtualEnvironmentPools instance can be instantiated.
func TestDataSourceVirtualEnvironmentPoolsInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentPools()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentPools")
	}
}

// TestDataSourceVirtualEnvironmentPoolsSchema tests the dataSourceVirtualEnvironmentPools schema.
func TestDataSourceVirtualEnvironmentPoolsSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentPools()

	attributeKeys := []string{
		mkDataSourceVirtualEnvironmentPoolsPoolIDs,
	}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in dataSourceVirtualEnvironmentPools.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in dataSourceVirtualEnvironmentPools.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
