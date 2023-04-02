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

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
)

// TestSecurityGroupInstantiation tests whether the SecurityGroup instance can be instantiated.
func TestSecurityGroupInstantiation(t *testing.T) {
	t.Parallel()
	require.NotNilf(t, SecurityGroup(), "Cannot instantiate SecurityGroup")
}

// TestSecurityGroupSchema tests the SecurityGroup Schema.
func TestSecurityGroupSchema(t *testing.T) {
	t.Parallel()
	s := SecurityGroup().Schema

	structure.AssertRequiredArguments(t, s, []string{
		mkSecurityGroupName,
	})

	structure.AssertOptionalArguments(t, s, []string{
		mkSecurityGroupComment,
	})

	structure.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkSecurityGroupName:    schema.TypeString,
		mkSecurityGroupComment: schema.TypeString,
	})

	structure.AssertNestedSchemaExistence(t, s, firewall.MkRule)
}
