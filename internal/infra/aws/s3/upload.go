package s3

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadObject(bucketName, objKey, filePath string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objKey),
		Body:   file,
	})
	if err != nil {
		return err
	}

	err = s3.NewObjectExistsWaiter(client).Wait(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objKey),
	}, time.Minute)
	if err != nil {
		return err
	}

	return nil
}
