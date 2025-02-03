//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package metrics_test

import (
	"testing"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceMetricsServer(t *testing.T) {
	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"create influxdb udp server & update it & again to default mtu", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_metrics_server" "acc_influxdb_server" {
					name   = "acc_example_influxdb_server"
					server = "192.168.3.2"
					port   = 18089
					type   = "influxdb"
					mtu    = 1000
				  }`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_metrics_server.acc_influxdb_server", map[string]string{
						"id":     "acc_example_influxdb_server",
						"name":   "acc_example_influxdb_server",
						"mtu":    "1000",
						"port":   "18089",
						"server": "192.168.3.2",
						"type":   "influxdb",
					}),
					test.NoResourceAttributesSet("proxmox_virtual_environment_metrics_server.acc_influxdb_server", []string{
						"disable",
						"timeout",
						"influx_api_path_prefix",
						"influx_bucket",
						"influx_db_proto",
						"influx_max_body_size",
						"influx_organization",
						"influx_token",
						"influx_verify",
						"graphite_path",
						"graphite_proto",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_metrics_server" "acc_influxdb_server" {
					name   			 = "acc_example_influxdb_server"
					server 			 = "192.168.3.2"
					port   			 = 18089
					type   			 = "influxdb"
					mtu    			 = 1000
					influx_bucket    = "xxxxx"
				  }`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_metrics_server.acc_influxdb_server", map[string]string{
						"id":            "acc_example_influxdb_server",
						"name":          "acc_example_influxdb_server",
						"mtu":           "1000",
						"port":          "18089",
						"server":        "192.168.3.2",
						"type":          "influxdb",
						"influx_bucket": "xxxxx",
					}),
					test.NoResourceAttributesSet("proxmox_virtual_environment_metrics_server.acc_influxdb_server", []string{
						"disable",
						"timeout",
						"influx_api_path_prefix",
						"influx_db_proto",
						"influx_max_body_size",
						"influx_organization",
						"influx_token",
						"influx_verify",
						"graphite_path",
						"graphite_proto",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_metrics_server" "acc_influxdb_server" {
					name   			 = "acc_example_influxdb_server"
					server 			 = "192.168.3.2"
					port   			 = 18089
					type   			 = "influxdb"
					influx_bucket    = "xxxxx"
				  }`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_virtual_environment_metrics_server.acc_influxdb_server", map[string]string{
						"id":            "acc_example_influxdb_server",
						"name":          "acc_example_influxdb_server",
						"port":          "18089",
						"server":        "192.168.3.2",
						"type":          "influxdb",
						"influx_bucket": "xxxxx",
					}),
					test.NoResourceAttributesSet("proxmox_virtual_environment_metrics_server.acc_influxdb_server", []string{
						"disable",
						"timeout",
						"mtu",
						"influx_api_path_prefix",
						"influx_db_proto",
						"influx_max_body_size",
						"influx_organization",
						"influx_token",
						"influx_verify",
						"graphite_path",
						"graphite_proto",
					}),
				),
			},
		}},
		{"create graphite udp metrics server & import it", []resource.TestStep{
			{
				ResourceName: "proxmox_virtual_environment_metrics_server.acc_graphite_server",
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_metrics_server" "acc_graphite_server" {
					name   = "acc_example_graphite_server"
					server = "192.168.3.2"
					port   = 18089
					type   = "graphite"
				  }`),
			},
			{
				ResourceName:      "proxmox_virtual_environment_metrics_server.acc_graphite_server",
				ImportState:       true,
				ImportStateVerify: true,
			},
		}},
		{"create graphite udp metrics server & test datasource", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_virtual_environment_metrics_server" "acc_graphite_server2" {
					name   = "acc_example_graphite_server2"
					server = "192.168.3.2"
					port   = 18089
					type   = "graphite"
				  }
				data "proxmox_virtual_environment_metrics_server" "acc_graphite_server2" {
					name = proxmox_virtual_environment_metrics_server.acc_graphite_server2.name
				  }`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("data.proxmox_virtual_environment_metrics_server.acc_graphite_server2", map[string]string{
						"id":     "acc_example_graphite_server2",
						"name":   "acc_example_graphite_server2",
						"port":   "18089",
						"server": "192.168.3.2",
						"type":   "graphite",
					}),
				),
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
