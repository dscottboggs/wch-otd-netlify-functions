/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package wch_otd_netlify_functions

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func BadRequest(message string) *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{StatusCode: 400, Body: message}
}

func InternalServerError(message string, err error) *events.APIGatewayProxyResponse {
	fmt.Printf("error %s: %v\n", message, err)
	return &events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       `{"error": "internal server error"}`,
	}
}
