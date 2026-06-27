//go:build acceptance || all

//testacc:tier=light
//testacc:resource=metrics

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package metrics_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/bpg/terraform-provider-proxmox/fwprovider/test"
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
				resource "proxmox_metrics_server" "acc_influxdb_server" {
					name   = "acc_example_influxdb_server"
					server = "192.168.3.2"
					port   = 18089
					type   = "influxdb"
					mtu    = 1000
				  }`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_metrics_server.acc_influxdb_server", map[string]string{
						"id":      "acc_example_influxdb_server",
						"name":    "acc_example_influxdb_server",
						"mtu":     "1000",
						"port":    "18089",
						"server":  "192.168.3.2",
						"type":    "influxdb",
						"disable": "false",
					}),
					test.NoResourceAttributesSet("proxmox_metrics_server.acc_influxdb_server", []string{
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
						"opentelemetry_proto",
						"opentelemetry_path",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_metrics_server" "acc_influxdb_server" {
					name   			 = "acc_example_influxdb_server"
					server 			 = "192.168.3.2"
					port   			 = 18089
					type   			 = "influxdb"
					mtu    			 = 1000
					influx_bucket    = "xxxxx"
				  }`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_metrics_server.acc_influxdb_server", map[string]string{
						"id":            "acc_example_influxdb_server",
						"name":          "acc_example_influxdb_server",
						"mtu":           "1000",
						"port":          "18089",
						"server":        "192.168.3.2",
						"type":          "influxdb",
						"influx_bucket": "xxxxx",
						"disable":       "false",
					}),
					test.NoResourceAttributesSet("proxmox_metrics_server.acc_influxdb_server", []string{
						"timeout",
						"influx_api_path_prefix",
						"influx_db_proto",
						"influx_max_body_size",
						"influx_organization",
						"influx_token",
						"influx_verify",
						"graphite_path",
						"graphite_proto",
						"opentelemetry_proto",
						"opentelemetry_path",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_metrics_server" "acc_influxdb_server" {
					name   			 = "acc_example_influxdb_server"
					server 			 = "192.168.3.2"
					port   			 = 18089
					type   			 = "influxdb"
					influx_bucket    = "xxxxx"
				  }`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_metrics_server.acc_influxdb_server", map[string]string{
						"id":            "acc_example_influxdb_server",
						"name":          "acc_example_influxdb_server",
						"port":          "18089",
						"server":        "192.168.3.2",
						"type":          "influxdb",
						"influx_bucket": "xxxxx",
						"disable":       "false",
					}),
					test.NoResourceAttributesSet("proxmox_metrics_server.acc_influxdb_server", []string{
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
						"opentelemetry_proto",
						"opentelemetry_path",
					}),
				),
			},
		}},
		{"create disabled influxdb server with verify toggle & round-trip bools", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_metrics_server" "acc_influxdb_bools" {
					name           = "acc_example_influxdb_bools"
					server         = "192.168.3.2"
					port           = 18090
					type           = "influxdb"
					disable        = true
					influx_verify  = false
				  }`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_metrics_server.acc_influxdb_bools", map[string]string{
						"disable":       "true",
						"influx_verify": "false",
						"port":          "18090",
						"type":          "influxdb",
					}),
				),
			},
			{
				Config: te.RenderConfig(`
				resource "proxmox_metrics_server" "acc_influxdb_bools" {
					name           = "acc_example_influxdb_bools"
					server         = "192.168.3.2"
					port           = 18090
					type           = "influxdb"
					disable        = false
					influx_verify  = true
				  }`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("proxmox_metrics_server.acc_influxdb_bools", map[string]string{
						"disable":       "false",
						"influx_verify": "true",
					}),
				),
			},
			{
				ResourceName:      "proxmox_metrics_server.acc_influxdb_bools",
				ImportState:       true,
				ImportStateVerify: true,
			},
		}},
		{"create graphite udp metrics server & import it", []resource.TestStep{
			{
				ResourceName: "proxmox_metrics_server.acc_graphite_server",
				Config: te.RenderConfig(`
				resource "proxmox_metrics_server" "acc_graphite_server" {
					name   = "acc_example_graphite_server"
					server = "192.168.3.2"
					port   = 18089
					type   = "graphite"
				  }`),
			},
			{
				ResourceName:      "proxmox_metrics_server.acc_graphite_server",
				ImportState:       true,
				ImportStateVerify: true,
			},
		}},
		{"create graphite udp metrics server & test datasource", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_metrics_server" "acc_graphite_server2" {
					name   = "acc_example_graphite_server2"
					server = "192.168.3.2"
					port   = 18089
					type   = "graphite"
				  }
				data "proxmox_metrics_server" "acc_graphite_server2" {
					name = proxmox_metrics_server.acc_graphite_server2.name
				  }`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("data.proxmox_metrics_server.acc_graphite_server2", map[string]string{
						"id":     "acc_example_graphite_server2",
						"name":   "acc_example_graphite_server2",
						"port":   "18089",
						"server": "192.168.3.2",
						"type":   "graphite",
					}),
				),
			},
		}},
		{"create opentelemetry metrics server & import it", []resource.TestStep{
			{
				// Skip this test until we have a way to test opentelemetry servers (i.e. setting up local otel collector)
				// Proxmox is trying to connect to the server when creating the resource.
				SkipFunc: func() (bool, error) {
					return true, nil
				},
				ResourceName: "proxmox_metrics_server.acc_otel_server",
				Config: te.RenderConfig(`
				resource "proxmox_metrics_server" "acc_otel_server" {
					name   = "acc_example_otel_server"
					server = "192.168.3.2"
					port   = 4318
					type   = "opentelemetry"
					opentelemetry_proto = "http"
				}`),
			},
			{
				// Skip this test until we have a way to test opentelemetry servers (i.e. setting up local otel collector)
				// Proxmox is trying to connect to the server when creating the resource.
				SkipFunc: func() (bool, error) {
					return true, nil
				},
				ResourceName:      "proxmox_metrics_server.acc_otel_server",
				ImportState:       true,
				ImportStateVerify: true,
			},
		}},
		{"influx_token survives create and refresh (regression)", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_metrics_server" "acc_influxdb_token" {
					name         = "acc_influxdb_token"
					server       = "192.168.3.2"
					port         = 18091
					type         = "influxdb"
					influx_token = "supersecrettoken"
				  }`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_metrics_server.acc_influxdb_token", "id", "acc_influxdb_token"),
					resource.TestCheckResourceAttrSet("proxmox_metrics_server.acc_influxdb_token", "influx_token"),
				),
			},
			// No-op refresh: influx_token must survive the round-trip through Read.
			{
				Config: te.RenderConfig(`
				resource "proxmox_metrics_server" "acc_influxdb_token" {
					name         = "acc_influxdb_token"
					server       = "192.168.3.2"
					port         = 18091
					type         = "influxdb"
					influx_token = "supersecrettoken"
				  }`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("proxmox_metrics_server.acc_influxdb_token", "influx_token"),
				),
			},
		}},
		{"create opentelemetry metrics server & test datasource", []resource.TestStep{
			{
				// Skip this test until we have a way to test opentelemetry servers (i.e. setting up local otel collector)
				// Proxmox is trying to connect to the server when creating the resource.
				SkipFunc: func() (bool, error) {
					return true, nil
				},
				Config: te.RenderConfig(`
				resource "proxmox_metrics_server" "acc_otel_server2" {
					name   = "acc_example_otel_server2"
					server = "192.168.3.2"
					port   = 4318
					type   = "opentelemetry"
					opentelemetry_proto = "https"
					opentelemetry_path  = "/v1/metrics"
				}
				data "proxmox_metrics_server" "acc_otel_server2" {
					name = proxmox_metrics_server.acc_otel_server2.name
				}`),
				Check: resource.ComposeTestCheckFunc(
					test.ResourceAttributes("data.proxmox_metrics_server.acc_otel_server2", map[string]string{
						"id":                  "acc_example_otel_server2",
						"name":                "acc_example_otel_server2",
						"port":                "4318",
						"server":              "192.168.3.2",
						"type":                "opentelemetry",
						"opentelemetry_proto": "https",
						"opentelemetry_path":  "/v1/metrics",
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

func TestAccResourceMetricsServerWriteOnly(t *testing.T) {
	t.Parallel()

	te := test.InitEnvironment(t)

	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{"influx_token_wo keeps the secret out of state", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_metrics_server" "acc_influxdb_token_wo" {
					name                    = "acc_influxdb_token_wo"
					server                  = "192.168.3.2"
					port                    = 18092
					type                    = "influxdb"
					influx_token_wo         = "supersecrettoken"
					influx_token_wo_version = 1
				  }`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_metrics_server.acc_influxdb_token_wo", "id", "acc_influxdb_token_wo"),
					resource.TestCheckResourceAttr("proxmox_metrics_server.acc_influxdb_token_wo", "influx_token_wo_version", "1"),
					// The write-only secret must never be persisted, and the legacy
					// influx_token mirror must stay unset when influx_token_wo is used.
					resource.TestCheckNoResourceAttr("proxmox_metrics_server.acc_influxdb_token_wo", "influx_token_wo"),
					resource.TestCheckNoResourceAttr("proxmox_metrics_server.acc_influxdb_token_wo", "influx_token"),
				),
			},
		}},
		{"influx_token_wo_version triggers rotation", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_metrics_server" "acc_influxdb_token_wo_rotate" {
					name                    = "acc_influxdb_token_wo_rotate"
					server                  = "192.168.3.2"
					port                    = 18093
					type                    = "influxdb"
					influx_token_wo         = "supersecrettoken"
					influx_token_wo_version = 1
				  }`),
				Check: resource.TestCheckResourceAttr("proxmox_metrics_server.acc_influxdb_token_wo_rotate", "influx_token_wo_version", "1"),
			},
			// Rotate the token by bumping the version counter.
			{
				Config: te.RenderConfig(`
				resource "proxmox_metrics_server" "acc_influxdb_token_wo_rotate" {
					name                    = "acc_influxdb_token_wo_rotate"
					server                  = "192.168.3.2"
					port                    = 18093
					type                    = "influxdb"
					influx_token_wo         = "rotatedsecrettoken"
					influx_token_wo_version = 2
				  }`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_metrics_server.acc_influxdb_token_wo_rotate", "influx_token_wo_version", "2"),
					resource.TestCheckNoResourceAttr("proxmox_metrics_server.acc_influxdb_token_wo_rotate", "influx_token"),
				),
			},
		}},
		{"influx_token_wo_version requires influx_token_wo", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_metrics_server" "acc_influxdb_token_wo_novalue" {
					name                    = "acc_influxdb_token_wo_novalue"
					server                  = "192.168.3.2"
					port                    = 18094
					type                    = "influxdb"
					influx_token_wo_version = 1
				  }`),
				ExpectError: regexp.MustCompile(`Attribute "influx_token_wo" must be specified`),
			},
		}},
		{"influx_token and influx_token_wo are mutually exclusive", []resource.TestStep{
			{
				Config: te.RenderConfig(`
				resource "proxmox_metrics_server" "acc_influxdb_token_conflict" {
					name            = "acc_influxdb_token_conflict"
					server          = "192.168.3.2"
					port            = 18095
					type            = "influxdb"
					influx_token    = "plaintext-token"
					influx_token_wo = "writeonly-token"
				  }`),
				ExpectError: regexp.MustCompile(`These attributes cannot be configured together`),
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.SkipBelow(tfversion.Version1_11_0),
				},
				ProtoV6ProviderFactories: te.AccProviders,
				Steps:                    tt.steps,
			})
		})
	}
}
