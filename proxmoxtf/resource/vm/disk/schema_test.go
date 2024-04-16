package disk

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

func TestDiskSchema(t *testing.T) {
	t.Parallel()

	s := Schema()

	diskSchema := test.AssertNestedSchemaExistence(t, s, MkDisk)

	test.AssertOptionalArguments(t, diskSchema, []string{
		mkDiskDatastoreID,
		mkDiskPathInDatastore,
		mkDiskFileFormat,
		mkDiskFileID,
		mkDiskSize,
	})

	test.AssertValueTypes(t, diskSchema, map[string]schema.ValueType{
		mkDiskDatastoreID:     schema.TypeString,
		mkDiskPathInDatastore: schema.TypeString,
		mkDiskFileFormat:      schema.TypeString,
		mkDiskFileID:          schema.TypeString,
		mkDiskSize:            schema.TypeInt,
	})

	diskSpeedSchema := test.AssertNestedSchemaExistence(
		t,
		diskSchema,
		mkDiskSpeed,
	)

	test.AssertOptionalArguments(t, diskSpeedSchema, []string{
		mkDiskSpeedRead,
		mkDiskSpeedReadBurstable,
		mkDiskSpeedWrite,
		mkDiskSpeedWriteBurstable,
	})

	test.AssertValueTypes(t, diskSpeedSchema, map[string]schema.ValueType{
		mkDiskSpeedRead:           schema.TypeInt,
		mkDiskSpeedReadBurstable:  schema.TypeInt,
		mkDiskSpeedWrite:          schema.TypeInt,
		mkDiskSpeedWriteBurstable: schema.TypeInt,
	})
}
