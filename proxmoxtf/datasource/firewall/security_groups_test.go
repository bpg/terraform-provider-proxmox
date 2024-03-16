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

// TestSecurityGroupsSchemaInstantiation tests whether the SecurityGroupsSchema instance can be instantiated.
func TestSecurityGroupsSchemaInstantiation(t *testing.T) {
	t.Parallel()
	require.NotNil(t, SecurityGroupsSchema(), "Cannot instantiate SecurityGroupsSchema")
}

// TestSecurityGroupsSchema tests the SecurityGroupsSchema.
func TestSecurityGroupsSchema(t *testing.T) {
	t.Parallel()

	s := SecurityGroupsSchema()

	test.AssertComputedAttributes(t, s, []string{
		mkSecurityGroupsSecurityGroupNames,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkSecurityGroupsSecurityGroupNames: schema.TypeList,
	})
}
