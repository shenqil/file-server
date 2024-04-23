package bucket

import "github.com/google/wire"

// 文件存储桶
var FileBucketName = "files"
var FileBucketLocation = "us-east-1"

// 头像存储桶
var AvatarBucketName = "avatars"
var AvatarBucketLocation = "us-east-1"

// BucketSet model 注入
var BucketSet = wire.NewSet(
	FileSet,
)
