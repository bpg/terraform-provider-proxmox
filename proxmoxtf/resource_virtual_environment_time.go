/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"fmt"
	"time"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	mkResourceVirtualEnvironmentTimeLocalTime = "local_time"
	mkResourceVirtualEnvironmentTimeNodeName  = "node_name"
	mkResourceVirtualEnvironmentTimeTimeZone  = "time_zone"
	mkResourceVirtualEnvironmentTimeUTCTime   = "utc_time"
)

func resourceVirtualEnvironmentTime() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentTimeLocalTime: {
				Type:        schema.TypeString,
				Description: "The local timestamp",
				Computed:    true,
			},
			mkResourceVirtualEnvironmentTimeNodeName: {
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
			},
			mkResourceVirtualEnvironmentTimeTimeZone: {
				Type:        schema.TypeString,
				Description: "The time zone",
				Required:    true,
			},
			mkResourceVirtualEnvironmentTimeUTCTime: {
				Type:        schema.TypeString,
				Description: "The UTC timestamp",
				Computed:    true,
			},
		},
		Create: resourceVirtualEnvironmentTimeCreate,
		Read:   resourceVirtualEnvironmentTimeRead,
		Update: resourceVirtualEnvironmentTimeUpdate,
		Delete: resourceVirtualEnvironmentTimeDelete,
	}
}

func resourceVirtualEnvironmentTimeCreate(d *schema.ResourceData, m interface{}) error {
	err := resourceVirtualEnvironmentTimeUpdate(d, m)

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentTimeNodeName).(string)

	d.SetId(fmt.Sprintf("%s_time", nodeName))

	return nil
}

func resourceVirtualEnvironmentTimeRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentTimeNodeName).(string)
	nodeTime, err := veClient.GetNodeTime(nodeName)

	if err != nil {
		return err
	}

	localLocation, err := time.LoadLocation(nodeTime.TimeZone)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s_time", nodeName))

	localTimeOffset := time.Time(nodeTime.LocalTime).Sub(time.Now().UTC())
	localTime := time.Time(nodeTime.LocalTime).Add(-localTimeOffset).In(localLocation)

	d.Set(mkDataSourceVirtualEnvironmentTimeLocalTime, localTime.Format(time.RFC3339))
	d.Set(mkDataSourceVirtualEnvironmentTimeTimeZone, nodeTime.TimeZone)
	d.Set(mkDataSourceVirtualEnvironmentTimeUTCTime, time.Time(nodeTime.UTCTime).Format(time.RFC3339))

	return nil
}

func resourceVirtualEnvironmentTimeUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentTimeNodeName).(string)
	timeZone := d.Get(mkResourceVirtualEnvironmentTimeTimeZone).(string)

	err = veClient.UpdateNodeTime(nodeName, &proxmox.VirtualEnvironmentNodeUpdateTimeRequestBody{
		TimeZone: timeZone,
	})

	if err != nil {
		return err
	}

	return resourceVirtualEnvironmentTimeRead(d, m)
}

func resourceVirtualEnvironmentTimeDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")

	return nil
}
