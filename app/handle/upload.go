package handle

import (
	"fileServer/app/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/tus/tusd/cmd/tusd/cli"
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/tus/tusd/pkg/s3store"
	"log"
)

// FileUpload .
func FileUpload() *tusd.Handler {

	minio := config.C.Minio

	// 根据 accessKey 创建一个 session
	newSession := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(minio.Region),
		Credentials: credentials.NewStaticCredentials(minio.AccessKey, minio.SecretKey, ""),
	}))

	// 新建一个 s3 配置
	s3Config := aws.NewConfig().
		WithEndpoint(minio.Endpoint).
		WithS3ForcePathStyle(true).
		WithDisableSSL(true)

	// 创建 S3 客户端
	s3Client := s3.New(newSession, s3Config)

	// 使用 tusd 连接 minio
	store := s3store.New(minio.Bucket, s3Client)

	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	// 创建 tusd 配置文件
	tusdConfig := tusd.Config{
		BasePath:                "/files/",
		StoreComposer:           composer,
		NotifyCompleteUploads:   true,
		NotifyTerminatedUploads: true,
		NotifyUploadProgress:    true,
		NotifyCreatedUploads:    true,
	}

	cli.SetupPreHooks(&tusdConfig)

	handler, err := tusd.NewHandler(tusdConfig)
	if err != nil {
		log.Printf("unable to create handler: %s", err)
	}

	eventHandler(handler)

	cli.SetupPostHooks(handler)

	return handler
}

func eventHandler(handler *tusd.Handler) {
	go func() {
		for {
			select {
			case info := <-handler.CompleteUploads:
				log.Println("file info --> ", info)
			}
		}
	}()
}
