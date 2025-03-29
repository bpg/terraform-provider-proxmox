/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package hardwaremapping

import (
	"fmt"

	proxmoxtypes "github.com/bpg/terraform-provider-proxmox/proxmox/types/hardwaremapping"
)

// ErrResourceMessageInvalidPath is the error message for an invalid Linux device path for a hardware mapping of the
// specified type.
// Extracting the message helps to reduce duplicated code and allows to use it in automated unit and acceptance tests.
//
//nolint:gochecknoglobals
var ErrResourceMessageInvalidPath = func(hmType proxmoxtypes.Type) string {
	return fmt.Sprintf("not a valid Linux device path for hardware mapping of type %q", hmType)
}
