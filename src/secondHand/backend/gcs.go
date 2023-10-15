package backend

import (
	"context"
	"fmt"
	"io"

	"secondHand/util"

	"cloud.google.com/go/storage"
)

var (
	GCSBackend *GoogleCloudStorageBackend
)

type GoogleCloudStorageBackend struct {
	client *storage.Client
	bucket string
}

func InitGCSBackend(config *util.GCSInfo) {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		panic(err)
	}

	GCSBackend = &GoogleCloudStorageBackend{
		client: client,
		bucket: config.Bucket,
	}
}

// SaveToGCS saves a file named objectName in google cloud storage
//
// A successful SaveToGCS returns file's url and error == nil.
// In contrast, a failed one returns an empty string and corresponding error.
func (backend *GoogleCloudStorageBackend) SaveToGCS(r io.Reader, objectName string) (string, error) {
	ctx := context.Background()
	object := backend.client.Bucket(backend.bucket).Object(objectName)
	wc := object.NewWriter(ctx)
	if _, err := io.Copy(wc, r); err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	// ACL: access control
	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}

	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", err
	}

	fmt.Printf("File is saved to GCS: %s\n", attrs.MediaLink)
	return attrs.MediaLink, nil
}

// DeleteFromGCS deletes a file in google cloud storage by its url
//
// A successful DeleteFromGCS returns error == nil. Otherwise, it returns corresponding error.
func (backend *GoogleCloudStorageBackend) DeleteFromGCS(url string) error {
	return nil
}
