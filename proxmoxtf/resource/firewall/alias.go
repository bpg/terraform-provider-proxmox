/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	dvAliasComment = ""

	mkAliasName    = "name"
	mkAliasCIDR    = "cidr"
	mkAliasComment = "comment"
)

func Alias() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkAliasName: {
				Type:        schema.TypeString,
				Description: "Alias name",
				Required:    true,
				ForceNew:    false,
			},
			mkAliasCIDR: {
				Type:        schema.TypeString,
				Description: "IP/CIDR block",
				Required:    true,
				ForceNew:    false,
			},
			mkAliasComment: {
				Type:        schema.TypeString,
				Description: "Alias comment",
				Optional:    true,
				Default:     dvAliasComment,
			},
		},
		CreateContext: aliasCreate,
		ReadContext:   aliasRead,
		UpdateContext: aliasUpdate,
		DeleteContext: aliasDelete,
	}
}

func aliasCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	comment := d.Get(mkAliasComment).(string)
	name := d.Get(mkAliasName).(string)
	cidr := d.Get(mkAliasCIDR).(string)

	body := &firewall.AliasCreateRequestBody{
		Comment: &comment,
		Name:    name,
		CIDR:    cidr,
	}

	err = veClient.API().Cluster().Firewall().CreateAlias(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)

	return aliasRead(ctx, d, m)
}

func aliasRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Id()
	alias, err := veClient.API().Cluster().Firewall().GetAlias(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	aliasMap := map[string]interface{}{
		mkAliasComment: alias.Comment,
		mkAliasName:    alias.Name,
		mkAliasCIDR:    alias.CIDR,
	}

	for key, val := range aliasMap {
		err = d.Set(key, val)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func aliasUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	comment := d.Get(mkAliasComment).(string)
	cidr := d.Get(mkAliasCIDR).(string)
	newName := d.Get(mkAliasName).(string)
	previousName := d.Id()

	body := &firewall.AliasUpdateRequestBody{
		ReName:  newName,
		CIDR:    cidr,
		Comment: &comment,
	}

	err = veClient.API().Cluster().Firewall().UpdateAlias(ctx, previousName, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newName)

	return aliasRead(ctx, d, m)
}

func aliasDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Id()
	err = veClient.API().Cluster().Firewall().DeleteAlias(ctx, name)

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
