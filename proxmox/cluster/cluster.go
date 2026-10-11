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

	"github.com/bpg/terraform-provider-proxmox/proxmox/api"
)

// ErrVMDoesNotExist is returned when the VM identifier cannot be found on any cluster node.
var ErrVMDoesNotExist = errors.New("unable to find VM identifier on any cluster node")

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

	return nil, ErrVMDoesNotExist
}

// FindContainerResource returns the cluster resource entry for the given container VMID, filtering
// out qemu entries so a container and a VM never collide even if VMID uniqueness is ever relaxed.
func FindContainerResource(resources []*ResourcesListResponseData, vmID int) (*ResourcesListResponseData, error) {
	for _, r := range resources {
		if r.Type == "lxc" && r.VMID == vmID {
			return r, nil
		}
	}

	return nil, ErrVMDoesNotExist
}

// GetContainerResource gets the cluster resource entry for the specified container, carrying its
// NodeName, Status and HaState.
func (c *Client) GetContainerResource(ctx context.Context, vmID int) (*ResourcesListResponseData, error) {
	allClusterVM, err := c.GetClusterResourcesVM(ctx)
	if err != nil {
		return nil, err
	}

	if allClusterVM == nil {
		return nil, api.ErrNoDataObjectInResponse
	}

	return FindContainerResource(allClusterVM, vmID)
}

// GetContainerNodeName gets the node hosting the specified container.
func (c *Client) GetContainerNodeName(ctx context.Context, vmID int) (*string, error) {
	res, err := c.GetContainerResource(ctx, vmID)
	if err != nil {
		return nil, err
	}

	return &res.NodeName, nil
}
