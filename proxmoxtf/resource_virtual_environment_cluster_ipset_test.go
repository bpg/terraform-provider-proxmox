/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"github.com/hashicorp/terraform/helper/schema"
	"testing"
)

// TestResourceVirtualEnvironmentIPSetInstantiation tests whether the resourceVirtualEnvironmentClusterIPSet
// instance can be instantiated.
func TestResourceVirtualEnvironmentIPSetInstantiation(t *testing.T) {
	s := resourceVirtualEnvironmentClusterIPSet()

	if s == nil {
		t.Fatalf("Cannot instantiate resourceVirtualEnvironmentAlias")
	}
}

// TestResourceVirtualEnvironmentIPSetSchema tests the resourceVirtualEnvironmentClusterIPSet schema.
func TestResourceVirtualEnvironmentIPSetSchema(t *testing.T) {
	s := resourceVirtualEnvironmentClusterIPSet()

	testRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentClusterIPSetName,
	})

	testOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentClusterIPSetCIDR,
		mkResourceVirtualEnvironmentClusterIPSetCIDRComment,
	})

	testValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentClusterIPSetName:        schema.TypeString,
		mkResourceVirtualEnvironmentClusterIPSetCIDR:        schema.TypeList,
		mkResourceVirtualEnvironmentClusterIPSetCIDRComment: schema.TypeString,
	})

	IPSetSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentClusterIPSetCIDR)

	testRequiredArguments(t, IPSetSchema, []string{
		mkResourceVirtualEnvironmentClusterIPSetCIDRName,
	})

	testOptionalArguments(t, IPSetSchema, []string{
		mkResourceVirtualEnvironmentClusterIPSetCIDRComment,
		mkResourceVirtualEnvironmentClusterIPSetCIDRNoMatch,
	})

	testValueTypes(t, IPSetSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentClusterIPSetCIDRName:    schema.TypeString,
		mkResourceVirtualEnvironmentClusterIPSetCIDRComment: schema.TypeString,
		mkResourceVirtualEnvironmentClusterIPSetCIDRNoMatch: schema.TypeBool,
	})

}
