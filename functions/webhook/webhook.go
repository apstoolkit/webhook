package main

import (
	"encoding/xml"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

func processRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("processRequest called", request.Body)
	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func makeHandlerfunc() func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return processRequest(request)
	}
}

type DocuSignEnvelopeInformation struct {
	XMLName        xml.Name       `xml:"DocuSignEnvelopeInformation"`
	EnvelopeStatus EnvelopeStatus `xml:"EnvelopeStatus"`
	TimeZone       string
	TimeZoneOffset string
}

type EnvelopeStatus struct {
	XMLName         xml.Name `xml:"EnvelopeStatus"`
	TimeGenerated   string   `xml:"TimeGenerated"`
	EnvelopeID      string
	Subject         string
	UserName        string
	Email           string
	Status          string
	Created         string
	Sent            string
	Delivered       string
	Signed          string
	Completed       string
	SigningLocation string
	SenderIPAddress string
}

func main() {
	handler := makeHandlerfunc()
	lambda.Start(handler)
}
