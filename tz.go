/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package wch_otd_api

import (
	"regexp"

	"github.com/aws/aws-lambda-go/events"
)

var validTzRegex *regexp.Regexp = regexp.MustCompile(`^(\w+/?)*\w+$`)

func ValidateTZ(untrusted string) (string, *events.APIGatewayProxyResponse) {
	if len(untrusted) <= 40 && validTzRegex.Match([]byte(untrusted)) {
		// validated
		return untrusted, nil
	} else {
		return "", BadRequest("invalid time zone")
	}
}
