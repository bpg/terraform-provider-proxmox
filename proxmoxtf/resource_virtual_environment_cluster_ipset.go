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
	dvResourceVirtualEnvironmentClusterIPSetCIDRComment = ""
	dvResourceVirtualEnvironmentClusterIPSetCIDRNoMatch = false

	mkResourceVirtualEnvironmentClusterIPSetName        = "name"
	mkResourceVirtualEnvironmentClusterIPSetCIDR        = "cidr"
	mkResourceVirtualEnvironmentClusterIPSetCIDRName    = "name"
	mkResourceVirtualEnvironmentClusterIPSetCIDRComment = "comment"
	mkResourceVirtualEnvironmentClusterIPSetCIDRNoMatch = "nomatch"
)

func resourceVirtualEnvironmentClusterIPSet() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentClusterIPSetName: {
				Type:        schema.TypeString,
				Description: "IPSet name",
				Required:    true,
				ForceNew:    false,
			},
			mkResourceVirtualEnvironmentClusterIPSetCIDR: {
				Type:        schema.TypeList,
				Description: "List of IP or Networks",
				Optional:    true,
				ForceNew:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentClusterIPSetCIDRName: {
							Type:        schema.TypeString,
							Description: "Network/IP specification in CIDR format",
							Required:    true,
							ForceNew:    true,
						},
						mkResourceVirtualEnvironmentClusterIPSetCIDRNoMatch: {
							Type:        schema.TypeBool,
							Description: "No match this IP/CIDR",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentClusterIPSetCIDRNoMatch,
							ForceNew:    true,
						},
						mkResourceVirtualEnvironmentClusterIPSetCIDRComment: {
							Type:        schema.TypeString,
							Description: "IP/CIDR comment",
							Optional:    true,
							Default:     dvResourceVirtualEnvironmentClusterIPSetCIDRComment,
							ForceNew:    true,
						},
					},
				},
			},
			mkResourceVirtualEnvironmentClusterIPSetCIDRComment: {
				Type:        schema.TypeString,
				Description: "IPSet comment",
				Optional:    true,
				Default:     dvResourceVirtualEnvironmentClusterIPSetCIDRComment,
			},
		},
		Create: resourceVirtualEnvironmentClusterIPSetCreate,
		Read:   resourceVirtualEnvironmentClusterIPSetRead,
		Update: resourceVirtualEnvironmentClusterIPSetUpdate,
		Delete: resourceVirtualEnvironmentClusterIPSetDelete,
	}
}

func resourceVirtualEnvironmentClusterIPSetCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	comment := d.Get(mkResourceVirtualEnvironmentClusterIPSetCIDRComment).(string)
	name := d.Get(mkResourceVirtualEnvironmentClusterIPSetName).(string)

	IPSets := d.Get(mkResourceVirtualEnvironmentClusterIPSetCIDR).([]interface{})
	IPSetsArray := make(proxmox.VirtualEnvironmentClusterIPSetContent, len(IPSets))

	for i, v := range IPSets {
		IPSetMap := v.(map[string]interface{})
		IPSetObject := proxmox.VirtualEnvironmentClusterIPSetGetResponseData{}

		cidr := IPSetMap[mkResourceVirtualEnvironmentClusterIPSetCIDRName].(string)
		noMatch := IPSetMap[mkResourceVirtualEnvironmentClusterIPSetCIDRNoMatch].(bool)
		comment := IPSetMap[mkResourceVirtualEnvironmentClusterIPSetCIDRComment].(string)

		IPSetObject.Comment = comment
		IPSetObject.CIDR = cidr

		if noMatch {
			noMatchBool := proxmox.CustomBool(true)
			IPSetObject.NoMatch = &noMatchBool
		}

		IPSetsArray[i] = IPSetObject
	}

	body := &proxmox.VirtualEnvironmentClusterIPSetCreateRequestBody{
		Comment: comment,
		Name:    name,
	}

	err = veClient.CreateIPSet(body)

	if err != nil {
		return err
	}

	for _, v := range IPSetsArray {
		err = veClient.AddCIDRToIPSet(name, &v)

		if err != nil {
			return err
		}
	}

	d.SetId(name)
	return resourceVirtualEnvironmentClusterIPSetRead(d, m)
}

func resourceVirtualEnvironmentClusterIPSetRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	name := d.Id()

	allIPSets, err := veClient.GetListIPSets()

	if err != nil {
		return err
	}

	for _, v := range allIPSets.Data {
		if v.Name == name {
			err = d.Set(mkResourceVirtualEnvironmentClusterIPSetName, v.Name)

			if err != nil {
				return err
			}

			err = d.Set(mkResourceVirtualEnvironmentClusterIPSetCIDRComment, v.Comment)

			if err != nil {
				return err
			}
		}
	}

	IPSet, err := veClient.GetListIPSetContent(name)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")
			return nil
		}

		return err
	}

	for key := range IPSet {
		d.Set(mkResourceVirtualEnvironmentClusterIPSetCIDR, IPSet[key])
	}

	return nil
}

func resourceVirtualEnvironmentClusterIPSetUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	comment := d.Get(mkResourceVirtualEnvironmentClusterIPSetCIDRComment).(string)
	newName := d.Get(mkResourceVirtualEnvironmentClusterIPSetName).(string)
	previousName := d.Id()

	body := &proxmox.VirtualEnvironmentClusterIPSetUpdateRequestBody{
		ReName:  previousName,
		Name:    newName,
		Comment: &comment,
	}

	err = veClient.UpdateIPSet(body)

	if err != nil {
		return err
	}

	d.SetId(newName)

	return resourceVirtualEnvironmentClusterIPSetRead(d, m)
}

func resourceVirtualEnvironmentClusterIPSetDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return nil
	}

	name := d.Id()

	IPSetContent, err := veClient.GetListIPSetContent(name)

	if err != nil {
		return err
	}

	// PVE requires content of IPSet be cleared before removal
	if len(IPSetContent) > 0 {
		for _, IPSet := range IPSetContent {
			err = veClient.DeleteIPSetContent(name, IPSet.CIDR)
			if err != nil {
				return err
			}
		}
	}

	err = veClient.DeleteIPSet(name)

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
