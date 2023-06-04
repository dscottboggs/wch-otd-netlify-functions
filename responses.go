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
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func BadRequest(message string) *events.APIGatewayProxyResponse {
	responseData := struct {
		E string `json:"error"`
	}{E: message}
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(responseData)
	if err != nil {
		return InternalServerError("failed to json-encode response for bad request", err)
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: 400,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       buf.String(),
	}
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
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       `{"error": "internal server error"}`,
	}
}
