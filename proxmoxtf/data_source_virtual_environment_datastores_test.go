/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"github.com/hashicorp/terraform/helper/schema"
	"testing"
)

// TestDataSourceVirtualEnvironmentDatastoresInstantiation tests whether the DataSourceVirtualEnvironmentDatastores instance can be instantiated.
func TestDataSourceVirtualEnvironmentDatastoresInstantiation(t *testing.T) {
	s := dataSourceVirtualEnvironmentDatastores()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentDatastores")
	}
}

// TestDataSourceVirtualEnvironmentDatastoresSchema tests the dataSourceVirtualEnvironmentDatastores schema.
func TestDataSourceVirtualEnvironmentDatastoresSchema(t *testing.T) {
	s := dataSourceVirtualEnvironmentDatastores()

	testRequiredArguments(t, s, []string{
		mkDataSourceVirtualEnvironmentDatastoresNodeName,
	})

	testComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentDatastoresActive,
		mkDataSourceVirtualEnvironmentDatastoresContentTypes,
		mkDataSourceVirtualEnvironmentDatastoresDatastoreIDs,
		mkDataSourceVirtualEnvironmentDatastoresEnabled,
		mkDataSourceVirtualEnvironmentDatastoresShared,
		mkDataSourceVirtualEnvironmentDatastoresSpaceAvailable,
		mkDataSourceVirtualEnvironmentDatastoresSpaceTotal,
		mkDataSourceVirtualEnvironmentDatastoresSpaceUsed,
		mkDataSourceVirtualEnvironmentDatastoresTypes,
	})

	testSchemaValueTypes(t, s, []string{
		mkDataSourceVirtualEnvironmentDatastoresActive,
		mkDataSourceVirtualEnvironmentDatastoresContentTypes,
		mkDataSourceVirtualEnvironmentDatastoresDatastoreIDs,
		mkDataSourceVirtualEnvironmentDatastoresEnabled,
		mkDataSourceVirtualEnvironmentDatastoresNodeName,
		mkDataSourceVirtualEnvironmentDatastoresShared,
		mkDataSourceVirtualEnvironmentDatastoresSpaceAvailable,
		mkDataSourceVirtualEnvironmentDatastoresSpaceTotal,
		mkDataSourceVirtualEnvironmentDatastoresSpaceUsed,
		mkDataSourceVirtualEnvironmentDatastoresTypes,
	}, []schema.ValueType{
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
		schema.TypeString,
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
		schema.TypeList,
	})
}
