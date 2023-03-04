/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/firewall"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	dvResourceVirtualEnvironmentClusterAliasComment = ""

	mkResourceVirtualEnvironmentClusterAliasName    = "name"
	mkResourceVirtualEnvironmentClusterAliasCIDR    = "cidr"
	mkResourceVirtualEnvironmentClusterAliasComment = "comment"
)

func ResourceVirtualEnvironmentClusterAlias() *schema.Resource {
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
		CreateContext: ResourceVirtualEnvironmentClusterAliasCreate,
		ReadContext:   ResourceVirtualEnvironmentClusterAliasRead,
		UpdateContext: ResourceVirtualEnvironmentClusterAliasUpdate,
		DeleteContext: ResourceVirtualEnvironmentClusterAliasDelete,
	}
}

func ResourceVirtualEnvironmentClusterAliasCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	comment := d.Get(mkResourceVirtualEnvironmentClusterAliasComment).(string)
	name := d.Get(mkResourceVirtualEnvironmentClusterAliasName).(string)
	cidr := d.Get(mkResourceVirtualEnvironmentClusterAliasCIDR).(string)

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

	return ResourceVirtualEnvironmentClusterAliasRead(ctx, d, m)
}

func ResourceVirtualEnvironmentClusterAliasRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
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

func ResourceVirtualEnvironmentClusterAliasUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	comment := d.Get(mkResourceVirtualEnvironmentClusterAliasComment).(string)
	cidr := d.Get(mkResourceVirtualEnvironmentClusterAliasCIDR).(string)
	newName := d.Get(mkResourceVirtualEnvironmentClusterAliasName).(string)
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

	return ResourceVirtualEnvironmentClusterAliasRead(ctx, d, m)
}

func ResourceVirtualEnvironmentClusterAliasDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
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
