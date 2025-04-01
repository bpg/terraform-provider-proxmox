/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

// RealmCreateRequestBody contains the data for a realm create request.
type RealmCreateRequestBody struct {
	Realm               string `json:"realm" url:"realm"`
	Type                string `json:"type" url:"type"` /*enum*/
	Acrvalues           string `json:"acr-values,omitempty" url:"acr-values,omitempty"`
	Autocreate          bool   `json:"autocreate,omitempty" url:"autocreate,omitempty,int"`
	Base_dn             string `json:"base_dn,omitempty" url:"base_dn,omitempty"`
	Bind_dn             string `json:"bind_dn,omitempty" url:"bind_dn,omitempty"`
	Capath              string `json:"capath,omitempty" url:"capath,omitempty"`
	Casesensitive       bool   `json:"case-sensitive,omitempty" url:"case-sensitive,omitempty,int"`
	Cert                string `json:"cert,omitempty" url:"cert,omitempty"`
	Certkey             string `json:"certkey,omitempty" url:"certkey,omitempty"`
	Checkconnection     bool   `json:"check-connection,omitempty" url:"check-connection,omitempty,int"`
	Clientid            string `json:"client-id,omitempty" url:"client-id,omitempty"`
	Clientkey           string `json:"client-key,omitempty" url:"client-key,omitempty"`
	Comment             string `json:"comment,omitempty" url:"comment,omitempty"`
	Default             bool   `json:"default,omitempty" url:"default,omitempty,int"`
	Domain              string `json:"domain,omitempty" url:"domain,omitempty"`
	Filter              string `json:"filter,omitempty" url:"filter,omitempty"`
	Group_classes       string `json:"group_classes,omitempty" url:"group_classes,omitempty"`
	Group_dn            string `json:"group_dn,omitempty" url:"group_dn,omitempty"`
	Group_filter        string `json:"group_filter,omitempty" url:"group_filter,omitempty"`
	Group_name_attr     string `json:"group_name_attr,omitempty" url:"group_name_attr,omitempty"`
	Issuerurl           string `json:"issuer-url,omitempty" url:"issuer-url,omitempty"`
	Mode                string `json:"mode,omitempty" url:"mode,omitempty"` /*enum*/
	Password            string `json:"password,omitempty" url:"password,omitempty"`
	Port                int    `json:"port,omitempty" url:"port,omitempty,int"`
	Prompt              string `json:"prompt,omitempty" url:"prompt,omitempty"`
	Scopes              string `json:"scopes,omitempty" url:"scopes,omitempty"`
	Secure              bool   `json:"secure,omitempty" url:"secure,omitempty,int"`
	Server1             string `json:"server1,omitempty" url:"server1,omitempty"`
	Server2             string `json:"server2,omitempty" url:"server2,omitempty"`
	Sslversion          string `json:"sslversion,omitempty" url:"sslversion,omitempty"` /*enum*/
	Syncdefaultsoptions string `json:"sync-defaults-options,omitempty" url:"sync-defaults-options,omitempty"`
	Sync_attributes     string `json:"sync_attributes,omitempty" url:"sync_attributes,omitempty"`
	Tfa                 string `json:"tfa,omitempty" url:"tfa,omitempty"` /*enum*/
	User_attr           string `json:"user_attr,omitempty" url:"user_attr,omitempty"`
	User_classes        string `json:"user_classes,omitempty" url:"user_classes,omitempty"`
	Usernameclaim       string `json:"username-claim,omitempty" url:"username-claim,omitempty"`
	Verify              bool   `json:"verify,omitempty" url:"verify,omitempty,int"`
}

// RealmGetResponseBody contains the body from an a realm get response.
type RealmGetResponseBody struct {
	Acrvalues           string `json:"acr-values,omitempty"`
	Autocreate          bool   `json:"autocreate,omitempty"`
	Base_dn             string `json:"base_dn,omitempty"`
	Bind_dn             string `json:"bind_dn,omitempty"`
	Capath              string `json:"capath,omitempty"`
	Casesensitive       bool   `json:"case-sensitive,omitempty"`
	Cert                string `json:"cert,omitempty"`
	Certkey             string `json:"certkey,omitempty"`
	Checkconnection     bool   `json:"check-connection,omitempty"`
	Clientid            string `json:"client-id,omitempty"`
	Clientkey           string `json:"client-key,omitempty"`
	Comment             string `json:"comment,omitempty"`
	Default             bool   `json:"default,omitempty"`
	Domain              string `json:"domain,omitempty"`
	Filter              string `json:"filter,omitempty"`
	Group_classes       string `json:"group_classes,omitempty"`
	Group_dn            string `json:"group_dn,omitempty"`
	Group_filter        string `json:"group_filter,omitempty"`
	Group_name_attr     string `json:"group_name_attr,omitempty"`
	Issuerurl           string `json:"issuer-url,omitempty"`
	Mode                string `json:"mode,omitempty"` /*enum*/
	Password            string `json:"password,omitempty"`
	Port                int    `json:"port,omitempty"`
	Prompt              string `json:"prompt,omitempty"`
	Scopes              string `json:"scopes,omitempty"`
	Secure              bool   `json:"secure,omitempty"`
	Server1             string `json:"server1,omitempty"`
	Server2             string `json:"server2,omitempty"`
	Sslversion          string `json:"sslversion,omitempty"` /*enum*/
	Syncdefaultsoptions string `json:"sync-defaults-options,omitempty"`
	Sync_attributes     string `json:"sync_attributes,omitempty"`
	Tfa                 string `json:"tfa,omitempty"` /*enum*/
	User_attr           string `json:"user_attr,omitempty"`
	User_classes        string `json:"user_classes,omitempty"`
	Usernameclaim       string `json:"username-claim,omitempty"`
	Verify              bool   `json:"verify,omitempty"`
}

// RealmListResponseBody contains the body from an a realm list response.
type RealmListResponseBody struct {
	Data []*RealmListResponseData `json:"data,omitempty"`
}

// RealmListResponseData contains the data from an a realm list response.
type RealmListResponseData struct {
	Realm   string `json:"realm" url:"realm"`
	Type    string `json:"type" url:"type"`
	Comment string `json:"comment,omitempty" url:"comment,omitempty"`
	Tfa     string `json:"tfa,omitempty" url:"tfa,omitempty"`
}

// RealmUpdateRequestBody contains the data for an a realm update request.
type RealmUpdateRequestBody struct {
	Acrvalues           string `json:"acr-values,omitempty" url:"acr-values,omitempty"`
	Autocreate          bool   `json:"autocreate,omitempty" url:"autocreate,omitempty,int"`
	Base_dn             string `json:"base_dn,omitempty" url:"base_dn,omitempty"`
	Bind_dn             string `json:"bind_dn,omitempty" url:"bind_dn,omitempty"`
	Capath              string `json:"capath,omitempty" url:"capath,omitempty"`
	Casesensitive       bool   `json:"case-sensitive,omitempty" url:"case-sensitive,omitempty,int"`
	Cert                string `json:"cert,omitempty" url:"cert,omitempty"`
	Certkey             string `json:"certkey,omitempty" url:"certkey,omitempty"`
	Checkconnection     bool   `json:"check-connection,omitempty" url:"check-connection,omitempty,int"`
	Clientid            string `json:"client-id,omitempty" url:"client-id,omitempty"`
	Clientkey           string `json:"client-key,omitempty" url:"client-key,omitempty"`
	Comment             string `json:"comment,omitempty" url:"comment,omitempty"`
	Default             bool   `json:"default,omitempty" url:"default,omitempty,int"`
	Domain              string `json:"domain,omitempty" url:"domain,omitempty"`
	Filter              string `json:"filter,omitempty" url:"filter,omitempty"`
	Group_classes       string `json:"group_classes,omitempty" url:"group_classes,omitempty"`
	Group_dn            string `json:"group_dn,omitempty" url:"group_dn,omitempty"`
	Group_filter        string `json:"group_filter,omitempty" url:"group_filter,omitempty"`
	Group_name_attr     string `json:"group_name_attr,omitempty" url:"group_name_attr,omitempty"`
	Issuerurl           string `json:"issuer-url,omitempty" url:"issuer-url,omitempty"`
	Mode                string `json:"mode,omitempty" url:"mode,omitempty"`
	Password            string `json:"password,omitempty" url:"password,omitempty"`
	Port                int    `json:"port,omitempty" url:"port,omitempty,int"`
	Prompt              string `json:"prompt,omitempty" url:"prompt,omitempty"`
	Scopes              string `json:"scopes,omitempty" url:"scopes,omitempty"`
	Secure              bool   `json:"secure,omitempty" url:"secure,omitempty,int"`
	Server1             string `json:"server1,omitempty" url:"server1,omitempty"`
	Server2             string `json:"server2,omitempty" url:"server2,omitempty"`
	Sslversion          string `json:"sslversion,omitempty" url:"sslversion,omitempty"`
	Syncdefaultsoptions string `json:"sync-defaults-options,omitempty" url:"sync-defaults-options,omitempty"`
	Sync_attributes     string `json:"sync_attributes,omitempty" url:"sync_attributes,omitempty"`
	Tfa                 string `json:"tfa,omitempty" url:"tfa,omitempty"`
	User_attr           string `json:"user_attr,omitempty" url:"user_attr,omitempty"`
	User_classes        string `json:"user_classes,omitempty" url:"user_classes,omitempty"`
	Usernameclaim       string `json:"username-claim,omitempty" url:"username-claim,omitempty"`
	Verify              bool   `json:"verify,omitempty" url:"verify,omitempty,int"`
}
