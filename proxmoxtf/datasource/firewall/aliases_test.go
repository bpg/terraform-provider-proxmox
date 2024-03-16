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

// TestAliasesSchemaInstantiation tests whether the AliasesSchema instance can be instantiated.
func TestAliasesSchemaInstantiation(t *testing.T) {
	t.Parallel()
	require.NotNil(t, AliasesSchema(), "Cannot instantiate AliasesSchema")
}

// TestAliasesSchema tests the AliasesSchema.
func TestAliasesSchema(t *testing.T) {
	t.Parallel()

	s := AliasesSchema()

	test.AssertComputedAttributes(t, s, []string{
		mkAliasesAliasNames,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkAliasesAliasNames: schema.TypeList,
	})
}
