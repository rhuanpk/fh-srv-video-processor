package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func getMessages(client *sqs.Client, queueURL string) ([]types.Message, error) {
	out, err := client.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: 10,
	})
	if err != nil {
		return nil, err
	}

	return out.Messages, nil
}

func ReceiveMessages(queueName string) (string, []types.Message, error) {
	client, err := getClient()
	if err != nil {
		return "", nil, err
	}

	queueURL, err := getQueueURL(client, queueName)
	if err != nil {
		return "", nil, err
	}

	messages, err := getMessages(client, queueURL)
	if err != nil {
		return "", nil, err
	}

	return queueURL, messages, nil
}

func DeleteMessage(queueURL, receiptHandle string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	if _, err := client.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	}); err != nil {
		return err
	}

	return nil
}
