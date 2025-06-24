/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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

	s := File().Schema

	test.AssertRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentFileDatastoreID,
		mkResourceVirtualEnvironmentFileNodeName,
	})

	test.AssertOptionalArguments(t, s, []string{
		mkResourceVirtualEnvironmentFileContentType,
		mkResourceVirtualEnvironmentFileSourceFile,
		mkResourceVirtualEnvironmentFileFileMode,
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
		mkResourceVirtualEnvironmentFileFileMode:             schema.TypeString,
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

func Test_fileParseVolumeID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		id      string
		want    fileVolumeID
		wantErr bool
	}{
		{"empty", "", fileVolumeID{}, true},
		{"missing datastore", "iso/file.ido", fileVolumeID{}, true},
		{"missing type", "local:/file.ido", fileVolumeID{}, true},
		{"missing file", "local:iso", fileVolumeID{}, true},
		{"missing file", "local:iso/", fileVolumeID{}, true},
		{"valid iso", "local:iso/file.iso", fileVolumeID{
			datastoreID: "local",
			contentType: "iso",
			fileName:    "file.iso",
		}, false},
		{"valid import", "local:import/file.qcow2", fileVolumeID{
			datastoreID: "local",
			contentType: "import",
			fileName:    "file.qcow2",
		}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := fileParseVolumeID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileParseVolumeID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fileParseVolumeID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileParseImportID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		id      string
		node    string
		volID   fileVolumeID
		wantErr bool
	}{
		{"empty", "", "", fileVolumeID{}, true},
		{"missing node", "local:iso/file.iso", "", fileVolumeID{}, true},
		{"missing node 2", "/local:iso/file.iso", "", fileVolumeID{}, true},
		{
			"valid", "pve/local:iso/file.iso",
			"pve",
			fileVolumeID{
				datastoreID: "local",
				contentType: "iso",
				fileName:    "file.iso",
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			node, volID, err := fileParseImportID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileParseImportID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if node != tt.node {
				t.Errorf("fileParseImportID() got node = %v, want %v", node, tt.node)
			}

			if !reflect.DeepEqual(volID, tt.volID) {
				t.Errorf("fileParseImportID() got volID = %v, want %v", volID, tt.volID)
			}
		})
	}
}
