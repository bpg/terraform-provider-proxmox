package network

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmoxtf/test"
)

func TestNetworkSchema(t *testing.T) {
	t.Parallel()

	s := Schema()

	test.AssertComputedAttributes(t, s, []string{
		mkIPv4Addresses,
		mkIPv6Addresses,
		mkMACAddresses,
		mkNetworkInterfaceNames,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkIPv4Addresses:         schema.TypeList,
		mkIPv6Addresses:         schema.TypeList,
		mkMACAddresses:          schema.TypeList,
		mkNetworkInterfaceNames: schema.TypeList,
	})

	deviceSchema := test.AssertNestedSchemaExistence(
		t,
		s,
		MkNetworkDevice,
	)

	test.AssertOptionalArguments(t, deviceSchema, []string{
		mkNetworkDeviceBridge,
		mkNetworkDeviceEnabled,
		mkNetworkDeviceMACAddress,
		mkNetworkDeviceModel,
		mkNetworkDeviceRateLimit,
		mkNetworkDeviceVLANID,
		mkNetworkDeviceMTU,
	})

	test.AssertValueTypes(t, deviceSchema, map[string]schema.ValueType{
		mkNetworkDeviceBridge:     schema.TypeString,
		mkNetworkDeviceEnabled:    schema.TypeBool,
		mkNetworkDeviceMACAddress: schema.TypeString,
		mkNetworkDeviceModel:      schema.TypeString,
		mkNetworkDeviceRateLimit:  schema.TypeFloat,
		mkNetworkDeviceVLANID:     schema.TypeInt,
		mkNetworkDeviceMTU:        schema.TypeInt,
	})
}
