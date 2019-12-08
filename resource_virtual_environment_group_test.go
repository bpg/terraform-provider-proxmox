/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestResourceVirtualEnvironmentGroupInstantiation tests whether the ResourceVirtualEnvironmentGroup instance can be instantiated.
func TestResourceVirtualEnvironmentGroupInstantiation(t *testing.T) {
	s := resourceVirtualEnvironmentGroup()

	if s == nil {
		t.Fatalf("Cannot instantiate resourceVirtualEnvironmentGroup")
	}
}

// TestResourceVirtualEnvironmentGroupSchema tests the resourceVirtualEnvironmentGroup schema.
func TestResourceVirtualEnvironmentGroupSchema(t *testing.T) {
	s := resourceVirtualEnvironmentGroup()

	attributeKeys := []string{
		mkResourceVirtualEnvironmentGroupMembers,
	}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in resourceVirtualEnvironmentGroup.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in resourceVirtualEnvironmentGroup.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
