package controllers

import (
	"flaver/api"
	"flaver/api/response"
	"flaver/lib/utils"
	"flaver/services/file"

	"github.com/gin-gonic/gin"
)

type FileController struct {
	fileService file.IFileService
}

func NewFileController() FileController {
	return FileController{
		fileService: file.NewFileServiceOption(
			file.WithStorageUtil(&utils.GcsUtil{}),
		),
	}
}

func (this FileController) UploadImage(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		api.SendResult(err, nil, c)
		return
	}
	files := form.File["imgs"]

	_type := c.PostForm("type")

	if len(files) == 0 {
		api.SendResult(api.InvalidArgument, nil, c)
		return
	}

	urls, err := this.fileService.UploadImage(_type, files)
	if err != nil {
		api.SendResult(err, nil, c)
		return
	}

	api.SendResult(nil, &response.UploadImg{Urls: urls}, c)
}
