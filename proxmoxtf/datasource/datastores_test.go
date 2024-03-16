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

// TestDatastoresInstantiation tests whether the Datastores instance can be instantiated.
func TestDatastoresInstantiation(t *testing.T) {
	t.Parallel()

	s := Datastores()
	if s == nil {
		t.Fatalf("Cannot instantiate Datastores")
	}
}

// TestDatastoresSchema tests the Datastores schema.
func TestDatastoresSchema(t *testing.T) {
	t.Parallel()

	s := Datastores().Schema

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
