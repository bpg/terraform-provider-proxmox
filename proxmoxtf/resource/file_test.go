/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

// TestFileInstantiation tests whether the File instance can be instantiated.
func TestFileInstantiation(t *testing.T) {
	t.Parallel()

	s := File()
	if s == nil {
		t.Fatalf("Cannot instantiate File")
	}
}

// TestFileSchema tests the File schema.
func TestFileSchema(t *testing.T) {
	t.Parallel()

	s := File()

	test.AssertRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentFileDatastoreID,
		mkResourceVirtualEnvironmentFileNodeName,
	})

	test.AssertOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentFileContentType,
		mkResourceVirtualEnvironmentFileSourceFile,
		mkResourceVirtualEnvironmentFileSourceRaw,
		mkResourceVirtualEnvironmentFileTimeoutUpload,
	})

	test.AssertComputedAttributes(t, s, []string{
		mkResourceVirtualEnvironmentFileFileModificationDate,
		mkResourceVirtualEnvironmentFileFileName,
		mkResourceVirtualEnvironmentFileFileSize,
		mkResourceVirtualEnvironmentFileFileTag,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentFileContentType:          schema.TypeString,
		mkResourceVirtualEnvironmentFileDatastoreID:          schema.TypeString,
		mkResourceVirtualEnvironmentFileFileModificationDate: schema.TypeString,
		mkResourceVirtualEnvironmentFileFileName:             schema.TypeString,
		mkResourceVirtualEnvironmentFileFileSize:             schema.TypeInt,
		mkResourceVirtualEnvironmentFileFileTag:              schema.TypeString,
		mkResourceVirtualEnvironmentFileNodeName:             schema.TypeString,
		mkResourceVirtualEnvironmentFileSourceFile:           schema.TypeList,
		mkResourceVirtualEnvironmentFileSourceRaw:            schema.TypeList,
		mkResourceVirtualEnvironmentFileTimeoutUpload:        schema.TypeInt,
	})

	sourceFileSchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentFileSourceFile)

	test.AssertRequiredArguments(t, sourceFileSchema, []string{
		mkResourceVirtualEnvironmentFileSourceFilePath,
	})

	test.AssertOptionalArguments(t, sourceFileSchema, []string{
		mkResourceVirtualEnvironmentFileSourceFileChanged,
		mkResourceVirtualEnvironmentFileSourceFileChecksum,
		mkResourceVirtualEnvironmentFileSourceFileFileName,
		mkResourceVirtualEnvironmentFileSourceFileInsecure,
	})

	test.AssertValueTypes(t, sourceFileSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentFileSourceFileChanged:  schema.TypeBool,
		mkResourceVirtualEnvironmentFileSourceFileChecksum: schema.TypeString,
		mkResourceVirtualEnvironmentFileSourceFileFileName: schema.TypeString,
		mkResourceVirtualEnvironmentFileSourceFileInsecure: schema.TypeBool,
		mkResourceVirtualEnvironmentFileSourceFilePath:     schema.TypeString,
	})

	sourceRawSchema := test.AssertNestedSchemaExistence(t, s, mkResourceVirtualEnvironmentFileSourceRaw)

	test.AssertRequiredArguments(t, sourceRawSchema, []string{
		mkResourceVirtualEnvironmentFileSourceRawData,
		mkResourceVirtualEnvironmentFileSourceRawFileName,
	})

	test.AssertOptionalArguments(t, sourceRawSchema, []string{
		mkResourceVirtualEnvironmentFileSourceRawResize,
	})

	test.AssertValueTypes(t, sourceRawSchema, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentFileSourceRawData:     schema.TypeString,
		mkResourceVirtualEnvironmentFileSourceRawFileName: schema.TypeString,
		mkResourceVirtualEnvironmentFileSourceRawResize:   schema.TypeInt,
	})
}

func Test_fileParseImportID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		value               string
		valid               bool
		expectedNodeName    string
		expectedDatastoreID string
		expectedVolumeID    string
	}{
		{"empty", "", false, "", "", ""},
		{"missing slash", "invalid", false, "", "", ""},
		{"missing parts", "invalid/invalid/invalid", false, "", "", ""},
		{"valid", "node/datastore_id/content_type/file_name", true, "node", "datastore_id", "content_type/file_name"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require := require.New(t)

			nodeName, datastoreID, volumeID, err := fileParseImportID(tt.value)

			if !tt.valid {
				require.Error(err)
				return
			}

			require.Nil(err)
			require.Equal(tt.expectedNodeName, nodeName)
			require.Equal(tt.expectedDatastoreID, datastoreID)
			require.Equal(tt.expectedVolumeID, volumeID)
		})
	}
}
