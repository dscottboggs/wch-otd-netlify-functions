/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package wch_otd_api

import (
	"testing"

	"github.com/dscottboggs/attest"
)

func TestValidateTZ(t *testing.T) {
	test := attest.New(t)
	tz, err := ValidateTZ("America/New_York")
	if err != nil {
		test.Logf("%#+v was not nil", err)
		test.Fail()
	}
	test.Equals(tz, "America/New_York")
	// longest entry in my tz db
	tz, err = ValidateTZ("posix/America/Argentina/ComodRivadavia")
	if err != nil {
		test.Logf("%#+v was not nil", err)
		test.Fail()
	}
	test.Equals(tz, "posix/America/Argentina/ComodRivadavia")
	tz, err = ValidateTZ("this should be invalid because it is too long")
	test.Equals("", tz)
	test.NotNil(err, "error for too-long time zone was nil")
	test.Equals("{\"error\":\"invalid time zone\"}\n", err.Body)
	test.Equals(err.StatusCode, 400)
	tz, err = ValidateTZ("invalid b/c spaces")
	test.NotNil(err, "error for time zone with spaces was nil")
	test.Equals("", tz)
	test.Equals("{\"error\":\"invalid time zone\"}\n", err.Body)
	test.Equals(err.StatusCode, 400)

}
