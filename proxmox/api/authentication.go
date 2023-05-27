/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package api

import (
	"net/http"
)

type Authenticator interface {
	// IsRoot returns true if the authenticator is configured to use the root
	IsRoot() bool

	// AuthenticateRequest adds authentication data to a new request.
	AuthenticateRequest(req *http.Request) error
}
