package config

import "os"

var (
	AWSRegion       string = os.Getenv("AWS_DEFAULT_REGION")
	AWSSQSQueueName string = os.Getenv("AWS_SQS_QUEUE_NAME")
	AWSSNSArnPrefix string = os.Getenv("AWS_SNS_ARN_PREFIX")
)
