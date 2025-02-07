/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

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
	mkDataSourceVirtualEnvironmentContainerName     = "name"
	mkDataSourceVirtualEnvironmentContainerNodeName = "node_name"
	mkDataSourceVirtualEnvironmentContainerTags     = "tags"
	mkDataSourceVirtualEnvironmentContainerTemplate = "template"
	mkDataSourceVirtualEnvironmentContainerStatus   = "status"
	mkDataSourceVirtualEnvironmentContainerVMID     = "vm_id"
)

// Container returns a resource for a single Proxmox Container.
func Container() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentContainerName: {
				Type:        schema.TypeString,
				Description: "The Container name",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentContainerNodeName: {
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
			},
			mkDataSourceVirtualEnvironmentContainerTags: {
				Type:        schema.TypeList,
				Description: "Tags of the Container",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceVirtualEnvironmentContainerTemplate: {
				Type:        schema.TypeBool,
				Description: "Is Container a template (true) or a regular Container (false)",
				Optional:    true,
			},
			mkDataSourceVirtualEnvironmentContainerStatus: {
				Type:        schema.TypeString,
				Description: "Status of the Container",
				Optional:    true,
			},
			mkDataSourceVirtualEnvironmentContainerVMID: {
				Type:        schema.TypeInt,
				Description: "The Container identifier",
				Required:    true,
			},
		},
		ReadContext: containerRead,
	}
}

// containerRead reads the data of a Container by ID.
func containerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)

	client, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkDataSourceVirtualEnvironmentContainerNodeName).(string)
	containerID := d.Get(mkDataSourceVirtualEnvironmentContainerVMID).(int)

	containerStatus, err := client.Node(nodeName).Container(containerID).GetContainerStatus(ctx)
	if err != nil {
		if errors.Is(err, api.ErrNoDataObjectInResponse) {
			d.SetId("")

			return nil
		}

		return diag.FromErr(err)
	}

	containerConfig, err := client.Node(nodeName).Container(containerID).GetContainer(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if containerStatus.Name != nil {
		err = d.Set(mkDataSourceVirtualEnvironmentContainerName, *containerStatus.Name)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentContainerName, "")
	}

	diags = append(diags, diag.FromErr(err)...)

	var tags []string

	if containerStatus.Tags != nil {
		for _, tag := range strings.Split(*containerStatus.Tags, ";") {
			t := strings.TrimSpace(tag)
			if len(t) > 0 {
				tags = append(tags, t)
			}
		}

		sort.Strings(tags)
	}

	err = d.Set(mkDataSourceVirtualEnvironmentContainerStatus, containerStatus.Status)
	diags = append(diags, diag.FromErr(err)...)

	if containerConfig.Template == nil {
		err = d.Set(mkDataSourceVirtualEnvironmentContainerTemplate, false)
	} else {
		err = d.Set(mkDataSourceVirtualEnvironmentContainerTemplate, *containerConfig.Template)
	}

	diags = append(diags, diag.FromErr(err)...)

	err = d.Set(mkDataSourceVirtualEnvironmentContainerTags, tags)
	diags = append(diags, diag.FromErr(err)...)

	d.SetId(strconv.Itoa(containerID))

	return diags
}
