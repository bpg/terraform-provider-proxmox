/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"github.com/hashicorp/terraform/helper/schema"
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

	testComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentGroupsComments,
		mkDataSourceVirtualEnvironmentGroupsGroupIDs,
	})

	testSchemaValueTypes(t, s, []string{
		mkDataSourceVirtualEnvironmentGroupsComments,
		mkDataSourceVirtualEnvironmentGroupsGroupIDs,
	}, []schema.ValueType{
		schema.TypeList,
		schema.TypeList,
	})
}
