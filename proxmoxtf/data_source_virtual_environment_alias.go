/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	dvDataVirtualEnvironmentClusterAliasComment 	           = ""

	mkDataSourceVirtualEnvironmentClusterAliasName             = "name"
	mkDataSourceVirtualEnvironmentClusterAliasCIDR             = "cidr"
	mkDataSourceVirtualEnvironmentClusterAliasComment          = "comment"

)

func dataSourceVirtualEnvironmentClusterAlias() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentClusterAliasName: {
				Type:        schema.TypeString,
				Description: "Alias name",
				Required:    true,
			},
			mkDataSourceVirtualEnvironmentClusterAliasCIDR: {
				Type:        schema.TypeString,
				Description: "IP/CIDR block",
				Computed:    true,
			},
			mkDataSourceVirtualEnvironmentClusterAliasComment: {
				Type:        schema.TypeString,
				Description: "Alias comment",
				Computed:    true,
			},

		},
		Read: dataSourceVirtualEnvironmentAliasRead,
	}
}

func dataSourceVirtualEnvironmentAliasRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	AliasID := d.Get(mkDataSourceVirtualEnvironmentClusterAliasName).(string)
	Alias, err := veClient.GetAlias(AliasID)

	if err != nil {
		return err
	}

	d.SetId(AliasID)

	d.Set(mkDataSourceVirtualEnvironmentClusterAliasCIDR, Alias.CIDR)

	if Alias.Comment != nil {
		d.Set(mkDataSourceVirtualEnvironmentClusterAliasComment, Alias.Comment)
	} else {
		d.Set(mkDataSourceVirtualEnvironmentClusterAliasComment, dvDataVirtualEnvironmentClusterAliasComment)
	}

	return nil
}
