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
	XMLName        xml.Name       `xml:"DocuSignEnvelopeInformation" json:"-"`
	EnvelopeStatus EnvelopeStatus `xml:"EnvelopeStatus"`
	DocumentPDFs   DocumentPDFs   `xml:"DocumentPDFs"`
	TimeZone       string
	TimeZoneOffset string
}

type EnvelopeStatus struct {
	XMLName           xml.Name `xml:"EnvelopeStatus" json:"-"`
	TimeGenerated     string   `xml:"TimeGenerated"`
	EnvelopeID        string
	Subject           string
	UserName          string
	Email             string
	Status            string
	Created           string
	Sent              string
	Delivered         string
	Signed            string
	Completed         string
	SigningLocation   string
	SenderIPAddress   string
	RecipientStatuses RecipientStatuses `xml:"RecipientStatuses"`
}

type DocumentPDFs struct {
	XMLName      xml.Name      `xml:"DocumentPDFs" json:"-"`
	DocumentPDFs []DocumentPDF `xml:"DocumentPDF"`
}

type DocumentPDF struct {
	XMLName      xml.Name `xml:"DocumentPDF" json:"-"`
	Name         string
	DocumentType string
}

type RecipientStatuses struct {
	XMLName           xml.Name          `xml:"RecipientStatuses" json:"-"`
	RecipientStatuses []RecipientStatus `xml:"RecipientStatus"`
}

type RecipientStatus struct {
	XMLName            xml.Name `xml:"RecipientStatus" json:"-"`
	Type               string
	EMail              string
	UserName           string
	RoutingOrder       string
	Sent               string
	Delivered          string
	Signed             string
	DeclineReason      string
	Status             string
	RecipientIPAddress string
	TabStatuses TabStatuses `xml:"TabStatuses"`
}

type TabStatuses struct {
	XMLName xml.Name `xml:"TabStatuses" json:"-"`
	TabStatuses []TabStatus `xml:"TabStatus"`
}

type TabStatus struct {
	XMLName xml.Name `xml:"TabStatus" json:"-"`
	TabType string
	Status string
	XPosition int
	YPosition int
	Signed string
	TabLabel string
	TabName string
	TabValue string
	DocumentID int
	PageNumber int
	OriginalValue string
	ValidationPattern string
	RoleName string
	ListValues string
	ListSelectedValue string
	ScaleValue float64
}

func main() {
	handler := makeHandlerfunc()
	lambda.Start(handler)
}
