/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/resource/validator"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

const (
	dvResourceVirtualEnvironmentFileContentType        = ""
	dvResourceVirtualEnvironmentFileSourceFileChanged  = false
	dvResourceVirtualEnvironmentFileSourceFileChecksum = ""
	dvResourceVirtualEnvironmentFileSourceFileFileName = ""
	dvResourceVirtualEnvironmentFileSourceFileInsecure = false
	dvResourceVirtualEnvironmentFileSourceRawResize    = 0

	mkResourceVirtualEnvironmentFileContentType          = "content_type"
	mkResourceVirtualEnvironmentFileDatastoreID          = "datastore_id"
	mkResourceVirtualEnvironmentFileFileModificationDate = "file_modification_date"
	mkResourceVirtualEnvironmentFileFileName             = "file_name"
	mkResourceVirtualEnvironmentFileFileSize             = "file_size"
	mkResourceVirtualEnvironmentFileFileTag              = "file_tag"
	mkResourceVirtualEnvironmentFileNodeName             = "node_name"
	mkResourceVirtualEnvironmentFileSourceFile           = "source_file"
	mkResourceVirtualEnvironmentFileSourceFilePath       = "path"
	mkResourceVirtualEnvironmentFileSourceFileChanged    = "changed"
	mkResourceVirtualEnvironmentFileSourceFileChecksum   = "checksum"
	mkResourceVirtualEnvironmentFileSourceFileFileName   = "file_name"
	mkResourceVirtualEnvironmentFileSourceFileInsecure   = "insecure"
	mkResourceVirtualEnvironmentFileSourceRaw            = "source_raw"
	mkResourceVirtualEnvironmentFileSourceRawData        = "data"
	mkResourceVirtualEnvironmentFileSourceRawFileName    = "file_name"
	mkResourceVirtualEnvironmentFileSourceRawResize      = "resize"
)

// File returns a resource that manages files on a node.
func File() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentFileContentType: {
				Type:             schema.TypeString,
				Description:      "The content type",
				Optional:         true,
				ForceNew:         true,
				Default:          dvResourceVirtualEnvironmentFileContentType,
				ValidateDiagFunc: validator.ContentType(),
			},
			mkResourceVirtualEnvironmentFileDatastoreID: {
				Type:        schema.TypeString,
				Description: "The datastore id",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentFileFileModificationDate: {
				Type:        schema.TypeString,
				Description: "The file modification date",
				Computed:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentFileFileName: {
				Type:        schema.TypeString,
				Description: "The file name",
				Computed:    true,
			},
			mkResourceVirtualEnvironmentFileFileSize: {
				Type:        schema.TypeInt,
				Description: "The file size in bytes",
				Computed:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentFileFileTag: {
				Type:        schema.TypeString,
				Description: "The file tag",
				Computed:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentFileNodeName: {
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentFileSourceFile: {
				Type:        schema.TypeList,
				Description: "The source file",
				Optional:    true,
				ForceNew:    true,
				DefaultFunc: func() (interface{}, error) {
					return make([]interface{}, 1), nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentFileSourceFilePath: {
							Type:        schema.TypeString,
							Description: "A path to a local file or a URL",
							Required:    true,
							ForceNew:    true,
						},
						mkResourceVirtualEnvironmentFileSourceFileChanged: {
							Type:        schema.TypeBool,
							Description: "Whether the source file has changed since the last run",
							Optional:    true,
							ForceNew:    true,
							Default:     dvResourceVirtualEnvironmentFileSourceFileChanged,
						},
						mkResourceVirtualEnvironmentFileSourceFileChecksum: {
							Type:        schema.TypeString,
							Description: "The SHA256 checksum of the source file",
							Optional:    true,
							ForceNew:    true,
							Default:     dvResourceVirtualEnvironmentFileSourceFileChecksum,
						},
						mkResourceVirtualEnvironmentFileSourceFileFileName: {
							Type:        schema.TypeString,
							Description: "The file name to use instead of the source file name",
							Optional:    true,
							ForceNew:    true,
							Default:     dvResourceVirtualEnvironmentFileSourceFileFileName,
						},
						mkResourceVirtualEnvironmentFileSourceFileInsecure: {
							Type:        schema.TypeBool,
							Description: "Whether to skip the TLS verification step for HTTPS sources",
							Optional:    true,
							ForceNew:    true,
							Default:     dvResourceVirtualEnvironmentFileSourceFileInsecure,
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentFileSourceRaw: {
				Type:        schema.TypeList,
				Description: "The raw source",
				Optional:    true,
				ForceNew:    true,
				DefaultFunc: func() (interface{}, error) {
					return make([]interface{}, 1), nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentFileSourceRawData: {
							Type:        schema.TypeString,
							Description: "The raw data",
							Required:    true,
							ForceNew:    true,
						},
						mkResourceVirtualEnvironmentFileSourceRawFileName: {
							Type:        schema.TypeString,
							Description: "The file name",
							Required:    true,
							ForceNew:    true,
						},
						mkResourceVirtualEnvironmentFileSourceRawResize: {
							Type:        schema.TypeInt,
							Description: "The number of bytes to resize the file to",
							Optional:    true,
							ForceNew:    true,
							Default:     dvResourceVirtualEnvironmentFileSourceRawResize,
						},
					},
				},
				MaxItems: 1,
				MinItems: 0,
			},
		},
		CreateContext: fileCreate,
		ReadContext:   fileRead,
		DeleteContext: fileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				node, datastore, volumeID, err := fileParseImportID(d.Id())
				if err != nil {
					return nil, err
				}

				d.SetId(volumeID)

				err = d.Set(mkResourceVirtualEnvironmentFileNodeName, node)
				if err != nil {
					return nil, fmt.Errorf("failed setting state during import: %w", err)
				}

				err = d.Set(mkResourceVirtualEnvironmentFileDatastoreID, datastore)
				if err != nil {
					return nil, fmt.Errorf("failed setting state during import: %w", err)
				}

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func fileParseImportID(id string) (string, string, string, error) {
	parts := strings.SplitN(id, "/", 4)

	if len(parts) != 4 || parts[0] == "" || parts[1] == "" || parts[2] == "" || parts[3] == "" {
		return "", "", "", fmt.Errorf("unexpected format of ID (%s), expected node/datastore_id/content_type/file_name", id)
	}

	return parts[0], parts[1], strings.Join(parts[2:], "/"), nil
}

func fileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	contentType, dg := fileGetContentType(d)
	diags = append(diags, dg...)

	datastoreID := d.Get(mkResourceVirtualEnvironmentFileDatastoreID).(string)
	fileName, err := fileGetFileName(d)
	diags = append(diags, diag.FromErr(err)...)

	nodeName := d.Get(mkResourceVirtualEnvironmentFileNodeName).(string)
	sourceFile := d.Get(mkResourceVirtualEnvironmentFileSourceFile).([]interface{})
	sourceRaw := d.Get(mkResourceVirtualEnvironmentFileSourceRaw).([]interface{})

	sourceFilePathLocal := ""

	// Determine if both source_data and source_file is specified as this is not supported.
	if len(sourceFile) > 0 && len(sourceRaw) > 0 {
		diags = append(diags, diag.Errorf(
			"please specify \"%s.%s\" or \"%s\" - not both",
			mkResourceVirtualEnvironmentFileSourceFile,
			mkResourceVirtualEnvironmentFileSourceFilePath,
			mkResourceVirtualEnvironmentFileSourceRaw,
		)...)
	}

	if diags.HasError() {
		return diags
	}

	// Determine if we're dealing with raw file data or a reference to a file or URL.
	// In case of a URL, we must first download the file before proceeding.
	// This is due to lack of support for chunked transfers in the Proxmox VE API.
	if len(sourceFile) > 0 {
		sourceFileBlock := sourceFile[0].(map[string]interface{})
		sourceFilePath := sourceFileBlock[mkResourceVirtualEnvironmentFileSourceFilePath].(string)
		sourceFileChecksum := sourceFileBlock[mkResourceVirtualEnvironmentFileSourceFileChecksum].(string)
		sourceFileInsecure := sourceFileBlock[mkResourceVirtualEnvironmentFileSourceFileInsecure].(bool)

		if fileIsURL(d) {
			tflog.Debug(ctx, "Downloading file from URL", map[string]interface{}{
				"url": sourceFilePath,
			})

			httpClient := http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: sourceFileInsecure,
					},
				},
			}

			res, err := httpClient.Get(sourceFilePath)
			if err != nil {
				return diag.FromErr(err)
			}

			defer utils.CloseOrLogError(ctx)(res.Body)

			tempDownloadedFile, err := os.CreateTemp("", "download")
			if err != nil {
				return diag.FromErr(err)
			}

			tempDownloadedFileName := tempDownloadedFile.Name()
			defer func(name string) {
				err := os.Remove(name)
				if err != nil {
					tflog.Error(ctx, "Failed to remove temporary file", map[string]interface{}{
						"error": err,
						"file":  name,
					})
				}
			}(tempDownloadedFileName)

			_, err = io.Copy(tempDownloadedFile, res.Body)
			diags = append(diags, diag.FromErr(err)...)
			err = tempDownloadedFile.Close()
			diags = append(diags, diag.FromErr(err)...)

			if diags.HasError() {
				return diags
			}

			sourceFilePathLocal = tempDownloadedFileName
		} else {
			sourceFilePathLocal = sourceFilePath
		}

		// Calculate the checksum of the source file now that it's available locally.
		if sourceFileChecksum != "" {
			file, err := os.Open(sourceFilePathLocal)
			if err != nil {
				return diag.FromErr(err)
			}

			h := sha256.New()
			_, err = io.Copy(h, file)
			diags = append(diags, diag.FromErr(err)...)
			err = file.Close()
			diags = append(diags, diag.FromErr(err)...)
			if diags.HasError() {
				return diags
			}

			calculatedChecksum := fmt.Sprintf("%x", h.Sum(nil))
			tflog.Debug(ctx, "Calculated checksum", map[string]interface{}{
				"source": sourceFilePath,
				"sha256": calculatedChecksum,
			})

			if sourceFileChecksum != calculatedChecksum {
				return diag.Errorf(
					"the calculated SHA256 checksum \"%s\" does not match source checksum \"%s\"",
					calculatedChecksum,
					sourceFileChecksum,
				)
			}
		}
	} else if len(sourceRaw) > 0 {
		sourceRawBlock := sourceRaw[0].(map[string]interface{})
		sourceRawData := sourceRawBlock[mkResourceVirtualEnvironmentFileSourceRawData].(string)
		sourceRawResize := sourceRawBlock[mkResourceVirtualEnvironmentFileSourceRawResize].(int)

		if sourceRawResize > 0 {
			if len(sourceRawData) <= sourceRawResize {
				sourceRawData = fmt.Sprintf(fmt.Sprintf("%%-%dv", sourceRawResize), sourceRawData)
			} else {
				return diag.Errorf("cannot resize %d bytes to %d bytes", len(sourceRawData), sourceRawResize)
			}
		}

		tempRawFile, err := os.CreateTemp("", "raw")
		if err != nil {
			return diag.FromErr(err)
		}

		tempRawFileName := tempRawFile.Name()
		_, err = io.Copy(tempRawFile, bytes.NewBufferString(sourceRawData))
		diags = append(diags, diag.FromErr(err)...)
		err = tempRawFile.Close()
		diags = append(diags, diag.FromErr(err)...)
		if diags.HasError() {
			return diags
		}

		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				tflog.Error(ctx, "Failed to remove temporary file", map[string]interface{}{
					"error": err,
					"file":  name,
				})
			}
		}(tempRawFileName)

		sourceFilePathLocal = tempRawFileName
	} else {
		return diag.Errorf(
			"please specify either \"%s.%s\" or \"%s\"",
			mkResourceVirtualEnvironmentFileSourceFile,
			mkResourceVirtualEnvironmentFileSourceFilePath,
			mkResourceVirtualEnvironmentFileSourceRaw,
		)
	}

	// Open the source file for reading in order to upload it.
	file, err := os.Open(sourceFilePathLocal)
	if err != nil {
		return diag.FromErr(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			tflog.Error(ctx, "Failed to close file", map[string]interface{}{
				"error": err,
			})
		}
	}(file)

	config := m.(proxmoxtf.ProviderConfiguration)

	capi, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	request := &api.FileUploadRequest{
		ContentType: *contentType,
		FileName:    *fileName,
		File:        file,
	}

	switch *contentType {
	case "iso", "vztmpl":
		_, err = capi.Node(nodeName).APIUpload(ctx, datastoreID, request)
	default:
		// For all other content types, we need to upload the file to the node's
		// datastore using SFTP.
		datastore, err2 := capi.Storage().GetDatastore(ctx, datastoreID)
		if err2 != nil {
			return diag.Errorf("failed to get datastore: %s", err2)
		}

		if datastore.Path == nil || *datastore.Path == "" {
			return diag.Errorf("failed to determine the datastore path")
		}

		sort.Strings(datastore.Content)

		_, found := slices.BinarySearch(datastore.Content, *contentType)
		if !found {
			diags = append(diags, diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary: fmt.Sprintf("the datastore %q does not support content type %q; supported content types are: %v",
						*datastore.Storage, *contentType, datastore.Content,
					),
				},
			}...)
		}

		remoteFileDir := *datastore.Path

		err = capi.SSH().NodeUpload(ctx, nodeName, remoteFileDir, request)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	volumeID, di := fileGetVolumeID(d)
	diags = append(diags, di...)
	if diags.HasError() {
		return diags
	}

	d.SetId(*volumeID)

	diags = append(diags, fileRead(ctx, d, m)...)

	if d.Id() == "" {
		diags = append(diags, diag.Errorf("failed to read file from %q", *volumeID)...)
	}

	return diags
}

func fileGetContentType(d *schema.ResourceData) (*string, diag.Diagnostics) {
	contentType := d.Get(mkResourceVirtualEnvironmentFileContentType).(string)
	sourceFile := d.Get(mkResourceVirtualEnvironmentFileSourceFile).([]interface{})
	sourceRaw := d.Get(mkResourceVirtualEnvironmentFileSourceRaw).([]interface{})

	sourceFilePath := ""

	if len(sourceFile) > 0 {
		sourceFileBlock := sourceFile[0].(map[string]interface{})
		sourceFilePath = sourceFileBlock[mkResourceVirtualEnvironmentFileSourceFilePath].(string)
	} else if len(sourceRaw) > 0 {
		sourceRawBlock := sourceRaw[0].(map[string]interface{})
		sourceFilePath = sourceRawBlock[mkResourceVirtualEnvironmentFileSourceRawFileName].(string)
	} else {
		return nil, diag.Errorf(
			"missing argument \"%s.%s\" or \"%s\"",
			mkResourceVirtualEnvironmentFileSourceFile,
			mkResourceVirtualEnvironmentFileSourceFilePath,
			mkResourceVirtualEnvironmentFileSourceRaw,
		)
	}

	if contentType == "" {
		if strings.HasSuffix(sourceFilePath, ".tar.gz") ||
			strings.HasSuffix(sourceFilePath, ".tar.xz") {
			contentType = "vztmpl"
		} else {
			ext := strings.TrimLeft(strings.ToLower(filepath.Ext(sourceFilePath)), ".")

			switch ext {
			case "img", "iso":
				contentType = "iso"
			case "yaml", "yml":
				contentType = "snippets"
			}
		}

		if contentType == "" {
			return nil, diag.Errorf(
				"cannot determine the content type of source \"%s\" - Please manually define the \"%s\" argument",
				sourceFilePath,
				mkResourceVirtualEnvironmentFileContentType,
			)
		}
	}

	ctValidator := validator.ContentType()
	diags := ctValidator(contentType, cty.GetAttrPath(mkResourceVirtualEnvironmentFileContentType))

	return &contentType, diags
}

func fileGetFileName(d *schema.ResourceData) (*string, error) {
	sourceFile := d.Get(mkResourceVirtualEnvironmentFileSourceFile).([]interface{})
	sourceRaw := d.Get(mkResourceVirtualEnvironmentFileSourceRaw).([]interface{})

	sourceFileFileName := ""
	sourceFilePath := ""

	if len(sourceFile) > 0 {
		sourceFileBlock := sourceFile[0].(map[string]interface{})
		sourceFileFileName = sourceFileBlock[mkResourceVirtualEnvironmentFileSourceFileFileName].(string)
		sourceFilePath = sourceFileBlock[mkResourceVirtualEnvironmentFileSourceFilePath].(string)
	} else if len(sourceRaw) > 0 {
		sourceRawBlock := sourceRaw[0].(map[string]interface{})
		sourceFileFileName = sourceRawBlock[mkResourceVirtualEnvironmentFileSourceRawFileName].(string)
	} else {
		return nil, fmt.Errorf(
			"missing argument \"%s.%s\"",
			mkResourceVirtualEnvironmentFileSourceFile,
			mkResourceVirtualEnvironmentFileSourceFilePath,
		)
	}

	if sourceFileFileName == "" {
		if fileIsURL(d) {
			downloadURL, err := url.ParseRequestURI(sourceFilePath)
			if err != nil {
				return nil, err
			}

			path := strings.Split(downloadURL.Path, "/")
			sourceFileFileName = path[len(path)-1]

			if sourceFileFileName == "" {
				return nil, fmt.Errorf(
					"failed to determine file name from the URL \"%s\"",
					sourceFilePath,
				)
			}
		} else {
			sourceFileFileName = filepath.Base(sourceFilePath)
		}
	}

	return &sourceFileFileName, nil
}

func fileGetVolumeID(d *schema.ResourceData) (*string, diag.Diagnostics) {
	fileName, err := fileGetFileName(d)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	datastoreID := d.Get(mkResourceVirtualEnvironmentFileDatastoreID).(string)
	contentType, diags := fileGetContentType(d)

	volumeID := fmt.Sprintf("%s:%s/%s", datastoreID, *contentType, *fileName)

	return &volumeID, diags
}

func fileIsURL(d *schema.ResourceData) bool {
	sourceFile := d.Get(mkResourceVirtualEnvironmentFileSourceFile).([]interface{})
	sourceFilePath := ""

	if len(sourceFile) > 0 {
		sourceFileBlock := sourceFile[0].(map[string]interface{})
		sourceFilePath = sourceFileBlock[mkResourceVirtualEnvironmentFileSourceFilePath].(string)
	} else {
		return false
	}

	return strings.HasPrefix(sourceFilePath, "http://") ||
		strings.HasPrefix(sourceFilePath, "https://")
}

func fileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	capi, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	datastoreID := d.Get(mkResourceVirtualEnvironmentFileDatastoreID).(string)
	nodeName := d.Get(mkResourceVirtualEnvironmentFileNodeName).(string)
	sourceFile := d.Get(mkResourceVirtualEnvironmentFileSourceFile).([]interface{})
	sourceFilePath := ""

	if len(sourceFile) == 0 {
		return nil
	}

	sourceFileBlock := sourceFile[0].(map[string]interface{})
	sourceFilePath = sourceFileBlock[mkResourceVirtualEnvironmentFileSourceFilePath].(string)

	list, err := capi.Node(nodeName).ListDatastoreFiles(ctx, datastoreID)
	if err != nil {
		return diag.FromErr(err)
	}

	fileIsURL := fileIsURL(d)
	fileName, err := fileGetFileName(d)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	for _, v := range list {
		if v.VolumeID == d.Id() {
			var fileModificationDate string
			var fileSize int64
			var fileTag string

			if fileIsURL {
				fileSize, fileModificationDate, fileTag, err = readURL(ctx, d, sourceFilePath)
			} else {
				fileModificationDate, fileSize, fileTag, err = readFile(ctx, sourceFilePath)
			}
			diags = append(diags, diag.FromErr(err)...)

			lastFileMD := d.Get(mkResourceVirtualEnvironmentFileFileModificationDate).(string)
			lastFileSize := int64(d.Get(mkResourceVirtualEnvironmentFileFileSize).(int))
			lastFileTag := d.Get(mkResourceVirtualEnvironmentFileFileTag).(string)

			err = d.Set(mkResourceVirtualEnvironmentFileFileModificationDate, fileModificationDate)
			diags = append(diags, diag.FromErr(err)...)
			err = d.Set(mkResourceVirtualEnvironmentFileFileName, *fileName)
			diags = append(diags, diag.FromErr(err)...)
			err = d.Set(mkResourceVirtualEnvironmentFileFileSize, fileSize)
			diags = append(diags, diag.FromErr(err)...)
			err = d.Set(mkResourceVirtualEnvironmentFileFileTag, fileTag)
			diags = append(diags, diag.FromErr(err)...)

			// just to make the logic easier to read
			changed := false
			if lastFileMD != "" && lastFileSize != 0 && lastFileTag != "" {
				changed = lastFileMD != fileModificationDate || lastFileSize != fileSize || lastFileTag != fileTag
			}

			sourceFileBlock[mkResourceVirtualEnvironmentFileSourceFileChanged] = changed
			err = d.Set(mkResourceVirtualEnvironmentFileSourceFile, sourceFile)
			diags = append(diags, diag.FromErr(err)...)

			if diags.HasError() {
				return diags
			}
			return nil
		}
	}

	d.SetId("")

	return nil
}

//nolint:nonamedreturns
func readFile(
	ctx context.Context,
	sourceFilePath string,
) (fileModificationDate string, fileSize int64, fileTag string, err error) {
	f, err := os.Open(sourceFilePath)
	if err != nil {
		return
	}

	defer func(f *os.File) {
		e := f.Close()
		if e != nil {
			tflog.Error(ctx, "failed to close the file", map[string]interface{}{
				"error": e.Error(),
			})
		}
	}(f)

	fileInfo, err := f.Stat()
	if err != nil {
		return
	}

	fileModificationDate = fileInfo.ModTime().UTC().Format(time.RFC3339)
	fileSize = fileInfo.Size()
	fileTag = fmt.Sprintf("%x-%x", fileInfo.ModTime().UTC().Unix(), fileInfo.Size())

	return fileModificationDate, fileSize, fileTag, nil
}

//nolint:nonamedreturns
func readURL(
	ctx context.Context,
	d *schema.ResourceData,
	sourceFilePath string,
) (fileSize int64, fileModificationDate string, fileTag string, err error) {
	res, err := http.Head(sourceFilePath)
	if err != nil {
		return
	}

	defer utils.CloseOrLogError(ctx)(res.Body)

	fileSize = res.ContentLength
	httpLastModified := res.Header.Get("Last-Modified")

	if httpLastModified != "" {
		var timeParsed time.Time
		timeParsed, err = time.Parse(time.RFC1123, httpLastModified)

		if err != nil {
			timeParsed, err = time.Parse(time.RFC1123Z, httpLastModified)
			if err != nil {
				return
			}
		}

		fileModificationDate = timeParsed.UTC().Format(time.RFC3339)
	} else {
		err = d.Set(mkResourceVirtualEnvironmentFileFileModificationDate, "")
		if err != nil {
			return
		}
	}

	httpTag := res.Header.Get("ETag")

	if httpTag != "" {
		httpTagParts := strings.Split(httpTag, "\"")

		if len(httpTagParts) > 1 {
			fileTag = httpTagParts[1]
		} else {
			fileTag = ""
		}
	} else {
		fileTag = ""
	}

	return
}

func fileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	capi, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	datastoreID := d.Get(mkResourceVirtualEnvironmentFileDatastoreID).(string)
	nodeName := d.Get(mkResourceVirtualEnvironmentFileNodeName).(string)

	err = capi.Node(nodeName).DeleteDatastoreFile(ctx, datastoreID, d.Id())

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
