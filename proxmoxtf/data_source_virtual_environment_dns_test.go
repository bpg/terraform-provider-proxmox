/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"github.com/hashicorp/terraform/helper/schema"
	"testing"
)

// TestDataSourceVirtualEnvironmentDNSInstantiation tests whether the DataSourceVirtualEnvironmentDNS instance can be instantiated.
func TestDataSourceVirtualEnvironmentDNSInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentDNS()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentDNS")
	}
}

// TestDataSourceVirtualEnvironmentDNSSchema tests the dataSourceVirtualEnvironmentDNS schema.
func TestDataSourceVirtualEnvironmentDNSSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentDNS()

	testRequiredArguments(t, s, []string{
		mkDataSourceVirtualEnvironmentDNSNodeName,
	})

	testComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentDNSDomain,
		mkDataSourceVirtualEnvironmentDNSServers,
	})

	testSchemaValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentDNSDomain:   schema.TypeString,
		mkDataSourceVirtualEnvironmentDNSNodeName: schema.TypeString,
		mkDataSourceVirtualEnvironmentDNSServers:  schema.TypeList,
	})
}
