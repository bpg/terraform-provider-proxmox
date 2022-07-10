/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"strings"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
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
		Create: resourceVirtualEnvironmentClusterAliasCreate,
		Read:   resourceVirtualEnvironmentClusterAliasRead,
		Update: resourceVirtualEnvironmentClusterAliasUpdate,
		Delete: resourceVirtualEnvironmentClusterAliasDelete,
	}
}

func resourceVirtualEnvironmentClusterAliasCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	comment := d.Get(mkResourceVirtualEnvironmentClusterAliasComment).(string)
	name := d.Get(mkResourceVirtualEnvironmentClusterAliasName).(string)
	cidr := d.Get(mkResourceVirtualEnvironmentClusterAliasCIDR).(string)

	body := &proxmox.VirtualEnvironmentClusterAliasCreateRequestBody{
		Comment: &comment,
		Name:    name,
		CIDR:    cidr,
	}

	err = veClient.CreateAlias(body)

	if err != nil {
		return err
	}

	d.SetId(name)

	return resourceVirtualEnvironmentClusterAliasRead(d, m)
}

func resourceVirtualEnvironmentClusterAliasRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	name := d.Id()
	alias, err := veClient.GetAlias(name)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")
			return nil
		}

		return err
	}

	aliasMap := map[string]interface{}{
		mkResourceVirtualEnvironmentClusterAliasComment: alias.Comment,
		mkResourceVirtualEnvironmentClusterAliasName:    alias.Name,
		mkResourceVirtualEnvironmentClusterAliasCIDR:    alias.CIDR,
	}

	for key, val := range aliasMap {
		err = d.Set(key, val)

		if err != nil {
			return err
		}
	}

	return nil
}

func resourceVirtualEnvironmentClusterAliasUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
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

	err = veClient.UpdateAlias(previousName, body)

	if err != nil {
		return err
	}

	d.SetId(newName)

	return resourceVirtualEnvironmentClusterAliasRead(d, m)
}

func resourceVirtualEnvironmentClusterAliasDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return nil
	}

	name := d.Id()
	err = veClient.DeleteAlias(name)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")
			return nil
		}

		return err
	}

	d.SetId("")

	return nil
}
