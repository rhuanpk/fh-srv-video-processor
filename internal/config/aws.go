package config

import "os"

var (
	AWSRegion       = os.Getenv("AWS_DEFAULT_REGION")
	AWSSQSQueueName = os.Getenv("AWS_SQS_QUEUE_NAME")
	AWSSNSArnPrefix = os.Getenv("AWS_SNS_ARN_PREFIX")
)
