/*
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package hardwaremapping

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

// EncodeText percent-encodes control characters, non-ASCII bytes, and the percent sign in free-form text.
// The Proxmox VE cluster mapping API stores descriptions raw and returns high-bit bytes double-encoded, so
// only printable ASCII survives a round-trip. This mirrors the encoding Proxmox VE itself applies to VM
// descriptions via PVE::Tools::encode_text.
func EncodeText(text string) string {
	var b strings.Builder

	for i := range len(text) {
		c := text[i]
		if c < 0x20 || c > 0x7e || c == '%' {
			fmt.Fprintf(&b, "%%%02X", c)
		} else {
			b.WriteByte(c)
		}
	}

	return b.String()
}

// DecodeText reverses EncodeText. Values not written by this provider may contain stray percent signs or
// invalid UTF-8 sequences; those are returned unchanged.
func DecodeText(text string) string {
	decoded, err := url.PathUnescape(text)
	if err != nil || !utf8.ValidString(decoded) {
		return text
	}

	return decoded
}
