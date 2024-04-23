package api

import (
	"fileServer/app/ginx"
	"fileServer/app/schema"
	"fileServer/app/service"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var FileSet = wire.NewSet(wire.Struct(new(FileApi), "*"))

// File 文件
type FileApi struct {
	FileSrv *service.FileServer
}

// Upload 上传文件
// @Tags File
// @Security ApiKeyAuth
// @Summary 上传文件
// @Accept multipart/form-data
// @Param file formData file true "file"
// @Success 200 {object} schema.IDResult
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /file-api/v1/files [post]
func (a *FileApi) Upload(c *gin.Context) {
	ctx := c.Request.Context()

	reader, err := c.Request.MultipartReader()
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	item := schema.File{
		Name:   "",
		Size:   -1,
		Type:   "",
		Reader: nil,
	}

	// Traverse form fields and file streams.
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			ginx.ResError(c, err)
			return
		}

		formName := part.FormName()

		if formName == "file" { // Handle file streams.
			// Get file name
			if item.Name == "" {
				item.Name = part.FileName()
			}
			item.Type = part.Header.Get("Content-Type")
			item.Reader = part
			break // 获取到文件流后直接退出，不接受后面内容
		}
	}

	result, err := a.FileSrv.Upload(ctx, item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	ginx.ResSuccess(c, result)
}
