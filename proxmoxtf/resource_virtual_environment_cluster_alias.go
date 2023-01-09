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
	dvResourceVirtualEnvironmentClusterAliasComment = ""

	mkResourceVirtualEnvironmentClusterAliasName    = "name"
	mkResourceVirtualEnvironmentClusterAliasCIDR    = "cidr"
	mkResourceVirtualEnvironmentClusterAliasComment = "comment"
)

func resourceVirtualEnvironmentClusterAlias() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentClusterAliasName: {
				Type:        schema.TypeString,
				Description: "Alias name",
				Required:    true,
				ForceNew:    false,
			},
			mkResourceVirtualEnvironmentClusterAliasCIDR: {
				Type:        schema.TypeString,
				Description: "IP/CIDR block",
				Required:    true,
				ForceNew:    false,
			},
			mkResourceVirtualEnvironmentClusterAliasComment: {
				Type:        schema.TypeString,
				Description: "Alias comment",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentClusterAliasComment,
			},
		},
		CreateContext: resourceVirtualEnvironmentClusterAliasCreate,
		ReadContext:   resourceVirtualEnvironmentClusterAliasRead,
		UpdateContext: resourceVirtualEnvironmentClusterAliasUpdate,
		DeleteContext: resourceVirtualEnvironmentClusterAliasDelete,
	}
}

func resourceVirtualEnvironmentClusterAliasCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	comment := d.Get(mkResourceVirtualEnvironmentClusterAliasComment).(string)
	name := d.Get(mkResourceVirtualEnvironmentClusterAliasName).(string)
	cidr := d.Get(mkResourceVirtualEnvironmentClusterAliasCIDR).(string)

	body := &proxmox.VirtualEnvironmentClusterAliasCreateRequestBody{
		Comment: &comment,
		Name:    name,
		CIDR:    cidr,
	}

	err = veClient.CreateAlias(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)

	return resourceVirtualEnvironmentClusterAliasRead(ctx, d, m)
}

func resourceVirtualEnvironmentClusterAliasRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Id()
	alias, err := veClient.GetAlias(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	aliasMap := map[string]interface{}{
		mkResourceVirtualEnvironmentClusterAliasComment: alias.Comment,
		mkResourceVirtualEnvironmentClusterAliasName:    alias.Name,
		mkResourceVirtualEnvironmentClusterAliasCIDR:    alias.CIDR,
	}

	for key, val := range aliasMap {
		err = d.Set(key, val)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceVirtualEnvironmentClusterAliasUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	comment := d.Get(mkResourceVirtualEnvironmentClusterAliasComment).(string)
	cidr := d.Get(mkResourceVirtualEnvironmentClusterAliasCIDR).(string)
	newName := d.Get(mkResourceVirtualEnvironmentClusterAliasName).(string)
	previousName := d.Id()

	body := &proxmox.VirtualEnvironmentClusterAliasUpdateRequestBody{
		ReName:  newName,
		CIDR:    cidr,
		Comment: &comment,
	}

	err = veClient.UpdateAlias(ctx, previousName, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newName)

	return resourceVirtualEnvironmentClusterAliasRead(ctx, d, m)
}

func resourceVirtualEnvironmentClusterAliasDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Id()
	err = veClient.DeleteAlias(ctx, name)

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
