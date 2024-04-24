package api

import (
	"fileServer/app/ginx"
	"fileServer/app/schema"
	"fileServer/app/service"
	"fileServer/util/errors"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var AvatarSet = wire.NewSet(wire.Struct(new(Avatar), "*"))

// File 文件
type Avatar struct {
	AvatarSrv *service.Avatar
}

// Upload 上传头像
// @Tags Avatar
// @Security ApiKeyAuth
// @Summary 上传头像
// @Accept multipart/form-data
// @Param name formData string false "文件名称(需要放在表单最前面)"
// @Param file formData file true "file"
// @Success 200 {object} schema.IDResult
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /file-api/v1/avatars [post]
func (a *Avatar) Upload(c *gin.Context) {
	ctx := c.Request.Context()

	userID := ginx.GetUserID(c)
	fileName := ""

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

		if formName == "name" { // 解析出文件名称
			value, err := io.ReadAll(part)
			if err != nil {
				ginx.ResError(c, errors.New400Response("form-data中name解析失败"))
				return
			}

			fileName = string(value)
		} else if formName == "file" { // 解析出文件流
			if fileName == "" {
				ext := filepath.Ext(part.FileName())
				item.Name = userID + "-avatar" + ext
			} else {
				item.Name = userID + "-" + fileName
			}

			item.Type = part.Header.Get("Content-Type")
			item.Reader = part
			defer part.Close()
			break // 获取到文件流后直接退出，不接受后面内容
		}
	}

	result, err := a.AvatarSrv.Upload(ctx, item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	ginx.ResSuccess(c, result)
}

// Get 下载文件
// @Tags Avatar
// @Security ApiKeyAuth
// @Summary 下载指定文件
// @Param name path string true "唯一标识"
// @Produce application/octet-stream
// @Success 200 {file} application/octet-stream
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /file-api/v1/avatars/{name} [get]
func (a *Avatar) Get(c *gin.Context) {
	ctx := c.Request.Context()
	object, err := a.AvatarSrv.Get(ctx, c.Param("name"))

	if err != nil {
		ginx.ResError(c, errors.New400Response(err.Error()))
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
		ginx.ResError(c, errors.New400Response(err.Error()))
		return
	}

	// 响应结束
	c.Status(http.StatusOK)
}

// Delete 删除文件
// @Tags Avatar
// @Security ApiKeyAuth
// @Summary 删除文件
// @Param name path string true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /file-api/v1/avatars/{name} [delete]
func (a *Avatar) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.AvatarSrv.Delete(ctx, c.Param("name"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
