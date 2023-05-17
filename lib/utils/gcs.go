package utils

import (
	"context"
	"flaver/globals"
	"io"
	"mime/multipart"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/api/option"
)

var (
	gcsClient     *storage.Client
	gcsClientMux  sync.Mutex
	gcsClientOnce sync.Once
)

type IStorageUtil interface {
	UploadFile(_type string, file multipart.File, wg *sync.WaitGroup, resultUrls chan string)
}

func newGcsClient() *storage.Client {
	client, err := storage.NewClient(context.Background(), option.WithCredentialsFile(globals.GetConfig().Gcs.GetCredentialFilePath()))
	if err != nil {
		globals.GetLogger().Fatalf("[gcs client startup] error: %v", err)
	}
	return client
}

func getGcsClient() *storage.Client {
	gcsClientMux.Lock()
	defer gcsClientMux.Unlock()
	gcsClientOnce.Do(func() {
		if gcsClient == nil {
			gcsClient = newGcsClient()
		}
	})

	return gcsClient
}

type GcsUtil struct {
}

func (this *GcsUtil) UploadFile(_type string, file multipart.File, wg *sync.WaitGroup, resultUrls chan string) {
	defer wg.Done()

	var err error
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	now := time.Now()
	objectName := _type + "/" + now.Format("2006-01-02") + "/" + uuid.NewV4().String()

	wc := getGcsClient().Bucket(globals.GetConfig().GetGcs().GetBucketName()).Object(objectName).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	if _, err = io.Copy(wc, file); err != nil {
		globals.GetLogger().Errorf("io.Copy error: ", err)
		return
	}
	if err = wc.Close(); err != nil {
		globals.GetLogger().Errorf("Writer.Close error: ", err)
		return
	}

	resultUrls <- globals.GetConfig().GetGcs().GetBaseUrl() + "/" + objectName
}
