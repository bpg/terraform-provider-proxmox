/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	})

	testOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentFileContentType,
		mkResourceVirtualEnvironmentFileSourceFile,
		mkResourceVirtualEnvironmentFileSourceRaw,
	})

	testComputedAttributes(t, s, []string{
		mkResourceVirtualEnvironmentFileFileModificationDate,
		mkResourceVirtualEnvironmentFileFileName,
		mkResourceVirtualEnvironmentFileFileSize,
		mkResourceVirtualEnvironmentFileFileTag,
	})

	testValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentFileContentType:          schema.TypeString,
		mkResourceVirtualEnvironmentFileDatastoreID:          schema.TypeString,
		mkResourceVirtualEnvironmentFileFileModificationDate: schema.TypeString,
		mkResourceVirtualEnvironmentFileFileName:             schema.TypeString,
		mkResourceVirtualEnvironmentFileFileSize:             schema.TypeInt,
		mkResourceVirtualEnvironmentFileFileTag:              schema.TypeString,
		mkResourceVirtualEnvironmentFileNodeName:             schema.TypeString,
		mkResourceVirtualEnvironmentFileSourceFile:           schema.TypeList,
		mkResourceVirtualEnvironmentFileSourceRaw:            schema.TypeList,
	})

	sourceFileSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentFileSourceFile)

	testRequiredArguments(t, sourceFileSchema, []string{
		mkResourceVirtualEnvironmentFileSourceFilePath,
	})

	testOptionalArguments(t, sourceFileSchema, []string{
		mkResourceVirtualEnvironmentFileSourceFileChanged,
		mkResourceVirtualEnvironmentFileSourceFileChecksum,
		mkResourceVirtualEnvironmentFileSourceFileFileName,
		mkResourceVirtualEnvironmentFileSourceFileInsecure,
	})

	testValueTypes(t, sourceFileSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentFileSourceFileChanged:  schema.TypeBool,
		mkResourceVirtualEnvironmentFileSourceFileChecksum: schema.TypeString,
		mkResourceVirtualEnvironmentFileSourceFileFileName: schema.TypeString,
		mkResourceVirtualEnvironmentFileSourceFileInsecure: schema.TypeBool,
		mkResourceVirtualEnvironmentFileSourceFilePath:     schema.TypeString,
	})

	sourceRawSchema := testNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentFileSourceRaw)

	testRequiredArguments(t, sourceRawSchema, []string{
		mkResourceVirtualEnvironmentFileSourceRawData,
		mkResourceVirtualEnvironmentFileSourceRawFileName,
	})

	testOptionalArguments(t, sourceRawSchema, []string{
		mkResourceVirtualEnvironmentFileSourceRawResize,
	})

	testValueTypes(t, sourceRawSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentFileSourceRawData:     schema.TypeString,
		mkResourceVirtualEnvironmentFileSourceRawFileName: schema.TypeString,
		mkResourceVirtualEnvironmentFileSourceRawResize:   schema.TypeInt,
	})
}
