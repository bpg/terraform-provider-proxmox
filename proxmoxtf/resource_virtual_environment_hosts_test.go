/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

// TestResourceVirtualEnvironmentHostsInstantiation tests whether the ResourceVirtualEnvironmentHosts instance can be instantiated.
func TestResourceVirtualEnvironmentHostsInstantiation(t *testing.T) {
	s := resourceVirtualEnvironmentHosts()

	if s == nil {
		t.Fatalf("Cannot instantiate resourceVirtualEnvironmentHosts")
	}
}

// TestResourceVirtualEnvironmentHostsSchema tests the resourceVirtualEnvironmentHosts schema.
func TestResourceVirtualEnvironmentHostsSchema(t *testing.T) {
	s := resourceVirtualEnvironmentHosts()

	testRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentHostsEntry,
		mkResourceVirtualEnvironmentHostsNodeName,
	})

	testComputedAttributes(t, s, []string{
		mkResourceVirtualEnvironmentHostsAddresses,
		mkResourceVirtualEnvironmentHostsDigest,
		mkResourceVirtualEnvironmentHostsEntries,
		mkResourceVirtualEnvironmentHostsHostnames,
	})

	testSchemaValueTypes(t, s, []string{
		mkResourceVirtualEnvironmentHostsAddresses,
		mkResourceVirtualEnvironmentHostsDigest,
		mkResourceVirtualEnvironmentHostsEntries,
		mkResourceVirtualEnvironmentHostsEntry,
		mkResourceVirtualEnvironmentHostsHostnames,
		mkResourceVirtualEnvironmentHostsNodeName,
	}, []schema.ValueType{
		schema.TypeList,
		schema.TypeString,
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
		schema.TypeString,
	})

	entriesSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentHostsEntries)

	testComputedAttributes(t, entriesSchema, []string{
		mkResourceVirtualEnvironmentHostsEntriesAddress,
		mkResourceVirtualEnvironmentHostsEntriesHostnames,
	})

	testSchemaValueTypes(t, entriesSchema, []string{
		mkResourceVirtualEnvironmentHostsEntriesAddress,
		mkResourceVirtualEnvironmentHostsEntriesHostnames,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeList,
	})

	entrySchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentHostsEntry)

	testRequiredArguments(t, entrySchema, []string{
		mkResourceVirtualEnvironmentHostsEntryAddress,
		mkResourceVirtualEnvironmentHostsEntryHostnames,
	})

	testSchemaValueTypes(t, entrySchema, []string{
		mkResourceVirtualEnvironmentHostsEntryAddress,
		mkResourceVirtualEnvironmentHostsEntryHostnames,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeList,
	})
}
