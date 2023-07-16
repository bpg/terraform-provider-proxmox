/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestIPSetSchemaInstantiation tests whether the IPSetSchema instance can be instantiated.
func TestIPSetSchemaInstantiation(t *testing.T) {
	t.Parallel()
	require.NotNil(t, IPSetSchema(), "Cannot instantiate IPSetSchema")
}

// TestIPSetSchema tests the IPSetSchema.
func TestIPSetSchema(t *testing.T) {
	t.Parallel()

	r := schema.Resource{Schema: IPSetSchema()}

	test.AssertRequiredArguments(t, &r, []string{
		mkIPSetName,
	})

	test.AssertComputedAttributes(t, &r, []string{
		mkIPSetCIDR,
		mkIPSetCIDRComment,
	})

	test.AssertValueTypes(t, &r, map[string]schema.ValueType{
		mkIPSetName:        schema.TypeString,
		mkIPSetCIDR:        schema.TypeList,
		mkIPSetCIDRComment: schema.TypeString,
	})

	cidr := test.AssertNestedSchemaExistence(t, &r, mkIPSetCIDR)

	test.AssertComputedAttributes(t, cidr, []string{
		mkIPSetCIDRName,
		mkIPSetCIDRNoMatch,
		mkIPSetCIDRComment,
	})

	test.AssertValueTypes(t, cidr, map[string]schema.ValueType{
		mkIPSetCIDRName:    schema.TypeString,
		mkIPSetCIDRNoMatch: schema.TypeBool,
		mkIPSetCIDRComment: schema.TypeString,
	})
}
