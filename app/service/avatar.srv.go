package service

import (
	"context"
	"fileServer/app/model/minio/bucket"
	"fileServer/app/schema"

	"github.com/google/wire"
)

// AvatarSet 上传头像注入
var AvatarSet = wire.NewSet(wire.Struct(new(Avatar), "*"))

// Avatar 头像
type Avatar struct {
	FileModel *bucket.FileBucket
}

func (a *Avatar) Upload(ctx context.Context, item schema.File) (*schema.IDResult, error) {
	info, err := a.FileModel.Upload(ctx, bucket.AvatarBucketName, item.Name, item.Reader, item.Size, item.Type)
	return schema.NewIDResult(info.Key), err
}
