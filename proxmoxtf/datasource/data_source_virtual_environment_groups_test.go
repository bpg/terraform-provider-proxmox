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

// TestDataSourceVirtualEnvironmentGroupsInstantiation tests whether the DataSourceVirtualEnvironmentGroups instance can be instantiated.
func TestDataSourceVirtualEnvironmentGroupsInstantiation(t *testing.T) {
	s := DataSourceVirtualEnvironmentGroups()

	if s == nil {
		t.Fatalf("Cannot instantiate DataSourceVirtualEnvironmentGroups")
	}
}

// TestDataSourceVirtualEnvironmentGroupsSchema tests the DataSourceVirtualEnvironmentGroups schema.
func TestDataSourceVirtualEnvironmentGroupsSchema(t *testing.T) {
	s := DataSourceVirtualEnvironmentGroups()

	test.AssertComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentGroupsComments,
		mkDataSourceVirtualEnvironmentGroupsGroupIDs,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentGroupsComments: schema.TypeList,
		mkDataSourceVirtualEnvironmentGroupsGroupIDs: schema.TypeList,
	})
}
