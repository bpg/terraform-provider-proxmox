/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestIPSetInstantiation tests whether the IPSet instance can be instantiated.
func TestIPSetInstantiation(t *testing.T) {
	t.Parallel()
	s := IPSet()

	if s == nil {
		t.Fatalf("Cannot instantiate IPSet")
	}
}

// TestIPSetSchema tests the IPSet schema.
func TestIPSetSchema(t *testing.T) {
	t.Parallel()
	s := IPSet()

	test.AssertRequiredArguments(t, s, []string{
		mkIPSetName,
	})

	test.AssertComputedAttributes(t, s, []string{
		mkIPSetCIDR,
		mkIPSetCIDRComment,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkIPSetName:        schema.TypeString,
		mkIPSetCIDR:        schema.TypeList,
		mkIPSetCIDRComment: schema.TypeString,
	})

	cirdSchema := test.AssertNestedSchemaExistence(t, s, mkIPSetCIDR)

	test.AssertComputedAttributes(t, cirdSchema, []string{
		mkIPSetCIDRName,
		mkIPSetCIDRNoMatch,
		mkIPSetCIDRComment,
	})

	test.AssertValueTypes(t, cirdSchema, map[string]schema.ValueType{
		mkIPSetCIDRName:    schema.TypeString,
		mkIPSetCIDRNoMatch: schema.TypeBool,
		mkIPSetCIDRComment: schema.TypeString,
	})
}
