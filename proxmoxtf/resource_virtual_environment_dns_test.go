/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

// TestResourceVirtualEnvironmentDNSInstantiation tests whether the ResourceVirtualEnvironmentDNS instance can be instantiated.
func TestResourceVirtualEnvironmentDNSInstantiation(t *testing.T) {
	s := resourceVirtualEnvironmentDNS()

	if s == nil {
		t.Fatalf("Cannot instantiate resourceVirtualEnvironmentDNS")
	}
}

// TestResourceVirtualEnvironmentDNSSchema tests the resourceVirtualEnvironmentDNS schema.
func TestResourceVirtualEnvironmentDNSSchema(t *testing.T) {
	s := resourceVirtualEnvironmentDNS()

	testRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentDNSDomain,
		mkResourceVirtualEnvironmentDNSNodeName,
		mkResourceVirtualEnvironmentDNSServers,
	})

	testSchemaValueTypes(t, s, []string{
		mkResourceVirtualEnvironmentDNSDomain,
		mkResourceVirtualEnvironmentDNSNodeName,
		mkResourceVirtualEnvironmentDNSServers,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeString,
		schema.TypeList,
	})
}
