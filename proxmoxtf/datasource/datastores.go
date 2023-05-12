/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datasource

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	mkDataSourceVirtualEnvironmentDatastoresActive         = "active"
	mkDataSourceVirtualEnvironmentDatastoresContentTypes   = "content_types"
	mkDataSourceVirtualEnvironmentDatastoresDatastoreIDs   = "datastore_ids"
	mkDataSourceVirtualEnvironmentDatastoresEnabled        = "enabled"
	mkDataSourceVirtualEnvironmentDatastoresNodeName       = "node_name"
	mkDataSourceVirtualEnvironmentDatastoresShared         = "shared"
	mkDataSourceVirtualEnvironmentDatastoresSpaceAvailable = "space_available"
	mkDataSourceVirtualEnvironmentDatastoresSpaceTotal     = "space_total"
	mkDataSourceVirtualEnvironmentDatastoresSpaceUsed      = "space_used"
	mkDataSourceVirtualEnvironmentDatastoresTypes          = "types"
)

// Datastores returns a resource for the Proxmox data store.
func Datastores() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentDatastoresActive: {
				Type:        schema.TypeList,
				Description: "Whether a datastore is active",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeBool},
			},
			mkDataSourceVirtualEnvironmentDatastoresContentTypes: {
				Type:        schema.TypeList,
				Description: "The allowed content types",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
			mkDataSourceVirtualEnvironmentDatastoresDatastoreIDs: {
				Type:        schema.TypeList,
				Description: "The datastore id",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			mkDataSourceVirtualEnvironmentDatastoresEnabled: {
				Type:        schema.TypeList,
				Description: "Whether a datastore is enabled",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeBool},
			},
			mkDataSourceVirtualEnvironmentDatastoresNodeName: {
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
			},
			mkDataSourceVirtualEnvironmentDatastoresShared: {
				Type:        schema.TypeList,
				Description: "Whether a datastore is shared",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeBool},
			},
			mkDataSourceVirtualEnvironmentDatastoresSpaceAvailable: {
				Type:        schema.TypeList,
				Description: "The available space in bytes",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			mkDataSourceVirtualEnvironmentDatastoresSpaceTotal: {
				Type:        schema.TypeList,
				Description: "The total space in bytes",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			mkDataSourceVirtualEnvironmentDatastoresSpaceUsed: {
				Type:        schema.TypeList,
				Description: "The used space in bytes",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			mkDataSourceVirtualEnvironmentDatastoresTypes: {
				Type:        schema.TypeList,
				Description: "The storage type",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		ReadContext: datastoresRead,
	}
}

func datastoresRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkDataSourceVirtualEnvironmentDatastoresNodeName).(string)
	list, err := veClient.ListDatastores(ctx, nodeName, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	active := make([]interface{}, len(list))
	contentTypes := make([]interface{}, len(list))
	datastoreIDs := make([]interface{}, len(list))
	enabled := make([]interface{}, len(list))
	shared := make([]interface{}, len(list))
	spaceAvailable := make([]interface{}, len(list))
	spaceTotal := make([]interface{}, len(list))
	spaceUsed := make([]interface{}, len(list))
	types := make([]interface{}, len(list))

	for i, v := range list {
		if v.Active != nil {
			active[i] = bool(*v.Active)
		} else {
			active[i] = true
		}

		if v.ContentTypes != nil {
			contentTypeList := []string(*v.ContentTypes)
			sort.Strings(contentTypeList)
			contentTypes[i] = contentTypeList
		} else {
			contentTypes[i] = []string{}
		}

		datastoreIDs[i] = v.ID

		if v.Enabled != nil {
			enabled[i] = bool(*v.Enabled)
		} else {
			enabled[i] = true
		}

		if v.Shared != nil {
			shared[i] = bool(*v.Shared)
		} else {
			shared[i] = true
		}

		if v.SpaceAvailable != nil {
			spaceAvailable[i] = *v.SpaceAvailable
		} else {
			spaceAvailable[i] = 0
		}

		if v.SpaceTotal != nil {
			spaceTotal[i] = *v.SpaceTotal
		} else {
			spaceTotal[i] = 0
		}

		if v.SpaceUsed != nil {
			spaceUsed[i] = *v.SpaceUsed
		} else {
			spaceUsed[i] = 0
		}

		types[i] = v.Type
	}

	d.SetId(fmt.Sprintf("%s_datastores", nodeName))

	err = d.Set(mkDataSourceVirtualEnvironmentDatastoresActive, active)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentDatastoresContentTypes, contentTypes)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentDatastoresDatastoreIDs, datastoreIDs)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentDatastoresEnabled, enabled)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentDatastoresShared, shared)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentDatastoresSpaceAvailable, spaceAvailable)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentDatastoresSpaceTotal, spaceTotal)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentDatastoresSpaceUsed, spaceUsed)
	diags = append(diags, diag.FromErr(err)...)
	err = d.Set(mkDataSourceVirtualEnvironmentDatastoresTypes, types)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}
