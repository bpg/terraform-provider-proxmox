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
)

const (
	mkAliasName    = "name"
	mkAliasCIDR    = "cidr"
	mkAliasComment = "comment"
)

func AliasSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
}

func AliasCreate(ctx context.Context, fw firewall.API, d *schema.ResourceData) diag.Diagnostics {
	comment := d.Get(mkAliasComment).(string)
	name := d.Get(mkAliasName).(string)
	cidr := d.Get(mkAliasCIDR).(string)

	body := &firewall.AliasCreateRequestBody{
		Comment: &comment,
		Name:    name,
		CIDR:    cidr,
	}

	err := fw.CreateAlias(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)

	return AliasRead(ctx, fw, d)
}

func AliasRead(ctx context.Context, fw firewall.API, d *schema.ResourceData) diag.Diagnostics {
	name := d.Id()
	alias, err := fw.GetAlias(ctx, name)
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

func AliasUpdate(ctx context.Context, fw firewall.API, d *schema.ResourceData) diag.Diagnostics {
	comment := d.Get(mkAliasComment).(string)
	cidr := d.Get(mkAliasCIDR).(string)
	newName := d.Get(mkAliasName).(string)
	previousName := d.Id()

	body := &firewall.AliasUpdateRequestBody{
		ReName:  newName,
		CIDR:    cidr,
		Comment: &comment,
	}

	err := fw.UpdateAlias(ctx, previousName, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newName)

	return AliasRead(ctx, fw, d)
}

func AliasDelete(ctx context.Context, fw firewall.API, d *schema.ResourceData) diag.Diagnostics {
	name := d.Id()
	err := fw.DeleteAlias(ctx, name)
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
