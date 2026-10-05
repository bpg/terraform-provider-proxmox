//go:build acceptance || all

//testacc:tier=heavy
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
	"regexp"
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

	proxy := newSharedSnippetStorageProxy(t, te.NodeName, te.DatastoreID, missingNode)

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

	config := te.RenderConfig(`
					resource "proxmox_virtual_environment_file" "test_removed_node" {
						content_type = "snippets"
						datastore_id = "{{.DatastoreID}}"
						node_name    = "{{.NodeName}}"
						source_raw {
							data      = "#cloud-config\nhostname: removed-node-test\n"
							file_name = "{{.SnippetFileName}}"
						}
					}`, WithInsecureEndpoint(proxy.endpoint))
	removedNodeConfig := strings.Replace(config, `node_name    = "`+te.NodeName+`"`, `node_name    = "`+missingNode+`"`, 1)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: te.AccProviders,
		CheckDestroy: func(*terraform.State) error {
			contentType := "snippets"

			files, err := nodeStorage.ListDatastoreFiles(ctx, &contentType)
			if err != nil {
				return err
			}

			for _, file := range files {
				if file.VolumeID == volumeID {
					return fmt.Errorf("snippet %q still exists after destroy", volumeID)
				}
			}

			if proxy.fallbackDeletes.Load() == 0 {
				return fmt.Errorf("deletion did not use the surviving node")
			}

			return nil
		},
		Steps: []resource.TestStep{
			{
				Config:      removedNodeConfig,
				ExpectError: regexp.MustCompile("error listing files from datastore.*HTTP 500"),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", volumeID),
					func(*terraform.State) error {
						if proxy.sharedResponses.Load() == 0 {
							return fmt.Errorf("provider did not receive simulated shared-storage metadata")
						}

						return checkFileExists()
					},
				),
			},

			{
				PreConfig: func() { proxy.removed.Store(true) },
				Config:    config,
				PlanOnly:  true,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "node_name", te.NodeName),
					resource.TestCheckResourceAttr(resourceName, "id", volumeID),
					func(*terraform.State) error {
						if proxy.fallbackReads.Load() == 0 {
							return fmt.Errorf("refresh did not use the surviving node")
						}

						return checkFileExists()
					},
				),
			},
			{
				ResourceName:  resourceName,
				ImportState:   true,
				ImportStateId: te.NodeName + "/" + volumeID,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 || states[0].ID != volumeID || states[0].Attributes["node_name"] != te.NodeName {
						return fmt.Errorf("import must retain the file ID and original node")
					}

					return checkFileExists()
				},
			},
		},
	})
}

type sharedSnippetStorageProxy struct {
	endpoint        string
	removed         atomic.Bool
	sharedResponses atomic.Int64
	fallbackReads   atomic.Int64
	fallbackDeletes atomic.Int64
}

// newSharedSnippetStorageProxy simulates shared storage and node removal while forwarding file operations to Proxmox.
func newSharedSnippetStorageProxy(t *testing.T, nodeName, datastoreID, missingNode string) *sharedSnippetStorageProxy {
	t.Helper()

	target, err := url.Parse(utils.GetAnyStringEnv("PROXMOX_VE_ENDPOINT"))
	require.NoError(t, err)

	state := &sharedSnippetStorageProxy{}
	const survivingNode = "test-surviving-node"

	transport := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	t.Cleanup(transport.CloseIdleConnections)

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			r.Out.Header.Del("Accept-Encoding")

			if strings.HasPrefix(r.In.URL.Path, "/api2/json/nodes/"+survivingNode+"/") {
				r.Out.URL.Path = strings.Replace(r.Out.URL.Path, "/nodes/"+survivingNode+"/", "/nodes/"+nodeName+"/", 1)
				r.Out.URL.Path = strings.Replace(r.Out.URL.Path, "UPID:"+survivingNode+":", "UPID:"+nodeName+":", 1)

				r.Out.URL.RawPath = ""
				if strings.Contains(r.In.URL.Path, "/content") {
					switch r.In.Method {
					case http.MethodGet:
						state.fallbackReads.Add(1)
					case http.MethodDelete:
						state.fallbackDeletes.Add(1)
					}
				}
			}
		},
		Transport: transport,
		ModifyResponse: func(resp *http.Response) error {
			path := strings.TrimPrefix(resp.Request.URL.Path, "/api2/json")
			if state.removed.Load() && resp.Request.Method == http.MethodDelete && resp.StatusCode == http.StatusOK {
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return fmt.Errorf("read deletion task: %w", err)
				}

				body = bytes.ReplaceAll(body, []byte("UPID:"+nodeName+":"), []byte("UPID:"+survivingNode+":"))
				resp.Body = io.NopCloser(bytes.NewReader(body))
				resp.ContentLength = int64(len(body))
				resp.Header.Set("Content-Length", strconv.Itoa(len(body)))

				return nil
			}

			if resp.Request.Method != http.MethodGet || resp.StatusCode != http.StatusOK {
				return nil
			}

			switch path {
			case "/nodes", "/storage", "/storage/" + datastoreID, "/nodes/" + nodeName + "/storage", "/cluster/resources":
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
				if state.removed.Load() && path == "/nodes" {
					var name string
					if err := json.Unmarshal(entry["node"], &name); err != nil {
						return fmt.Errorf("decode node name: %w", err)
					}

					if name == nodeName {
						entry["node"] = json.RawMessage(`"` + survivingNode + `"`)
					}
				}
				var id string
				if raw, ok := entry["storage"]; ok {
					if err := json.Unmarshal(raw, &id); err != nil {
						return fmt.Errorf("transform shared-storage metadata: %w", err)
					}
				}

				if entry != nil && (single || id == datastoreID) {
					entry["shared"] = json.RawMessage("1")

					state.sharedResponses.Add(1)
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

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api2/json")
		if strings.HasPrefix(path, "/nodes/"+missingNode+"/") ||
			(state.removed.Load() && strings.HasPrefix(path, "/nodes/"+nodeName+"/")) {
			http.Error(w, "node "+missingNode+" does not exist", http.StatusInternalServerError)
			return
		}

		if state.removed.Load() && path == "/storage/"+datastoreID {
			http.Error(w, "Datastore.Allocate required", http.StatusForbidden)
			return
		}

		proxy.ServeHTTP(w, r)
	}))
	t.Cleanup(server.Close)

	state.endpoint = server.URL

	return state
}
