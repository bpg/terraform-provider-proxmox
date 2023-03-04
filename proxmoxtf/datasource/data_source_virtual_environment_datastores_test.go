/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datasource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestDataSourceVirtualEnvironmentDatastoresInstantiation tests whether the DataSourceVirtualEnvironmentDatastores instance can be instantiated.
func TestDataSourceVirtualEnvironmentDatastoresInstantiation(t *testing.T) {
	s := DataSourceVirtualEnvironmentDatastores()

	if s == nil {
		t.Fatalf("Cannot instantiate DataSourceVirtualEnvironmentDatastores")
	}
}

// TestDataSourceVirtualEnvironmentDatastoresSchema tests the DataSourceVirtualEnvironmentDatastores schema.
func TestDataSourceVirtualEnvironmentDatastoresSchema(t *testing.T) {
	s := DataSourceVirtualEnvironmentDatastores()

	test.AssertRequiredArguments(t, s, []string{
		mkDataSourceVirtualEnvironmentDatastoresNodeName,
	})

	test.AssertComputedAttributes(t, s, []string{
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

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentDatastoresActive:         schema.TypeList,
		mkDataSourceVirtualEnvironmentDatastoresContentTypes:   schema.TypeList,
		mkDataSourceVirtualEnvironmentDatastoresDatastoreIDs:   schema.TypeList,
		mkDataSourceVirtualEnvironmentDatastoresEnabled:        schema.TypeList,
		mkDataSourceVirtualEnvironmentDatastoresNodeName:       schema.TypeString,
		mkDataSourceVirtualEnvironmentDatastoresShared:         schema.TypeList,
		mkDataSourceVirtualEnvironmentDatastoresSpaceAvailable: schema.TypeList,
		mkDataSourceVirtualEnvironmentDatastoresSpaceTotal:     schema.TypeList,
		mkDataSourceVirtualEnvironmentDatastoresSpaceUsed:      schema.TypeList,
		mkDataSourceVirtualEnvironmentDatastoresTypes:          schema.TypeList,
	})
}
