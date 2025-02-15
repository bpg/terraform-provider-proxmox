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

// TestContainerInstantiation tests whether the Container instance can be instantiated.
func TestContainerInstantiation(t *testing.T) {
	t.Parallel()

	s := Container()

	if s == nil {
		t.Fatalf("Cannot instantiate Container")
	}
}

// TestContainerSchema tests the Container schema.
func TestContainerSchema(t *testing.T) {
	t.Parallel()

	s := Container().Schema

	test.AssertComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentContainerName,
		mkDataSourceVirtualEnvironmentContainerTags,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentContainerName:     schema.TypeString,
		mkDataSourceVirtualEnvironmentContainerNodeName: schema.TypeString,
		mkDataSourceVirtualEnvironmentContainerTags:     schema.TypeList,
		mkDataSourceVirtualEnvironmentContainerTemplate: schema.TypeBool,
		mkDataSourceVirtualEnvironmentContainerStatus:   schema.TypeString,
		mkDataSourceVirtualEnvironmentContainerVMID:     schema.TypeInt,
	})
}
