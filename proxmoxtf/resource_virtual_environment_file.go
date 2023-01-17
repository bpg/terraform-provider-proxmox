/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

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
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
)

const (
	dvResourceVirtualEnvironmentFileContentType        = ""
	dvResourceVirtualEnvironmentFileSourceData         = ""
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

func resourceVirtualEnvironmentFile() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentFileContentType: {
				Type:             schema.TypeString,
				Description:      "The content type",
				Optional:         true,
				ForceNew:         true,
				Default:          dvResourceVirtualEnvironmentFileContentType,
				ValidateDiagFunc: getContentTypeValidator(),
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
		CreateContext: resourceVirtualEnvironmentFileCreate,
		ReadContext:   resourceVirtualEnvironmentFileRead,
		DeleteContext: resourceVirtualEnvironmentFileDelete,
	}
}

func resourceVirtualEnvironmentFileCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	contentType, dg := resourceVirtualEnvironmentFileGetContentType(d)
	diags = append(diags, dg...)

	datastoreID := d.Get(mkResourceVirtualEnvironmentFileDatastoreID).(string)
	fileName, err := resourceVirtualEnvironmentFileGetFileName(d)
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

		if resourceVirtualEnvironmentFileIsURL(d) {
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
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					tflog.Error(ctx, "Failed to close body", map[string]interface{}{
						"error": err,
					})
				}
			}(res.Body)

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

	body := &proxmox.VirtualEnvironmentDatastoreUploadRequestBody{
		ContentType: *contentType,
		DatastoreID: datastoreID,
		FileName:    *fileName,
		FileReader:  file,
		NodeName:    nodeName,
	}

	_, err = veClient.UploadFileToDatastore(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	volumeID, diags := resourceVirtualEnvironmentFileGetVolumeID(d)
	if diags.HasError() {
		return diags
	}

	d.SetId(*volumeID)

	return resourceVirtualEnvironmentFileRead(ctx, d, m)
}

func resourceVirtualEnvironmentFileGetContentType(
	d *schema.ResourceData,
) (*string, diag.Diagnostics) {
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

	ctValidator := getContentTypeValidator()
	diags := ctValidator(contentType, cty.GetAttrPath(mkResourceVirtualEnvironmentFileContentType))

	return &contentType, diags
}

func resourceVirtualEnvironmentFileGetFileName(d *schema.ResourceData) (*string, error) {
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
		if resourceVirtualEnvironmentFileIsURL(d) {
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

func resourceVirtualEnvironmentFileGetVolumeID(d *schema.ResourceData) (*string, diag.Diagnostics) {
	fileName, err := resourceVirtualEnvironmentFileGetFileName(d)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	datastoreID := d.Get(mkResourceVirtualEnvironmentFileDatastoreID).(string)
	contentType, diags := resourceVirtualEnvironmentFileGetContentType(d)

	volumeID := fmt.Sprintf("%s:%s/%s", datastoreID, *contentType, *fileName)

	return &volumeID, diags
}

func resourceVirtualEnvironmentFileIsURL(d *schema.ResourceData) bool {
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

func resourceVirtualEnvironmentFileRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
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

	list, err := veClient.ListDatastoreFiles(ctx, nodeName, datastoreID)
	if err != nil {
		return diag.FromErr(err)
	}

	fileIsURL := resourceVirtualEnvironmentFileIsURL(d)
	fileName, err := resourceVirtualEnvironmentFileGetFileName(d)
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

			sourceFileBlock[mkResourceVirtualEnvironmentFileSourceFileChanged] = lastFileMD != fileModificationDate ||
				lastFileSize != fileSize ||
				lastFileTag != fileTag
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

func readFile(
	ctx context.Context,
	sourceFilePath string,
) (fileModificationDate string, fileSize int64, fileTag string, err error) {
	f, err := os.Open(sourceFilePath)
	if err != nil {
		return
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			tflog.Error(ctx, "failed to close the file", map[string]interface{}{
				"error": err.Error(),
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

func readURL(
	ctx context.Context,
	d *schema.ResourceData,
	sourceFilePath string,
) (fileSize int64, fileModificationDate string, fileTag string, err error) {
	res, err := http.Head(sourceFilePath)
	if err != nil {
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			tflog.Error(ctx, "failed to close the response body", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}(res.Body)

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

func resourceVirtualEnvironmentFileDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	datastoreID := d.Get(mkResourceVirtualEnvironmentFileDatastoreID).(string)
	nodeName := d.Get(mkResourceVirtualEnvironmentFileNodeName).(string)

	err = veClient.DeleteDatastoreFile(ctx, nodeName, datastoreID, d.Id())

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
