package s3

import (
	"context"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"extractor/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
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

func DownloadObjects(bucketName, pathSuffix string) (*objectMetadata, []string, error) {
	var (
		objectMetadata objectMetadata
		videosPaths    []string
	)

	client, err := getClient()
	if err != nil {
		return nil, nil, err
	}

	objKeys, err := listObjectsPaginator(client, bucketName, pathSuffix)
	if err != nil {
		return nil, nil, err
	}

	var obj *s3.GetObjectOutput
	for _, objKey := range objKeys {
		if !strings.Contains(config.VideoExtensions, filepath.Ext(objKey)) {
			continue
		}

		videoName := filepath.Join(config.ExtractorFolderTmp, filepath.Base(objKey))
		videosPaths = append(videosPaths, videoName)

		obj, err = client.GetObject(context.Background(), &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objKey),
		})
		if err != nil {
			return nil, nil, err
		}
		defer obj.Body.Close()

		file, err := os.Create(videoName)
		if err != nil {
			return nil, nil, err
		}
		defer file.Close()

		if _, err := io.Copy(file, obj.Body); err != nil {
			return nil, nil, err
		}
	}

	if obj != nil {
		objectMetadata.UserEmail = obj.Metadata["email"]
		objectMetadata.VideoID = obj.Metadata["id"]
	}

	return &objectMetadata, videosPaths, nil
}
