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

// TestDNSInstantiation tests whether the DNS instance can be instantiated.
func TestDNSInstantiation(t *testing.T) {
	t.Parallel()

	s := DNS()
	if s == nil {
		t.Fatalf("Cannot instantiate DNS")
	}
}

// TestDNSSchema tests the DNS schema.
func TestDNSSchema(t *testing.T) {
	t.Parallel()

	s := DNS()

	test.AssertRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentDNSDomain,
		mkResourceVirtualEnvironmentDNSNodeName,
	})

	test.AssertOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentDNSServers,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentDNSDomain:   schema.TypeString,
		mkResourceVirtualEnvironmentDNSNodeName: schema.TypeString,
		mkResourceVirtualEnvironmentDNSServers:  schema.TypeList,
	})
}
