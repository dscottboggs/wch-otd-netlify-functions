package wch_otd_api

import "github.com/aws/aws-lambda-go/events"

func CheckForPrefetch(request events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	if request.HTTPMethod == "OPTIONS" {
		/** TODO:
		 * it may be a good idea to log the User-Agent and preflight headers
		 * (https://developer.mozilla.org/en-US/docs/Glossary/Preflight_request)
		 * (especially origin) here.
		 */
		headers := CORSHeaders()
		headers["Content-Type"] = "application/json"
		return &events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    headers,
		}
	}
	return nil
}

func CORSHeaders() map[string]string {
	return map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, OPTIONS",
	}
}
