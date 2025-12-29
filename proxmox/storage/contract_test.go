package storage

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

type capturedRequest struct {
	Method   string
	Path     string
	BodyKeys []string
}

func TestClient_Contract_CreateAndUpdateStorage_RequestKeys(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string

		storageID   string
		createReq   any
		updateReq   any
		createKeys  []string
		updateKeys  []string
		updateNoKey []string
	}

	const (
		nodeName = "node-1"
	)

	nodes := types.CustomCommaSeparatedList{nodeName}
	contentImages := types.CustomCommaSeparatedList{"images"}

	disableFalse := types.CustomBool(false)
	sharedFalse := types.CustomBool(false)

	keepDaily := 7
	maxProtected := types.CustomInt64(5)
	backups := DataStoreWithBackups{
		MaxProtectedBackups: &maxProtected,
		KeepDaily:           &keepDaily,
	}

	cases := []testCase{
		{
			name:      "directory",
			storageID: "dir-test",
			createReq: DirectoryStorageCreateRequest{
				DataStoreCommonImmutableFields: DataStoreCommonImmutableFields{
					ID:   ptr("dir-test"),
					Type: ptr("dir"),
				},
				DirectoryStorageMutableFields: DirectoryStorageMutableFields{
					DataStoreCommonMutableFields: DataStoreCommonMutableFields{
						Nodes:        &nodes,
						ContentTypes: &contentImages,
						Disable:      &disableFalse,
					},
					Backups: backups,
					Shared:  &sharedFalse,
				},
				DirectoryStorageImmutableFields: DirectoryStorageImmutableFields{
					Path: ptr("/var/lib/vz"),
				},
			},
			updateReq: DirectoryStorageUpdateRequest{
				DirectoryStorageMutableFields: DirectoryStorageMutableFields{
					DataStoreCommonMutableFields: DataStoreCommonMutableFields{
						Nodes:        &nodes,
						ContentTypes: &contentImages,
						Disable:      &disableFalse,
					},
					Backups: backups,
					Shared:  &sharedFalse,
				},
			},
			createKeys: []string{
				"content",
				"disable",
				"max-protected-backups",
				"nodes",
				"path",
				"prune-backups",
				"shared",
				"storage",
				"type",
			},
			updateKeys: []string{
				"content",
				"disable",
				"max-protected-backups",
				"nodes",
				"prune-backups",
				"shared",
			},
			updateNoKey: []string{"path", "storage", "type"},
		},
		{
			name:      "nfs",
			storageID: "nfs-test",
			createReq: NFSStorageCreateRequest{
				DataStoreCommonImmutableFields: DataStoreCommonImmutableFields{
					ID:   ptr("nfs-test"),
					Type: ptr("nfs"),
				},
				NFSStorageMutableFields: NFSStorageMutableFields{
					DataStoreCommonMutableFields: DataStoreCommonMutableFields{
						Nodes:        &nodes,
						ContentTypes: &contentImages,
						Disable:      &disableFalse,
					},
					Backups: backups,
				},
				NFSStorageImmutableFields: NFSStorageImmutableFields{
					Server: ptr("127.0.0.1"),
					Export: ptr("/export"),
				},
			},
			updateReq: NFSStorageUpdateRequest{
				NFSStorageMutableFields: NFSStorageMutableFields{
					DataStoreCommonMutableFields: DataStoreCommonMutableFields{
						Nodes:        &nodes,
						ContentTypes: &contentImages,
						Disable:      &disableFalse,
					},
					Backups: backups,
				},
			},
			createKeys: []string{
				"content",
				"disable",
				"export",
				"max-protected-backups",
				"nodes",
				"prune-backups",
				"server",
				"storage",
				"type",
			},
			updateKeys: []string{
				"content",
				"disable",
				"max-protected-backups",
				"nodes",
				"prune-backups",
			},
			updateNoKey: []string{"export", "server", "storage", "type"},
		},
		{
			name:      "cifs",
			storageID: "cifs-test",
			createReq: CIFSStorageCreateRequest{
				DataStoreCommonImmutableFields: DataStoreCommonImmutableFields{
					ID:   ptr("cifs-test"),
					Type: ptr("cifs"),
				},
				CIFSStorageMutableFields: CIFSStorageMutableFields{
					DataStoreCommonMutableFields: DataStoreCommonMutableFields{
						Nodes:        &nodes,
						ContentTypes: &contentImages,
						Disable:      &disableFalse,
					},
					Backups: backups,
				},
				CIFSStorageImmutableFields: CIFSStorageImmutableFields{
					Server:   ptr("127.0.0.1"),
					Username: ptr("user"),
					Password: ptr("pass"),
					Share:    ptr("share"),
				},
			},
			updateReq: CIFSStorageUpdateRequest{
				CIFSStorageMutableFields: CIFSStorageMutableFields{
					DataStoreCommonMutableFields: DataStoreCommonMutableFields{
						Nodes:        &nodes,
						ContentTypes: &contentImages,
						Disable:      &disableFalse,
					},
					Backups: backups,
				},
			},
			createKeys: []string{
				"content",
				"disable",
				"max-protected-backups",
				"nodes",
				"password",
				"prune-backups",
				"server",
				"share",
				"storage",
				"type",
				"username",
			},
			updateKeys: []string{
				"content",
				"disable",
				"max-protected-backups",
				"nodes",
				"prune-backups",
			},
			updateNoKey: []string{"password", "server", "share", "storage", "type", "username"},
		},
		{
			name:      "pbs",
			storageID: "pbs-test",
			createReq: PBSStorageCreateRequest{
				DataStoreCommonImmutableFields: DataStoreCommonImmutableFields{
					ID:   ptr("pbs-test"),
					Type: ptr("pbs"),
				},
				PBSStorageMutableFields: PBSStorageMutableFields{
					DataStoreCommonMutableFields: DataStoreCommonMutableFields{
						Nodes:        &nodes,
						ContentTypes: &contentImages,
						Disable:      &disableFalse,
					},
					Backups: backups,
				},
				PBSStorageImmutableFields: PBSStorageImmutableFields{
					Server:    ptr("127.0.0.1"),
					Datastore: ptr("ds"),
				},
			},
			updateReq: PBSStorageUpdateRequest{
				PBSStorageMutableFields: PBSStorageMutableFields{
					DataStoreCommonMutableFields: DataStoreCommonMutableFields{
						Nodes:        &nodes,
						ContentTypes: &contentImages,
						Disable:      &disableFalse,
					},
					Backups: backups,
				},
			},
			createKeys: []string{
				"content",
				"datastore",
				"disable",
				"max-protected-backups",
				"nodes",
				"prune-backups",
				"server",
				"storage",
				"type",
			},
			updateKeys: []string{
				"content",
				"disable",
				"max-protected-backups",
				"nodes",
				"prune-backups",
			},
			updateNoKey: []string{"datastore", "server", "storage", "type"},
		},
		{
			name:      "lvm",
			storageID: "lvm-test",
			createReq: LVMStorageCreateRequest{
				DataStoreCommonImmutableFields: DataStoreCommonImmutableFields{
					ID:   ptr("lvm-test"),
					Type: ptr("lvm"),
				},
				LVMStorageMutableFields: LVMStorageMutableFields{
					DataStoreCommonMutableFields: DataStoreCommonMutableFields{
						Nodes:        &nodes,
						ContentTypes: &contentImages,
						Disable:      &disableFalse,
					},
					WipeRemovedVolumes: types.CustomBool(false),
				},
				LVMStorageImmutableFields: LVMStorageImmutableFields{
					VolumeGroup: ptr("vg0"),
				},
			},
			updateReq: LVMStorageUpdateRequest{
				LVMStorageMutableFields: LVMStorageMutableFields{
					DataStoreCommonMutableFields: DataStoreCommonMutableFields{
						Nodes:        &nodes,
						ContentTypes: &contentImages,
						Disable:      &disableFalse,
					},
					WipeRemovedVolumes: types.CustomBool(false),
				},
			},
			createKeys: []string{
				"content",
				"disable",
				"nodes",
				"saferemove",
				"storage",
				"type",
				"vgname",
			},
			updateKeys: []string{
				"content",
				"disable",
				"nodes",
				"saferemove",
			},
			updateNoKey: []string{"storage", "type", "vgname"},
		},
		{
			name:      "lvmthin",
			storageID: "lvmthin-test",
			createReq: LVMThinStorageCreateRequest{
				DataStoreCommonImmutableFields: DataStoreCommonImmutableFields{
					ID:   ptr("lvmthin-test"),
					Type: ptr("lvmthin"),
				},
				LVMThinStorageMutableFields: LVMThinStorageMutableFields{
					DataStoreCommonMutableFields: DataStoreCommonMutableFields{
						Nodes:        &nodes,
						ContentTypes: &contentImages,
						Disable:      &disableFalse,
					},
				},
				LVMThinStorageImmutableFields: LVMThinStorageImmutableFields{
					VolumeGroup: ptr("vg0"),
					ThinPool:    ptr("data"),
				},
			},
			updateReq: LVMThinStorageUpdateRequest{
				LVMThinStorageMutableFields: LVMThinStorageMutableFields{
					DataStoreCommonMutableFields: DataStoreCommonMutableFields{
						Nodes:        &nodes,
						ContentTypes: &contentImages,
						Disable:      &disableFalse,
					},
				},
			},
			createKeys: []string{
				"content",
				"disable",
				"nodes",
				"storage",
				"thinpool",
				"type",
				"vgname",
			},
			updateKeys: []string{
				"content",
				"disable",
				"nodes",
			},
			updateNoKey: []string{"storage", "thinpool", "type", "vgname"},
		},
		{
			name:      "zfspool",
			storageID: "zfs-test",
			createReq: ZFSStorageCreateRequest{
				DataStoreCommonImmutableFields: DataStoreCommonImmutableFields{
					ID:   ptr("zfs-test"),
					Type: ptr("zfspool"),
				},
				ZFSStorageMutableFields: ZFSStorageMutableFields{
					DataStoreCommonMutableFields: DataStoreCommonMutableFields{
						Nodes:        &nodes,
						ContentTypes: &contentImages,
						Disable:      &disableFalse,
					},
					ThinProvision: types.CustomBool(true),
				},
				ZFSStorageImmutableFields: ZFSStorageImmutableFields{
					ZFSPool: ptr("rpool/data"),
				},
			},
			updateReq: ZFSStorageUpdateRequest{
				ZFSStorageMutableFields: ZFSStorageMutableFields{
					DataStoreCommonMutableFields: DataStoreCommonMutableFields{
						Nodes:        &nodes,
						ContentTypes: &contentImages,
						Disable:      &disableFalse,
					},
					ThinProvision: types.CustomBool(true),
				},
			},
			createKeys: []string{
				"content",
				"disable",
				"nodes",
				"pool",
				"sparse",
				"storage",
				"type",
			},
			updateKeys: []string{
				"content",
				"disable",
				"nodes",
				"sparse",
			},
			updateNoKey: []string{"pool", "storage", "type"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server, captures := newStorageContractServer()
			t.Cleanup(server.Close)

			storageClient := newStorageClientForTest(t, server.URL)
			ctx := context.Background()

			_, err := storageClient.CreateDatastore(ctx, tc.createReq)
			if err != nil {
				t.Fatalf("CreateDatastore returned error: %v", err)
			}

			if err := storageClient.UpdateDatastore(ctx, tc.storageID, tc.updateReq); err != nil {
				t.Fatalf("UpdateDatastore returned error: %v", err)
			}

			got := captures.all()
			if len(got) != 2 {
				t.Fatalf("expected 2 requests, got %d", len(got))
			}

			assertRequest(t, got[0], http.MethodPost, "/api2/json/storage", tc.createKeys, nil)
			assertRequest(t, got[1], http.MethodPut, "/api2/json/storage/"+tc.storageID, tc.updateKeys, tc.updateNoKey)
		})
	}
}

type requestCaptures struct {
	mu  sync.Mutex
	req []capturedRequest
}

func (c *requestCaptures) add(r capturedRequest) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.req = append(c.req, r)
}

func (c *requestCaptures) all() []capturedRequest {
	c.mu.Lock()
	defer c.mu.Unlock()

	out := make([]capturedRequest, len(c.req))
	copy(out, c.req)

	return out
}

func newStorageContractServer() (*httptest.Server, *requestCaptures) {
	captures := &requestCaptures{}

	mux := http.NewServeMux()
	mux.HandleFunc("/api2/json/storage", func(w http.ResponseWriter, r *http.Request) {
		captures.add(captureRequest(r))

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		storeID := r.PostFormValue("storage")
		storeType := r.PostFormValue("type")

		if err := json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"storage": storeID,
				"type":    storeType,
				"config":  map[string]any{},
			},
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("/api2/json/storage/", func(w http.ResponseWriter, r *http.Request) {
		captures.add(captureRequest(r))

		switch r.Method {
		case http.MethodPut, http.MethodDelete, http.MethodGet:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return httptest.NewTLSServer(mux), captures
}

func captureRequest(r *http.Request) capturedRequest {
	cr := capturedRequest{
		Method: r.Method,
		Path:   r.URL.Path,
	}

	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		if err := r.ParseForm(); err == nil {
			for k := range r.PostForm {
				cr.BodyKeys = append(cr.BodyKeys, k)
			}
		}

		sort.Strings(cr.BodyKeys)
	}

	return cr
}

func newStorageClientForTest(t *testing.T, endpoint string) *Client {
	t.Helper()

	conn, err := api.NewConnection(endpoint, true, "1.2")
	if err != nil {
		t.Fatalf("NewConnection returned error: %v", err)
	}

	creds, err := api.NewCredentials("", "", "", "user@pve!token=abcd", "", "")
	if err != nil {
		t.Fatalf("NewCredentials returned error: %v", err)
	}

	c, err := api.NewClient(creds, conn)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	return &Client{Client: c}
}

func assertRequest(
	t *testing.T,
	got capturedRequest,
	wantMethod string,
	wantPath string,
	wantKeys []string,
	wantNoKeys []string,
) {
	t.Helper()

	if got.Method != wantMethod {
		t.Fatalf("unexpected method: got %q want %q", got.Method, wantMethod)
	}

	if got.Path != wantPath {
		t.Fatalf("unexpected path: got %q want %q", got.Path, wantPath)
	}

	assertKeysPresent(t, got.BodyKeys, wantKeys)
	assertKeysAbsent(t, got.BodyKeys, wantNoKeys)
}

func assertKeysPresent(t *testing.T, gotKeys []string, wantKeys []string) {
	t.Helper()

	for _, want := range wantKeys {
		if !containsString(gotKeys, want) {
			t.Fatalf("missing key %q in request body keys: %v", want, gotKeys)
		}
	}
}

func assertKeysAbsent(t *testing.T, gotKeys []string, wantNoKeys []string) {
	t.Helper()

	for _, want := range wantNoKeys {
		if containsString(gotKeys, want) {
			t.Fatalf("unexpected key %q in request body keys: %v", want, gotKeys)
		}
	}
}

func containsString(haystack []string, needle string) bool {
	return slices.Contains(haystack, needle)
}

func ptr(s string) *string {
	v := strings.Clone(s)
	return &v
}
