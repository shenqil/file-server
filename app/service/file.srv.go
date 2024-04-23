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
	info, err := a.FileModel.Upload(ctx, bucket.FileBucketName, item.Name, item.Reader, item.Size, item.Type)
	return schema.NewIDResult(info.Key), err
}

func (a *FileServer) Get(ctx context.Context, fileName string) (*schema.File, error) {
	reader, err := a.FileModel.Get(ctx, bucket.FileBucketName, fileName)
	if err != nil {
		return nil, err
	}

	stat, err := reader.Stat()
	if err != nil {
		return nil, err
	}

	return &schema.File{
		Name:   fileName,
		Size:   stat.Size,
		Type:   stat.ContentType,
		Reader: reader,
	}, err
}

func (a *FileServer) Delete(ctx context.Context, fileName string) error {
	return a.FileModel.Remove(ctx, bucket.FileBucketName, fileName)
}
