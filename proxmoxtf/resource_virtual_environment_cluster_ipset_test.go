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
		mkResourceVirtualEnvironmentClusterIPSet,
		mkResourceVirtualEnvironmentClusterIPSetComment,
	})

	testValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentClusterIPSetName: 		schema.TypeString,
		mkResourceVirtualEnvironmentClusterIPSet:			schema.TypeList,
		mkResourceVirtualEnvironmentClusterIPSetComment:   	schema.TypeString,
	})

	IPSetSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentClusterIPSet)

	testRequiredArguments(t, IPSetSchema, []string{
		mkResourceVirtualEnvironmentClusterIPSetCIDR,
	})

	testOptionalArguments(t, IPSetSchema, []string{
		mkResourceVirtualEnvironmentClusterIPSetComment,
		mkResourceVirtualEnvironmentClusterIPSetNoMatch,
	})

	testValueTypes(t, IPSetSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentClusterIPSetCIDR: schema.TypeString,
		mkResourceVirtualEnvironmentClusterIPSetComment: schema.TypeString,
		mkResourceVirtualEnvironmentClusterIPSetNoMatch: schema.TypeBool,
	})

}
