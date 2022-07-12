/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	mkDataSourceVirtualEnvironmentTimeLocalTime = "local_time"
	mkDataSourceVirtualEnvironmentTimeNodeName  = "node_name"
	mkDataSourceVirtualEnvironmentTimeTimeZone  = "time_zone"
	mkDataSourceVirtualEnvironmentTimeUTCTime   = "utc_time"
)

func dataSourceVirtualEnvironmentTime() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentTimeLocalTime: {
				Type:        schema.TypeString,
				Description: "The local timestamp",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentTimeNodeName: {
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
			},
			mkDataSourceVirtualEnvironmentTimeTimeZone: {
				Type:        schema.TypeString,
				Description: "The time zone",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentTimeUTCTime: {
				Type:        schema.TypeString,
				Description: "The UTC timestamp",
				Computed:    true,
			},
		},
		ReadContext: dataSourceVirtualEnvironmentTimeRead,
	}
}

func dataSourceVirtualEnvironmentTimeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkDataSourceVirtualEnvironmentTimeNodeName).(string)
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

	return diags
}
