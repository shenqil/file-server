package common

import (
	"path/filepath"

	"github.com/google/uuid"
)

// 生成唯一的文件名，保留原始文件的后缀名
func GenerateUniqueFilename(filePath string) string {
	// 使用UUID生成唯一标识符
	uuid := uuid.New()
	// 获取原始文件的后缀名
	ext := filepath.Ext(filePath)
	// 将UUID和后缀名拼接作为唯一文件名
	return uuid.String() + ext
}
