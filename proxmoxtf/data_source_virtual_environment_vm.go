/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	mkDataSourceVirtualEnvironmentVMName     = "name"
	mkDataSourceVirtualEnvironmentVMNodeName = "node_name"
	mkDataSourceVirtualEnvironmentVMTags     = "tags"
	mkDataSourceVirtualEnvironmentVMVMID     = "vm_id"
)

func dataSourceVirtualEnvironmentVM() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentVMName: {
				Type:        schema.TypeString,
				Description: "The VM name",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentVMNodeName: {
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
			},
			mkDataSourceVirtualEnvironmentVMTags: {
				Type:        schema.TypeList,
				Description: "Tags of the VM",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceVirtualEnvironmentVMVMID: {
				Type:        schema.TypeInt,
				Description: "The VM identifier",
				Required:    true,
			},
		},
		ReadContext: dataSourceVirtualEnvironmentVMRead,
	}
}

// dataSourceVirtualEnvironmentVMRead reads the data of a VM by ID.
func dataSourceVirtualEnvironmentVMRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkDataSourceVirtualEnvironmentVMNodeName).(string)
	vmID := d.Get(mkDataSourceVirtualEnvironmentVMVMID).(int)

	vmStatus, err := veClient.GetVMStatus(ctx, nodeName, vmID)
	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") ||
			(strings.Contains(err.Error(), "HTTP 500") && strings.Contains(err.Error(), "does not exist")) {
			d.SetId("")

			return nil
		}

		return diag.FromErr(err)
	}

	if vmStatus.Name != nil {
		err = d.Set(mkDataSourceVirtualEnvironmentVMName, *vmStatus.Name)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentVMName, "")
	}
	diags = append(diags, diag.FromErr(err)...)

	var tags []string
	if vmStatus.Tags != nil {
		for _, tag := range strings.Split(*vmStatus.Tags, ";") {
			t := strings.TrimSpace(tag)
			if len(t) > 0 {
				tags = append(tags, t)
			}
		}
		sort.Strings(tags)
	}
	err = d.Set(mkDataSourceVirtualEnvironmentVMTags, tags)
	diags = append(diags, diag.FromErr(err)...)

	d.SetId(strconv.Itoa(vmID))

	return diags
}
