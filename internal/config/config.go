package config

import "os"

func init() {
	if AWSRegion = os.Getenv("AWS_REGION"); AWSRegion == "" {
		if AWSRegion = os.Getenv("AWS_DEFAULT_REGION"); AWSRegion == "" {
			AWSRegion = "us-east-1"
		}
	}
}
