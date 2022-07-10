/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	testValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentHostsAddresses: schema.TypeList,
		mkResourceVirtualEnvironmentHostsDigest:    schema.TypeString,
		mkResourceVirtualEnvironmentHostsEntries:   schema.TypeList,
		mkResourceVirtualEnvironmentHostsEntry:     schema.TypeList,
		mkResourceVirtualEnvironmentHostsHostnames: schema.TypeList,
		mkResourceVirtualEnvironmentHostsNodeName:  schema.TypeString,
	})

	entriesSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentHostsEntries)

	testComputedAttributes(t, entriesSchema, []string{
		mkResourceVirtualEnvironmentHostsEntriesAddress,
		mkResourceVirtualEnvironmentHostsEntriesHostnames,
	})

	testValueTypes(t, entriesSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentHostsEntriesAddress:   schema.TypeString,
		mkResourceVirtualEnvironmentHostsEntriesHostnames: schema.TypeList,
	})

	entrySchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentHostsEntry)

	testRequiredArguments(t, entrySchema, []string{
		mkResourceVirtualEnvironmentHostsEntryAddress,
		mkResourceVirtualEnvironmentHostsEntryHostnames,
	})

	testValueTypes(t, entrySchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentHostsEntryAddress:   schema.TypeString,
		mkResourceVirtualEnvironmentHostsEntryHostnames: schema.TypeList,
	})
}
