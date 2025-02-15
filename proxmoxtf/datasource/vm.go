/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

//nolint:dupl
package datasource

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	mkDataSourceVirtualEnvironmentVMName     = "name"
	mkDataSourceVirtualEnvironmentVMNodeName = "node_name"
	mkDataSourceVirtualEnvironmentVMTags     = "tags"
	mkDataSourceVirtualEnvironmentVMTemplate = "template"
	mkDataSourceVirtualEnvironmentVMStatus   = "status"
	mkDataSourceVirtualEnvironmentVMVMID     = "vm_id"
)

// VM returns a resource for a single Proxmox VM.
func VM() *schema.Resource {
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
			mkDataSourceVirtualEnvironmentVMTemplate: {
				Type:        schema.TypeBool,
				Description: "Is VM a template (true) or a regular VM (false)",
				Optional:    true,
			},
			mkDataSourceVirtualEnvironmentVMStatus: {
				Type:        schema.TypeString,
				Description: "Status of the VM",
				Optional:    true,
			},
			mkDataSourceVirtualEnvironmentVMVMID: {
				Type:        schema.TypeInt,
				Description: "The VM identifier",
				Required:    true,
			},
		},
		ReadContext: vmRead,
	}
}

// vmRead reads the data of a VM by ID.
func vmRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)

	client, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkDataSourceVirtualEnvironmentVMNodeName).(string)
	vmID := d.Get(mkDataSourceVirtualEnvironmentVMVMID).(int)

	vmStatus, err := client.Node(nodeName).VM(vmID).GetVMStatus(ctx)
	if err != nil {
		if errors.Is(err, api.ErrNoDataObjectInResponse) {
			d.SetId("")

			return nil
		}

		return diag.FromErr(err)
	}

	vmConfig, err := client.Node(nodeName).VM(vmID).GetVM(ctx)
	if err != nil {
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

	err = d.Set(mkDataSourceVirtualEnvironmentVMStatus, vmStatus.Status)
	diags = append(diags, diag.FromErr(err)...)

	if vmConfig.Template == nil {
		err = d.Set(mkDataSourceVirtualEnvironmentVMTemplate, false)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentVMTemplate, *vmConfig.Template)
	}

	diags = append(diags, diag.FromErr(err)...)

	err = d.Set(mkDataSourceVirtualEnvironmentVMTags, tags)
	diags = append(diags, diag.FromErr(err)...)

	d.SetId(strconv.Itoa(vmID))

	return diags
}
