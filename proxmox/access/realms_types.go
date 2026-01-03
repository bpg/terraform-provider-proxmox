/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"github.com/bpg/terraform-provider-proxmox/proxmox/types"
)

// RealmCreateRequestBody contains the data for creating a realm.
type RealmCreateRequestBody struct {
	Realm string `json:"realm" url:"realm"`
	Type  string `json:"type"  url:"type"` // "ldap", "ad", "openid", "pam", "pve"

	// LDAP/AD specific
	Server1          *string           `json:"server1,omitempty"               url:"server1,omitempty"`               // Primary LDAP server
	Server2          *string           `json:"server2,omitempty"               url:"server2,omitempty"`               // Fallback LDAP server
	Port             *int              `json:"port,omitempty"                  url:"port,omitempty"`                  // Default: 389 (636 for ldaps)
	BaseDN           *string           `json:"base_dn,omitempty"               url:"base_dn,omitempty"`               // LDAP base DN
	BindDN           *string           `json:"bind_dn,omitempty"               url:"bind_dn,omitempty"`               // LDAP bind DN
	BindPassword     *string           `json:"password,omitempty"              url:"password,omitempty"`              // Bind password
	UserAttr         *string           `json:"user_attr,omitempty"             url:"user_attr,omitempty"`             // Default: "uid"
	Secure           *types.CustomBool `json:"secure,omitempty"                url:"secure,omitempty,int"`            // Use LDAPS
	Verify           *types.CustomBool `json:"verify,omitempty"                url:"verify,omitempty,int"`            // Verify SSL cert
	CaPath           *string           `json:"capath,omitempty"                url:"capath,omitempty"`                // CA cert path
	CertPath         *string           `json:"cert,omitempty"                  url:"cert,omitempty"`                  // Client cert
	CertKeyPath      *string           `json:"certkey,omitempty"               url:"certkey,omitempty"`               // Client key
	Filter           *string           `json:"filter,omitempty"                url:"filter,omitempty"`                // LDAP filter
	GroupDN          *string           `json:"group_dn,omitempty"              url:"group_dn,omitempty"`              // Group base DN
	GroupFilter      *string           `json:"group_filter,omitempty"          url:"group_filter,omitempty"`          // Group filter
	GroupClasses     *string           `json:"group_classes,omitempty"         url:"group_classes,omitempty"`         // Group objectClasses
	GroupNameAttr    *string           `json:"group_name_attr,omitempty"       url:"group_name_attr,omitempty"`       // Group name attribute
	Mode             *string           `json:"mode,omitempty"                  url:"mode,omitempty"`                  // LDAP mode: ldap, ldaps, ldap+starttls
	SSLVersion       *string           `json:"sslversion,omitempty"            url:"sslversion,omitempty"`            // SSL/TLS version
	UserClasses      *string           `json:"user_classes,omitempty"          url:"user_classes,omitempty"`          // User objectClasses
	SyncAttributes   *string           `json:"sync_attributes,omitempty"       url:"sync_attributes,omitempty"`       // Attributes to sync
	SyncDefaultsOpts *string           `json:"sync-defaults-options,omitempty" url:"sync-defaults-options,omitempty"` // Default sync options

	// OpenID specific
	ClientID      *string           `json:"client-id,omitempty"      url:"client-id,omitempty"`      // OpenID client ID
	ClientKey     *string           `json:"client-key,omitempty"     url:"client-key,omitempty"`     // OpenID client secret
	IssuerURL     *string           `json:"issuer-url,omitempty"     url:"issuer-url,omitempty"`     // OpenID issuer URL
	UsernameClaim *string           `json:"username-claim,omitempty" url:"username-claim,omitempty"` // Username claim
	Scopes        *string           `json:"scopes,omitempty"         url:"scopes,omitempty"`         // OAuth scopes
	Prompt        *string           `json:"prompt,omitempty"         url:"prompt,omitempty"`         // OAuth prompt
	ACRValues     *string           `json:"acr-values,omitempty"     url:"acr-values,omitempty"`     // ACR values
	AutoCreate    *types.CustomBool `json:"autocreate,omitempty"     url:"autocreate,omitempty,int"` // Auto-create users

	// General options
	Comment       *string           `json:"comment,omitempty"        url:"comment,omitempty"`            // Description
	Default       *types.CustomBool `json:"default,omitempty"        url:"default,omitempty,int"`        // Default realm
	TFA           *string           `json:"tfa,omitempty"            url:"tfa,omitempty"`                // TFA config
	CaseSensitive *types.CustomBool `json:"case-sensitive,omitempty" url:"case-sensitive,omitempty,int"` // Case-sensitive usernames
	Domain        *string           `json:"domain,omitempty"         url:"domain,omitempty"`             // AD domain
}

// RealmUpdateRequestBody contains the data for updating a realm.
type RealmUpdateRequestBody struct {
	// Note: Realm and Type are not included (part of URL path)
	Server1          *string           `json:"server1,omitempty"               url:"server1,omitempty"`
	Server2          *string           `json:"server2,omitempty"               url:"server2,omitempty"`
	Port             *int              `json:"port,omitempty"                  url:"port,omitempty,int"`
	BaseDN           *string           `json:"base_dn,omitempty"               url:"base_dn,omitempty"`
	BindDN           *string           `json:"bind_dn,omitempty"               url:"bind_dn,omitempty"`
	BindPassword     *string           `json:"password,omitempty"              url:"password,omitempty"`
	UserAttr         *string           `json:"user_attr,omitempty"             url:"user_attr,omitempty"`
	Secure           *types.CustomBool `json:"secure,omitempty"                url:"secure,omitempty,int"`
	Verify           *types.CustomBool `json:"verify,omitempty"                url:"verify,omitempty,int"`
	CaPath           *string           `json:"capath,omitempty"                url:"capath,omitempty"`
	CertPath         *string           `json:"cert,omitempty"                  url:"cert,omitempty"`
	CertKeyPath      *string           `json:"certkey,omitempty"               url:"certkey,omitempty"`
	Filter           *string           `json:"filter,omitempty"                url:"filter,omitempty"`
	GroupDN          *string           `json:"group_dn,omitempty"              url:"group_dn,omitempty"`
	GroupFilter      *string           `json:"group_filter,omitempty"          url:"group_filter,omitempty"`
	GroupClasses     *string           `json:"group_classes,omitempty"         url:"group_classes,omitempty"`
	GroupNameAttr    *string           `json:"group_name_attr,omitempty"       url:"group_name_attr,omitempty"`
	Mode             *string           `json:"mode,omitempty"                  url:"mode,omitempty"`
	SSLVersion       *string           `json:"sslversion,omitempty"            url:"sslversion,omitempty"`
	UserClasses      *string           `json:"user_classes,omitempty"          url:"user_classes,omitempty"`
	SyncAttributes   *string           `json:"sync_attributes,omitempty"       url:"sync_attributes,omitempty"`
	SyncDefaultsOpts *string           `json:"sync-defaults-options,omitempty" url:"sync-defaults-options,omitempty"`
	ClientID         *string           `json:"client-id,omitempty"             url:"client-id,omitempty"`
	ClientKey        *string           `json:"client-key,omitempty"            url:"client-key,omitempty"`
	IssuerURL        *string           `json:"issuer-url,omitempty"            url:"issuer-url,omitempty"`
	UsernameClaim    *string           `json:"username-claim,omitempty"        url:"username-claim,omitempty"`
	Scopes           *string           `json:"scopes,omitempty"                url:"scopes,omitempty"`
	Prompt           *string           `json:"prompt,omitempty"                url:"prompt,omitempty"`
	ACRValues        *string           `json:"acr-values,omitempty"            url:"acr-values,omitempty"`
	AutoCreate       *types.CustomBool `json:"autocreate,omitempty"            url:"autocreate,omitempty,int"`
	Comment          *string           `json:"comment,omitempty"               url:"comment,omitempty"`
	Default          *types.CustomBool `json:"default,omitempty"               url:"default,omitempty,int"`
	TFA              *string           `json:"tfa,omitempty"                   url:"tfa,omitempty"`
	CaseSensitive    *types.CustomBool `json:"case-sensitive,omitempty"        url:"case-sensitive,omitempty,int"`
	Domain           *string           `json:"domain,omitempty"                url:"domain,omitempty"`
	Delete           *string           `json:"delete,omitempty"                url:"delete,omitempty"` // Comma-separated list of properties to delete
}

// RealmGetResponseBody contains the response for GET /access/domains/{realm}.
type RealmGetResponseBody struct {
	Data *RealmGetResponseData `json:"data,omitempty"`
}

// RealmGetResponseData contains realm configuration data.
type RealmGetResponseData struct {
	Realm            string            `json:"realm"`
	Type             string            `json:"type"`
	Server1          *string           `json:"server1,omitempty"`
	Server2          *string           `json:"server2,omitempty"`
	Port             *int              `json:"port,omitempty"`
	BaseDN           *string           `json:"base_dn,omitempty"`
	BindDN           *string           `json:"bind_dn,omitempty"`
	UserAttr         *string           `json:"user_attr,omitempty"`
	Secure           *types.CustomBool `json:"secure,omitempty"`
	Verify           *types.CustomBool `json:"verify,omitempty"`
	CaPath           *string           `json:"capath,omitempty"`
	CertPath         *string           `json:"cert,omitempty"`
	CertKeyPath      *string           `json:"certkey,omitempty"`
	Filter           *string           `json:"filter,omitempty"`
	GroupDN          *string           `json:"group_dn,omitempty"`
	GroupFilter      *string           `json:"group_filter,omitempty"`
	GroupClasses     *string           `json:"group_classes,omitempty"`
	GroupNameAttr    *string           `json:"group_name_attr,omitempty"`
	Mode             *string           `json:"mode,omitempty"`
	SSLVersion       *string           `json:"sslversion,omitempty"`
	UserClasses      *string           `json:"user_classes,omitempty"`
	SyncAttributes   *string           `json:"sync_attributes,omitempty"`
	SyncDefaultsOpts *string           `json:"sync-defaults-options,omitempty"`
	ClientID         *string           `json:"client-id,omitempty"`
	ClientKey        *string           `json:"client-key,omitempty"`
	IssuerURL        *string           `json:"issuer-url,omitempty"`
	UsernameClaim    *string           `json:"username-claim,omitempty"`
	Scopes           *string           `json:"scopes,omitempty"`
	Prompt           *string           `json:"prompt,omitempty"`
	ACRValues        *string           `json:"acr-values,omitempty"`
	AutoCreate       *types.CustomBool `json:"autocreate,omitempty"`
	Comment          *string           `json:"comment,omitempty"`
	Default          *types.CustomBool `json:"default,omitempty"`
	TFA              *string           `json:"tfa,omitempty"`
	CaseSensitive    *types.CustomBool `json:"case-sensitive,omitempty"`
	Domain           *string           `json:"domain,omitempty"`
	Digest           *string           `json:"digest,omitempty"` // Config digest for change detection
}

// RealmListResponseBody contains the response for GET /access/domains.
type RealmListResponseBody struct {
	Data []*RealmListResponseData `json:"data,omitempty"`
}

// RealmListResponseData contains a realm list entry.
type RealmListResponseData struct {
	Realm   string  `json:"realm"`
	Type    string  `json:"type"`
	Comment *string `json:"comment,omitempty"`
	TFA     *string `json:"tfa,omitempty"`
}

// RealmSyncRequestBody contains the request for POST /access/domains/{realm}/sync.
type RealmSyncRequestBody struct {
	Scope          *string           `json:"scope,omitempty"           url:"scope,omitempty"`           // "users", "groups", "both"
	RemoveVanished *string           `json:"remove-vanished,omitempty" url:"remove-vanished,omitempty"` // "acl", "properties", "entry"
	EnableNew      *types.CustomBool `json:"enable-new,omitempty"      url:"enable-new,omitempty,int"`  // Enable new users
	Full           *types.CustomBool `json:"full,omitempty"            url:"full,omitempty,int"`        // Full sync
	Purge          *types.CustomBool `json:"purge,omitempty"           url:"purge,omitempty,int"`       // Purge removed entries
	DryRun         *types.CustomBool `json:"dry-run,omitempty"         url:"dry-run,omitempty,int"`     // Test mode
}
