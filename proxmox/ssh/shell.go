/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package ssh

import "strings"

// ShellQuote returns s wrapped in single quotes so that a POSIX shell treats it as a single literal word.
func ShellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, `'`, `'"'"'`) + "'"
}
