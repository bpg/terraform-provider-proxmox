/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestResourceVirtualEnvironmentRoleInstantiation tests whether the ResourceVirtualEnvironmentRole instance can be instantiated.
func TestResourceVirtualEnvironmentRoleInstantiation(t *testing.T) {
	s := resourceVirtualEnvironmentRole()

	if s == nil {
		t.Fatalf("Cannot instantiate resourceVirtualEnvironmentRole")
	}
}

// TestResourceVirtualEnvironmentRoleSchema tests the resourceVirtualEnvironmentRole schema.
func TestResourceVirtualEnvironmentRoleSchema(t *testing.T) {
	s := resourceVirtualEnvironmentRole()

	testRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentRolePrivileges,
		mkResourceVirtualEnvironmentRoleRoleID,
	})

	testValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentRolePrivileges: schema.TypeSet,
		mkResourceVirtualEnvironmentRoleRoleID:     schema.TypeString,
	})
}
