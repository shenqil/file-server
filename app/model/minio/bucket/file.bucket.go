package bucket

import (
	"context"
	"io"

	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
)

var FileSet = wire.NewSet(wire.Struct(new(FileBucket), "*"))

type FileBucket struct {
	MinioClient *minio.Client
}

func (a *FileBucket) Upload(ctx context.Context, bucketName string, fileName string, reader io.Reader, size int64, contentType string) (info minio.UploadInfo, err error) {
	return a.MinioClient.PutObject(ctx, bucketName, fileName, reader, size, minio.PutObjectOptions{ContentType: contentType})
}

func (a *FileBucket) Get(ctx context.Context, bucketName string, fileName string) (*minio.Object, error) {
	return a.MinioClient.GetObject(ctx, bucketName, fileName, minio.GetObjectOptions{})
}

func (a *FileBucket) Remove(ctx context.Context, bucketName string, fileName string) error {
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}
	return a.MinioClient.RemoveObject(context.Background(), bucketName, fileName, opts)
}
