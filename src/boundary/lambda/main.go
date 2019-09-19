package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/webdevwilson/serverless-multicloud/src/boundary"
	"log"
	"net/http"
	"net/url"
)

func HandleRequest(ctx context.Context, apiGatewayReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("[INFO] HTTP PATH: %s", apiGatewayReq.Path)

	// Log the request event
	evt, err := json.Marshal(apiGatewayReq)
	if err != nil {
		log.Printf("[WARN] Error marshalling event")
	}
	log.Printf("[INFO] Request Event: %s", evt)

	// Create the handler
	pathPrefix := fmt.Sprintf("/%s", apiGatewayReq.RequestContext.Stage)
	handler := boundary.NewRequestHandler(pathPrefix)

	// Create an HTTP request from the event
	req, err := getRequestFromLambda(apiGatewayReq)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	resp := newLambdaResponseWriter()

	handler.ServeHTTP(resp, req)

	// Return the lambda response event
	response := events.APIGatewayProxyResponse{
		StatusCode:	resp.statusCode,
		Body:		string(resp.body),
		Headers: 	resp.flattenHeaders(),
	}

	// Log the response event
	evt, err = json.Marshal(response)
	if err != nil {
		log.Printf("[WARN] Error marshalling event %s", err)
	}
	log.Printf("[INFO] Response Event: %s", evt)

	return response, nil
}

func getRequestFromLambda(req events.APIGatewayProxyRequest) (*http.Request, error) {

	var path string
	if req.RequestContext.Stage == "" {
		path = req.Path
	} else {
		path = fmt.Sprintf("/%s%s", req.RequestContext.Stage, req.Path)
	}
	url, err := url.Parse(path)

	if err != nil {
		return nil, err
	}

	return &http.Request{
		Method:           req.HTTPMethod,
		URL:              url,
	}, nil
}

type lambdaResponseWriter struct {
	statusCode int
	headers http.Header
	body []byte
}

// type checking lambdaResponseWriter
var _ http.ResponseWriter = &lambdaResponseWriter{}

func newLambdaResponseWriter() *lambdaResponseWriter {
	return &lambdaResponseWriter{
		statusCode: -1,
		headers: make(map[string][]string),
		body: nil,
	}
}

func (l *lambdaResponseWriter) Header() http.Header {
	return l.headers
}

func (l *lambdaResponseWriter) Write(data []byte) (int, error) {
	// mimic standard net/http lib from documentation
	l.Header().Set("Content-Type", http.DetectContentType(data))
	if l.statusCode == -1 {
		l.WriteHeader(http.StatusOK)
	}
	l.body = data
	return len(data), nil
}

func (l *lambdaResponseWriter) WriteHeader(statusCode int) {
	l.statusCode = statusCode
}

// Lambda API does not allow for multiple headers like go standard lib
func (l *lambdaResponseWriter) flattenHeaders() map[string]string {
	var headers = make(map[string]string)
	for k, v := range l.headers {
		headers[k] = v[0]
	}
	return headers
}

func main() {
	lambda.Start(HandleRequest)
}