/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

// TestDataSourceVirtualEnvironmentVersionInstantiation tests whether the DataSourceVirtualEnvironmentVersion instance can be instantiated.
func TestDataSourceVirtualEnvironmentVersionInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentVersion()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentVersion")
	}
}

// TestDataSourceVirtualEnvironmentVersionSchema tests the dataSourceVirtualEnvironmentVersion schema.
func TestDataSourceVirtualEnvironmentVersionSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentVersion()

	testComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentVersionKeyboardLayout,
		mkDataSourceVirtualEnvironmentVersionRelease,
		mkDataSourceVirtualEnvironmentVersionRepositoryID,
		mkDataSourceVirtualEnvironmentVersionVersion,
	})

	testSchemaValueTypes(t, s, []string{
		mkDataSourceVirtualEnvironmentVersionKeyboardLayout,
		mkDataSourceVirtualEnvironmentVersionRelease,
		mkDataSourceVirtualEnvironmentVersionRepositoryID,
		mkDataSourceVirtualEnvironmentVersionVersion,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeString,
		schema.TypeString,
		schema.TypeString,
	})
}
