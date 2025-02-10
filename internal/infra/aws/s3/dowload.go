package s3

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"extractor/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func listObjectsPaginator(client s3.ListObjectsV2APIClient, bucketName, pathSuffix string) ([]string, error) {
	var objKeys []string

	unescapedPath, err := url.PathUnescape(pathSuffix)
	if err != nil {
		return nil, err
	}

	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(unescapedPath),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}

		for _, content := range page.Contents {
			objKeys = append(objKeys, aws.ToString(content.Key))
		}
	}

	return objKeys, nil
}

func DownloadObjects(bucketName, pathSuffix string) ([]string, error) {
	var videosPaths []string

	client, err := getClient()
	if err != nil {
		return nil, err
	}

	objKeys, err := listObjectsPaginator(client, bucketName, pathSuffix)
	if err != nil {
		return nil, err
	}

	for _, objKey := range objKeys {
		if !strings.Contains(config.VideoExtensions, filepath.Ext(objKey)) {
			continue
		}

		videoName := filepath.Join(config.ExtractorFolderTmp, filepath.Base(objKey))
		videosPaths = append(videosPaths, videoName)

		file, err := os.Create(videoName)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		if _, err := manager.NewDownloader(client).Download(context.Background(), file, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objKey),
		}); err != nil {
			return nil, err
		}
	}

	return videosPaths, nil
}
