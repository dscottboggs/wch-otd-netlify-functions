/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * This file contains some convenient constructor functions for AWS's
 * *events.APIGatewayProxyResponse.
 */

package wch_otd_api

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func BadRequest(message string) *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{StatusCode: 400, Body: message}
}

/**
 * Return a 500-status error response, logging the given error but not
 * transmitting information about it. Basically this gives us a quick way to
 * handle an error which doesn't offer any risk of exposing sensitive
 * information to the client.
 */
func InternalServerError(message string, err error) *events.APIGatewayProxyResponse {
	fmt.Printf("error %s: %v\n", message, err)
	return &events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       `{"error": "internal server error"}`,
	}
}
