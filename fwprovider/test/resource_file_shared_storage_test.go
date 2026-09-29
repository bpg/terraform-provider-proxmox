//go:build acceptance || all

//testacc:tier=light
//testacc:resource=file

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

// TestAccResourceFileReadAfterNodeRemoval uses real snippets with simulated shared-storage metadata.
func TestAccResourceFileReadAfterNodeRemoval(t *testing.T) {
	te := InitEnvironment(t)
	ctx := context.Background()
	nodeStorage := te.NodeStorageClient()
	datastores, err := nodeStorage.ListDatastores(ctx, &storage.DatastoreListRequestBody{ID: &te.DatastoreID})
	require.NoError(t, err)
	require.Len(t, datastores, 1, "test datastore must be accessible on the test node")

	datastore := datastores[0]
	require.Equal(t, te.DatastoreID, datastore.ID)

	require.NotNil(t, datastore.ContentTypes)
	require.Contains(t, []string(*datastore.ContentTypes), "snippets")

	missingNode := "test-removed-node"
	nodes, err := te.ClusterClient().GetClusterResources(ctx, "node")
	require.NoError(t, err)
	require.NotEmpty(t, nodes, "node discovery must be permitted for this test")

	for _, node := range nodes {
		require.NotEqual(t, missingNode, node.NodeName, "the stale node must not exist in the cluster")
	}

	endpoint, sharedResponses := newSharedSnippetStorageProxy(t, te.NodeName, te.DatastoreID)

	fileName := SafeResourceName("removed-node-snippet") + ".yaml"
	volumeID := fmt.Sprintf("%s:snippets/%s", te.DatastoreID, fileName)
	resourceName := "proxmox_virtual_environment_file.test_removed_node"

	checkFileExists := func() error {
		contentType := "snippets"

		files, err := nodeStorage.ListDatastoreFiles(ctx, &contentType)
		if err != nil {
			return err
		}

		for _, file := range files {
			if file.VolumeID == volumeID {
				return nil
			}
		}

		return fmt.Errorf("snippet %q is missing from shared storage on surviving node %q", volumeID, te.NodeName)
	}

	te.AddTemplateVars(map[string]any{
		"SnippetFileName": fileName,
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: te.RenderConfig(`
					resource "proxmox_virtual_environment_file" "test_removed_node" {
						content_type = "snippets"
						datastore_id = "{{.DatastoreID}}"
						node_name    = "{{.NodeName}}"
						source_raw {
							data      = "#cloud-config\nhostname: removed-node-test\n"
							file_name = "{{.SnippetFileName}}"
						}
					}`, WithInsecureEndpoint(endpoint)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", volumeID),
					func(*terraform.State) error {
						if sharedResponses.Load() == 0 {
							return fmt.Errorf("provider did not receive simulated shared-storage metadata")
						}

						return checkFileExists()
					},
				),
			},
			{
				ResourceName:  resourceName,
				ImportState:   true,
				ImportStateId: missingNode + "/" + volumeID,
				// Leave the original state intact so cleanup can use the surviving node
				ImportStatePersist: false,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected one refreshed snippet, got %d", len(states))
					}

					if states[0].ID != volumeID {
						return fmt.Errorf("refresh changed snippet ID: got %q, want %q", states[0].ID, volumeID)
					}

					if states[0].Attributes["file_name"] != fileName {
						return fmt.Errorf("refresh did not recover snippet file_name %q", fileName)
					}

					return checkFileExists()
				},
			},
		},
	})
}

// newSharedSnippetStorageProxy only changes the selected datastore's shared flag. Node discovery,
// missing-node errors, file contents, uploads and deletes all come from the real Proxmox instance.
func newSharedSnippetStorageProxy(t *testing.T, nodeName, datastoreID string) (string, *atomic.Int64) {
	t.Helper()

	target, err := url.Parse(utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT"))
	require.NoError(t, err)

	var sharedResponses atomic.Int64

	transport := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	t.Cleanup(transport.CloseIdleConnections)

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			r.Out.Header.Del("Accept-Encoding")
		},
		Transport: transport,
		ModifyResponse: func(resp *http.Response) error {
			path := strings.TrimPrefix(resp.Request.URL.Path, "/api2/json")
			if resp.Request.Method != http.MethodGet || resp.StatusCode != http.StatusOK {
				return nil
			}

			switch path {
			case "/storage", "/storage/" + datastoreID, "/nodes/" + nodeName + "/storage", "/cluster/resources":
			default:
				return nil
			}

			defer resp.Body.Close()

			var body map[string]json.RawMessage
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				return fmt.Errorf("decode storage metadata: %w", err)
			}

			var entries []map[string]json.RawMessage

			single := path == "/storage/"+datastoreID
			if single {
				var entry map[string]json.RawMessage
				if err := json.Unmarshal(body["data"], &entry); err != nil {
					return fmt.Errorf("transform shared-storage metadata: %w", err)
				}

				entries = append(entries, entry)
			} else if err := json.Unmarshal(body["data"], &entries); err != nil {
				return fmt.Errorf("transform shared-storage metadata: %w", err)
			}

			for _, entry := range entries {
				var id string
				if raw, ok := entry["storage"]; ok {
					if err := json.Unmarshal(raw, &id); err != nil {
						return fmt.Errorf("transform shared-storage metadata: %w", err)
					}
				}

				if entry != nil && (single || id == datastoreID) {
					entry["shared"] = json.RawMessage("1")

					sharedResponses.Add(1)
				}
			}

			var data any = entries
			if single {
				data = entries[0]
			}

			encodedData, err := json.Marshal(data)
			if err != nil {
				return fmt.Errorf("transform shared-storage metadata: %w", err)
			}

			body["data"] = encodedData

			encoded, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("transform shared-storage metadata: %w", err)
			}

			resp.Body = io.NopCloser(bytes.NewReader(encoded))
			resp.ContentLength = int64(len(encoded))
			resp.Header.Set("Content-Length", strconv.Itoa(len(encoded)))

			return nil
		},
	}

	server := httptest.NewTLSServer(proxy)
	t.Cleanup(server.Close)

	return server.URL, &sharedResponses
}
