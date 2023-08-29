/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package api

import "fmt"

// Error is a sentinel error type for API errors.
type Error string

func (err Error) Error() string {
	return string(err)
}

// ErrNoDataObjectInResponse is returned when the server does not include a data object in the response.
const ErrNoDataObjectInResponse Error = "the server did not include a data object in the response"

// HTTPError is a generic error type for HTTP errors.
type HTTPError struct {
	Code    int
	Message string
}

func (err *HTTPError) Error() string {
	return fmt.Sprintf("received an HTTP %d response - Reason: %s", err.Code, err.Message)
}
