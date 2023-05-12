/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datasource

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
		mkDataSourceVirtualEnvironmentHostsNodeName,
	})

	test.AssertComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentHostsAddresses,
		mkDataSourceVirtualEnvironmentHostsDigest,
		mkDataSourceVirtualEnvironmentHostsEntries,
		mkDataSourceVirtualEnvironmentHostsHostnames,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentHostsAddresses: schema.TypeList,
		mkDataSourceVirtualEnvironmentHostsDigest:    schema.TypeString,
		mkDataSourceVirtualEnvironmentHostsEntries:   schema.TypeList,
		mkDataSourceVirtualEnvironmentHostsHostnames: schema.TypeList,
		mkDataSourceVirtualEnvironmentHostsNodeName:  schema.TypeString,
	})

	entriesSchema := test.AssertNestedSchemaExistence(t, s, mkDataSourceVirtualEnvironmentHostsEntries)

	test.AssertComputedAttributes(t, entriesSchema, []string{
		mkDataSourceVirtualEnvironmentHostsEntriesAddress,
		mkDataSourceVirtualEnvironmentHostsEntriesHostnames,
	})

	test.AssertValueTypes(t, entriesSchema, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentHostsEntriesAddress:   schema.TypeString,
		mkDataSourceVirtualEnvironmentHostsEntriesHostnames: schema.TypeList,
	})
}
