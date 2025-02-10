package sns

import (
	"context"

	"extractor/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func Publish(email, objectURL string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	topicArn := config.AWSSNSArnPrefix + ":" + email
	message := "Seu vídeo foi processado e está disponível em: " + objectURL
	subject := "Processador de Vídeo"

	if _, err := client.Publish(context.Background(), &sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Message:  aws.String(message),
		Subject:  aws.String(subject),
	}); err != nil {
		return err
	}

	return nil
}
