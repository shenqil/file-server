package minio

import (
	"bytes"
	"fileServer/app/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func Upload(fileName string, readSeeker *bytes.Reader) error {
	_, err := S3Client.PutObject(&s3.PutObjectInput{
		Body:   readSeeker,
		Bucket: aws.String(config.C.Minio.Bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return err
	}

	return nil
}
