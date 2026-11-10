/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cloudinitbootstrap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bpg/terraform-provider-proxmox/proxmox"
	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/storage"
	"github.com/bpg/terraform-provider-proxmox/proxmox/nodes/vms"
	"github.com/bpg/terraform-provider-proxmox/proxmox/retry"
)

// cloudInitResultPath is written by cloud-init when its final stage finishes; /run is cleared on reboot,
// so a result there always comes from the boot with our seed.
const cloudInitResultPath = "/run/cloud-init/result.json"

const cloudInitPollInterval = 10 * time.Second

// uploadISO builds the seed ISO and uploads it to the datastore as fileName.
func uploadISO(ctx context.Context, client proxmox.Client, ds *storage.Client, files map[string]string, fileName string) error {
	// The ISO holds the (sensitive) seed contents: keep it in a private 0700 directory.
	dir, err := os.MkdirTemp(client.TempDir(), "cloudinit-iso-")
	if err != nil {
		return fmt.Errorf("failed to create ISO directory: %w", err)
	}

	defer os.RemoveAll(dir)

	isoPath, err := buildISO(dir, files)
	if err != nil {
		return err
	}

	// Uploads run in parallel with other heavy API work and pveproxy sometimes drops the connection (EOF).
	// Retrying is safe: the file name is content-addressed and PVE only moves a complete upload into place.
	op := retry.NewAPICallOperation("cloud-init ISO upload",
		retry.WithAttempts(4),
		retry.WithBaseDelay(2*time.Second),
		retry.WithRetryIf(isRetryableUploadError),
	)

	return op.Do(ctx, func() error {
		f, err := os.Open(isoPath)
		if err != nil {
			return fmt.Errorf("failed to open ISO: %w", err)
		}

		defer f.Close()

		_, err = ds.APIUpload(ctx, &api.FileUploadRequest{
			ContentType: storage.ContentTypeISO,
			FileName:    fileName,
			File:        f,
		}, dir)

		return err
	})
}

func isRetryableUploadError(err error) bool {
	return retry.IsTransientAPIError(err) ||
		errors.Is(err, io.EOF) ||
		errors.Is(err, io.ErrUnexpectedEOF) ||
		errors.Is(err, syscall.ECONNRESET)
}

type cloudInitResult struct {
	V1 struct {
		Errors []string `json:"errors"`
	} `json:"v1"`
}

func checkCloudInitResult(content string) error {
	var res cloudInitResult
	if err := json.Unmarshal([]byte(content), &res); err != nil {
		return fmt.Errorf("unable to parse the cloud-init result: %w", err)
	}

	if len(res.V1.Errors) > 0 {
		return fmt.Errorf("cloud-init finished with errors: %s", strings.Join(res.V1.Errors, "; "))
	}

	return nil
}

// isPermanentAgentError reports errors that polling can't outlast: a missing privilege or no agent configured.
// "Agent not running" and "file not found" are expected while the guest boots and are retried.
func isPermanentAgentError(err error) bool {
	var httpErr *api.HTTPError
	if errors.As(err, &httpErr) && httpErr.Code >= 400 && httpErr.Code < 500 {
		return true
	}

	return strings.Contains(err.Error(), "No QEMU guest agent configured")
}

// waitForCloudInit polls the guest until cloud-init has written its result file.
func waitForCloudInit(ctx context.Context, vm *vms.Client, timeout time.Duration) error {
	pollCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var (
		content string
		lastErr error
	)

	op := retry.NewPollOperation("cloud-init result",
		retry.WithBaseDelay(cloudInitPollInterval),
		retry.WithRetryIf(func(err error) bool { return !isPermanentAgentError(err) }),
	)

	err := op.DoPoll(pollCtx, func() error {
		c, err := vm.GetFileFromAgent(pollCtx, cloudInitResultPath)
		if err != nil {
			// Keep the agent's error, not the deadline that cut the last request short.
			if pollCtx.Err() == nil || lastErr == nil {
				lastErr = err
			}

			tflog.Debug(ctx, "waiting for cloud-init to finish", map[string]any{"vm_id": vm.VMID, "error": err.Error()})

			return err
		}

		content = c

		return nil
	})

	switch {
	case err == nil:
		return checkCloudInitResult(content)
	case ctx.Err() != nil:
		return fmt.Errorf("cancelled while waiting for cloud-init to finish: %w", ctx.Err())
	case pollCtx.Err() != nil && lastErr == nil:
		return fmt.Errorf("timed out after %s waiting for cloud-init to finish", timeout)
	case pollCtx.Err() != nil:
		return fmt.Errorf("timed out after %s waiting for cloud-init to finish (last error: %w)", timeout, lastErr)
	default:
		return fmt.Errorf("unable to read the cloud-init result: %w", err)
	}
}
