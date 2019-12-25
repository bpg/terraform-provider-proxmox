/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

// TestResourceVirtualEnvironmentFileInstantiation tests whether the ResourceVirtualEnvironmentFile instance can be instantiated.
func TestResourceVirtualEnvironmentFileInstantiation(t *testing.T) {
	s := resourceVirtualEnvironmentFile()

	if s == nil {
		t.Fatalf("Cannot instantiate resourceVirtualEnvironmentFile")
	}
}

// TestResourceVirtualEnvironmentFileSchema tests the resourceVirtualEnvironmentFile schema.
func TestResourceVirtualEnvironmentFileSchema(t *testing.T) {
	s := resourceVirtualEnvironmentFile()

	testRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentFileDatastoreID,
		mkResourceVirtualEnvironmentFileNodeName,
		mkResourceVirtualEnvironmentFileSource,
	})

	testOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentFileContentType,
		mkResourceVirtualEnvironmentFileOverrideFileName,
		mkResourceVirtualEnvironmentFileSourceChanged,
		mkResourceVirtualEnvironmentFileSourceChecksum,
		mkResourceVirtualEnvironmentFileSourceInsecure,
	})

	testComputedAttributes(t, s, []string{
		mkResourceVirtualEnvironmentFileFileModificationDate,
		mkResourceVirtualEnvironmentFileFileName,
		mkResourceVirtualEnvironmentFileFileSize,
		mkResourceVirtualEnvironmentFileFileTag,
	})

	testSchemaValueTypes(t, s, []string{
		mkResourceVirtualEnvironmentFileContentType,
		mkResourceVirtualEnvironmentFileDatastoreID,
		mkResourceVirtualEnvironmentFileFileModificationDate,
		mkResourceVirtualEnvironmentFileFileName,
		mkResourceVirtualEnvironmentFileFileSize,
		mkResourceVirtualEnvironmentFileFileTag,
		mkResourceVirtualEnvironmentFileOverrideFileName,
		mkResourceVirtualEnvironmentFileSourceChanged,
		mkResourceVirtualEnvironmentFileNodeName,
		mkResourceVirtualEnvironmentFileSource,
		mkResourceVirtualEnvironmentFileSourceChecksum,
		mkResourceVirtualEnvironmentFileSourceInsecure,
	}, []schema.ValueType{
		schema.TypeString,
		schema.TypeString,
		schema.TypeString,
		schema.TypeString,
		schema.TypeInt,
		schema.TypeString,
		schema.TypeString,
		schema.TypeBool,
		schema.TypeString,
		schema.TypeString,
		schema.TypeString,
		schema.TypeBool,
	})
}
