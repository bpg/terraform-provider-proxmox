/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	mkResourceVirtualEnvironmentTimeLocalTime = "local_time"
	mkResourceVirtualEnvironmentTimeNodeName  = "node_name"
	mkResourceVirtualEnvironmentTimeTimeZone  = "time_zone"
	mkResourceVirtualEnvironmentTimeUTCTime   = "utc_time"
)

// Time returns a resource that manages time settings for a node.
func Time() *schema.Resource {
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
		CreateContext: timeCreate,
		ReadContext:   timeRead,
		UpdateContext: timeUpdate,
		DeleteContext: timeDelete,
	}
}

func timeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diags := timeUpdate(ctx, d, m)
	if diags.HasError() {
		return diags
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentTimeNodeName).(string)

	d.SetId(fmt.Sprintf("%s_time", nodeName))

	return nil
}

func timeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentTimeNodeName).(string)
	nodeTime, err := api.Node(nodeName).GetTime(ctx)
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

	err = d.Set(mkResourceVirtualEnvironmentTimeLocalTime, localTime.Format(time.RFC3339))
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkResourceVirtualEnvironmentTimeTimeZone, nodeTime.TimeZone)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(
		mkResourceVirtualEnvironmentTimeUTCTime,
		time.Time(nodeTime.UTCTime).Format(time.RFC3339),
	)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}

func timeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentTimeNodeName).(string)
	timeZone := d.Get(mkResourceVirtualEnvironmentTimeTimeZone).(string)

	err = api.Node(nodeName).UpdateTime(
		ctx,
		&nodes.UpdateTimeRequestBody{
			TimeZone: timeZone,
		},
	)
	if err != nil {
		return diag.FromErr(err)
	}

	return timeRead(ctx, d, m)
}

func timeDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}
