package sdn

import (
	"os"
	"testing"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/subnets"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/vnets"
	"github.com/bpg/terraform-provider-proxmox/proxmox/cluster/sdn/zones"
	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
)

const (
	testZoneID          = "testzone"
	testVnetID          = "testvnet"
	testSubnetCIDR      = "10.10.0.0/24"
	testSubnetCanonical = "testzone-10.10.0.0-24"
	testGateway         = "10.10.0.1"
	testDNS             = "10.10.0.53"
	testDHCPStart       = "10.10.0.10"
	testDHCPEnd         = "10.10.0.100"
)

type testClients struct {
	zone   *zones.Client
	vnet   *vnets.Client
	subnet *subnets.Client
}

func getTestClients(t *testing.T) *testClients {
	t.Helper()

	apiToken := os.Getenv("PVE_TOKEN")

	url := os.Getenv("PVE_URL")
	if apiToken == "" || url == "" {
		t.Skip("PVE_TOKEN and PVE_URL must be set")
	}

	conn, err := api.NewConnection(url, true, "")
	if err != nil {
		t.Fatalf("connection error: %v", err)
	}

	creds := api.Credentials{TokenCredentials: &api.TokenCredentials{APIToken: apiToken}}

	client, err := api.NewClient(creds, conn)
	if err != nil {
		t.Fatalf("client error: %v", err)
	}

	return &testClients{
		zone:   &zones.Client{Client: client},
		vnet:   &vnets.Client{Client: client},
		subnet: &subnets.Client{Client: client},
	}
}

func TestSDNLifecycle(t *testing.T) {
	clients := getTestClients(t)

	t.Run("Create Zone", func(t *testing.T) {
		t.Parallel()

		err := clients.zone.CreateZone(t.Context(), &zones.ZoneRequestData{
			ZoneData: zones.ZoneData{
				ID:     testZoneID,
				Type:   ptr.Ptr("vlan"),
				IPAM:   ptr.Ptr("pve"),
				Bridge: ptr.Ptr("vmbr0"),
				MTU:    ptr.Ptr(int64(1500)),
				Nodes:  ptr.Ptr("pvenode1"),
			},
		})
		if err != nil {
			t.Fatalf("CreateZone failed: %v", err)
		}
	})

	t.Run("Get Zone", func(t *testing.T) {
		t.Parallel()

		zone, err := clients.zone.GetZone(t.Context(), testZoneID)
		if err != nil {
			t.Fatalf("GetZone failed: %v", err)
		}

		t.Logf("Zone: %+v", zone)
	})

	t.Run("Update Zone", func(t *testing.T) {
		t.Parallel()

		err := clients.zone.UpdateZone(t.Context(), &zones.ZoneRequestData{
			ZoneData: zones.ZoneData{
				ID:     testZoneID,
				Nodes:  ptr.Ptr("updatednode"),
				Bridge: ptr.Ptr("vmbr1"), // simulate a VLAN-related update.
			},
		})
		if err != nil {
			t.Fatalf("UpdateZone failed: %v", err)
		}
	})

	t.Run("Create VNet", func(t *testing.T) {
		t.Parallel()

		err := clients.vnet.CreateVnet(t.Context(), &vnets.VnetRequestData{
			VnetData: vnets.VnetData{
				ID:           testVnetID,
				Zone:         ptr.Ptr(testZoneID),
				Alias:        ptr.Ptr("TestVNet"),
				IsolatePorts: ptr.Ptr(int64(0)),
				Type:         ptr.Ptr("vnet"),
				Tag:          ptr.Ptr(int64(100)),
				VlanAware:    ptr.Ptr(int64(0)),
			},
		})
		if err != nil {
			t.Fatalf("CreateVnet failed: %v", err)
		}
	})

	t.Run("Get VNet", func(t *testing.T) {
		t.Parallel()

		vnet, err := clients.vnet.GetVnet(t.Context(), testVnetID)
		if err != nil {
			t.Fatalf("GetVnet failed: %v", err)
		}

		t.Logf("VNet: %+v", vnet)
	})

	t.Run("Update VNet", func(t *testing.T) {
		t.Parallel()

		err := clients.vnet.UpdateVnet(t.Context(), &vnets.VnetRequestData{
			VnetData: vnets.VnetData{
				ID:    testVnetID,
				Alias: ptr.Ptr("UpdatedAlias"),
			},
		})
		if err != nil {
			t.Fatalf("UpdateVnet failed: %v", err)
		}
	})

	t.Run("Create Subnet", func(t *testing.T) {
		t.Parallel()

		ptr := &subnets.SubnetData{
			ID:            testSubnetCIDR,
			Vnet:          ptr.Ptr(testVnetID),
			Type:          ptr.Ptr("subnet"),
			Gateway:       ptr.Ptr(testGateway),
			DHCPDNSServer: ptr.Ptr(testDNS),
			DHCPRange: subnets.DHCPRangeList{
				{StartAddress: testDHCPStart, EndAddress: testDHCPEnd},
			},
			SNAT: ptr.Ptr(int64(1)),
		}
		req := &subnets.SubnetRequestData{
			EncodedSubnetData: *ptr.ToEncoded(),
		}

		err := clients.subnet.CreateSubnet(t.Context(), testVnetID, req)
		if err != nil {
			t.Fatalf("CreateSubnet failed: %v", err)
		}
	})

	t.Run("Get Subnet", func(t *testing.T) {
		t.Parallel()

		subnet, err := clients.subnet.GetSubnet(t.Context(), testVnetID, testSubnetCanonical)
		if err != nil {
			t.Fatalf("GetSubnet failed: %v", err)
		}

		t.Logf("Subnet: %+v", subnet)
	})

	t.Run("Update Subnet", func(t *testing.T) {
		t.Parallel()

		ptr := &subnets.SubnetData{
			ID:      testSubnetCanonical,
			Vnet:    ptr.Ptr(testVnetID),
			Gateway: ptr.Ptr("10.10.0.254"),
		}
		req := &subnets.SubnetRequestData{
			EncodedSubnetData: *ptr.ToEncoded(),
		}

		err := clients.subnet.UpdateSubnet(t.Context(), testVnetID, req)
		if err != nil {
			t.Fatalf("UpdateSubnet failed: %v", err)
		}
	})

	t.Run("Delete Subnet", func(t *testing.T) {
		t.Parallel()

		err := clients.subnet.DeleteSubnet(t.Context(), testVnetID, testSubnetCanonical)
		if err != nil {
			t.Fatalf("DeleteSubnet failed: %v", err)
		}
	})

	t.Run("Delete VNet", func(t *testing.T) {
		t.Parallel()

		err := clients.vnet.DeleteVnet(t.Context(), testVnetID)
		if err != nil {
			t.Fatalf("DeleteVnet failed: %v", err)
		}
	})

	t.Run("Delete Zone", func(t *testing.T) {
		t.Parallel()

		err := clients.zone.DeleteZone(t.Context(), testZoneID)
		if err != nil {
			t.Fatalf("DeleteZone failed: %v", err)
		}
	})
}
