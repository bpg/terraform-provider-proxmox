/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestResourceVirtualEnvironmentRoleInstantiation tests whether the ResourceVirtualEnvironmentRole instance can be instantiated.
func TestResourceVirtualEnvironmentRoleInstantiation(t *testing.T) {
	s := ResourceVirtualEnvironmentRole()

	if s == nil {
		t.Fatalf("Cannot instantiate ResourceVirtualEnvironmentRole")
	}
}

// TestResourceVirtualEnvironmentRoleSchema tests the ResourceVirtualEnvironmentRole schema.
func TestResourceVirtualEnvironmentRoleSchema(t *testing.T) {
	s := ResourceVirtualEnvironmentRole()

	test.AssertRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentRolePrivileges,
		mkResourceVirtualEnvironmentRoleRoleID,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentRolePrivileges: schema.TypeSet,
		mkResourceVirtualEnvironmentRoleRoleID:     schema.TypeString,
	})
}
