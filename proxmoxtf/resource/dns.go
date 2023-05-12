/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes"
	"github.com/bpg/terraform-provider-proxmox/proxmoxtf"
)

const (
	mkResourceVirtualEnvironmentDNSDomain   = "domain"
	mkResourceVirtualEnvironmentDNSNodeName = "node_name"
	mkResourceVirtualEnvironmentDNSServers  = "servers"
)

// DNS returns a resource that manages DNS settings for a node.
func DNS() *schema.Resource {
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
		CreateContext: dnsCreate,
		ReadContext:   dnsRead,
		UpdateContext: dnsUpdate,
		DeleteContext: dnsDelete,
	}
}

func dnsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diags := dnsUpdate(ctx, d, m)
	if diags.HasError() {
		return diags
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentDNSNodeName).(string)

	d.SetId(fmt.Sprintf("%s_dns", nodeName))

	return nil
}

func dnsGetUpdateBody(d *schema.ResourceData) *nodes.DNSUpdateRequestBody {
	domain := d.Get(mkResourceVirtualEnvironmentDNSDomain).(string)
	servers := d.Get(mkResourceVirtualEnvironmentDNSServers).([]interface{})

	body := &nodes.DNSUpdateRequestBody{
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

	return body
}

func dnsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentDNSNodeName).(string)
	api := veClient.API().Node(nodeName)

	dns, err := api.GetDNS(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if dns.SearchDomain != nil {
		err = d.Set(mkResourceVirtualEnvironmentDNSDomain, *dns.SearchDomain)
	} else {
		err = d.Set(mkResourceVirtualEnvironmentDNSDomain, "")
	}
	diags = append(diags, diag.FromErr(err)...)

	var servers []interface{}

	if dns.Server1 != nil {
		servers = append(servers, *dns.Server1)
	}

	if dns.Server2 != nil {
		servers = append(servers, *dns.Server2)
	}

	if dns.Server3 != nil {
		servers = append(servers, *dns.Server3)
	}

	err = d.Set(mkResourceVirtualEnvironmentDNSServers, servers)
	diags = append(diags, diag.FromErr(err)...)

	return diags
}

func dnsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(proxmoxtf.ProviderConfiguration)
	veClient, err := config.GetVEClient()
	if err != nil {
		return diag.FromErr(err)
	}

	nodeName := d.Get(mkResourceVirtualEnvironmentDNSNodeName).(string)
	api := veClient.API().Node(nodeName)

	body := dnsGetUpdateBody(d)

	err = api.UpdateDNS(ctx, body)
	if err != nil {
		return diag.FromErr(err)
	}

	return dnsRead(ctx, d, m)
}

func dnsDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}
