package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"log"
)

func processRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("processRequest called", request.Body)
	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func makeHandlerfunc()func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return processRequest(request)
	}
}


func main() {
	handler := makeHandlerfunc()
	lambda.Start(handler)
}
