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

// TestIPSetsSchemaInstantiation tests whether the IPSetsSchema instance can be instantiated.
func TestIPSetsSchemaInstantiation(t *testing.T) {
	t.Parallel()
	require.NotNil(t, IPSetsSchema(), "Cannot instantiate IPSetsSchema")
}

// TestIPSetsSchema tests the IPSetsSchema.
func TestIPSetsSchema(t *testing.T) {
	t.Parallel()

	r := schema.Resource{Schema: IPSetsSchema()}

	test.AssertComputedAttributes(t, &r, []string{
		mkIPSetsIPSetNames,
	})

	test.AssertValueTypes(t, &r, map[string]schema.ValueType{
		mkIPSetsIPSetNames: schema.TypeList,
	})
}
