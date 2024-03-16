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

// TestVMsInstantiation tests whether the dataSourceVirtualEnvironmentVMs instance can be instantiated.
func TestVMsInstantiation(t *testing.T) {
	t.Parallel()

	s := VMs()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentVMs")
	}
}

// TestVMsSchema tests the dataSourceVirtualEnvironmentVMs schema.
func TestVMsSchema(t *testing.T) {
	t.Parallel()

	s := VMs().Schema

	test.AssertComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentVMs,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentVMNodeName: schema.TypeString,
		mkDataSourceVirtualEnvironmentVMTags:     schema.TypeList,
		mkDataSourceVirtualEnvironmentVMs:        schema.TypeList,
	})

	vmsSchema := test.AssertNestedSchemaExistence(t, s, mkDataSourceVirtualEnvironmentVMs)

	test.AssertComputedAttributes(t, vmsSchema, []string{
		mkDataSourceVirtualEnvironmentVMName,
		mkDataSourceVirtualEnvironmentVMTags,
	})

	test.AssertValueTypes(t, vmsSchema, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentVMName:     schema.TypeString,
		mkDataSourceVirtualEnvironmentVMNodeName: schema.TypeString,
		mkDataSourceVirtualEnvironmentVMTags:     schema.TypeList,
		mkDataSourceVirtualEnvironmentVMVMID:     schema.TypeInt,
	})
}
