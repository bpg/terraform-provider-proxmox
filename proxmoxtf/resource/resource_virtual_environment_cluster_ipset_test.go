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

// TestResourceVirtualEnvironmentIPSetInstantiation tests whether the ResourceVirtualEnvironmentFirewallIPSet
// instance can be instantiated.
func TestResourceVirtualEnvironmentIPSetInstantiation(t *testing.T) {
	s := ResourceVirtualEnvironmentFirewallIPSet()

	if s == nil {
		t.Fatalf("Cannot instantiate ResourceVirtualEnvironmentAlias")
	}
}

// TestResourceVirtualEnvironmentIPSetSchema tests the ResourceVirtualEnvironmentFirewallIPSet schema.
func TestResourceVirtualEnvironmentIPSetSchema(t *testing.T) {
	s := ResourceVirtualEnvironmentFirewallIPSet()

	test.AssertRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentFirewallIPSetName,
	})

	test.AssertOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentFirewallIPSetCIDR,
		mkResourceVirtualEnvironmentFirewallIPSetCIDRComment,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentFirewallIPSetName:        schema.TypeString,
		mkResourceVirtualEnvironmentFirewallIPSetCIDR:        schema.TypeList,
		mkResourceVirtualEnvironmentFirewallIPSetCIDRComment: schema.TypeString,
	})

	IPSetSchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentFirewallIPSetCIDR)

	test.AssertRequiredArguments(t, IPSetSchema, []string{
		mkResourceVirtualEnvironmentFirewallIPSetCIDRName,
	})

	test.AssertOptionalArguments(t, IPSetSchema, []string{
		mkResourceVirtualEnvironmentFirewallIPSetCIDRComment,
		mkResourceVirtualEnvironmentFirewallIPSetCIDRNoMatch,
	})

	test.AssertValueTypes(t, IPSetSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentFirewallIPSetCIDRName:    schema.TypeString,
		mkResourceVirtualEnvironmentFirewallIPSetCIDRComment: schema.TypeString,
		mkResourceVirtualEnvironmentFirewallIPSetCIDRNoMatch: schema.TypeBool,
	})
}
