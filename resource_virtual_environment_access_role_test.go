/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestResourceVirtualEnvironmentAccessRoleInstantiation tests whether the ResourceVirtualEnvironmentAccessRole instance can be instantiated.
func TestResourceVirtualEnvironmentAccessRoleInstantiation(t *testing.T) {
	s := resourceVirtualEnvironmentAccessRole()

	if s == nil {
		t.Fatalf("Cannot instantiate resourceVirtualEnvironmentAccessRole")
	}
}

// TestResourceVirtualEnvironmentAccessRoleSchema tests the resourceVirtualEnvironmentAccessRole schema.
func TestResourceVirtualEnvironmentAccessRoleSchema(t *testing.T) {
	s := resourceVirtualEnvironmentAccessRole()

	attributeKeys := []string{}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in resourceVirtualEnvironmentAccessRole.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in resourceVirtualEnvironmentAccessRole.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
