/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestResourceVirtualEnvironmentIPSetInstantiation tests whether the resourceVirtualEnvironmentFirewallIPSet
// instance can be instantiated.
func TestResourceVirtualEnvironmentIPSetInstantiation(t *testing.T) {
	s := resourceVirtualEnvironmentFirewallIPSet()

	if s == nil {
		t.Fatalf("Cannot instantiate resourceVirtualEnvironmentAlias")
	}
}

// TestResourceVirtualEnvironmentIPSetSchema tests the resourceVirtualEnvironmentFirewallIPSet schema.
func TestResourceVirtualEnvironmentIPSetSchema(t *testing.T) {
	s := resourceVirtualEnvironmentFirewallIPSet()

	testRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentFirewallIPSetName,
	})

	testOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentFirewallIPSetCIDR,
		mkResourceVirtualEnvironmentFirewallIPSetCIDRComment,
	})

	testValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentFirewallIPSetName:        schema.TypeString,
		mkResourceVirtualEnvironmentFirewallIPSetCIDR:        schema.TypeList,
		mkResourceVirtualEnvironmentFirewallIPSetCIDRComment: schema.TypeString,
	})

	IPSetSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentFirewallIPSetCIDR)

	testRequiredArguments(t, IPSetSchema, []string{
		mkResourceVirtualEnvironmentFirewallIPSetCIDRName,
	})

	testOptionalArguments(t, IPSetSchema, []string{
		mkResourceVirtualEnvironmentFirewallIPSetCIDRComment,
		mkResourceVirtualEnvironmentFirewallIPSetCIDRNoMatch,
	})

	testValueTypes(t, IPSetSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentFirewallIPSetCIDRName:    schema.TypeString,
		mkResourceVirtualEnvironmentFirewallIPSetCIDRComment: schema.TypeString,
		mkResourceVirtualEnvironmentFirewallIPSetCIDRNoMatch: schema.TypeBool,
	})
}
