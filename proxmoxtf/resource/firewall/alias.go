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

	"github.com/bpg/terraform-provider-proxmox/proxmox/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/structure"
)

const (
	mkAliasName    = "name"
	mkAliasCIDR    = "cidr"
	mkAliasComment = "comment"
)

// Alias returns a resource to manage aliases.
func Alias() *schema.Resource {
	s := map[string]*schema.Schema{
		mkAliasName: {
			Type:        schema.TypeString,
			Description: "Alias name",
			Required:    true,
		},
		mkAliasCIDR: {
			Type:        schema.TypeString,
			Description: "IP/CIDR block",
			Required:    true,
		},
		mkAliasComment: {
			Type:        schema.TypeString,
			Description: "Alias comment",
			Optional:    true,
			Default:     "",
		},
	}

	structure.MergeSchema(s, selectorSchema())

	return &schema.Resource{
		Schema:        s,
		CreateContext: selectFirewallAPI(aliasCreate),
		ReadContext:   selectFirewallAPI(aliasRead),
		UpdateContext: selectFirewallAPI(aliasUpdate),
		DeleteContext: selectFirewallAPI(aliasDelete),
	}
}

func aliasCreate(ctx context.Context, api firewall.API, d *schema.ResourceData) diag.Diagnostics {
	comment := d.Get(mkAliasComment).(string)
	name := d.Get(mkAliasName).(string)
	cidr := d.Get(mkAliasCIDR).(string)

	body := &firewall.AliasCreateRequestBody{
		Comment: &comment,
		Name:    name,
		CIDR:    cidr,
	}

	err := api.CreateAlias(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)

	return aliasRead(ctx, api, d)
}

func aliasRead(ctx context.Context, api firewall.API, d *schema.ResourceData) diag.Diagnostics {
	name := d.Id()
	alias, err := api.GetAlias(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), "no such alias") {
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

func aliasUpdate(ctx context.Context, api firewall.API, d *schema.ResourceData) diag.Diagnostics {
	comment := d.Get(mkAliasComment).(string)
	cidr := d.Get(mkAliasCIDR).(string)
	newName := d.Get(mkAliasName).(string)
	previousName := d.Id()

	body := &firewall.AliasUpdateRequestBody{
		ReName:  newName,
		CIDR:    cidr,
		Comment: &comment,
	}

	err := api.UpdateAlias(ctx, previousName, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newName)

	return aliasRead(ctx, api, d)
}

func aliasDelete(ctx context.Context, api firewall.API, d *schema.ResourceData) diag.Diagnostics {
	name := d.Id()
	err := api.DeleteAlias(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), "no such alias") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
