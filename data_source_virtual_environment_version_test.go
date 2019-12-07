/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestDataSourceVirtualEnvironmentVersionInstantiation tests whether the DataSourceVirtualEnvironmentVersion instance can be instantiated.
func TestDataSourceVirtualEnvironmentVersionInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentVersion()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentVersion")
	}
}

// TestDataSourceVirtualEnvironmentVersionSchema tests the dataSourceVirtualEnvironmentVersion schema.
func TestDataSourceVirtualEnvironmentVersionSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentVersion()

	attributeKeys := []string{
		mkDataSourceVirtualEnvironmentVersionRelease,
		mkDataSourceVirtualEnvironmentVersionRepositoryID,
		mkDataSourceVirtualEnvironmentVersionVersion,
	}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in dataSourceVirtualEnvironmentVersion.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in dataSourceVirtualEnvironmentVersion.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
