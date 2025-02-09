package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func getQueueURL(client sqs.ListQueuesAPIClient, queueName string) (string, error) {
	paginator := sqs.NewListQueuesPaginator(client, &sqs.ListQueuesInput{
		QueueNamePrefix: aws.String(queueName),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return "", err
		}

		return page.QueueUrls[0], nil
	}

	return "", nil
}
