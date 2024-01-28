/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/pools"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	dvResourceVirtualEnvironmentPoolComment = ""

	mkResourceVirtualEnvironmentPoolComment            = "comment"
	mkResourceVirtualEnvironmentPoolMembers            = "members"
	mkResourceVirtualEnvironmentPoolMembersDatastoreID = "datastore_id"
	mkResourceVirtualEnvironmentPoolMembersID          = "id"
	mkResourceVirtualEnvironmentPoolMembersNodeName    = "node_name"
	mkResourceVirtualEnvironmentPoolMembersType        = "type"
	mkResourceVirtualEnvironmentPoolMembersVMID        = "vm_id"
	mkResourceVirtualEnvironmentPoolPoolID             = "pool_id"
)

// Pool returns a resource that manages pools.
func Pool() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentPoolComment: {
				Type:        schema.TypeString,
				Description: "The pool comment",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentPoolComment,
			},
			mkResourceVirtualEnvironmentPoolMembers: {
				Type:        schema.TypeList,
				Description: "The pool members",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentPoolMembersDatastoreID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The datastore id",
						},
						mkResourceVirtualEnvironmentPoolMembersID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The member id",
						},
						mkResourceVirtualEnvironmentPoolMembersNodeName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The node name",
						},
						mkResourceVirtualEnvironmentPoolMembersType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The member type",
						},
						mkResourceVirtualEnvironmentPoolMembersVMID: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The virtual machine id",
						},
					},
				},
			},
			mkResourceVirtualEnvironmentPoolPoolID: {
				Type:        schema.TypeString,
				Description: "The pool id",
				Required:    true,
				ForceNew:    true,
			},
		},
		CreateContext: poolCreate,
		ReadContext:   poolRead,
		UpdateContext: poolUpdate,
		DeleteContext: poolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				d.SetId(d.Id())
				err := d.Set(mkResourceVirtualEnvironmentPoolPoolID, d.Id())
				if err != nil {
					return nil, fmt.Errorf("failed setting state during import: %w", err)
				}

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func poolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	comment := d.Get(mkResourceVirtualEnvironmentPoolComment).(string)
	poolID := d.Get(mkResourceVirtualEnvironmentPoolPoolID).(string)

	body := &pools.PoolCreateRequestBody{
		Comment: &comment,
		ID:      poolID,
	}

	err = api.Pool().CreatePool(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(poolID)

	return poolRead(ctx, d, m)
}

func poolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	poolID := d.Id()
	pool, err := api.Pool().GetPool(ctx, poolID)
	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if pool.Comment != nil {
		err = d.Set(mkResourceVirtualEnvironmentPoolComment, pool.Comment)
	} else {
		err = d.Set(mkResourceVirtualEnvironmentPoolComment, "")
	}
	diags = append(diags, diag.FromErr(err)...)

	members := make([]interface{}, len(pool.Members))

	for i, v := range pool.Members {
		values := map[string]interface{}{}

		values[mkResourceVirtualEnvironmentPoolMembersID] = v.ID
		values[mkResourceVirtualEnvironmentPoolMembersNodeName] = v.Node

		if v.DatastoreID != nil {
			values[mkResourceVirtualEnvironmentPoolMembersDatastoreID] = v.DatastoreID
		} else {
			values[mkResourceVirtualEnvironmentPoolMembersDatastoreID] = ""
		}

		values[mkResourceVirtualEnvironmentPoolMembersType] = v.Type

		if v.VMID != nil {
			values[mkResourceVirtualEnvironmentPoolMembersVMID] = v.VMID
		} else {
			values[mkResourceVirtualEnvironmentPoolMembersVMID] = 0
		}

		members[i] = values
	}

	err = d.Set(mkResourceVirtualEnvironmentPoolMembers, members)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}

func poolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	comment := d.Get(mkResourceVirtualEnvironmentPoolComment).(string)
	poolID := d.Id()

	body := &pools.PoolUpdateRequestBody{
		Comment: &comment,
	}

	err = api.Pool().UpdatePool(ctx, poolID, body)
	if err != nil {
		return diag.FromErr(err)
	}

	return poolRead(ctx, d, m)
}

func poolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	api, err := config.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}

	poolID := d.Id()
	err = api.Pool().DeletePool(ctx, poolID)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")

			return nil
		}

		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
