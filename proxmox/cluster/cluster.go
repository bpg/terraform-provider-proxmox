/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cluster

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

const (
	getVMIDStep = 1
)

var (
	//nolint:gochecknoglobals
	getVMIDCounter = -1
	//nolint:gochecknoglobals
	getVMIDCounterMutex = &sync.Mutex{}
)

// GetNextID retrieves the next free VM identifier for the cluster.
func (c *Client) GetNextID(ctx context.Context, vmID *int) (*int, error) {
	reqBody := &NextIDRequestBody{
		VMID: vmID,
	}

	resBody := &NextIDResponseBody{}

	err := c.DoRequest(ctx, http.MethodGet, "cluster/nextid", reqBody, resBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving next VM ID: %w", err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return (*int)(resBody.Data), nil
}

// GetVMID retrieves the next available VM identifier.
func (c *Client) GetVMID(ctx context.Context) (*int, error) {
	getVMIDCounterMutex.Lock()
	defer getVMIDCounterMutex.Unlock()

	if getVMIDCounter < 0 {
		nextVMID, err := c.GetNextID(ctx, nil)
		if err != nil {
			return nil, err
		}

		if nextVMID == nil {
			return nil, errors.New("unable to retrieve the next available VM identifier")
		}

		getVMIDCounter = *nextVMID + getVMIDStep

		tflog.Debug(ctx, "next VM identifier", map[string]interface{}{
			"id": *nextVMID,
		})

		return nextVMID, nil
	}

	vmID := getVMIDCounter

	for vmID <= 2147483637 {
		_, err := c.GetNextID(ctx, &vmID)
		if err != nil {
			vmID += getVMIDStep

			continue
		}

		getVMIDCounter = vmID + getVMIDStep

		tflog.Debug(ctx, "next VM identifier", map[string]interface{}{
			"id": vmID,
		})

		return &vmID, nil
	}

	return nil, errors.New("unable to determine the next available VM identifier")
}

// GetClusterResources retrieves current resources for cluster.
func (c *Client) GetClusterResources(ctx context.Context, resourceType string) ([]*ResourcesListResponseData, error) {
	reqBody := &ResourcesListRequestBody{
		Type: resourceType,
	}
	resBody := &ResourcesListBody{}

	err := c.DoRequest(ctx, http.MethodGet, "cluster/resources", reqBody, resBody)
	if err != nil {
		return nil, fmt.Errorf("failed to get resources list of type (\"%s\") for cluster: %w", resourceType, err)
	}

	if resBody.Data == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return resBody.Data, nil
}

// GetClusterResourcesVM retrieves current VM resources for cluster.
func (c *Client) GetClusterResourcesVM(ctx context.Context) ([]*ResourcesListResponseData, error) {
	return c.GetClusterResources(ctx, "vm")
}

// GetVMNodeName gets node for specified vmID.
func (c *Client) GetVMNodeName(ctx context.Context, vmID int) (*string, error) {
	allClusterVM, err := c.GetClusterResourcesVM(ctx)
	if err != nil {
		return nil, err
	}

	if allClusterVM == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	for _, v := range allClusterVM {
		if v.VMID == vmID {
			return &v.NodeName, nil
		}
	}

	return nil, errors.New("unable to determine node name for VM identifier")
}
