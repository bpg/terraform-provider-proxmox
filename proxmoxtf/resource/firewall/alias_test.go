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

// TestAliasSchemaInstantiation tests whether the AliasSchema instance can be instantiated.
func TestAliasSchemaInstantiation(t *testing.T) {
	t.Parallel()
	require.NotNilf(t, Alias(), "Cannot instantiate AliasSchema")
}

// TestAliasSchema tests the AliasSchema.
func TestAliasSchema(t *testing.T) {
	t.Parallel()
	s := Alias().Schema

	structure.AssertRequiredArguments(t, s, []string{
		mkAliasName,
		mkAliasCIDR,
	})

	structure.AssertOptionalArguments(t, s, []string{
		mkAliasComment,
	})

	structure.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkAliasName:    schema.TypeString,
		mkAliasCIDR:    schema.TypeString,
		mkAliasComment: schema.TypeString,
	})
}
