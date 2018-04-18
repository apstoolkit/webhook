package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"encoding/xml"
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

type DocuSignEnvelopeInformation struct {
	XMLName xml.Name `xml:"DocuSignEnvelopeInformation"`
	EnvelopeStatus EnvelopeStatus `xml:"EnvelopeStatus"`
}

type EnvelopeStatus struct {
	XMLName xml.Name `xml:"EnvelopeStatus"`
	TimeGenerated string `xml:"TimeGenerated"`
	EnvelopeID string
	Subject string
}


func main() {
	handler := makeHandlerfunc()
	lambda.Start(handler)
}
