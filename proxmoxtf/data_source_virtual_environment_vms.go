/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
)

const (
	mkDataSourceVirtualEnvironmentVMs = "vms"
)

func dataSourceVirtualEnvironmentVMs() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentVMNodeName: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The node name",
			},
			mkDataSourceVirtualEnvironmentVMTags: {
				Type:        schema.TypeList,
				Description: "Tags of the VM to match",
				Computed:    true,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceVirtualEnvironmentVMs: {
				Type:        schema.TypeList,
				Description: "VMs",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceVirtualEnvironmentVM().Schema,
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

	nodeNames, err := getNodeNames(ctx, d, veClient)
	if err != nil {
		return diag.FromErr(err)
	}

	tagsData := d.Get(mkResourceVirtualEnvironmentVMTags).([]interface{})
	var filterTags []string
	for i := 0; i < len(tagsData); i++ {
		tag := strings.TrimSpace(tagsData[i].(string))
		if len(tag) > 0 {
			filterTags = append(filterTags, tag)
		}
	}
	sort.Strings(filterTags)

	var vms []interface{}
	for _, nodeName := range nodeNames {
		listData, err := veClient.ListVMs(ctx, nodeName)
		if err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}

		sort.Slice(listData, func(i, j int) bool {
			return listData[i].VMID < listData[j].VMID
		})

		for _, data := range listData {
			vm := map[string]interface{}{
				mkDataSourceVirtualEnvironmentVMNodeName: nodeName,
				mkDataSourceVirtualEnvironmentVMVMID:     data.VMID,
			}

			if data.Name != nil {
				vm[mkDataSourceVirtualEnvironmentVMName] = *data.Name
			} else {
				vm[mkDataSourceVirtualEnvironmentVMName] = ""
			}

			var tags []string
			if data.Tags != nil && *data.Tags != "" {
				tags = strings.Split(*data.Tags, ";")
				sort.Strings(tags)
				vm[mkDataSourceVirtualEnvironmentVMTags] = tags
			}

			if len(filterTags) > 0 {
				match := true
				for _, tag := range filterTags {
					if !slices.Contains(tags, tag) {
						match = false
						break
					}
				}
				if !match {
					continue
				}
			}

			vms = append(vms, vm)
		}
	}

	err = d.Set(mkDataSourceVirtualEnvironmentVMs, vms)
	diags = append(diags, diag.FromErr(err)...)

	d.SetId(uuid.New().String())

	return diags
}

func getNodeNames(
	ctx context.Context,
	d *schema.ResourceData,
	veClient *proxmox.VirtualEnvironmentClient,
) ([]string, error) {
	var nodeNames []string
	nodeName := d.Get(mkDataSourceVirtualEnvironmentVMNodeName).(string)
	if nodeName != "" {
		nodeNames = append(nodeNames, nodeName)
	} else {
		nodes, err := veClient.ListNodes(ctx)
		if err != nil {
			return nil, fmt.Errorf("error listing nodes: %w", err)
		}

		for _, node := range nodes {
			nodeNames = append(nodeNames, node.Name)
		}
	}

	sort.Strings(nodeNames)
	return nodeNames, nil
}
