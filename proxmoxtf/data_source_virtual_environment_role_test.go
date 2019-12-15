/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

// TestDataSourceVirtualEnvironmentRoleInstantiation tests whether the DataSourceVirtualEnvironmentRole instance can be instantiated.
func TestDataSourceVirtualEnvironmentRoleInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentRole()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentRole")
	}
}

// TestDataSourceVirtualEnvironmentRoleSchema tests the dataSourceVirtualEnvironmentRole schema.
func TestDataSourceVirtualEnvironmentRoleSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentRole()

	testRequiredArguments(t, s, []string{
		mkDataSourceVirtualEnvironmentRoleID,
	})

	testComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentRolePrivileges,
	})

	testSchemaValueTypes(t, s, []string{
		mkDataSourceVirtualEnvironmentRoleID,
		mkDataSourceVirtualEnvironmentRolePrivileges,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeSet,
	})
}
