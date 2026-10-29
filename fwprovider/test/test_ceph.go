/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/utils"
)

// RequireRootSSH skips the test unless SSH as root is configured: PROXMOX_VE_SSH_USERNAME=root with an agent, a
// private key or a password.
func (e *Environment) RequireRootSSH() {
	e.t.Helper()

	if utils.GetAnyStringEnv("PROXMOX_VE_SSH_USERNAME") != "root" ||
		(!utils.GetAnyBoolEnv("PROXMOX_VE_SSH_AGENT") &&
			utils.GetAnyStringEnv("PROXMOX_VE_SSH_PRIVATE_KEY", "PROXMOX_VE_SSH_PASSWORD") == "") {
		e.t.Skip("Skipping: needs SSH as root (PROXMOX_VE_SSH_USERNAME=root with an agent, key or password)")
	}
}

// CreateCephMDS creates a Ceph metadata server on the test node, waits until it's in standby, and removes it when the
// test ends. A CephFS needs an MDS in standby to be created, and the cluster may have none.
func (e *Environment) CreateCephMDS() string {
	e.t.Helper()

	name := SafeResourceName("ceph-fs-mds")
	path := e.nodePath("ceph/mds/" + url.PathEscape(name))

	e.t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		if err := e.runTask(ctx, http.MethodDelete, path, nil); err != nil {
			e.t.Errorf("removing test MDS %q: %v", name, err)
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	require.NoError(e.t, e.runTask(ctx, http.MethodPost, path, nil), "creating test MDS %q", name)

	type mdsEntry struct {
		Name  string `json:"name"`
		State string `json:"state"`
	}

	for {
		resp := &struct {
			Data []mdsEntry `json:"data"`
		}{}
		require.NoError(e.t, e.Client().DoRequest(ctx, http.MethodGet, e.nodePath("ceph/mds"), nil, resp))

		if slices.ContainsFunc(resp.Data, func(d mdsEntry) bool { return d.Name == name && d.State == "up:standby" }) {
			return name
		}

		select {
		case <-ctx.Done():
			e.t.Fatalf("test MDS %q didn't reach up:standby", name)
		case <-time.After(2 * time.Second):
		}
	}
}

// RemoveCephFS removes a CephFS and its pools the way the admin guide's "Destroy CephFS" does by hand: `ceph fs fail`,
// which the API doesn't offer, then `pveceph fs destroy --remove-pools`, over SSH as root. It then checks with the API
// that the filesystem and its pools are gone. A filesystem that doesn't exist is fine.
func (e *Environment) RemoveCephFS(ctx context.Context, name string) error {
	pools := []string{name + "_metadata", name + "_data"}

	out, err := e.SSHClient().ExecuteNodeCommands(ctx, e.NodeName, []string{
		fmt.Sprintf("ceph fs fail %s 2>&1 || true", name),
		fmt.Sprintf("pveceph fs destroy %s --remove-pools 1 2>&1 || true", name),
	})
	if err != nil {
		return fmt.Errorf("removing filesystem %q over SSH: %w\n%s", name, err, out)
	}

	_, err = e.NodeClient().Ceph().FS().Get(ctx, name)
	if err == nil {
		return fmt.Errorf("filesystem %q still exists after removal:\n%s", name, out)
	}

	if !errors.Is(err, api.ErrResourceDoesNotExist) {
		return fmt.Errorf("checking filesystem %q: %w", name, err)
	}

	for _, pool := range pools {
		_, err := e.NodeClient().Ceph().Pool().Get(ctx, pool)
		if err == nil {
			return fmt.Errorf("pool %q still exists after removing CephFS %q:\n%s", pool, name, out)
		}

		if !errors.Is(err, api.ErrResourceDoesNotExist) {
			return fmt.Errorf("checking pool %q: %w", pool, err)
		}
	}

	return nil
}

// runTask sends a request that starts a Proxmox task and waits for the task.
func (e *Environment) runTask(ctx context.Context, method, path string, body any) error {
	resp := &struct {
		Data *string `json:"data"`
	}{}
	if err := e.Client().DoRequest(ctx, method, path, body, resp); err != nil {
		return err
	}

	if resp.Data == nil {
		return fmt.Errorf("%s %s: no task ID in the response", method, path)
	}

	return e.NodeClient().Tasks().WaitForTask(ctx, *resp.Data).Err()
}

func (e *Environment) nodePath(path string) string {
	return fmt.Sprintf("nodes/%s/%s", url.PathEscape(e.NodeName), path)
}
