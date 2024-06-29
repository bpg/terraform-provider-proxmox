/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package apt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

//nolint:gochecknoglobals
var (
	// ErrValueConversion indicates an error while converting a value for a Proxmox VE API APT entity.
	ErrValueConversion = func(format string, attrs ...any) error {
		return function.NewFuncError(fmt.Sprintf(format, attrs...))
	}
)
