/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	mkResourceVirtualEnvironmentFileDatastoreID = "datastore_id"
	mkResourceVirtualEnvironmentFileFileName    = "file_name"
	mkResourceVirtualEnvironmentFileNodeName    = "node_name"
	mkResourceVirtualEnvironmentFileSource      = "source"
	mkResourceVirtualEnvironmentFileTemplate    = "template"
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
				Description: "The file name to use in the datastore",
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
	fileName := d.Get(mkResourceVirtualEnvironmentFileFileName).(string)
	nodeName := d.Get(mkResourceVirtualEnvironmentFileNodeName).(string)
	source := d.Get(mkResourceVirtualEnvironmentFileSource).(string)
	template := d.Get(mkResourceVirtualEnvironmentFileTemplate).(bool)

	if fileName == "" {
		fileName = filepath.Base(source)
	}

	file, err := os.Open(source)

	if err != nil {
		return err
	}

	defer file.Close()

	contentType := "iso"

	if template {
		contentType = "vztmpl"
	}

	body := &proxmox.VirtualEnvironmentDatastoreUploadRequestBody{
		ContentType: contentType,
		DatastoreID: datastoreID,
		FileName:    fileName,
		FileReader:  file,
		NodeName:    nodeName,
	}

	_, err = veClient.UploadFileToDatastore(body)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s:%s:%s", nodeName, datastoreID, fileName))

	return resourceVirtualEnvironmentFileRead(d, m)
}

func resourceVirtualEnvironmentFileGetVolumeID(d *schema.ResourceData, m interface{}) string {
	datastoreID := d.Get(mkResourceVirtualEnvironmentFileDatastoreID).(string)
	fileName := d.Get(mkResourceVirtualEnvironmentFileFileName).(string)
	source := d.Get(mkResourceVirtualEnvironmentFileSource).(string)
	template := d.Get(mkResourceVirtualEnvironmentFileTemplate).(bool)

	if fileName == "" {
		fileName = filepath.Base(source)
	}

	contentType := "iso"

	if template {
		contentType = "vztmpl"
	}

	return fmt.Sprintf("%s:%s/%s", datastoreID, contentType, fileName)
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

	volumeID := resourceVirtualEnvironmentFileGetVolumeID(d, m)

	for _, v := range list {
		if v.VolumeID == volumeID {
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
	volumeID := resourceVirtualEnvironmentFileGetVolumeID(d, m)

	err = veClient.DeleteDatastoreFile(nodeName, datastoreID, volumeID)

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
