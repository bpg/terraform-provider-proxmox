/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"github.com/hashicorp/terraform/helper/schema"
	"testing"
)

// TestDataSourceVirtualEnvironmentHostsInstantiation tests whether the DataSourceVirtualEnvironmentHosts instance can be instantiated.
func TestDataSourceVirtualEnvironmentHostsInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentHosts()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentHosts")
	}
}

// TestDataSourceVirtualEnvironmentHostsSchema tests the dataSourceVirtualEnvironmentHosts schema.
func TestDataSourceVirtualEnvironmentHostsSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentHosts()

	testRequiredArguments(t, s, []string{
		mkDataSourceVirtualEnvironmentHostsNodeName,
	})

	testComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentHostsAddresses,
		mkDataSourceVirtualEnvironmentHostsDigest,
		mkDataSourceVirtualEnvironmentHostsHostnames,
	})

	testSchemaValueTypes(t, s, []string{
		mkDataSourceVirtualEnvironmentHostsAddresses,
		mkDataSourceVirtualEnvironmentHostsDigest,
		mkDataSourceVirtualEnvironmentHostsHostnames,
		mkDataSourceVirtualEnvironmentHostsNodeName,
	}, []schema.ValueType{
		schema.TypeList,
		schema.TypeString,
		schema.TypeList,
		schema.TypeString,
	})
}
