package bucket

import (
	"context"
	"io"

	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
)

var AvatarSet = wire.NewSet(wire.Struct(new(AvatarBucket), "*"))

var AvatarBucketName = "avatars"
var AvatarBucketLocation = "us-east-1"

type AvatarBucket struct {
	MinioClient *minio.Client
}

func (a *AvatarBucket) Upload(ctx context.Context, fileName string, reader io.Reader, size int64, contentType string) (info minio.UploadInfo, err error) {
	return a.MinioClient.PutObject(ctx, AvatarBucketName, fileName, reader, size, minio.PutObjectOptions{ContentType: contentType})
}
