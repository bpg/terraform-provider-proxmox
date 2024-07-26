/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package account

// ACMEAccountListResponseBody contains the body from a ACME account list response.
type ACMEAccountListResponseBody struct {
	Data []*ACMEAccountListResponseData `json:"data,omitempty"`
}

// ACMEAccountListResponseData contains the data from a ACME account list response.
type ACMEAccountListResponseData struct {
	Name string `json:"name"`
}

// ACMEAccountGetResponseBody contains the body from a ACME account get response.
type ACMEAccountGetResponseBody struct {
	Data *ACMEAccountGetResponseData `json:"data,omitempty"`
}

// ACMEAccountGetResponseData contains the data from a ACME account get response.
type ACMEAccountGetResponseData struct {
	// Account is the ACME account data.
	Account map[string]interface{} `json:"account"`
	// Directory is the URL of the ACME CA directory endpoint.
	Directory string `json:"directory"`
	// Location is the location of the ACME account.
	Location string `json:"location"`
	// TOS is the terms of service URL.
	TOS string `json:"tos"`
}

// ACMEAccountCreateRequestBody contains the body for creating a new ACME account.
type ACMEAccountCreateRequestBody struct {
	// Contact is the contact email addresses.
	Contact string `url:"contact"`
	// Directory is the URL of the ACME CA directory endpoint.
	Directory string `url:"directory,omitempty"`
	// EABHMACKey is the HMAC key for External Account Binding.
	EABHMACKey string `url:"eab-hmac-key,omitempty"`
	// EABKID is the Key Identifier for External Account Binding.
	EABKID string `url:"eab-kid,omitempty"`
	// Name is the ACME account config file name.
	Name string `url:"name,omitempty"`
	// TOS is the URL of CA TermsOfService - setting this indicates agreement.
	TOS string `url:"tos_url,omitempty"`
}

// ACMEAccountCreateResponseBody contains the body from an ACME account create request.
type ACMEAccountCreateResponseBody struct {
	Data *string `json:"data,omitempty"`
}

// ACMEAccountUpdateRequestBody contains the body for updating an existing ACME account.
type ACMEAccountUpdateRequestBody struct {
	// Contact is the contact email addresses.
	Contact string `url:"contact,omitempty"`
	// Name is the ACME account config file name.
	Name string `url:"name,omitempty"`
}

// ACMEAccountUpdateResponseBody contains the body from an ACME account update request.
type ACMEAccountUpdateResponseBody struct {
	Data *string `json:"data,omitempty"`
}

// ACMEAccountDeleteResponseBody contains the body from an ACME account delete request.
type ACMEAccountDeleteResponseBody struct {
	Data *string `json:"data,omitempty"`
}
