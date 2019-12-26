/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const (
	dvResourceVirtualEnvironmentVMCPUCores                = 1
	dvResourceVirtualEnvironmentVMCPUSockets              = 1
	dvResourceVirtualEnvironmentVMDiskDatastoreID         = "local-lvm"
	dvResourceVirtualEnvironmentVMDiskFileID              = ""
	dvResourceVirtualEnvironmentVMDiskSize                = 8
	dvResourceVirtualEnvironmentVMMemoryDedicated         = 256
	dvResourceVirtualEnvironmentVMMemoryFloating          = 0
	dvResourceVirtualEnvironmentVMMemoryShared            = 0
	dvResourceVirtualEnvironmentVMName                    = ""
	dvResourceVirtualEnvironmentVMNetworkDeviceBridge     = "vmbr0"
	dvResourceVirtualEnvironmentVMNetworkDeviceMACAddress = ""
	dvResourceVirtualEnvironmentVMNetworkDeviceModel      = "virtio"
	dvResourceVirtualEnvironmentVMNetworkDeviceVLANID     = -1
	dvResourceVirtualEnvironmentVMVMID                    = -1

	mkResourceVirtualEnvironmentVMCPU                     = "cpu"
	mkResourceVirtualEnvironmentVMCPUCores                = "cores"
	mkResourceVirtualEnvironmentVMCPUSockets              = "sockets"
	mkResourceVirtualEnvironmentVMDisk                    = "disk"
	mkResourceVirtualEnvironmentVMDiskDatastoreID         = "datastore_id"
	mkResourceVirtualEnvironmentVMDiskFileID              = "file_id"
	mkResourceVirtualEnvironmentVMDiskSize                = "size"
	mkResourceVirtualEnvironmentVMMemory                  = "memory"
	mkResourceVirtualEnvironmentVMMemoryDedicated         = "dedicated"
	mkResourceVirtualEnvironmentVMMemoryFloating          = "floating"
	mkResourceVirtualEnvironmentVMMemoryShared            = "shared"
	mkResourceVirtualEnvironmentVMName                    = "name"
	mkResourceVirtualEnvironmentVMNetworkDevice           = "network_device"
	mkResourceVirtualEnvironmentVMNetworkDeviceBridge     = "bridge"
	mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress = "mac_address"
	mkResourceVirtualEnvironmentVMNetworkDeviceModel      = "model"
	mkResourceVirtualEnvironmentVMNetworkDeviceVLANID     = "vlan_id"
	mkResourceVirtualEnvironmentVMVMID                    = "vm_id"
)

func resourceVirtualEnvironmentVM() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentVMCPU: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The CPU allocation",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					defaultList := make([]interface{}, 1)
					defaultMap := make(map[string]interface{})

					defaultMap[mkResourceVirtualEnvironmentVMCPUCores] = dvResourceVirtualEnvironmentVMCPUCores
					defaultMap[mkResourceVirtualEnvironmentVMCPUSockets] = dvResourceVirtualEnvironmentVMCPUSockets

					defaultList[0] = defaultMap

					return defaultList, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMCPUCores: {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "The number of CPU cores",
							Default:      dvResourceVirtualEnvironmentVMCPUCores,
							ValidateFunc: validation.IntBetween(1, 2304),
						},
						mkResourceVirtualEnvironmentVMCPUSockets: {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "The number of CPU sockets",
							Default:      dvResourceVirtualEnvironmentVMCPUSockets,
							ValidateFunc: validation.IntBetween(1, 16),
						},
					},
				},
				MaxItems: 1,
				MinItems: 1,
			},
			mkResourceVirtualEnvironmentVMDisk: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The disk devices",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					defaultList := make([]interface{}, 1)
					defaultMap := make(map[string]interface{})

					defaultMap[mkResourceVirtualEnvironmentVMDiskDatastoreID] = dvResourceVirtualEnvironmentVMDiskDatastoreID
					defaultMap[mkResourceVirtualEnvironmentVMDiskFileID] = dvResourceVirtualEnvironmentVMDiskFileID
					defaultMap[mkResourceVirtualEnvironmentVMDiskSize] = dvResourceVirtualEnvironmentVMDiskSize

					defaultList[0] = defaultMap

					return defaultList, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMDiskDatastoreID: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The datastore id",
							Default:     dvResourceVirtualEnvironmentVMDiskDatastoreID,
						},
						mkResourceVirtualEnvironmentVMDiskFileID: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The file id for a disk image",
							Default:     dvResourceVirtualEnvironmentVMDiskFileID,
						},
						mkResourceVirtualEnvironmentVMDiskSize: {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The disk size in gigabytes",
							Default:     dvResourceVirtualEnvironmentVMDiskSize,
						},
					},
				},
				MaxItems: 14,
				MinItems: 1,
			},
			mkResourceVirtualEnvironmentVMMemory: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The memory allocation",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					defaultList := make([]interface{}, 1)
					defaultMap := make(map[string]interface{})

					defaultMap[mkResourceVirtualEnvironmentVMMemoryDedicated] = dvResourceVirtualEnvironmentVMMemoryDedicated
					defaultMap[mkResourceVirtualEnvironmentVMMemoryFloating] = dvResourceVirtualEnvironmentVMMemoryFloating
					defaultMap[mkResourceVirtualEnvironmentVMMemoryShared] = dvResourceVirtualEnvironmentVMMemoryShared

					defaultList[0] = defaultMap

					return defaultList, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMMemoryDedicated: {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "The dedicated memory in megabytes",
							Default:      dvResourceVirtualEnvironmentVMMemoryDedicated,
							ValidateFunc: validation.IntBetween(64, 268435456),
						},
						mkResourceVirtualEnvironmentVMMemoryFloating: {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "The floating memory in megabytes (balloon)",
							Default:      dvResourceVirtualEnvironmentVMMemoryFloating,
							ValidateFunc: validation.IntBetween(0, 268435456),
						},
						mkResourceVirtualEnvironmentVMMemoryShared: {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "The shared memory in megabytes",
							Default:      dvResourceVirtualEnvironmentVMMemoryShared,
							ValidateFunc: validation.IntBetween(0, 268435456),
						},
					},
				},
				MaxItems: 1,
				MinItems: 1,
			},
			mkResourceVirtualEnvironmentVMName: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name",
				Default:     dvResourceVirtualEnvironmentVMName,
			},
			mkResourceVirtualEnvironmentVMNetworkDevice: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The network devices",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					defaultList := make([]interface{}, 1)
					defaultMap := make(map[string]interface{})

					defaultMap[mkResourceVirtualEnvironmentVMNetworkDeviceBridge] = dvResourceVirtualEnvironmentVMNetworkDeviceBridge
					defaultMap[mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress] = dvResourceVirtualEnvironmentVMNetworkDeviceMACAddress
					defaultMap[mkResourceVirtualEnvironmentVMNetworkDeviceModel] = dvResourceVirtualEnvironmentVMNetworkDeviceModel
					defaultMap[mkResourceVirtualEnvironmentVMNetworkDeviceVLANID] = dvResourceVirtualEnvironmentVMNetworkDeviceVLANID

					defaultList[0] = defaultMap

					return defaultList, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentVMNetworkDeviceBridge: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The bridge",
							Default:     dvResourceVirtualEnvironmentVMNetworkDeviceBridge,
						},
						mkResourceVirtualEnvironmentVMNetworkDeviceMACAddress: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The MAC address",
							Default:     dvResourceVirtualEnvironmentVMNetworkDeviceMACAddress,
						},
						mkResourceVirtualEnvironmentVMNetworkDeviceModel: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The model",
							Default:     dvResourceVirtualEnvironmentVMNetworkDeviceModel,
						},
						mkResourceVirtualEnvironmentVMNetworkDeviceVLANID: {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The VLAN identifier",
							Default:     dvResourceVirtualEnvironmentVMNetworkDeviceVLANID,
							ValidateFunc: func(i interface{}, k string) (s []string, es []error) {
								min := 1

								v, ok := i.(int)

								if !ok {
									es = append(es, fmt.Errorf("expected type of %s to be int", k))
									return
								}

								if v != -1 {
									if v < min {
										es = append(es, fmt.Errorf("expected %s to be at least %d, got %d", k, min, v))
										return
									}
								}

								return
							},
						},
					},
				},
				MaxItems: 8,
				MinItems: 0,
			},
			mkResourceVirtualEnvironmentVMVMID: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The identifier",
				Default:     dvResourceVirtualEnvironmentVMVMID,
				ValidateFunc: func(i interface{}, k string) (s []string, es []error) {
					min := 100
					max := 2000000

					v, ok := i.(int)

					if !ok {
						es = append(es, fmt.Errorf("expected type of %s to be int", k))
						return
					}

					if v != -1 {
						if v < min || v > max {
							es = append(es, fmt.Errorf("expected %s to be in the range (%d - %d), got %d", k, min, max, v))
							return
						}
					}

					return
				},
			},
		},
		Create: resourceVirtualEnvironmentVMCreate,
		Read:   resourceVirtualEnvironmentVMRead,
		Update: resourceVirtualEnvironmentVMUpdate,
		Delete: resourceVirtualEnvironmentVMDelete,
	}
}

func resourceVirtualEnvironmentVMCreate(d *schema.ResourceData, m interface{}) error {
	/*
		config := m.(providerConfiguration)
		veClient, err := config.GetVEClient()

		if err != nil {
			return err
		}

		d.SetId("")
	*/

	return resourceVirtualEnvironmentVMRead(d, m)
}

func resourceVirtualEnvironmentVMRead(d *schema.ResourceData, m interface{}) error {
	/*
		config := m.(providerConfiguration)
		veClient, err := config.GetVEClient()

		if err != nil {
			return err
		}
	*/

	return nil
}

func resourceVirtualEnvironmentVMUpdate(d *schema.ResourceData, m interface{}) error {
	/*
		config := m.(providerConfiguration)
		veClient, err := config.GetVEClient()

		if err != nil {
			return err
		}
	*/

	return resourceVirtualEnvironmentVMRead(d, m)
}

func resourceVirtualEnvironmentVMDelete(d *schema.ResourceData, m interface{}) error {
	/*
		config := m.(providerConfiguration)
		veClient, err := config.GetVEClient()

		if err != nil {
			return err
		}

		d.SetId("")
	*/

	return nil
}
