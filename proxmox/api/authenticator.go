/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package api

import (
	"context"
	"net/http"
)

// Authenticator is an interface for adding authentication data to a request.
// The authenticator is also aware of the authentication context, e.g. if it
// is configured to use the root user.
type Authenticator interface {
	// IsRoot returns true if the authenticator is configured to use the root
	IsRoot(ctx context.Context) bool

	// IsRootTicket returns true if the authenticator is configured to use the root directly using a login ticket.
	// (root using token is weaker, cannot change VM arch)
	IsRootTicket(ctx context.Context) bool

	// AuthenticateRequest adds authentication data to a new request.
	AuthenticateRequest(ctx context.Context, req *http.Request) error
}
