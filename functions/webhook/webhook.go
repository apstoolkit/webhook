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
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/apstoolkit/webhook/awsctx"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"os"
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

func processRequest(awsContext *awsctx.AWSContext, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	sendMessageInput := sqs.SendMessageInput{
		MessageBody: aws.String(string(jsonBytes)),
		QueueUrl: aws.String(os.Getenv("SQS_URL")),
	}

	_, err = awsContext.SQSSvc.SendMessage(&sendMessageInput)
	if err != nil {
		log.Println("Error queuing status", err.Error())
		return events.APIGatewayProxyResponse{StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func makeHandlerfunc(awsContent *awsctx.AWSContext) func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return processRequest(awsContent,request)
	}
}



func main() {
	var awsContext awsctx.AWSContext

	sess := session.New()
	svc := sqs.New(sess)
	awsContext.SQSSvc = svc

	handler := makeHandlerfunc(&awsContext)
	lambda.Start(handler)
}
