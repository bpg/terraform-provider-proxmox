/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func dataSourceVirtualEnvironmentDatastores() *schema.Resource {
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
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
		Read: dataSourceVirtualEnvironmentDatastoresRead,
	}
}

func dataSourceVirtualEnvironmentDatastoresRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkDataSourceVirtualEnvironmentDatastoresNodeName).(string)
	list, err := veClient.ListDatastores(nodeName, nil)

	if err != nil {
		return err
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

	d.Set(mkDataSourceVirtualEnvironmentDatastoresActive, active)
	d.Set(mkDataSourceVirtualEnvironmentDatastoresContentTypes, contentTypes)
	d.Set(mkDataSourceVirtualEnvironmentDatastoresDatastoreIDs, datastoreIDs)
	d.Set(mkDataSourceVirtualEnvironmentDatastoresEnabled, enabled)
	d.Set(mkDataSourceVirtualEnvironmentDatastoresShared, shared)
	d.Set(mkDataSourceVirtualEnvironmentDatastoresSpaceAvailable, spaceAvailable)
	d.Set(mkDataSourceVirtualEnvironmentDatastoresSpaceTotal, spaceTotal)
	d.Set(mkDataSourceVirtualEnvironmentDatastoresSpaceUsed, spaceUsed)
	d.Set(mkDataSourceVirtualEnvironmentDatastoresTypes, types)

	return nil
}
