package wch_otd_api

import "github.com/aws/aws-lambda-go/events"

func CheckForPrefetch(request events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
	if request.HTTPMethod == "HEAD" {
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
		"Access-Control-Allow-Methods": "GET",
	}
}
