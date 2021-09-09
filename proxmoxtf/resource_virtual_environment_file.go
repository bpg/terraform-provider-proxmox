/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Type:         schema.TypeString,
				Description:  "The content type",
				Optional:     true,
				ForceNew:     true,
				Default:      dvResourceVirtualEnvironmentFileContentType,
				ValidateFunc: getContentTypeValidator(),
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
		Create: resourceVirtualEnvironmentFileCreate,
		Read:   resourceVirtualEnvironmentFileRead,
		Delete: resourceVirtualEnvironmentFileDelete,
	}
}

func resourceVirtualEnvironmentFileCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	contentType, err := resourceVirtualEnvironmentFileGetContentType(d, m)

	if err != nil {
		return err
	}

	datastoreID := d.Get(mkResourceVirtualEnvironmentFileDatastoreID).(string)
	fileName, err := resourceVirtualEnvironmentFileGetFileName(d, m)

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentFileNodeName).(string)
	sourceFile := d.Get(mkResourceVirtualEnvironmentFileSourceFile).([]interface{})
	sourceRaw := d.Get(mkResourceVirtualEnvironmentFileSourceRaw).([]interface{})

	sourceFilePathLocal := ""

	// Determine if both source_data and source_file is specified as this is not supported.
	if len(sourceFile) > 0 && len(sourceRaw) > 0 {
		return fmt.Errorf(
			"Please specify \"%s.%s\" or \"%s\" - not both",
			mkResourceVirtualEnvironmentFileSourceFile,
			mkResourceVirtualEnvironmentFileSourceFilePath,
			mkResourceVirtualEnvironmentFileSourceRaw,
		)
	}

	// Determine if we're dealing with raw file data or a reference to a file or URL.
	// In case of a URL, we must first download the file before proceeding.
	// This is due to lack of support for chunked transfers in the Proxmox VE API.
	if len(sourceFile) > 0 {
		sourceFileBlock := sourceFile[0].(map[string]interface{})
		sourceFilePath := sourceFileBlock[mkResourceVirtualEnvironmentFileSourceFilePath].(string)
		sourceFileChecksum := sourceFileBlock[mkResourceVirtualEnvironmentFileSourceFileChecksum].(string)
		sourceFileInsecure := sourceFileBlock[mkResourceVirtualEnvironmentFileSourceFileInsecure].(bool)

		if resourceVirtualEnvironmentFileIsURL(d, m) {
			log.Printf("[DEBUG] Downloading file from '%s'", sourceFilePath)

			httpClient := http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: sourceFileInsecure,
					},
				},
			}

			res, err := httpClient.Get(sourceFilePath)

			if err != nil {
				return err
			}

			defer res.Body.Close()

			tempDownloadedFile, err := ioutil.TempFile("", "download")

			if err != nil {
				return err
			}

			tempDownloadedFileName := tempDownloadedFile.Name()
			_, err = io.Copy(tempDownloadedFile, res.Body)

			if err != nil {
				tempDownloadedFile.Close()

				return err
			}

			tempDownloadedFile.Close()

			defer os.Remove(tempDownloadedFileName)

			sourceFilePathLocal = tempDownloadedFileName
		} else {
			sourceFilePathLocal = sourceFilePath
		}

		// Calculate the checksum of the source file now that it's available locally.
		if sourceFileChecksum != "" {
			file, err := os.Open(sourceFilePathLocal)

			if err != nil {
				return err
			}

			h := sha256.New()
			_, err = io.Copy(h, file)

			if err != nil {
				file.Close()

				return err
			}

			file.Close()

			calculatedChecksum := fmt.Sprintf("%x", h.Sum(nil))

			log.Printf("[DEBUG] The calculated SHA256 checksum for source \"%s\" is \"%s\"", sourceFilePath, calculatedChecksum)

			if sourceFileChecksum != calculatedChecksum {
				return fmt.Errorf("The calculated SHA256 checksum \"%s\" does not match source checksum \"%s\"", calculatedChecksum, sourceFileChecksum)
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
				return fmt.Errorf("Cannot resize %d bytes to %d bytes", len(sourceRawData), sourceRawResize)
			}
		}

		tempRawFile, err := ioutil.TempFile("", "raw")

		if err != nil {
			return err
		}

		tempRawFileName := tempRawFile.Name()
		_, err = io.Copy(tempRawFile, bytes.NewBufferString(sourceRawData))

		if err != nil {
			tempRawFile.Close()

			return err
		}

		tempRawFile.Close()

		defer os.Remove(tempRawFileName)

		sourceFilePathLocal = tempRawFileName
	} else {
		return fmt.Errorf(
			"Please specify either \"%s.%s\" or \"%s\"",
			mkResourceVirtualEnvironmentFileSourceFile,
			mkResourceVirtualEnvironmentFileSourceFilePath,
			mkResourceVirtualEnvironmentFileSourceRaw,
		)
	}

	// Open the source file for reading in order to upload it.
	file, err := os.Open(sourceFilePathLocal)

	if err != nil {
		return err
	}

	defer file.Close()

	body := &proxmox.VirtualEnvironmentDatastoreUploadRequestBody{
		ContentType: *contentType,
		DatastoreID: datastoreID,
		FileName:    *fileName,
		FileReader:  file,
		NodeName:    nodeName,
	}

	_, err = veClient.UploadFileToDatastore(body)

	if err != nil {
		return err
	}

	volumeID, err := resourceVirtualEnvironmentFileGetVolumeID(d, m)

	if err != nil {
		return err
	}

	d.SetId(*volumeID)

	return resourceVirtualEnvironmentFileRead(d, m)
}

func resourceVirtualEnvironmentFileGetContentType(d *schema.ResourceData, m interface{}) (*string, error) {
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
		return nil, fmt.Errorf(
			"Missing argument \"%s.%s\" or \"%s\"",
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
			return nil, fmt.Errorf(
				"Cannot determine the content type of source \"%s\" - Please manually define the \"%s\" argument",
				sourceFilePath,
				mkResourceVirtualEnvironmentFileContentType,
			)
		}
	}

	ctValidator := getContentTypeValidator()
	_, errs := ctValidator(contentType, mkResourceVirtualEnvironmentFileContentType)

	if len(errs) > 0 {
		return nil, errs[0]
	}

	return &contentType, nil
}

func resourceVirtualEnvironmentFileGetFileName(d *schema.ResourceData, m interface{}) (*string, error) {
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
			"Missing argument \"%s.%s\"",
			mkResourceVirtualEnvironmentFileSourceFile,
			mkResourceVirtualEnvironmentFileSourceFilePath,
		)
	}

	if sourceFileFileName == "" {
		if resourceVirtualEnvironmentFileIsURL(d, m) {
			downloadURL, err := url.ParseRequestURI(sourceFilePath)

			if err != nil {
				return nil, err
			}

			path := strings.Split(downloadURL.Path, "/")
			sourceFileFileName = path[len(path)-1]

			if sourceFileFileName == "" {
				return nil, fmt.Errorf("Failed to determine file name from the URL \"%s\"", sourceFilePath)
			}
		} else {
			sourceFileFileName = filepath.Base(sourceFilePath)
		}
	}

	return &sourceFileFileName, nil
}

func resourceVirtualEnvironmentFileGetVolumeID(d *schema.ResourceData, m interface{}) (*string, error) {
	fileName, err := resourceVirtualEnvironmentFileGetFileName(d, m)

	if err != nil {
		return nil, err
	}

	datastoreID := d.Get(mkResourceVirtualEnvironmentFileDatastoreID).(string)
	contentType, err := resourceVirtualEnvironmentFileGetContentType(d, m)

	if err != nil {
		return nil, err
	}

	volumeID := fmt.Sprintf("%s:%s/%s", datastoreID, *contentType, *fileName)

	return &volumeID, nil
}

func resourceVirtualEnvironmentFileIsURL(d *schema.ResourceData, m interface{}) bool {
	sourceFile := d.Get(mkResourceVirtualEnvironmentFileSourceFile).([]interface{})
	sourceFilePath := ""

	if len(sourceFile) > 0 {
		sourceFileBlock := sourceFile[0].(map[string]interface{})
		sourceFilePath = sourceFileBlock[mkResourceVirtualEnvironmentFileSourceFilePath].(string)
	} else {
		return false
	}

	return strings.HasPrefix(sourceFilePath, "http://") || strings.HasPrefix(sourceFilePath, "https://")
}

func resourceVirtualEnvironmentFileRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
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

	list, err := veClient.ListDatastoreFiles(nodeName, datastoreID)

	if err != nil {
		return err
	}

	fileIsURL := resourceVirtualEnvironmentFileIsURL(d, m)
	fileName, err := resourceVirtualEnvironmentFileGetFileName(d, m)

	if err != nil {
		return err
	}

	for _, v := range list {
		if v.VolumeID == d.Id() {
			var fileModificationDate string
			var fileSize int64
			var fileTag string

			if fileIsURL {
				res, err := http.Head(sourceFilePath)

				if err != nil {
					return err
				}

				defer res.Body.Close()

				fileSize = res.ContentLength
				httpLastModified := res.Header.Get("Last-Modified")

				if httpLastModified != "" {
					timeParsed, err := time.Parse(time.RFC1123, httpLastModified)

					if err != nil {
						timeParsed, err = time.Parse(time.RFC1123Z, httpLastModified)

						if err != nil {
							return err
						}
					}

					fileModificationDate = timeParsed.UTC().Format(time.RFC3339)
				} else {
					d.Set(mkResourceVirtualEnvironmentFileFileModificationDate, "")
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
			} else {
				f, err := os.Open(sourceFilePath)

				if err != nil {
					return err
				}

				defer f.Close()

				fileInfo, err := f.Stat()

				if err != nil {
					return err
				}

				fileModificationDate = fileInfo.ModTime().UTC().Format(time.RFC3339)
				fileSize = fileInfo.Size()
				fileTag = fmt.Sprintf("%x-%x", fileInfo.ModTime().UTC().Unix(), fileInfo.Size())
			}

			lastFileModificationDate := d.Get(mkResourceVirtualEnvironmentFileFileModificationDate).(string)
			lastFileSize := int64(d.Get(mkResourceVirtualEnvironmentFileFileSize).(int))
			lastFileTag := d.Get(mkResourceVirtualEnvironmentFileFileTag).(string)

			d.Set(mkResourceVirtualEnvironmentFileFileModificationDate, fileModificationDate)
			d.Set(mkResourceVirtualEnvironmentFileFileName, *fileName)
			d.Set(mkResourceVirtualEnvironmentFileFileSize, fileSize)
			d.Set(mkResourceVirtualEnvironmentFileFileTag, fileTag)
			d.Set(mkResourceVirtualEnvironmentFileSourceFileChanged, lastFileModificationDate != fileModificationDate || lastFileSize != fileSize || lastFileTag != fileTag)

			return nil
		}
	}

	d.SetId("")

	return nil
}

func resourceVirtualEnvironmentFileDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	datastoreID := d.Get(mkResourceVirtualEnvironmentFileDatastoreID).(string)
	nodeName := d.Get(mkResourceVirtualEnvironmentFileNodeName).(string)

	err = veClient.DeleteDatastoreFile(nodeName, datastoreID, d.Id())

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}

		return err
	}

	d.SetId("")

	return nil
}
