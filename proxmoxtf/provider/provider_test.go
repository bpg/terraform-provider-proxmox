/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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
		mkProviderUsername,
		mkProviderPassword,
		mkProviderEndpoint,
		mkProviderInsecure,
		mkProviderMinTLS,
		mkProviderOTP,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkProviderUsername: schema.TypeString,
		mkProviderPassword: schema.TypeString,
		mkProviderEndpoint: schema.TypeString,
		mkProviderInsecure: schema.TypeBool,
		mkProviderMinTLS:   schema.TypeString,
		mkProviderOTP:      schema.TypeString,
	})

	providerSSHSchema := test.AssertNestedSchemaExistence(t, s, mkProviderSSH)

	// do not limit number of nodes in the cluster
	test.AssertListMaxItems(t, providerSSHSchema, mkProviderSSHNode, 0)
}
