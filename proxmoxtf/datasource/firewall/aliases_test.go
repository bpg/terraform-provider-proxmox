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

// TestAliasesSchemaInstantiation tests whether the AliasesSchema instance can be instantiated.
func TestAliasesSchemaInstantiation(t *testing.T) {
	t.Parallel()
	s := AliasesSchema()

	if s == nil {
		t.Fatalf("Cannot instantiate Aliases")
	}
}

// TestAliasesSchema tests the Aliases schema.
func TestAliasesSchema(t *testing.T) {
	t.Parallel()
	s := AliasesSchema()

	structure.AssertComputedAttributes(t, s, []string{
		mkAliasesAliasNames,
	})

	structure.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkAliasesAliasNames: schema.TypeList,
	})
}
