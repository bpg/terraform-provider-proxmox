/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestDataSourceVirtualEnvironmentDatastoresInstantiation tests whether the DataSourceVirtualEnvironmentDatastores instance can be instantiated.
func TestDataSourceVirtualEnvironmentDatastoresInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentDatastores()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentDatastores")
	}
}

// TestDataSourceVirtualEnvironmentDatastoresSchema tests the dataSourceVirtualEnvironmentDatastores schema.
func TestDataSourceVirtualEnvironmentDatastoresSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentDatastores()

	attributeKeys := []string{
		mkDataSourceVirtualEnvironmentDatastoresActive,
		mkDataSourceVirtualEnvironmentDatastoresContentTypes,
		mkDataSourceVirtualEnvironmentDatastoresDatastoreIDs,
		mkDataSourceVirtualEnvironmentDatastoresEnabled,
		mkDataSourceVirtualEnvironmentDatastoresShared,
		mkDataSourceVirtualEnvironmentDatastoresSpaceAvailable,
		mkDataSourceVirtualEnvironmentDatastoresSpaceTotal,
		mkDataSourceVirtualEnvironmentDatastoresSpaceUsed,
		mkDataSourceVirtualEnvironmentDatastoresTypes,
	}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in dataSourceVirtualEnvironmentDatastores.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in dataSourceVirtualEnvironmentDatastores.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
