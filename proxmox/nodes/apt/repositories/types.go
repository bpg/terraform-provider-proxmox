/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package repositories

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// baseData contains common data for APT repository API calls.
type baseData struct {
	Digest *string `json:"digest,omitempty" url:"digest,omitempty"`
}

// file contains the data of a parsed APT repository file.
type file struct {
	// FileType is the format of the file.
	FileType string `json:"file-type"`

	// Path is the path to the repository file.
	Path string `json:"path"`

	// Repositories is the list of parsed repositories.
	Repositories []*repo `json:"repositories"`
}

// repo contains the data of an APT repository from a parsed file.
type repo struct {
	// Comment is the associated comment.
	Comment *string `json:"Comment,omitempty"`

	// Components is the list of repository components.
	Components []string `json:"Components,omitempty"`

	// Enabled indicates whether the repository is enabled.
	Enabled types.CustomBool `json:"Enabled"`

	// FileTpe is the format of the defining file.
	FileType string `json:"FileType"`

	// PackageTypes is the list of package types.
	PackageTypes []string `json:"Types"`

	// Suites is the list of package distributions.
	Suites []string `json:"Suites"`

	// URIs is the list of repository URIs.
	URIs []string `json:"URIs"`
}

// standardRepo contains the data for an APT standard repository.
type standardRepo struct {
	// Description is the description of the APT standard repository.
	Description *string `json:"description,omitempty"`

	// Handle is the pre-defined handle of the APT standard repository.
	Handle string `json:"handle"`

	// Name is the human-readable name of the APT standard repository.
	Name string `json:"Name"`

	// Status is the activation status of the APT standard repository.
	// Can be either 0 (disabled) or 1 (enabled).
	Status *int64 `json:"status,omitempty"`
}

// AddRequestBody contains the body for an APT repository PUT request to add a standard repository.
type AddRequestBody struct {
	baseData

	// Handle is the pre-defined handle of the APT standard repository.
	Handle string `json:"handle" url:"handle"`

	// Node is the name of the target Proxmox VE node.
	Node string `json:"node" url:"node"`
}

// GetResponseBody is the body from an APT repository GET response.
type GetResponseBody struct {
	Data *GetResponseData `json:"data,omitempty"`
}

// GetResponseData contains the data from an APT repository GET response.
type GetResponseData struct {
	baseData

	// Files contains the APT repository files.
	Files []*file `json:"files,omitempty"`

	// StandardRepos contains the APT standard repositories.
	StandardRepos []*standardRepo `json:"standard-repos,omitempty"`
}

// ModifyRequestBody contains the body for an APT repository POST request to modify a repository.
type ModifyRequestBody struct {
	baseData

	// Enabled indicates the activation status of the APT repository.
	// Must either be 0 (disabled) or 1 (enabled).
	Enabled types.CustomBool `json:"enabled" url:"enabled,int"`

	// Index is the index of the APT repository within the defining repository source file.
	Index int64 `json:"handle" url:"index"`

	// Path is the absolute path of the defining source file for the APT repository.
	Path string `json:"path" url:"path"`
}
