package sdn

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

const (
	mkSubnetID             = "subnet"
	mkSubnetType           = "type"
	mkSubnetVnet           = "vnet"
	mkSubnetDhcpDnsServer  = "DhcpDnsServer"
	mkSubnetDhcpRange      = "DhcpRange"
	mkSubnetDnsZonePrefix  = "DnsZonePrefix"
	mkSubnetGateway        = "gateway"
	mkSubnetSnat           = "snat"
	mkSubnetDeleteSettings = "deleteSettings"
	mkSubnetDigest         = "digest"
)

func Subnet() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkSubnetID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet value",
			},
			mkSubnetType: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subnet type",
			},
		},
	}
}
