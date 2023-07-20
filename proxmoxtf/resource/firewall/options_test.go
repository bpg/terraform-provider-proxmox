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

// TestOptionsInstantiation tests whether the Options instance can be instantiated.
func TestOptionsInstantiation(t *testing.T) {
	t.Parallel()
	require.NotNilf(t, Options(), "Cannot instantiate Options")
}

// TestOptionsSchema tests the Options Schema.
func TestOptionsSchema(t *testing.T) {
	t.Parallel()

	s := Options()

	test.AssertOptionalArguments(t, s, []string{
		mkDHCP,
		mkEnabled,
		mkIPFilter,
		mkLogLevelIN,
		mkLogLevelOUT,
		mkMACFilter,
		mkNDP,
		mkPolicyIn,
		mkPolicyOut,
		mkRadv,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkDHCP:        schema.TypeBool,
		mkEnabled:     schema.TypeBool,
		mkIPFilter:    schema.TypeBool,
		mkLogLevelIN:  schema.TypeString,
		mkLogLevelOUT: schema.TypeString,
		mkMACFilter:   schema.TypeBool,
		mkNDP:         schema.TypeBool,
		mkPolicyIn:    schema.TypeString,
		mkPolicyOut:   schema.TypeString,
		mkRadv:        schema.TypeBool,
	})
}
