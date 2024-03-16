package network

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// GetNetworkDeviceObjects returns a list of network devices from the resource data.
func GetNetworkDeviceObjects(d *schema.ResourceData) (vms.CustomNetworkDevices, error) {
	networkDevice := d.Get(MkNetworkDevice).([]interface{})
	networkDeviceObjects := make(vms.CustomNetworkDevices, len(networkDevice))

	for i, networkDeviceEntry := range networkDevice {
		block := networkDeviceEntry.(map[string]interface{})

		bridge := block[mkNetworkDeviceBridge].(string)
		disconnected := types.CustomBool(block[mkNetworkDeviceDisconnected].(bool))
		enabled := block[mkNetworkDeviceEnabled].(bool)
		firewall := types.CustomBool(block[mkNetworkDeviceFirewall].(bool))
		macAddress := block[mkNetworkDeviceMACAddress].(string)
		model := block[mkNetworkDeviceModel].(string)
		mtu := block[mkNetworkDeviceMTU].(int)
		queues := block[mkNetworkDeviceQueues].(int)
		rateLimit := block[mkNetworkDeviceRateLimit].(float64)
		trunks := block[mkNetworkDeviceTrunks].(string)
		vlanID := block[mkNetworkDeviceVLANID].(int)

		device := vms.CustomNetworkDevice{
			Enabled:  enabled,
			Firewall: &firewall,
			Model:    model,
		}

		if bridge != "" {
			device.Bridge = &bridge
		}

		if disconnected {
			device.LinkDown = &disconnected
		}

		if macAddress != "" {
			device.MACAddress = &macAddress
		}

		if queues != 0 {
			device.Queues = &queues
		}

		if rateLimit != 0 {
			device.RateLimit = &rateLimit
		}

		if vlanID != 0 {
			device.Tag = &vlanID
		}

		if trunks != "" {
			splitTrunks := strings.Split(trunks, ";")

			var trunksAsInt []int

			for _, numStr := range splitTrunks {
				num, err := strconv.Atoi(numStr)
				if err != nil {
					return nil, fmt.Errorf("error parsing trunks: %w", err)
				}

				trunksAsInt = append(trunksAsInt, num)
			}

			device.Trunks = trunksAsInt
		}

		if mtu != 0 {
			device.MTU = &mtu
		}

		networkDeviceObjects[i] = device
	}

	return networkDeviceObjects, nil
}

// ReadNetworkDeviceObjects reads the network device objects from the resource data.
func ReadNetworkDeviceObjects(d *schema.ResourceData, vmConfig *vms.GetResponseData) diag.Diagnostics {
	var diags diag.Diagnostics

	// Compare the network devices to those stored in the state.
	currentNetworkDeviceList := d.Get(MkNetworkDevice).([]interface{})

	macAddresses := make([]interface{}, MaxNetworkDevices)
	networkDeviceLast := -1
	networkDeviceList := make([]interface{}, MaxNetworkDevices)
	networkDeviceObjects := []*vms.CustomNetworkDevice{
		vmConfig.NetworkDevice0,
		vmConfig.NetworkDevice1,
		vmConfig.NetworkDevice2,
		vmConfig.NetworkDevice3,
		vmConfig.NetworkDevice4,
		vmConfig.NetworkDevice5,
		vmConfig.NetworkDevice6,
		vmConfig.NetworkDevice7,
		vmConfig.NetworkDevice8,
		vmConfig.NetworkDevice9,
		vmConfig.NetworkDevice10,
		vmConfig.NetworkDevice11,
		vmConfig.NetworkDevice12,
		vmConfig.NetworkDevice13,
		vmConfig.NetworkDevice14,
		vmConfig.NetworkDevice15,
		vmConfig.NetworkDevice16,
		vmConfig.NetworkDevice17,
		vmConfig.NetworkDevice18,
		vmConfig.NetworkDevice19,
		vmConfig.NetworkDevice20,
		vmConfig.NetworkDevice21,
		vmConfig.NetworkDevice22,
		vmConfig.NetworkDevice23,
		vmConfig.NetworkDevice24,
		vmConfig.NetworkDevice25,
		vmConfig.NetworkDevice26,
		vmConfig.NetworkDevice27,
		vmConfig.NetworkDevice28,
		vmConfig.NetworkDevice29,
		vmConfig.NetworkDevice30,
		vmConfig.NetworkDevice31,
	}

	for ni, nd := range networkDeviceObjects {
		networkDevice := map[string]interface{}{}

		if nd != nil {
			networkDeviceLast = ni

			if nd.Bridge != nil {
				networkDevice[mkNetworkDeviceBridge] = *nd.Bridge
			} else {
				networkDevice[mkNetworkDeviceBridge] = ""
			}

			networkDevice[mkNetworkDeviceEnabled] = nd.Enabled

			if nd.LinkDown != nil {
				networkDevice[mkNetworkDeviceDisconnected] = *nd.LinkDown
			} else {
				networkDevice[mkNetworkDeviceDisconnected] = false
			}

			if nd.Firewall != nil {
				networkDevice[mkNetworkDeviceFirewall] = *nd.Firewall
			} else {
				networkDevice[mkNetworkDeviceFirewall] = false
			}

			if nd.MACAddress != nil {
				macAddresses[ni] = *nd.MACAddress
			} else {
				macAddresses[ni] = ""
			}

			networkDevice[mkNetworkDeviceMACAddress] = macAddresses[ni]
			networkDevice[mkNetworkDeviceModel] = nd.Model

			if nd.Queues != nil {
				networkDevice[mkNetworkDeviceQueues] = *nd.Queues
			} else {
				networkDevice[mkNetworkDeviceQueues] = 0
			}

			if nd.RateLimit != nil {
				networkDevice[mkNetworkDeviceRateLimit] = *nd.RateLimit
			} else {
				networkDevice[mkNetworkDeviceRateLimit] = 0
			}

			if nd.Tag != nil {
				networkDevice[mkNetworkDeviceVLANID] = nd.Tag
			} else {
				networkDevice[mkNetworkDeviceVLANID] = 0
			}

			if nd.Trunks != nil {
				networkDevice[mkNetworkDeviceTrunks] = strings.Trim(
					strings.Join(strings.Fields(fmt.Sprint(nd.Trunks)), ";"), "[]")
			} else {
				networkDevice[mkNetworkDeviceTrunks] = ""
			}

			if nd.MTU != nil {
				networkDevice[mkNetworkDeviceMTU] = nd.MTU
			} else {
				networkDevice[mkNetworkDeviceMTU] = 0
			}
		} else {
			macAddresses[ni] = ""
			networkDevice[mkNetworkDeviceEnabled] = false
		}

		networkDeviceList[ni] = networkDevice
	}

	if len(currentNetworkDeviceList) > 0 || networkDeviceLast > -1 {
		err := d.Set(MkNetworkDevice, networkDeviceList[:networkDeviceLast+1])
		diags = append(diags, diag.FromErr(err)...)
	}

	err := d.Set(mkMACAddresses, macAddresses[0:len(currentNetworkDeviceList)])
	diags = append(diags, diag.FromErr(err)...)

	return diags
}

// ReadNetworkValues reads the network values from the resource data.
func ReadNetworkValues(
	ctx context.Context,
	d *schema.ResourceData,
	vmAPI *vms.Client,
	started bool,
	vmConfig *vms.GetResponseData,
	agentTimeout time.Duration,
) diag.Diagnostics {
	var diags diag.Diagnostics

	var ipv4Addresses []interface{}

	var ipv6Addresses []interface{}

	var networkInterfaceNames []interface{}

	if started {
		if vmConfig.Agent != nil && vmConfig.Agent.Enabled != nil && *vmConfig.Agent.Enabled {
			var macAddresses []interface{}

			networkInterfaces, err := vmAPI.WaitForNetworkInterfacesFromVMAgent(ctx, int(agentTimeout.Seconds()), 5, true)
			if err == nil && networkInterfaces.Result != nil {
				ipv4Addresses = make([]interface{}, len(*networkInterfaces.Result))
				ipv6Addresses = make([]interface{}, len(*networkInterfaces.Result))
				macAddresses = make([]interface{}, len(*networkInterfaces.Result))
				networkInterfaceNames = make([]interface{}, len(*networkInterfaces.Result))

				for ri, rv := range *networkInterfaces.Result {
					var rvIPv4Addresses []interface{}

					var rvIPv6Addresses []interface{}

					if rv.IPAddresses != nil {
						for _, ip := range *rv.IPAddresses {
							switch ip.Type {
							case "ipv4":
								rvIPv4Addresses = append(rvIPv4Addresses, ip.Address)
							case "ipv6":
								rvIPv6Addresses = append(rvIPv6Addresses, ip.Address)
							}
						}
					}

					ipv4Addresses[ri] = rvIPv4Addresses
					ipv6Addresses[ri] = rvIPv6Addresses
					macAddresses[ri] = strings.ToUpper(rv.MACAddress)
					networkInterfaceNames[ri] = rv.Name
				}
			}

			err = d.Set(mkMACAddresses, macAddresses)
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	e := d.Set(mkIPv4Addresses, ipv4Addresses)
	diags = append(diags, diag.FromErr(e)...)
	e = d.Set(mkIPv6Addresses, ipv6Addresses)
	diags = append(diags, diag.FromErr(e)...)
	e = d.Set(mkNetworkInterfaceNames, networkInterfaceNames)
	diags = append(diags, diag.FromErr(e)...)

	return diags
}
