package wch_otd_netlify_functions

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
