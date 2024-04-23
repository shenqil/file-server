package bucket

import (
	"context"
	"io"

	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
)

var FileSet = wire.NewSet(wire.Struct(new(FileBucket), "*"))

var FileBucketName = "files"
var FileBucketLocation = "us-east-1"

type FileBucket struct {
	MinioClient *minio.Client
}

func (a *FileBucket) Upload(ctx context.Context, fileName string, reader io.Reader, size int64, contentType string) (info minio.UploadInfo, err error) {
	return a.MinioClient.PutObject(ctx, FileBucketName, fileName, reader, size, minio.PutObjectOptions{ContentType: contentType})
}
