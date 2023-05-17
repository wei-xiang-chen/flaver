package file

import (
	"flaver/globals"
	"flaver/lib/utils"
	"fmt"
	"mime/multipart"
	"sync"
)

type FileService struct {
	storageUtil utils.IStorageUtil
}

type FileServiceOption func(*FileService)

func NewFileServiceOption(options ...func(*FileService)) IFileService {
	service := FileService{}

	for _, option := range options {
		option(&service)
	}

	return &service
}

func WithStorageUtil(util utils.IStorageUtil) FileServiceOption {
	return func(service *FileService) {
		service.storageUtil = util
	}
}

func (this *FileService) UploadImage(_type string, fileHeaders []*multipart.FileHeader) ([]string, error) {
	fmt.Println("UploadImageService")
	urlChan := make(chan string, len(fileHeaders))
	resultUrls := make([]string, 0)

	wg := new(sync.WaitGroup)
	wg.Add(len(fileHeaders))

	for _, header := range fileHeaders {
		file, err := header.Open()
		if err != nil {
			globals.GetLogger().Errorf("fileHeader.Open: ", err)
		}

		go this.storageUtil.UploadFile(_type, file, wg, urlChan)
	}
	wg.Wait()
	close(urlChan)

	for {
		url, ok := <-urlChan
		if !ok {
			break
		}
		resultUrls = append(resultUrls, url)
	}

	return resultUrls, nil
}
