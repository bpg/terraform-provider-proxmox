/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
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

	attributeKeys := []string{}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in resourceVirtualEnvironmentRole.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in resourceVirtualEnvironmentRole.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
