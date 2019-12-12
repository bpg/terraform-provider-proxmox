/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	mkResourceVirtualEnvironmentFileDatastoreID      = "datastore_id"
	mkResourceVirtualEnvironmentFileFileName         = "file_name"
	mkResourceVirtualEnvironmentFileOverrideFileName = "override_file_name"
	mkResourceVirtualEnvironmentFileNodeName         = "node_name"
	mkResourceVirtualEnvironmentFileSource           = "source"
	mkResourceVirtualEnvironmentFileTemplate         = "template"
)

func resourceVirtualEnvironmentFile() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentFileDatastoreID: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The datastore id",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentFileFileName: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The datastore file name",
				Computed:    true,
			},
			mkResourceVirtualEnvironmentFileOverrideFileName: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The file name to use in the datastore (leave undefined to use source file name)",
				Optional:    true,
				ForceNew:    true,
				Default:     "",
			},
			mkResourceVirtualEnvironmentFileNodeName: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentFileSource: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The path to a file",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentFileTemplate: &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Whether this is a container template",
				Required:    true,
				ForceNew:    true,
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

	datastoreID := d.Get(mkResourceVirtualEnvironmentFileDatastoreID).(string)
	fileName, err := resourceVirtualEnvironmentFileGetFileName(d, m)

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentFileNodeName).(string)
	source := d.Get(mkResourceVirtualEnvironmentFileSource).(string)
	template := d.Get(mkResourceVirtualEnvironmentFileTemplate).(bool)

	var sourceReader io.Reader

	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		log.Printf("[DEBUG] Downloading file '%s'", source)

		res, err := http.Get(source)

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

		file, err := os.Open(tempDownloadedFileName)

		if err != nil {
			return err
		}

		defer file.Close()

		sourceReader = file
	} else {
		file, err := os.Open(source)

		if err != nil {
			return err
		}

		defer file.Close()

		sourceReader = file
	}

	contentType := "iso"

	if template {
		contentType = "vztmpl"
	}

	body := &proxmox.VirtualEnvironmentDatastoreUploadRequestBody{
		ContentType: contentType,
		DatastoreID: datastoreID,
		FileName:    *fileName,
		FileReader:  sourceReader,
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

func resourceVirtualEnvironmentFileGetFileName(d *schema.ResourceData, m interface{}) (*string, error) {
	fileName := d.Get(mkResourceVirtualEnvironmentFileOverrideFileName).(string)
	source := d.Get(mkResourceVirtualEnvironmentFileSource).(string)

	if fileName == "" {
		if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
			downloadURL, err := url.ParseRequestURI(source)

			if err != nil {
				return nil, err
			}

			path := strings.Split(downloadURL.Path, "/")
			fileName = path[len(path)-1]

			if fileName == "" {
				return nil, errors.New("Failed to determine file name from source URL")
			}
		} else {
			fileName = filepath.Base(source)
		}
	}

	return &fileName, nil
}

func resourceVirtualEnvironmentFileGetVolumeID(d *schema.ResourceData, m interface{}) (*string, error) {
	contentType := "iso"
	fileName, err := resourceVirtualEnvironmentFileGetFileName(d, m)

	if err != nil {
		return nil, err
	}

	datastoreID := d.Get(mkResourceVirtualEnvironmentFileDatastoreID).(string)
	template := d.Get(mkResourceVirtualEnvironmentFileTemplate).(bool)

	if template {
		contentType = "vztmpl"
	}

	volumeID := fmt.Sprintf("%s:%s/%s", datastoreID, contentType, *fileName)

	return &volumeID, nil
}

func resourceVirtualEnvironmentFileRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	datastoreID := d.Get(mkResourceVirtualEnvironmentFileDatastoreID).(string)
	nodeName := d.Get(mkResourceVirtualEnvironmentFileNodeName).(string)

	list, err := veClient.ListDatastoreFiles(nodeName, datastoreID)

	if err != nil {
		return err
	}

	for _, v := range list {
		if v.VolumeID == d.Id() {
			fileName, _ := resourceVirtualEnvironmentFileGetFileName(d, m)

			d.Set(mkResourceVirtualEnvironmentFileFileName, *fileName)

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
