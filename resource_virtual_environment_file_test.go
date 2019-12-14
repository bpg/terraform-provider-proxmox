/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

// TestResourceVirtualEnvironmentFileInstantiation tests whether the ResourceVirtualEnvironmentFile instance can be instantiated.
func TestResourceVirtualEnvironmentFileInstantiation(t *testing.T) {
	s := resourceVirtualEnvironmentFile()

	if s == nil {
		t.Fatalf("Cannot instantiate resourceVirtualEnvironmentFile")
	}
}

// TestResourceVirtualEnvironmentFileSchema tests the resourceVirtualEnvironmentFile schema.
func TestResourceVirtualEnvironmentFileSchema(t *testing.T) {
	s := resourceVirtualEnvironmentFile()

	testRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentFileDatastoreID,
		mkResourceVirtualEnvironmentFileNodeName,
		mkResourceVirtualEnvironmentFileSource,
		mkResourceVirtualEnvironmentFileTemplate,
	})

	testOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentFileOverrideFileName,
	})

	testComputedAttributes(t, s, []string{
		mkResourceVirtualEnvironmentFileFileName,
	})

	testSchemaValueTypes(t, s, []string{
		mkResourceVirtualEnvironmentFileDatastoreID,
		mkResourceVirtualEnvironmentFileFileName,
		mkResourceVirtualEnvironmentFileOverrideFileName,
		mkResourceVirtualEnvironmentFileNodeName,
		mkResourceVirtualEnvironmentFileSource,
		mkResourceVirtualEnvironmentFileTemplate,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeString,
		schema.TypeString,
		schema.TypeString,
		schema.TypeString,
		schema.TypeBool,
	})
}
