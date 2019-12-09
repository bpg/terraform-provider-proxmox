/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestResourceVirtualEnvironmentUserInstantiation tests whether the ResourceVirtualEnvironmentUser instance can be instantiated.
func TestResourceVirtualEnvironmentUserInstantiation(t *testing.T) {
	s := resourceVirtualEnvironmentUser()

	if s == nil {
		t.Fatalf("Cannot instantiate resourceVirtualEnvironmentUser")
	}
}

// TestResourceVirtualEnvironmentUserSchema tests the resourceVirtualEnvironmentUser schema.
func TestResourceVirtualEnvironmentUserSchema(t *testing.T) {
	s := resourceVirtualEnvironmentUser()

	attributeKeys := []string{}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in resourceVirtualEnvironmentUser.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in resourceVirtualEnvironmentUser.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
