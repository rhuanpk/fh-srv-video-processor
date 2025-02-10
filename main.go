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

	"github.com/aws/aws-sdk-go-v2/aws"
)

func main() {
	service := video.NewService(zipper.NewService())

	// REMOVE
	config.AWSRegion = "us-east-1"
	config.AWSSQSQueueName = "hack-sqs-queue"
	config.AWSSNSArnPrefix = "arn:aws:sns:us-east-1:086085977965"

	for {
		log.Println("pulling for aws sqs messages")

		queueURL, messages, err := sqs.ReceiveMessages(config.AWSSQSQueueName)
		if err != nil {
			log.Println("error in receive messages:", err)
		}

		if len(messages) > 0 {
			log.Println("retrieved", len(messages), "messages")
		}

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
				videosPaths, err := s3.DownloadObjects(record.S3.Bucket.Name, record.S3.Object.Key)
				if err != nil {
					log.Println("error in dowload object:", err)
					continue
				}
				if len(videosPaths) <= 0 {
					continue
				}
				log.Println("download objects:", strings.TrimPrefix(record.S3.Object.Key, "videos/"))

				// call srv status antes

				zipsPaths, err := service.Process(videosPaths, config.FrameInterval, config.FrameHighQuality)
				if err != nil {
					log.Println("error in process videos:", err)
					continue
				}

				var objPublicLink string
				for _, zipPath := range zipsPaths {
					unescapedKey, err := url.PathUnescape(record.S3.Object.Key)
					if err != nil {
						log.Println("error in unescape object key:", err)
						continue
					}
					zipPathBase := filepath.Base(zipPath)
					s3FileName := filepath.Join(filepath.Dir(unescapedKey), zipPathBase)

					objPublicLink, err = s3.UploadObject(record.S3.Bucket.Name, s3FileName, zipPath)
					if err != nil {
						log.Println("error in upload object:", err)
						continue
					}
					log.Println("upload object:", zipPathBase)
					log.Println("object public link:", objPublicLink)
				}

				// call srv status depois

				email := regexp.MustCompile(`=(.*)$`).FindStringSubmatch(record.UserIdentity.PrincipalID)[1]
				snsTopicID, err := sns.Publish(regexp.MustCompile(`[[:punct:]]`).ReplaceAllString(email, "_"), objPublicLink)
				if err != nil {
					log.Println("error in publish sns topic:", err)
					continue
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
