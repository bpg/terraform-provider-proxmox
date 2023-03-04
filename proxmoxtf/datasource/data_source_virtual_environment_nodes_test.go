/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datasource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestDataSourceVirtualEnvironmentNodesInstantiation tests whether the DataSourceVirtualEnvironmentNodes instance can be instantiated.
func TestDataSourceVirtualEnvironmentNodesInstantiation(t *testing.T) {
	s := DataSourceVirtualEnvironmentNodes()

	if s == nil {
		t.Fatalf("Cannot instantiate DataSourceVirtualEnvironmentNodes")
	}
}

// TestDataSourceVirtualEnvironmentNodesSchema tests the DataSourceVirtualEnvironmentNodes schema.
func TestDataSourceVirtualEnvironmentNodesSchema(t *testing.T) {
	s := DataSourceVirtualEnvironmentNodes()

	test.AssertComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentNodesCPUCount,
		mkDataSourceVirtualEnvironmentNodesCPUUtilization,
		mkDataSourceVirtualEnvironmentNodesMemoryAvailable,
		mkDataSourceVirtualEnvironmentNodesMemoryUsed,
		mkDataSourceVirtualEnvironmentNodesNames,
		mkDataSourceVirtualEnvironmentNodesOnline,
		mkDataSourceVirtualEnvironmentNodesSSLFingerprints,
		mkDataSourceVirtualEnvironmentNodesSupportLevels,
		mkDataSourceVirtualEnvironmentNodesUptime,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentNodesCPUCount:        schema.TypeList,
		mkDataSourceVirtualEnvironmentNodesCPUUtilization:  schema.TypeList,
		mkDataSourceVirtualEnvironmentNodesMemoryAvailable: schema.TypeList,
		mkDataSourceVirtualEnvironmentNodesMemoryUsed:      schema.TypeList,
		mkDataSourceVirtualEnvironmentNodesNames:           schema.TypeList,
		mkDataSourceVirtualEnvironmentNodesOnline:          schema.TypeList,
		mkDataSourceVirtualEnvironmentNodesSSLFingerprints: schema.TypeList,
		mkDataSourceVirtualEnvironmentNodesSupportLevels:   schema.TypeList,
		mkDataSourceVirtualEnvironmentNodesUptime:          schema.TypeList,
	})
}
