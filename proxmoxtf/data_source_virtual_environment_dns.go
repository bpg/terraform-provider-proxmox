/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

const (
	mkDataSourceVirtualEnvironmentDNSDomain   = "domain"
	mkDataSourceVirtualEnvironmentDNSNodeName = "node_name"
	mkDataSourceVirtualEnvironmentDNSServers  = "servers"
)

func dataSourceVirtualEnvironmentDNS() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentDNSDomain: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The DNS search domain",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentDNSNodeName: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The node name",
				Required:    true,
			},
			mkDataSourceVirtualEnvironmentDNSServers: &schema.Schema{
				Type:        schema.TypeList,
				Description: "The DNS servers",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Read: dataSourceVirtualEnvironmentDNSRead,
	}
}

func dataSourceVirtualEnvironmentDNSRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	nodeName := d.Get(mkDataSourceVirtualEnvironmentDNSNodeName).(string)
	dns, err := veClient.GetDNS(nodeName)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s_dns", nodeName))

	if dns.SearchDomain != nil {
		d.Set(mkDataSourceVirtualEnvironmentDNSDomain, *dns.SearchDomain)
	} else {
		d.Set(mkDataSourceVirtualEnvironmentDNSDomain, "")
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

	d.Set(mkDataSourceVirtualEnvironmentDNSServers, servers)

	return nil
}
