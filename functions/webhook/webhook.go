package main

import (
	"encoding/xml"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
	"encoding/json"
	"strings"
	"errors"
	"strconv"
)

func ipOctects(ip string) ([]int, error) {
	ipparts := strings.Split(ip, ".")
	if len(ipparts) != 4 {
		return nil, errors.New("Split ip does not have 4 octets")
	}

	var octets []int

	for _, part := range ipparts {
		o, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}

		octets = append(octets,o)
	}

	return octets, nil
}

// This link documents the IPs to whitelist for Docusign callbacks
// https://trust.docusign.com/en-us/trust-certifications/whitelist/
// Note this set of whitelisted addresses are all IPv4

func isKnownSingleAddress(address string) bool {
	knownIps := []string{
		"54.149.21.90",
		"54.69.114.54",
		"52.25.122.31",
		"52.25.145.215",
		"52.26.192.160",
		"52.24.91.157",
		"52.27.126.9",
		"52.11.152.229",
	}

	for _, ip := range knownIps {
		if address == ip {
			return true
		}
	}

	return false
}

func isCallerDocusign(xForwardedFor string) (bool, error) {

	parts := strings.Split(xForwardedFor, ",")
	if len(parts) == 0 {
		return false, errors.New("Expected parts to x-forwarded-for")
	}

	//First x-forwarded-for address is the originating address
	if isKnownSingleAddress(parts[0]) {
		return true, nil
	}

	octets, err := ipOctects(parts[0])
	if err != nil {
		return false, err
	}

	if octets[0] == 162 && octets[1] == 248 {
		return (octets[2] >= 184 && octets[2] <= 187) && (octets[3] >= 1 && octets[3] <= 254), nil
	}

	return false, nil
}

func processRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("processRequest headers", request.Headers)

	xForwardedFor := request.Headers["X-Forwarded-For"]
	isDocusign, err := isCallerDocusign(xForwardedFor)
	if err != nil {
		log.Println(err.Error())
	}

	if isDocusign == false {
		log.Println("Callback from a non-docusign user")
		return events.APIGatewayProxyResponse{StatusCode: http.StatusForbidden}, nil
	}

	var envInfo DocuSignEnvelopeInformation
	err = xml.Unmarshal([]byte(request.Body), &envInfo)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}

	jsonBytes, err := json.Marshal(&envInfo)
	if err == nil {
		log.Println(string(jsonBytes))
	}

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
