/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
		CreateContext: resourceVirtualEnvironmentTimeCreate,
		ReadContext:   resourceVirtualEnvironmentTimeRead,
		UpdateContext: resourceVirtualEnvironmentTimeUpdate,
		DeleteContext: resourceVirtualEnvironmentTimeDelete,
	}
}

func resourceVirtualEnvironmentTimeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diags := resourceVirtualEnvironmentTimeUpdate(ctx, d, m)
	if diags.HasError() {
		return diags
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentTimeNodeName).(string)

	d.SetId(fmt.Sprintf("%s_time", nodeName))

	return nil
}

func resourceVirtualEnvironmentTimeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentTimeNodeName).(string)
	nodeTime, err := veClient.GetNodeTime(ctx, nodeName)
	if err != nil {
		return diag.FromErr(err)
	}

	localLocation, err := time.LoadLocation(nodeTime.TimeZone)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s_time", nodeName))

	localTimeOffset := time.Time(nodeTime.LocalTime).Sub(time.Now().UTC())
	localTime := time.Time(nodeTime.LocalTime).Add(-localTimeOffset).In(localLocation)

	err = d.Set(mkDataSourceVirtualEnvironmentTimeLocalTime, localTime.Format(time.RFC3339))
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentTimeTimeZone, nodeTime.TimeZone)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentTimeUTCTime, time.Time(nodeTime.UTCTime).Format(time.RFC3339))
	diags = append(diags, diag.FromErr(err)...)

	return nil
}

func resourceVirtualEnvironmentTimeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentTimeNodeName).(string)
	timeZone := d.Get(mkResourceVirtualEnvironmentTimeTimeZone).(string)

	err = veClient.UpdateNodeTime(ctx, nodeName, &proxmox.VirtualEnvironmentNodeUpdateTimeRequestBody{
		TimeZone: timeZone,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceVirtualEnvironmentTimeRead(ctx, d, m)
}

func resourceVirtualEnvironmentTimeDelete(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}
