/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
)

// TestIPSetSchemaInstantiation tests whether the IPSetSchema instance can be instantiated.
func TestIPSetSchemaInstantiation(t *testing.T) {
	t.Parallel()
	s := IPSetSchema()

	if s == nil {
		t.Fatalf("Cannot instantiate IPSet")
	}
}

// TestIPSetSchema tests the IPSet schema.
func TestIPSetSchema(t *testing.T) {
	t.Parallel()
	s := IPSetSchema()

	structure.AssertRequiredArguments(t, s, []string{
		mkIPSetName,
	})

	structure.AssertComputedAttributes(t, s, []string{
		mkIPSetCIDR,
		mkIPSetCIDRComment,
	})

	structure.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkIPSetName:        schema.TypeString,
		mkIPSetCIDR:        schema.TypeList,
		mkIPSetCIDRComment: schema.TypeString,
	})

	cirdSchema := structure.AssertNestedSchemaExistence(t, s, mkIPSetCIDR).Schema

	structure.AssertComputedAttributes(t, cirdSchema, []string{
		mkIPSetCIDRName,
		mkIPSetCIDRNoMatch,
		mkIPSetCIDRComment,
	})

	structure.AssertValueTypes(t, cirdSchema, map[string]schema.ValueType{
		mkIPSetCIDRName:    schema.TypeString,
		mkIPSetCIDRNoMatch: schema.TypeBool,
		mkIPSetCIDRComment: schema.TypeString,
	})
}
