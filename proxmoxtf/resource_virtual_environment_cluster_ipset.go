/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

const (
	dvResourceVirtualEnvironmentClusterIPSetComment	= ""
	dvResourceVirtualEnvironmentClusterIPSetNoMatch	= false

	mkResourceVirtualEnvironmentClusterIPSet        = "ipset"
	mkResourceVirtualEnvironmentClusterIPSetCIDR    = "cidr"
	mkResourceVirtualEnvironmentClusterIPSetName    = "name"
	mkResourceVirtualEnvironmentClusterIPSetComment = "comment"
	mkResourceVirtualEnvironmentClusterIPSetNoMatch = "nomatch"
)

func resourceVirtualEnvironmentClusterIPSet() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentClusterIPSetName: {
				Type: schema.TypeString,
				Description: "IPSet name",
				Required: true,
				ForceNew: false,
			},
			mkResourceVirtualEnvironmentClusterIPSet: {
				Type: schema.TypeList,
				Description: "List of IP or Networks",
				Optional: true,
				ForceNew: true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						mkResourceVirtualEnvironmentClusterIPSetCIDR: {
							Type: schema.TypeString,
							Description: "Network/IP specification in CIDR format",
							Required: true,
							ForceNew: true,
						},
						mkResourceVirtualEnvironmentClusterIPSetNoMatch: {
							Type: schema.TypeBool,
							Description: "No match this IP/CIDR",
							Optional: true,
							Default: dvResourceVirtualEnvironmentClusterIPSetNoMatch,
							ForceNew: true,
						},
						mkResourceVirtualEnvironmentClusterIPSetComment: {
							Type: schema.TypeString,
							Description: "IP/CIDR comment",
							Optional: true,
							Default: dvResourceVirtualEnvironmentClusterIPSetComment,
							ForceNew: true,
						},
					},
				},
			},
			mkResourceVirtualEnvironmentClusterIPSetComment: {
				Type: schema.TypeString,
				Description: "IPSet comment",
				Optional: true,
				Default: dvResourceVirtualEnvironmentClusterIPSetComment,
			},
		},
		Create: resourceVirtualEnvironmentClusterIPSetCreate,
		Read: resourceVirtualEnvironmentClusterIPSetRead,
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

	comment := d.Get(mkResourceVirtualEnvironmentClusterIPSetComment).(string)
	name := d.Get(mkResourceVirtualEnvironmentClusterIPSetName).(string)

	IPSets := d.Get(mkResourceVirtualEnvironmentClusterIPSet).([]interface{})
	IPSetsArray := make(proxmox.VirtualEnvironmentClusterIPSetContent, len(IPSets))

	for i, v := range IPSets {
		IPSetMap := v.(map[string]interface{})
		IPSetObject := proxmox.VirtualEnvironmentClusterIPSetGetResponseData{}

		cidr := IPSetMap[mkResourceVirtualEnvironmentClusterIPSetCIDR].(string)
		noMatch := IPSetMap[mkResourceVirtualEnvironmentClusterIPSetNoMatch].(bool)
		comment := IPSetMap[mkResourceVirtualEnvironmentClusterIPSetComment].(string)


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
		Name: name,
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

			err = d.Set(mkResourceVirtualEnvironmentClusterIPSetComment, v.Comment)

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

	for key, _ := range IPSet {
		d.Set(mkResourceVirtualEnvironmentClusterIPSet, IPSet[key])
	}

	return nil
}

func resourceVirtualEnvironmentClusterIPSetUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	comment := d.Get(mkResourceVirtualEnvironmentClusterIPSetComment).(string)
	newName := d.Get(mkResourceVirtualEnvironmentClusterIPSetName).(string)
	previousName := d.Id()

	body := &proxmox.VirtualEnvironmentClusterIPSetUpdateRequestBody{
		ReName: previousName,
		Name: newName,
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














