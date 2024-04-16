package cpu

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

func TestCPUSchema(t *testing.T) {
	t.Parallel()

	s := Schema()

	cpuSchema := test.AssertNestedSchemaExistence(t, s, MkCPU)

	test.AssertOptionalArguments(t, cpuSchema, []string{
		mkCPUArchitecture,
		mkCPUCores,
		mkCPUFlags,
		mkCPUHotplugged,
		mkCPUNUMA,
		mkCPUSockets,
		mkCPUType,
		mkCPUUnits,
	})

	test.AssertValueTypes(t, cpuSchema, map[string]schema.ValueType{
		mkCPUArchitecture: schema.TypeString,
		mkCPUCores:        schema.TypeInt,
		mkCPUFlags:        schema.TypeList,
		mkCPUHotplugged:   schema.TypeInt,
		mkCPUNUMA:         schema.TypeBool,
		mkCPUSockets:      schema.TypeInt,
		mkCPUType:         schema.TypeString,
		mkCPUUnits:        schema.TypeInt,
	})

}
