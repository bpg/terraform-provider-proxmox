/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestResourceVirtualEnvironmentPoolInstantiation tests whether the ResourceVirtualEnvironmentPool instance can be instantiated.
func TestResourceVirtualEnvironmentPoolInstantiation(t *testing.T) {
	s := resourceVirtualEnvironmentPool()

	if s == nil {
		t.Fatalf("Cannot instantiate resourceVirtualEnvironmentPool")
	}
}

// TestResourceVirtualEnvironmentPoolSchema tests the resourceVirtualEnvironmentPool schema.
func TestResourceVirtualEnvironmentPoolSchema(t *testing.T) {
	s := resourceVirtualEnvironmentPool()

	attributeKeys := []string{}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in resourceVirtualEnvironmentPool.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in resourceVirtualEnvironmentPool.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
