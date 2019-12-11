/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestDataSourceVirtualEnvironmentNodesInstantiation tests whether the DataSourceVirtualEnvironmentNodes instance can be instantiated.
func TestDataSourceVirtualEnvironmentNodesInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentNodes()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentNodes")
	}
}

// TestDataSourceVirtualEnvironmentNodesSchema tests the dataSourceVirtualEnvironmentNodes schema.
func TestDataSourceVirtualEnvironmentNodesSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentNodes()

	attributeKeys := []string{
		mkDataSourceVirtualEnvironmentNodesCPUCount,
		mkDataSourceVirtualEnvironmentNodesCPUUtilization,
		mkDataSourceVirtualEnvironmentNodesMemoryAvailable,
		mkDataSourceVirtualEnvironmentNodesMemoryUsed,
		mkDataSourceVirtualEnvironmentNodesNames,
		mkDataSourceVirtualEnvironmentNodesOnline,
		mkDataSourceVirtualEnvironmentNodesSSLFingerprints,
		mkDataSourceVirtualEnvironmentNodesSupportLevels,
		mkDataSourceVirtualEnvironmentNodesUptime,
	}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in dataSourceVirtualEnvironmentNodes.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in dataSourceVirtualEnvironmentNodes.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
