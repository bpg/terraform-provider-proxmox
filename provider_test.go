/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
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

	testOptionalArguments(t, s, []string{
		mkProviderVirtualEnvironment,
	})

	testSchemaValueTypes(t, s, []string{
		mkProviderVirtualEnvironment,
	}, []schema.ValueType{
		schema.TypeList,
	})

	veSchema := testNestedSchemaExistence(t, s, mkProviderVirtualEnvironment)

	testRequiredArguments(t, veSchema, []string{
		mkProviderVirtualEnvironmentEndpoint,
		mkProviderVirtualEnvironmentPassword,
		mkProviderVirtualEnvironmentUsername,
	})

	testOptionalArguments(t, veSchema, []string{
		mkProviderVirtualEnvironmentInsecure,
	})

	testSchemaValueTypes(t, veSchema, []string{
		mkProviderVirtualEnvironmentEndpoint,
		mkProviderVirtualEnvironmentInsecure,
		mkProviderVirtualEnvironmentPassword,
		mkProviderVirtualEnvironmentUsername,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeBool,
		schema.TypeString,
		schema.TypeString,
	})
}
