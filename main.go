package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"extractor/internal/config"
	"extractor/internal/infra/aws/s3"
	"extractor/internal/infra/aws/sns"
	"extractor/internal/infra/aws/sqs"
	video "extractor/internal/resouce/video/handler"
	zipper "extractor/internal/resouce/zipper/handler"
	"extractor/pkg/request"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func main() {
	service := video.NewService(zipper.NewService())

	for {
		log.Println("pulling for aws sqs messages")

		queueURL, messages, err := sqs.ReceiveMessages(config.AWSSQSQueueName)
		if err != nil {
			log.Println("error in receive messages:", err)
		}

		if len(messages) > 0 {
			log.Println("retrieved", len(messages), "messages")
		}

	messagesLoop:
		for _, message := range messages {
			var body sqs.Event
			if err := json.Unmarshal([]byte(*message.Body), &body); err != nil {
				log.Println("error in unmarshal message:", err)
				continue
			}

			if err := os.MkdirAll(config.ExtractorFolderTmp, 0775); err != nil {
				log.Println("error in create base folder:", err)
				continue
			}

			for _, record := range body.Records {
				objMetadata, videosPaths, err := s3.DownloadObjects(record.S3.Bucket.Name, record.S3.Object.Key)
				if err != nil {
					log.Println("error in dowload object:", err)
					continue messagesLoop
				}
				if len(videosPaths) <= 0 {
					if err := sqs.DeleteMessage(queueURL, aws.ToString(message.ReceiptHandle)); err != nil {
						log.Println("error in delete message:", err)
						continue
					}
					log.Println("delete message:", aws.ToString(message.MessageId))
					continue messagesLoop
				}
				log.Println("download objects:", strings.TrimPrefix(record.S3.Object.Key, "videos/"))

				if err := request.Post(config.APISrvStatusURL, "application/json", map[string]any{
					"id":     objMetadata.VideoID,
					"status": "EM_PROCESSAMENTO",
				}); err != nil {
					log.Println("error in request status service:", err)
				}

				zipsPaths, err := service.Process(videosPaths, config.FrameInterval, config.FrameHighQuality)
				if err != nil {
					log.Println("error in process videos:", err)
					continue messagesLoop
				}

				var objPublicLink string
				for _, zipPath := range zipsPaths {
					unescapedKey, err := url.PathUnescape(record.S3.Object.Key)
					if err != nil {
						log.Println("error in unescape object key:", err)
						continue messagesLoop
					}
					zipPathBase := filepath.Base(zipPath)
					s3FileName := filepath.Join(filepath.Dir(unescapedKey), zipPathBase)

					objPublicLink, err = s3.UploadObject(record.S3.Bucket.Name, s3FileName, zipPath)
					if err != nil {
						log.Println("error in upload object:", err)
						continue messagesLoop
					}
					log.Println("upload object:", zipPathBase)
					log.Println("object public link:", objPublicLink)
				}

				if err := request.Post(config.APISrvStatusURL, "application/json", map[string]any{
					"id":     objMetadata.VideoID,
					"status": "FINALIZADO",
				}); err != nil {
					log.Println("error in request status service:", err)
				}

				snsTopicID, err := sns.Publish(regexp.MustCompile(`[[:punct:]]`).ReplaceAllString(objMetadata.UserEmail, "_"), objPublicLink)
				if err != nil {
					log.Println("error in publish sns topic:", err)
					continue messagesLoop
				}

				log.Println("send sns topic message:", snsTopicID)
			}

			if err := os.RemoveAll(config.ExtractorFolderTmp); err != nil {
				log.Println("error in remove tmp folder:", err)
			}

			if err := sqs.DeleteMessage(queueURL, aws.ToString(message.ReceiptHandle)); err != nil {
				log.Println("error in delete message:", err)
				continue
			}
			log.Println("delete message:", aws.ToString(message.MessageId))
		}

		time.Sleep(time.Second * 3)
	}
}
