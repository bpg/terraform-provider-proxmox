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

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
)

// TestIPSetInstantiation tests whether the IPSet
// instance can be instantiated.
func TestIPSetInstantiation(t *testing.T) {
	t.Parallel()
	require.NotNilf(t, IPSet(), "Cannot instantiate IPSet")
}

// TestIPSetSchema tests the IPSet Schema.
func TestIPSetSchema(t *testing.T) {
	t.Parallel()

	s := IPSet().Schema

	structure.AssertRequiredArguments(t, s, []string{
		mkIPSetName,
	})

	structure.AssertOptionalArguments(t, s, []string{
		mkSelectorVMID,
		mkSelectorNodeName,
		mkIPSetCIDR,
		mkIPSetCIDRComment,
	})

	structure.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkIPSetName:        schema.TypeString,
		mkIPSetCIDR:        schema.TypeList,
		mkIPSetCIDRComment: schema.TypeString,
	})

	nested := structure.AssertNestedSchemaExistence(t, s, mkIPSetCIDR).Schema

	structure.AssertRequiredArguments(t, nested, []string{
		mkIPSetCIDRName,
	})

	structure.AssertOptionalArguments(t, nested, []string{
		mkIPSetCIDRComment,
		mkIPSetCIDRNoMatch,
	})

	structure.AssertValueTypes(t, nested, map[string]schema.ValueType{
		mkIPSetCIDRName:    schema.TypeString,
		mkIPSetCIDRComment: schema.TypeString,
		mkIPSetCIDRNoMatch: schema.TypeBool,
	})
}
