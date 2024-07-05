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

// TestVMInstantiation tests whether the VM instance can be instantiated.
func TestVMInstantiation(t *testing.T) {
	t.Parallel()

	s := VM()

	if s == nil {
		t.Fatalf("Cannot instantiate VM")
	}
}

// TestVMSchema tests the VM schema.
func TestVMSchema(t *testing.T) {
	t.Parallel()

	s := VM().Schema

	test.AssertComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentVMName,
		mkDataSourceVirtualEnvironmentVMTags,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentVMName:     schema.TypeString,
		mkDataSourceVirtualEnvironmentVMNodeName: schema.TypeString,
		mkDataSourceVirtualEnvironmentVMTags:     schema.TypeList,
		mkDataSourceVirtualEnvironmentVMTemplate: schema.TypeBool,
		mkDataSourceVirtualEnvironmentVMStatus:   schema.TypeString,
		mkDataSourceVirtualEnvironmentVMVMID:     schema.TypeInt,
	})
}
