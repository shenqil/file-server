package schema

import "io"

// File 上传文件信息
type File struct {
	Name   string        // 文件名称
	Size   int64         // 文件大小
	Type   string        // 文件类型 (MIME 类型)
	Path   string        // 文件路径
	Reader io.ReadCloser // 文件流 (实现了 io.ReadSeeker 接口)
}
