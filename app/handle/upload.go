package handle

import (
	"fileServer/app/config"
	"fileServer/app/pkg/minio"
	"github.com/tus/tusd/cmd/tusd/cli"
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/tus/tusd/pkg/s3store"
	"log"
)

// FileUpload .
func FileUpload() *tusd.Handler {
	cfg := config.C.Minio

	// 使用 tusd 连接 minio
	store := s3store.New(cfg.Bucket, minio.S3Client)

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
