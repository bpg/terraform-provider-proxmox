/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"context"
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func resourceVirtualEnvironmentPool() *schema.Resource {
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
		CreateContext: resourceVirtualEnvironmentPoolCreate,
		ReadContext:   resourceVirtualEnvironmentPoolRead,
		UpdateContext: resourceVirtualEnvironmentPoolUpdate,
		DeleteContext: resourceVirtualEnvironmentPoolDelete,
	}
}

func resourceVirtualEnvironmentPoolCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	comment := d.Get(mkResourceVirtualEnvironmentPoolComment).(string)
	poolID := d.Get(mkResourceVirtualEnvironmentPoolPoolID).(string)

	body := &proxmox.VirtualEnvironmentPoolCreateRequestBody{
		Comment: &comment,
		ID:      poolID,
	}

	err = veClient.CreatePool(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(poolID)

	return resourceVirtualEnvironmentPoolRead(ctx, d, m)
}

func resourceVirtualEnvironmentPoolRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	poolID := d.Id()
	pool, err := veClient.GetPool(ctx, poolID)
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

func resourceVirtualEnvironmentPoolUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	comment := d.Get(mkResourceVirtualEnvironmentPoolComment).(string)
	poolID := d.Id()

	body := &proxmox.VirtualEnvironmentPoolUpdateRequestBody{
		Comment: &comment,
	}

	err = veClient.UpdatePool(ctx, poolID, body)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceVirtualEnvironmentPoolRead(ctx, d, m)
}

func resourceVirtualEnvironmentPoolDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	poolID := d.Id()
	err = veClient.DeletePool(ctx, poolID)

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
