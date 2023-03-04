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

// TestProviderInstantiation() tests whether the Provider instance can be instantiated.
func TestProviderInstantiation(t *testing.T) {
	s := Provider()

	if s == nil {
		t.Fatalf("Cannot instantiate Provider")
	}
}

// TestProviderSchema() tests the Provider schema.
func TestProviderSchema(t *testing.T) {
	s := &schema.Resource{
		Schema: Provider().Schema,
	}

	test.AssertOptionalArguments(t, s, []string{
		mkProviderVirtualEnvironment,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkProviderVirtualEnvironment: schema.TypeList,
	})

	veSchema := test.AssertNestedSchemaExistence(t, s, mkProviderVirtualEnvironment)

	test.AssertOptionalArguments(t, veSchema, []string{
		mkProviderVirtualEnvironmentEndpoint,
		mkProviderVirtualEnvironmentInsecure,
		mkProviderVirtualEnvironmentOTP,
		mkProviderVirtualEnvironmentPassword,
		mkProviderVirtualEnvironmentUsername,
	})

	test.AssertValueTypes(t, veSchema, map[string]schema.ValueType{
		mkProviderVirtualEnvironmentEndpoint: schema.TypeString,
		mkProviderVirtualEnvironmentInsecure: schema.TypeBool,
		mkProviderVirtualEnvironmentOTP:      schema.TypeString,
		mkProviderVirtualEnvironmentPassword: schema.TypeString,
		mkProviderVirtualEnvironmentUsername: schema.TypeString,
	})
}
