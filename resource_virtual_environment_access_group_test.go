/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestResourceVirtualEnvironmentAccessGroupInstantiation tests whether the ResourceVirtualEnvironmentAccessGroup instance can be instantiated.
func TestResourceVirtualEnvironmentAccessGroupInstantiation(t *testing.T) {
	s := resourceVirtualEnvironmentAccessGroup()

	if s == nil {
		t.Fatalf("Cannot instantiate resourceVirtualEnvironmentAccessGroup")
	}
}

// TestResourceVirtualEnvironmentAccessGroupSchema tests the resourceVirtualEnvironmentAccessGroup schema.
func TestResourceVirtualEnvironmentAccessGroupSchema(t *testing.T) {
	s := resourceVirtualEnvironmentAccessGroup()

	attributeKeys := []string{
		mkResourceVirtualEnvironmentAccessGroupMembers,
	}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in resourceVirtualEnvironmentAccessGroup.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in resourceVirtualEnvironmentAccessGroup.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
