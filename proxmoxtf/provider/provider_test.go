/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestProviderInstantiation() tests whether the ProxmoxVirtualEnvironment instance can be instantiated.
func TestProviderInstantiation(t *testing.T) {
	t.Parallel()

	s := ProxmoxVirtualEnvironment()
	if s == nil {
		t.Fatalf("Cannot instantiate ProxmoxVirtualEnvironment")
	}
}

// TestProviderSchema() tests the ProxmoxVirtualEnvironment schema.
func TestProviderSchema(t *testing.T) {
	t.Parallel()

	s := ProxmoxVirtualEnvironment().Schema

	test.AssertOptionalArguments(t, s, []string{
		mkProviderEndpoint,
		mkProviderAPIHeaders,
		mkProviderInsecure,
		mkProviderMinTLS,
		mkProviderAuthTicket,
		mkProviderCSRFPreventionToken,
		mkProviderOTP,
		mkProviderUsername,
		mkProviderPassword,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkProviderEndpoint:            schema.TypeString,
		mkProviderAPIHeaders:          schema.TypeMap,
		mkProviderInsecure:            schema.TypeBool,
		mkProviderMinTLS:              schema.TypeString,
		mkProviderAuthTicket:          schema.TypeString,
		mkProviderCSRFPreventionToken: schema.TypeString,
		mkProviderOTP:                 schema.TypeString,
		mkProviderUsername:            schema.TypeString,
		mkProviderPassword:            schema.TypeString,
	})

	// header values are credentials
	require.True(t, s[mkProviderAPIHeaders].Sensitive)

	providerSSHSchema := test.AssertNestedSchemaExistence(t, s, mkProviderSSH)

	// do not limit number of nodes in the cluster
	test.AssertListMaxItems(t, providerSSHSchema, mkProviderSSHNode, 0)
}

// TestIsAPIHeadersSetWithoutRawConfig ensures the presence check falls back to "not set" - and so to
// the environment variable - rather than panicking when no raw configuration is available.
func TestIsAPIHeadersSetWithoutRawConfig(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ProxmoxVirtualEnvironment().Schema, map[string]any{})

	require.False(t, isAPIHeadersSet(d))
}
