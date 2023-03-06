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

// TestIPSetInstantiation tests whether the IPSet
// instance can be instantiated.
func TestIPSetInstantiation(t *testing.T) {
	t.Parallel()
	require.NotNilf(t, IPSet(), "Cannot instantiate IPSet")
}

// TestIPSetSchema tests the IPSet schema.
func TestIPSetSchema(t *testing.T) {
	t.Parallel()
	s := IPSet()

	test.AssertRequiredArguments(t, s, []string{
		mkIPSetName,
	})

	test.AssertOptionalArguments(t, s, []string{
		mkIPSetCIDR,
		mkIPSetCIDRComment,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkIPSetName:        schema.TypeString,
		mkIPSetCIDR:        schema.TypeList,
		mkIPSetCIDRComment: schema.TypeString,
	})

	IPSetSchema := test.AssertNestedSchemaExistence(t, s, mkIPSetCIDR)

	test.AssertRequiredArguments(t, IPSetSchema, []string{
		mkIPSetCIDRName,
	})

	test.AssertOptionalArguments(t, IPSetSchema, []string{
		mkIPSetCIDRComment,
		mkIPSetCIDRNoMatch,
	})

	test.AssertValueTypes(t, IPSetSchema, map[string]schema.ValueType{
		mkIPSetCIDRName:    schema.TypeString,
		mkIPSetCIDRComment: schema.TypeString,
		mkIPSetCIDRNoMatch: schema.TypeBool,
	})
}
