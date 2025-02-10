package s3

import (
	"context"
	"extractor/internal/config"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadObject(bucketName, objKey, filePath string) (string, error) {
	client, err := getClient()
	if err != nil {
		return "", err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objKey),
		Body:   file,
		// ACL:    types.ObjectCannedACLPublicRead,
	}); err != nil {
		return "", err
	}

	if err = s3.NewObjectExistsWaiter(client).Wait(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objKey),
	}, time.Minute); err != nil {
		return "", err
	}

	publicLink := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, config.AWSRegion, objKey)
	return publicLink, nil
}
