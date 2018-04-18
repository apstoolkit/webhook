package awsctx

import "github.com/aws/aws-sdk-go/service/sqs/sqsiface"

type AWSContext struct {
	SQSSvc sqsiface.SQSAPI
}
