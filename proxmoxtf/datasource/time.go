/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datasource

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	mkDataSourceVirtualEnvironmentTimeLocalTime = "local_time"
	mkDataSourceVirtualEnvironmentTimeNodeName  = "node_name"
	mkDataSourceVirtualEnvironmentTimeTimeZone  = "time_zone"
	mkDataSourceVirtualEnvironmentTimeUTCTime   = "utc_time"
)

// Time returns a resource for the Proxmox time.
func Time() *schema.Resource {
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
		ReadContext: timeRead,
	}
}

func timeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkDataSourceVirtualEnvironmentTimeNodeName).(string)
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

	err = d.Set(mkDataSourceVirtualEnvironmentTimeLocalTime, localTime.Format(time.RFC3339))
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentTimeTimeZone, nodeTime.TimeZone)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(
		mkDataSourceVirtualEnvironmentTimeUTCTime,
		time.Time(nodeTime.UTCTime).Format(time.RFC3339),
	)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}
