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

// TestContainersInstantiation tests whether the dataSourceVirtualEnvironmentContainers instance can be instantiated.
func TestContainersInstantiation(t *testing.T) {
	t.Parallel()

	s := Containers()

	if s == nil {
		t.Fatalf("Cannot instantiate dataSourceVirtualEnvironmentContainers")
	}
}

// TestContainersSchema tests the dataSourceVirtualEnvironmentContainers schema.
func TestContainersSchema(t *testing.T) {
	t.Parallel()

	s := Containers().Schema

	test.AssertComputedAttributes(t, s, []string{
		mkDataSourceVirtualEnvironmentContainers,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentContainerNodeName: schema.TypeString,
		mkDataSourceVirtualEnvironmentContainerTags:     schema.TypeList,
		mkDataSourceFilter:                              schema.TypeList,
		mkDataSourceVirtualEnvironmentContainers:        schema.TypeList,
	})

	containersSchema := test.AssertNestedSchemaExistence(t, s, mkDataSourceVirtualEnvironmentContainers)

	test.AssertComputedAttributes(t, containersSchema, []string{
		mkDataSourceVirtualEnvironmentContainerName,
		mkDataSourceVirtualEnvironmentContainerTags,
	})

	test.AssertValueTypes(t, containersSchema, map[string]schema.ValueType{
		mkDataSourceVirtualEnvironmentContainerName:     schema.TypeString,
		mkDataSourceVirtualEnvironmentContainerNodeName: schema.TypeString,
		mkDataSourceVirtualEnvironmentContainerTags:     schema.TypeList,
		mkDataSourceVirtualEnvironmentContainerVMID:     schema.TypeInt,
	})

	filterSchema := test.AssertNestedSchemaExistence(t, s, mkDataSourceFilter)
	test.AssertValueTypes(t, filterSchema, map[string]schema.ValueType{
		mkDataSourceFilterName:   schema.TypeString,
		mkDataSourceFilterValues: schema.TypeList,
		mkDataSourceFilterRegex:  schema.TypeBool,
	})
}
