package service

import (
	"context"
	"fileServer/app/model/minio/bucket"
	"fileServer/app/schema"

	"github.com/google/wire"
)

// FileSet 注入File
var FileSet = wire.NewSet(wire.Struct(new(FileServer), "*"))

// File 文件
type FileServer struct {
	FileModel *bucket.FileBucket
}

func (a *FileServer) Upload(ctx context.Context, item schema.File) (*schema.IDResult, error) {
	info, err := a.FileModel.Upload(ctx, item.Name, item.Reader, item.Size, item.Type)
	return schema.NewIDResult(info.Key), err
}
