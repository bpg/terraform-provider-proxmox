/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"context"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const mkDataSourceVirtualEnvironmentVMs = "vms"

func dataSourceVirtualEnvironmentVMs() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentVMNodeName: {
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
			},
			mkDataSourceVirtualEnvironmentVMs: {
				Type:        schema.TypeList,
				Description: "VMs",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkDataSourceVirtualEnvironmentVMName: {
							Type:        schema.TypeString,
							Description: "The VM name",
							Computed:    true,
						},
						mkDataSourceVirtualEnvironmentVMTags: {
							Type:        schema.TypeList,
							Description: "Tags of the VM",
							Computed:    true,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						mkDataSourceVirtualEnvironmentVMVMID: {
							Type:        schema.TypeInt,
							Description: "The VM identifier",
							Required:    true,
						},
					},
				},
			},
		},
		ReadContext: dataSourceVirtualEnvironmentVMsRead,
	}
}

// dataSourceVirtualEnvironmentVMRead reads the data of a VM by ID.
func dataSourceVirtualEnvironmentVMsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkDataSourceVirtualEnvironmentVMNodeName).(string)

	listData, err := veClient.ListVMs(ctx, nodeName)
	if err != nil {
		return diag.FromErr(err)
	}

	vms := make([]interface{}, len(listData))
	for i, data := range listData {
		vm := map[string]interface{}{
			mkDataSourceVirtualEnvironmentVMVMID: data.VMID,
		}

		if data.Name != nil {
			vm[mkDataSourceVirtualEnvironmentVMName] = *data.Name
		}

		if data.Tags != nil && *data.Tags != "" {
			tags := strings.Split(*data.Tags, ";")
			sort.Strings(tags)
			vm[mkDataSourceVirtualEnvironmentVMTags] = tags
		}

		vms[i] = vm
	}

	err = d.Set(mkDataSourceVirtualEnvironmentVMs, vms)
	diags = append(diags, diag.FromErr(err)...)

	d.SetId(nodeName)

	return diags
}
