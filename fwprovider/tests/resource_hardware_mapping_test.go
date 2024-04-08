/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package tests

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/bpg/terraform-provider-proxmox/fwprovider"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/types"
	"github.com/bpg/terraform-provider-proxmox/fwprovider/validators"
	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

const (
	accTestHardwareMappingNamePCI = "proxmox_virtual_environment_hardware_mapping_pci.test"
	accTestHardwareMappingNameUSB = "proxmox_virtual_environment_hardware_mapping_usb.test"
)

type accTestHardwareMappingFakeData struct {
	Comments        []string `fake:"{sentence:3}"         fakesize:"2"`
	MapComments     []string `fake:"{sentence:3}"         fakesize:"2"`
	MapDeviceIDs    []string `fake:"{linuxdeviceid}"      fakesize:"2"`
	MapIOMMUGroups  []uint   `fake:"{number:1,20}"        fakesize:"2"`
	MapPathsPCI     []string `fake:"{linuxdevicepathpci}" fakesize:"2"`
	MapPathsUSB     []string `fake:"{linuxdevicepathusb}" fakesize:"2"`
	MapSubsystemIDs []string `fake:"{linuxdeviceid}"      fakesize:"2"`
	MediatedDevices bool     `fake:"{bool}"`
	Names           []string `fake:"{noun}"               fakesize:"2"`
}

func testAccResourceHardwareMappingInit(t *testing.T) (*accTestHardwareMappingFakeData, *testEnvironment) {
	t.Helper()

	// Register a new custom function to generate random Linux device IDs.
	gofakeit.AddFuncLookup(
		"linuxdeviceid", gofakeit.Info{
			Category:    "custom",
			Description: "Random Linux device ID",
			Example:     "8086:5916",
			Output:      "string",
			Generate: func(f *gofakeit.Faker, _ *gofakeit.MapParams, _ *gofakeit.Info) (any, error) {
				return f.Regex(proxmoxtypes.HardwareMappingDeviceIDAttrValueRegEx.String()), nil
			},
		},
	)
	// Register a new custom function to generate random Linux PCI device paths.
	gofakeit.AddFuncLookup(
		"linuxdevicepathpci", gofakeit.Info{
			Category:    "custom",
			Description: "Random Linux PCI device path",
			Example:     "0000:00:02.0",
			Output:      "string",
			Generate: func(f *gofakeit.Faker, _ *gofakeit.MapParams, _ *gofakeit.Info) (any, error) {
				return f.Regex(types.HardwareMappingPathPCIValueRegEx.String()), nil
			},
		},
	)
	// Register a new custom function to generate random Linux USB device paths.
	gofakeit.AddFuncLookup(
		"linuxdevicepathusb", gofakeit.Info{
			Category:    "custom",
			Description: "Random Linux USB device path",
			Example:     "1-5.2",
			Output:      "string",
			Generate: func(f *gofakeit.Faker, _ *gofakeit.MapParams, _ *gofakeit.Info) (any, error) {
				return f.Regex(types.HardwareMappingPathUSBValueRegEx.String()), nil
			},
		},
	)

	te := initTestEnvironment(t)

	var data accTestHardwareMappingFakeData

	if err := gofakeit.Struct(&data); err != nil {
		t.Fatalf("could not create fake data for hardware mapping: %s", err)
	}

	return &data, te
}

// TestAccResourceHardwareMappingPCIValidInput runs tests for PCI hardware mapping resource definitions with valid input
// where all possible attributes are
// specified.
// All implementations of the [github.com/hashicorp/terraform-plugin-framework/resource.Resource] interface are tested
// in sequential steps.
func TestAccResourceHardwareMappingPCIValidInput(t *testing.T) {
	data, te := testAccResourceHardwareMappingInit(t)

	resource.Test(
		t, resource.TestCase{
			ProtoV6ProviderFactories: te.accProviders,
			Steps: []resource.TestStep{
				// Test the "Create" and "Read" implementations where all possible attributes are specified.
				{
					Config: fmt.Sprintf(
						`
					resource "proxmox_virtual_environment_hardware_mapping_pci" "test" {
						comment = "%s"
						name    = "%s"
						map     = [
							{
								comment      = "%s"
								id           = "%s"
								iommu_group  = %d
								node         = "%s"
								path         = "%s"
								subsystem_id = "%s"
							},
						]
						mediated_devices = %t
					}
					`,
						data.Comments[0],
						data.Names[0],
						data.MapComments[0],
						data.MapDeviceIDs[0],
						data.MapIOMMUGroups[0],
						te.nodeName,
						data.MapPathsPCI[0],
						data.MapSubsystemIDs[0],
						data.MediatedDevices,
					),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(accTestHardwareMappingNamePCI, "comment", data.Comments[0]),
						resource.TestCheckResourceAttrSet(accTestHardwareMappingNamePCI, "id"),
						resource.TestCheckTypeSetElemNestedAttrs(
							accTestHardwareMappingNamePCI, "map.*", map[string]string{
								"comment":      data.MapComments[0],
								"id":           data.MapDeviceIDs[0],
								"iommu_group":  strconv.Itoa(int(data.MapIOMMUGroups[0])),
								"node":         te.nodeName,
								"path":         data.MapPathsPCI[0],
								"subsystem_id": data.MapSubsystemIDs[0],
							},
						),
						resource.TestCheckResourceAttr(
							accTestHardwareMappingNamePCI,
							"mediated_devices",
							strconv.FormatBool(data.MediatedDevices),
						),
						resource.TestCheckResourceAttr(accTestHardwareMappingNamePCI, "name", data.Names[0]),
					),
				},

				// Test the "ImportState" implementation.
				{
					ImportState:       true,
					ImportStateId:     data.Names[0],
					ImportStateVerify: true,
					ResourceName:      accTestHardwareMappingNamePCI,
				},

				// Test the "Update" implementation where all possible attributes are specified.
				{
					Config: fmt.Sprintf(
						`
					resource "proxmox_virtual_environment_hardware_mapping_pci" "test" {
						comment = "%s"
						name    = "%s"
						map     = [
							{
								comment      = "%s"
								id           = "%s"
								iommu_group  = %d
								node         = "%s"
								path         = "%s"
								subsystem_id = "%s"
							},
						]
						mediated_devices = %t
					}
					`,
						data.Comments[1],
						data.Names[0],
						data.MapComments[1],
						data.MapDeviceIDs[0],
						data.MapIOMMUGroups[1],
						te.nodeName,
						data.MapPathsPCI[1],
						data.MapSubsystemIDs[1],
						!data.MediatedDevices,
					),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(accTestHardwareMappingNamePCI, "comment", data.Comments[1]),
						resource.TestCheckResourceAttrSet(accTestHardwareMappingNamePCI, "id"),
						resource.TestCheckTypeSetElemNestedAttrs(
							accTestHardwareMappingNamePCI, "map.*", map[string]string{
								"comment":      data.MapComments[1],
								"id":           data.MapDeviceIDs[0],
								"iommu_group":  strconv.Itoa(int(data.MapIOMMUGroups[1])),
								"node":         te.nodeName,
								"path":         data.MapPathsPCI[1],
								"subsystem_id": data.MapSubsystemIDs[1],
							},
						),
						resource.TestCheckResourceAttr(
							accTestHardwareMappingNamePCI,
							"mediated_devices",
							strconv.FormatBool(!data.MediatedDevices),
						),
						resource.TestCheckResourceAttr(accTestHardwareMappingNamePCI, "name", data.Names[0]),
					),
				},
			},
		},
	)
}

// TestAccResourceHardwareMappingPCIValidInputMinimal runs tests for PCI hardware mapping resource definitions with
// valid input that only have the minimum
// amount of attributes set to test computed and default values within the resulting plan and state. The last step sets
// the undefined values to test the update
// logic.
// All implementations of the [github.com/hashicorp/terraform-plugin-framework/resource.Resource] interface are tested
// in sequential steps.
func TestAccResourceHardwareMappingPCIValidInputMinimal(t *testing.T) {
	data, te := testAccResourceHardwareMappingInit(t)

	resource.Test(
		t, resource.TestCase{
			ProtoV6ProviderFactories: te.accProviders,
			Steps: []resource.TestStep{
				// Test the "Create" and "Read" implementations with only the minimum amount of attributes being set.
				{
					Config: fmt.Sprintf(
						`
					resource "proxmox_virtual_environment_hardware_mapping_pci" "test" {
						name    = "%s"
						map     = [
							{
								id   = "%s"
								node = "%s"
								path = "%s"
							},
						]
					}
					`,
						data.Names[0],
						data.MapDeviceIDs[0],
						te.nodeName,
						data.MapPathsPCI[0],
					),
					ConfigStateChecks: []statecheck.StateCheck{
						// Optional attributes should all be unset.
						statecheck.ExpectKnownValue(
							accTestHardwareMappingNamePCI,
							tfjsonpath.New("map").AtSliceIndex(0),
							knownvalue.MapPartial(
								map[string]knownvalue.Check{
									"comment":      knownvalue.Null(),
									"iommu_group":  knownvalue.Null(),
									"subsystem_id": knownvalue.Null(),
								},
							),
						),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet(accTestHardwareMappingNamePCI, "id"),
						resource.TestCheckTypeSetElemNestedAttrs(
							accTestHardwareMappingNamePCI, "map.*", map[string]string{
								"id":   data.MapDeviceIDs[0],
								"node": te.nodeName,
								"path": data.MapPathsPCI[0],
							},
						),
						resource.TestCheckResourceAttr(accTestHardwareMappingNamePCI, "name", data.Names[0]),
					),
				},

				// Test the "ImportState" implementation.
				{
					ImportState:       true,
					ImportStateId:     data.Names[0],
					ImportStateVerify: true,
					ResourceName:      accTestHardwareMappingNamePCI,
				},

				// Test the "Update" implementation by setting all previously undefined attributes.
				{
					Config: fmt.Sprintf(
						`
					resource "proxmox_virtual_environment_hardware_mapping_pci" "test" {
						comment = "%s"
						name    = "%s"
						map     = [
							{
								comment      = "%s"
								id           = "%s"
								iommu_group  = %d
								node         = "%s"
								path         = "%s"
								subsystem_id = "%s"
							},
						]
						mediated_devices = %t
					}
					`,
						data.Comments[1],
						data.Names[0],
						data.MapComments[1],
						data.MapDeviceIDs[0],
						data.MapIOMMUGroups[1],
						te.nodeName,
						data.MapPathsPCI[1],
						data.MapSubsystemIDs[1],
						!data.MediatedDevices,
					),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(accTestHardwareMappingNamePCI, "comment", data.Comments[1]),
						resource.TestCheckResourceAttrSet(accTestHardwareMappingNamePCI, "id"),
						resource.TestCheckTypeSetElemNestedAttrs(
							accTestHardwareMappingNamePCI, "map.*", map[string]string{
								"comment":      data.MapComments[1],
								"id":           data.MapDeviceIDs[0],
								"iommu_group":  strconv.Itoa(int(data.MapIOMMUGroups[1])),
								"node":         te.nodeName,
								"path":         data.MapPathsPCI[1],
								"subsystem_id": data.MapSubsystemIDs[1],
							},
						),
						resource.TestCheckResourceAttr(
							accTestHardwareMappingNamePCI,
							"mediated_devices",
							strconv.FormatBool(!data.MediatedDevices),
						),
						resource.TestCheckResourceAttr(accTestHardwareMappingNamePCI, "name", data.Names[0]),
					),
				},
			},
		},
	)
}

// TestAccResourceHardwareMappingPCIInvalidInput runs tests for PCI hardware mapping resource definitions with invalid
// input where all possible attributes are
// specified.
// Only the "Create" method implementation of the [github.com/hashicorp/terraform-plugin-framework/resource.Resource]
// interface is tested in sequential steps.
func TestAccResourceHardwareMappingPCIInvalidInput(t *testing.T) {
	data, te := testAccResourceHardwareMappingInit(t)

	resource.Test(
		t, resource.TestCase{
			ProtoV6ProviderFactories: te.accProviders,
			Steps: []resource.TestStep{
				// Test the "Create" method implementation where all possible attributes are specified, but an error is expected
				// when using an invalid device path.
				{
					ExpectError: regexp.MustCompile(
						fmt.Sprintf(
							// The error line is, for whatever reason, broken down into multiple lines in acceptance tests, so we need
							// to capture newline characters.
							// Note that the regular expression syntax used by Go does not capture newlines with the "." matcher,
							// so we need to enable the "s" flag that enabled "."
							// to match "\n".
							// References:
							//   1. https://pkg.go.dev/regexp/syntax
							`(?s).*%s(?s).*`,
							fwprovider.HardwareMappingResourceErrMessageInvalidPath(proxmoxtypes.HardwareMappingTypePCI),
						),
					),
					Config: fmt.Sprintf(
						`
					resource "proxmox_virtual_environment_hardware_mapping_pci" "test" {
						comment = "%s"
						name    = "%s"
						map     = [
							{
								comment      = "%s"
								id           = "%s"
								iommu_group  = %d
								node         = "%s"
								# Only valid Linux PCI device paths should pass the verification.
								path         = "wxyz:1337"
								subsystem_id = "%s"
							},
						]
						mediated_devices = %t
					}
					`,
						data.Comments[0],
						data.Names[0],
						data.Comments[1],
						data.MapDeviceIDs[0],
						data.MapIOMMUGroups[0],
						te.nodeName,
						data.MapSubsystemIDs[0],
						data.MediatedDevices,
					),
				},
			},
		},
	)

	resource.Test(
		t, resource.TestCase{
			ProtoV6ProviderFactories: te.accProviders,
			Steps: []resource.TestStep{
				// Test the "Create" method implementation where all possible attributes are specified, but an error is expected
				// when using an invalid device subsystem
				// ID.
				{
					ExpectError: regexp.MustCompile(fmt.Sprintf(`.*%s.*`, validators.HardwareMappingDeviceIDValidatorErrMessage)),
					Config: fmt.Sprintf(
						`
					resource "proxmox_virtual_environment_hardware_mapping_pci" "test" {
						comment = "%s"
						name    = "%s"
						map     = [
							{
								comment      = "%s"
								id           = "%s"
								iommu_group  = %d
								node         = "%s"
								path         = "%s"
								# Only valid Linux device subsystem IDs should pass the verification.
								subsystem_id = "x1y2:1337"
							},
						]
						mediated_devices = %t
					}
					`,
						data.Comments[0],
						data.Names[0],
						data.Comments[1],
						data.MapDeviceIDs[0],
						data.MapIOMMUGroups[0],
						te.nodeName,
						data.MapPathsPCI[0],
						data.MediatedDevices,
					),
				},
			},
		},
	)
}

// TestAccResourceHardwareMappingUSBValidInput runs tests for USB hardware mapping resource definitions with valid input
// where all possible attributes are
// specified.
// All implementations of the [github.com/hashicorp/terraform-plugin-framework/resource.Resource] interface are tested
// in sequential steps.
func TestAccResourceHardwareMappingUSBValidInput(t *testing.T) {
	data, te := testAccResourceHardwareMappingInit(t)

	resource.Test(
		t, resource.TestCase{
			ProtoV6ProviderFactories: te.accProviders,
			Steps: []resource.TestStep{
				// Test the "Create" and "Read" implementations where all possible attributes are specified.
				{
					Config: fmt.Sprintf(
						`
					resource "proxmox_virtual_environment_hardware_mapping_usb" "test" {
						comment = "%s"
						name    = "%s"
						map     = [
							{
								comment = "%s"
								id      = "%s"
								node    = "%s"
								path    = "%s"
							},
						]
					}
					`,
						data.Comments[0],
						data.Names[0],
						data.MapComments[0],
						data.MapDeviceIDs[0],
						te.nodeName,
						data.MapPathsUSB[0],
					),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(accTestHardwareMappingNameUSB, "comment", data.Comments[0]),
						resource.TestCheckResourceAttrSet(accTestHardwareMappingNameUSB, "id"),
						resource.TestCheckTypeSetElemNestedAttrs(
							accTestHardwareMappingNameUSB, "map.*", map[string]string{
								"comment": data.MapComments[0],
								"id":      data.MapDeviceIDs[0],
								"node":    te.nodeName,
								"path":    data.MapPathsUSB[0],
							},
						),
						resource.TestCheckResourceAttr(accTestHardwareMappingNameUSB, "name", data.Names[0]),
					),
				},

				// Test the "ImportState" implementation and ensure that PCI-only attributes are not set.
				{
					ImportState:       true,
					ImportStateId:     data.Names[0],
					ImportStateVerify: true,
					ResourceName:      accTestHardwareMappingNameUSB,
				},

				// Test the "Update" implementation where all possible attributes are specified.
				{
					Config: fmt.Sprintf(
						`
					resource "proxmox_virtual_environment_hardware_mapping_usb" "test" {
						comment = "%s"
						name    = "%s"
						map     = [
							{
								comment = "%s"
								id      = "%s"
								node    = "%s"
								path    = "%s"
							},
						]
					}
					`,
						data.Comments[1],
						data.Names[0],
						data.MapComments[1],
						data.MapDeviceIDs[0],
						te.nodeName,
						data.MapPathsUSB[1],
					),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(accTestHardwareMappingNameUSB, "comment", data.Comments[1]),
						resource.TestCheckResourceAttrSet(accTestHardwareMappingNameUSB, "id"),
						resource.TestCheckTypeSetElemNestedAttrs(
							accTestHardwareMappingNameUSB, "map.*", map[string]string{
								"comment": data.MapComments[1],
								"id":      data.MapDeviceIDs[0],
								"node":    te.nodeName,
								"path":    data.MapPathsUSB[1],
							},
						),
						resource.TestCheckResourceAttr(accTestHardwareMappingNameUSB, "name", data.Names[0]),
					),
				},
			},
		},
	)
}

// TestAccResourceHardwareMappingUSBValidInputMinimal runs tests for USB hardware mapping resource definitions with
// valid input that only have the minimum
// amount of attributes set to test computed and default values within the resulting plan and state. The last step sets
// the undefined values to test the update
// logic.
// All implementations of the [github.com/hashicorp/terraform-plugin-framework/resource.Resource] interface are tested
// in sequential steps.
func TestAccResourceHardwareMappingUSBValidInputMinimal(t *testing.T) {
	data, te := testAccResourceHardwareMappingInit(t)

	resource.Test(
		t, resource.TestCase{
			ProtoV6ProviderFactories: te.accProviders,
			Steps: []resource.TestStep{
				// Test the "Create" and "Read" implementations with only the minimum amount of attributes being set.
				{
					Config: fmt.Sprintf(
						`
					resource "proxmox_virtual_environment_hardware_mapping_usb" "test" {
						name    = "%s"
						map     = [
							{
								id   = "%s"
								node = "%s"
							},
						]
					}
					`,
						data.Names[0],
						data.MapDeviceIDs[0],
						te.nodeName,
					),
					ConfigStateChecks: []statecheck.StateCheck{
						// Optional attributes should all be unset.
						statecheck.ExpectKnownValue(
							accTestHardwareMappingNameUSB,
							tfjsonpath.New("map").AtSliceIndex(0),
							knownvalue.MapPartial(
								map[string]knownvalue.Check{
									"path": knownvalue.Null(),
								},
							),
						),
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet(accTestHardwareMappingNameUSB, "id"),
						resource.TestCheckTypeSetElemNestedAttrs(
							accTestHardwareMappingNameUSB, "map.*", map[string]string{
								"id":   data.MapDeviceIDs[0],
								"node": te.nodeName,
							},
						),
						resource.TestCheckResourceAttr(accTestHardwareMappingNameUSB, "name", data.Names[0]),
					),
				},

				// Test the "Update" implementation by setting all previously undefined attributes.
				{
					Config: fmt.Sprintf(
						`
					resource "proxmox_virtual_environment_hardware_mapping_usb" "test" {
						comment = "%s"
						name    = "%s"
						map     = [
							{
								comment = "%s"
								id      = "%s"
								node    = "%s"
								path    = "%s"
							},
						]
					}
					`,
						data.Comments[0],
						data.Names[0],
						data.Comments[1],
						data.MapDeviceIDs[1],
						te.nodeName,
						data.MapPathsUSB[0],
					),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(accTestHardwareMappingNameUSB, "comment", data.Comments[0]),
						resource.TestCheckResourceAttrSet(accTestHardwareMappingNameUSB, "id"),
						resource.TestCheckTypeSetElemNestedAttrs(
							accTestHardwareMappingNameUSB, "map.*", map[string]string{
								"comment": data.Comments[1],
								"id":      data.MapDeviceIDs[1],
								"node":    te.nodeName,
								"path":    data.MapPathsUSB[0],
							},
						),
						resource.TestCheckResourceAttr(accTestHardwareMappingNameUSB, "name", data.Names[0]),
					),
				},
			},
		},
	)
}

// TestAccResourceHardwareMappingUSBInvalidInput runs tests for USB hardware mapping resource definitions where all
// possible attributes are specified.
// Only the "Create" method implementation of the [github.com/hashicorp/terraform-plugin-framework/resource.Resource]
// interface is tested in sequential steps.
func TestAccResourceHardwareMappingUSBInvalidInput(t *testing.T) {
	data, te := testAccResourceHardwareMappingInit(t)

	resource.Test(
		t, resource.TestCase{
			ProtoV6ProviderFactories: te.accProviders,
			Steps: []resource.TestStep{
				// Test the "Create" method implementation where all possible attributes are specified, but an error is expected
				// when using an invalid device path.
				{
					ExpectError: regexp.MustCompile(
						fmt.Sprintf(
							// The error line is, for whatever reason, broken down into multiple lines in acceptance tests, so we need
							// to capture newline characters.
							// Note that the regular expression syntax used by Go does not capture newlines with the "." matcher,
							// so we need to enable the "s" flag that enabled "."
							// to match "\n".
							// References:
							//   1. https://pkg.go.dev/regexp/syntax
							`(?s).*%s(?s).*`,
							fwprovider.HardwareMappingResourceErrMessageInvalidPath(proxmoxtypes.HardwareMappingTypeUSB),
						),
					),
					Config: fmt.Sprintf(
						`
					resource "proxmox_virtual_environment_hardware_mapping_usb" "test" {
						comment = "%s"
						name    = "%s"
						map     = [
							{
								comment = "%s"
								id      = "%s"
								node    = "%s"
								# Only valid Linux USB device paths should pass the verification.
								path    = "xyz3:1337foobar"
							},
						]
					}
					`,
						data.Comments[0],
						data.Names[0],
						data.Comments[1],
						data.MapDeviceIDs[0],
						te.nodeName,
					),
				},
			},
		},
	)
}
