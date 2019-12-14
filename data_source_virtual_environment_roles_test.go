/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"testing"
)

// TestDataSourceVirtualEnvironmentRolesInstantiation tests whether the DataSourceVirtualEnvironmentRoles instance can be instantiated.
func TestDataSourceVirtualEnvironmentRolesInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentRoles()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentRoles")
	}
}

// TestDataSourceVirtualEnvironmentRolesSchema tests the dataSourceVirtualEnvironmentRoles schema.
func TestDataSourceVirtualEnvironmentRolesSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentRoles()

	testComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentRolesPrivileges,
		mkDataSourceVirtualEnvironmentRolesRoleIDs,
		mkDataSourceVirtualEnvironmentRolesSpecial,
	})

	testSchemaValueTypes(t, s, []string{
		mkDataSourceVirtualEnvironmentRolesPrivileges,
		mkDataSourceVirtualEnvironmentRolesRoleIDs,
		mkDataSourceVirtualEnvironmentRolesSpecial,
	}, []schema.ValueType{
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
	})
}
