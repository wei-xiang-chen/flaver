package file

import "mime/multipart"

type IFileService interface {
	UploadImage(_type string, fileHeaders []*multipart.FileHeader) ([]string, error)
}

type FileType string

const (
	UserAvatar FileType = "user_avatar"
	Post       FileType = "post"
)
