/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package applier

// ApplyResponseBody represents the response of PUT /cluster/sdn.
// PVE typically returns a UPID string in `data` for async tasks.
// We keep it here even if ApplyConfig currently ignores it, so the
// API can easily be extended later to surface it.
type ApplyResponseBody struct {
	Data *string `json:"data"`
}
