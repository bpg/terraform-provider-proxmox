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

// TestDataSourceVirtualEnvironmentRoleInstantiation tests whether the DataSourceVirtualEnvironmentRole instance can be instantiated.
func TestDataSourceVirtualEnvironmentRoleInstantiation(t *testing.T) {
	s := DataSourceVirtualEnvironmentRole()

	if s == nil {
		t.Fatalf("Cannot instantiate DataSourceVirtualEnvironmentRole")
	}
}

// TestDataSourceVirtualEnvironmentRoleSchema tests the DataSourceVirtualEnvironmentRole schema.
func TestDataSourceVirtualEnvironmentRoleSchema(t *testing.T) {
	s := DataSourceVirtualEnvironmentRole()

	test.AssertRequiredArguments(t, s, []string{
		mkDataSourceVirtualEnvironmentRoleID,
	})

	test.AssertComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentRolePrivileges,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentRoleID:         schema.TypeString,
		mkDataSourceVirtualEnvironmentRolePrivileges: schema.TypeSet,
	})
}
