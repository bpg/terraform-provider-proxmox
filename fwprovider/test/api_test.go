//go:build acceptance || all

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package test

import (
	"log"
	"testing"
)

func TestAPI(t *testing.T) {
	t.Parallel()

	te := InitEnvironment(t)

	acc, err := te.AccessClient().ListRealms(te.t.Context())
	if err != nil {
		log.Printf("%s", err)
	}

	for _, obj := range acc {
		log.Printf("%s", obj)
	}

	_ = acc

}
