/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package repositories

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

//nolint:gochecknoglobals
var (
	// ErrCephVersionNameIllegal indicates an error for an illegal Ceph major version name.
	ErrCephVersionNameIllegal = func(name string) error {
		return function.NewFuncError(fmt.Sprintf("illegal Ceph major version name %q", name))
	}

	// ErrCephVersionNameMarshal indicates an error while marshalling a Ceph major version name.
	ErrCephVersionNameMarshal = function.NewFuncError("cannot marshal Ceph major version name")

	// ErrCephVersionNameUnmarshal indicates an error while unmarshalling a Ceph major version name.
	ErrCephVersionNameUnmarshal = function.NewFuncError("cannot unmarshal Ceph major version name")

	// ErrStandardRepoHandleKindIllegal indicates an error for an illegal APT standard repository handle.
	ErrStandardRepoHandleKindIllegal = func(handle string) error {
		return function.NewFuncError(fmt.Sprintf("illegal APT standard repository handle kind %q", handle))
	}

	// ErrStandardRepoHandleKindMarshal indicates an error while marshalling an APT standard repository handle kind.
	ErrStandardRepoHandleKindMarshal = function.NewFuncError("cannot marshal APT standard repository handle kind")

	// ErrStandardRepoHandleKindUnmarshal indicates an error while unmarshalling an APT standard repository handle kind.
	ErrStandardRepoHandleKindUnmarshal = function.NewFuncError("cannot unmarshal APT standard repository handle kind")
)
