package api

import (
	"fileServer/app/ginx"
	"fileServer/app/schema"
	"fileServer/app/service"
	"fileServer/util/errors"
	"io"
	"net/http"
	"strconv"
	"strings"

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
			ginx.ResError(c, errors.New400Response("form-data中未找到file"))
			break
		}

		if err != nil {
			ginx.ResError(c, errors.New400Response(err.Error()))
			return
		}

		formName := part.FormName()

		if formName == "file" {
			if item.Name == "" {
				item.Name = part.FileName()
			}
			item.Type = part.Header.Get("Content-Type")
			item.Reader = part
			defer part.Close()
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

// Get 下载文件
// @Tags File
// @Security ApiKeyAuth
// @Summary 下载指定文件
// @Param fileName path string true "唯一标识"
// @Produce application/octet-stream
// @Success 200 {file} application/octet-stream
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /file-api/v1/files/{fileName} [get]
func (a *FileApi) Get(c *gin.Context) {
	ctx := c.Request.Context()
	object, err := a.FileSrv.Get(ctx, c.Param("fileName"))

	if err != nil {
		ginx.ResError(c, err)
		return
	}
	defer object.Reader.Close()

	c.Header("Content-Type", object.Type)
	c.Header("X-File-Size", strconv.FormatInt(object.Size, 10))
	// 设置响应头
	if strings.HasPrefix(object.Type, "image/") {
		// 如果是图片类型，则支持预览
		c.Header("Content-Disposition", "inline")
	} else {
		// 否则作为附件下载
		c.Header("Content-Disposition", "attachment; filename="+object.Name)
	}
	// 将图片写入响应体，并获取文件大小
	_, err = io.Copy(c.Writer, object.Reader)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	// 响应结束
	c.Status(http.StatusOK)
}

// Delete 删除文件
// @Tags File
// @Security ApiKeyAuth
// @Summary 删除文件
// @Param fileName path string true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /file-api/v1/files/{fileName} [delete]
func (a *FileApi) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.FileSrv.Delete(ctx, c.Param("fileName"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
