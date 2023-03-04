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

// TestHostsInstantiation tests whether the Hosts instance can be instantiated.
func TestHostsInstantiation(t *testing.T) {
	t.Parallel()
	s := Hosts()

	if s == nil {
		t.Fatalf("Cannot instantiate Hosts")
	}
}

// TestHostsSchema tests the Hosts schema.
func TestHostsSchema(t *testing.T) {
	t.Parallel()
	s := Hosts()

	test.AssertRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentHostsEntry,
		mkResourceVirtualEnvironmentHostsNodeName,
	})

	test.AssertComputedAttributes(t, s, []string{
		mkResourceVirtualEnvironmentHostsAddresses,
		mkResourceVirtualEnvironmentHostsDigest,
		mkResourceVirtualEnvironmentHostsEntries,
		mkResourceVirtualEnvironmentHostsHostnames,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentHostsAddresses: schema.TypeList,
		mkResourceVirtualEnvironmentHostsDigest:    schema.TypeString,
		mkResourceVirtualEnvironmentHostsEntries:   schema.TypeList,
		mkResourceVirtualEnvironmentHostsEntry:     schema.TypeList,
		mkResourceVirtualEnvironmentHostsHostnames: schema.TypeList,
		mkResourceVirtualEnvironmentHostsNodeName:  schema.TypeString,
	})

	entriesSchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentHostsEntries)

	test.AssertComputedAttributes(t, entriesSchema, []string{
		mkResourceVirtualEnvironmentHostsEntriesAddress,
		mkResourceVirtualEnvironmentHostsEntriesHostnames,
	})

	test.AssertValueTypes(t, entriesSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentHostsEntriesAddress:   schema.TypeString,
		mkResourceVirtualEnvironmentHostsEntriesHostnames: schema.TypeList,
	})

	entrySchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentHostsEntry)

	test.AssertRequiredArguments(t, entrySchema, []string{
		mkResourceVirtualEnvironmentHostsEntryAddress,
		mkResourceVirtualEnvironmentHostsEntryHostnames,
	})

	test.AssertValueTypes(t, entrySchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentHostsEntryAddress:   schema.TypeString,
		mkResourceVirtualEnvironmentHostsEntryHostnames: schema.TypeList,
	})
}
