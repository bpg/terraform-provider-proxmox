/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"fmt"

	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	mkResourceVirtualEnvironmentDNSDomain   = "domain"
	mkResourceVirtualEnvironmentDNSNodeName = "node_name"
	mkResourceVirtualEnvironmentDNSServers  = "servers"
)

func resourceVirtualEnvironmentDNS() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentDNSDomain: {
				Type:        schema.TypeString,
				Description: "The DNS search domain",
				Required:    true,
			},
			mkResourceVirtualEnvironmentDNSNodeName: {
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
				ForceNew:    true,
			},
			mkResourceVirtualEnvironmentDNSServers: {
				Type:        schema.TypeList,
				Description: "The DNS servers",
				Optional:    true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 0,
				MaxItems: 3,
			},
		},
		Create: resourceVirtualEnvironmentDNSCreate,
		Read:   resourceVirtualEnvironmentDNSRead,
		Update: resourceVirtualEnvironmentDNSUpdate,
		Delete: resourceVirtualEnvironmentDNSDelete,
	}
}

func resourceVirtualEnvironmentDNSCreate(d *schema.ResourceData, m interface{}) error {
	err := resourceVirtualEnvironmentDNSUpdate(d, m)

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentDNSNodeName).(string)

	d.SetId(fmt.Sprintf("%s_dns", nodeName))

	return nil
}

func resourceVirtualEnvironmentDNSGetUpdateBody(d *schema.ResourceData, m interface{}) (*proxmox.VirtualEnvironmentDNSUpdateRequestBody, error) {
	domain := d.Get(mkResourceVirtualEnvironmentDNSDomain).(string)
	servers := d.Get(mkResourceVirtualEnvironmentDNSServers).([]interface{})

	body := &proxmox.VirtualEnvironmentDNSUpdateRequestBody{
		SearchDomain: &domain,
	}

	for i, server := range servers {
		s := server.(string)

		switch i {
		case 0:
			body.Server1 = &s
		case 1:
			body.Server2 = &s
		case 2:
			body.Server3 = &s
		}
	}

	return body, nil
}

func resourceVirtualEnvironmentDNSRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentDNSNodeName).(string)
	dns, err := veClient.GetDNS(nodeName)

	if err != nil {
		return err
	}

	if dns.SearchDomain != nil {
		d.Set(mkResourceVirtualEnvironmentDNSDomain, *dns.SearchDomain)
	} else {
		d.Set(mkResourceVirtualEnvironmentDNSDomain, "")
	}

	servers := []interface{}{}

	if dns.Server1 != nil {
		servers = append(servers, *dns.Server1)
	}

	if dns.Server2 != nil {
		servers = append(servers, *dns.Server2)
	}

	if dns.Server3 != nil {
		servers = append(servers, *dns.Server3)
	}

	d.Set(mkResourceVirtualEnvironmentDNSServers, servers)

	return nil
}

func resourceVirtualEnvironmentDNSUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentDNSNodeName).(string)

	body, err := resourceVirtualEnvironmentDNSGetUpdateBody(d, m)

	if err != nil {
		return err
	}

	err = veClient.UpdateDNS(nodeName, body)

	if err != nil {
		return err
	}

	return resourceVirtualEnvironmentDNSRead(d, m)
}

func resourceVirtualEnvironmentDNSDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")

	return nil
}
