/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cluster

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/rogpeppe/go-internal/lockedfile"

	"github.com/bpg/terraform-provider-proxmox/proxmox/helpers/ptr"
)

const (
	idGeneratorLockFile         = "terraform-provider-proxmox-id-gen.lock"
	idGeneratorSequenceFile     = "terraform-provider-proxmox-id-gen.seq"
	idGeneratorContentionWindow = 5 * time.Second
)

// IDGenerator is responsible for generating unique identifiers for VMs and Containers.
type IDGenerator struct {
	client *Client
	config IDGeneratorConfig
}

// IDGeneratorConfig is the configuration for the IDGenerator.
type IDGeneratorConfig struct {
	RandomIDs    bool
	RandomIDStat int
	RandomIDEnd  int

	lockFName string
	seqFName  string
}

// NewIDGenerator creates a new IDGenerator with the given parameters.
func NewIDGenerator(client *Client, config IDGeneratorConfig) IDGenerator {
	if config.RandomIDStat == 0 {
		config.RandomIDStat = 10000
	}

	if config.RandomIDEnd == 0 {
		config.RandomIDEnd = 99999
	}

	config.lockFName = filepath.Join(os.TempDir(), idGeneratorLockFile)
	config.seqFName = filepath.Join(os.TempDir(), idGeneratorSequenceFile)

	unlock, err := lockedfile.MutexAt(config.lockFName).Lock()
	if err == nil {
		defer unlock()

		// delete the sequence file if it is older than 10 seconds
		// this is to prevent the sequence file from growing indefinitely,
		// while giving some protection against parallel runs of the provider
		// that might interfere with each other and reset the sequence at the same time
		stat, err := os.Stat(config.seqFName)
		if err == nil && time.Since(stat.ModTime()) > idGeneratorContentionWindow {
			_ = os.Remove(config.seqFName)
		}
	}

	return IDGenerator{client, config}
}

// NextID returns the next available VM identifier.
func (g IDGenerator) NextID(ctx context.Context) (int, error) {
	// lock the ID generator to prevent concurrent access
	// it should be unlocked only when the new ID is successfully
	// retrieved (and optionally written to the sequence file)
	unlock, err := lockedfile.MutexAt(g.config.lockFName).Lock()
	if err != nil {
		return -1, fmt.Errorf("unable to lock the ID generator: %w", err)
	}

	defer unlock()

	ctx, cancel := context.WithTimeout(ctx, idGeneratorContentionWindow+time.Second)
	defer cancel()

	var newID *int

	var errs []error

	id, err := retry.DoWithData(func() (*int, error) {
		if g.config.RandomIDs {
			//nolint:gosec
			newID = ptr.Ptr(rand.Intn(g.config.RandomIDEnd-g.config.RandomIDStat) + g.config.RandomIDStat)
		} else if newID == nil {
			newID, err = nextSequentialID(g.config.seqFName)
			if err != nil {
				return nil, err
			}
		}

		return g.client.GetNextID(ctx, newID)
	},
		retry.OnRetry(func(_ uint, err error) {
			if strings.Contains(err.Error(), "already exists") && newID != nil {
				newID, err = g.client.GetNextID(ctx, nil)
			}

			errs = append(errs, err)
		}),
		retry.Context(ctx),
		retry.UntilSucceeded(),
		retry.DelayType(retry.FixedDelay),
		retry.Delay(200*time.Millisecond),
	)
	if err != nil {
		errs = append(errs, err)
		return -1, fmt.Errorf("unable to retrieve the next available VM identifier: %w", errors.Join(errs...))
	}

	if !g.config.RandomIDs {
		var b bytes.Buffer

		_, _ = fmt.Fprintf(&b, "%d", *id)

		if err := lockedfile.Write(g.config.seqFName, &b, 0o666); err != nil {
			return -1, fmt.Errorf("unable to write the ID generator file: %w", err)
		}
	}

	return *id, nil
}

func nextSequentialID(seqFName string) (*int, error) {
	buf, err := lockedfile.Read(seqFName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil //nolint:nilnil
		}

		return nil, fmt.Errorf("unable to read the ID generator sequence file: %w", err)
	}

	id, err := strconv.Atoi(string(buf))
	if err != nil {
		return nil, fmt.Errorf("unable to parse the ID generator file: %w", err)
	}

	return ptr.Ptr(id + 1), nil
}
