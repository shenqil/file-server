package minio

import (
	"fileServer/app/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var S3Client *s3.S3
var newSession *session.Session

func init() {
	cfg := config.C.Minio

	// 根据 accessKey 创建一个 session
	newSession = session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, ""),
	}))

	// 新建一个 s3 配置
	s3Config := aws.NewConfig().
		WithEndpoint(cfg.Endpoint).
		WithS3ForcePathStyle(true).
		WithDisableSSL(true)

	// 创建 S3 客户端
	S3Client = s3.New(newSession, s3Config)

	// 判断是否存在 对应的存储桶
	result, err := S3Client.ListBuckets(nil)
	if err != nil {
		panic(err)
		return
	}

	isExist := false
	for _, b := range result.Buckets {
		if *b.Name == cfg.Bucket {
			isExist = true
			break
		}
	}

	// 不存在 则创建存储桶
	if !isExist {
		cParams := &s3.CreateBucketInput{
			Bucket: aws.String(cfg.Bucket), // 必须
		}

		// 调用CreateBucket创建一个新的存储桶。
		_, err := S3Client.CreateBucket(cParams)
		if err != nil {
			panic(err.Error())
			return
		}
	}

}
