/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
)

// TestResourceVirtualEnvironmentFileInstantiation tests whether the ResourceVirtualEnvironmentFile instance can be instantiated.
func TestResourceVirtualEnvironmentFileInstantiation(t *testing.T) {
	s := resourceVirtualEnvironmentFile()

	if s == nil {
		t.Fatalf("Cannot instantiate resourceVirtualEnvironmentFile")
	}
}

// TestResourceVirtualEnvironmentFileSchema tests the resourceVirtualEnvironmentFile schema.
func TestResourceVirtualEnvironmentFileSchema(t *testing.T) {
	s := resourceVirtualEnvironmentFile()

	attributeKeys := []string{}

	for _, v := range attributeKeys {
		if s.Schema[v] == nil {
			t.Fatalf("Error in resourceVirtualEnvironmentFile.Schema: Missing attribute \"%s\"", v)
		}

		if s.Schema[v].Computed != true {
			t.Fatalf("Error in resourceVirtualEnvironmentFile.Schema: Attribute \"%s\" is not computed", v)
		}
	}
}
