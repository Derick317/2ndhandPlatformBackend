package backend

import (
	"context"
	"fmt"
	"io"
	"secondHand/constants"
	"secondHand/util"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var (
	gcsBackend *GoogleCloudStorageBackend
)

type GoogleCloudStorageBackend struct {
	client *storage.Client
	bucket string
}

func InitGCSBackend() {
	var client *storage.Client
	var err error
	if constants.DEPLOYED {
		client, err = storage.NewClient(context.Background())
	} else {
		client, err = storage.NewClient(context.Background(),
			option.WithCredentialsFile(util.MustGetenv("GOOGLE_APPLICATION_CREDENTIALS")))
	}

	if err != nil {
		panic(err)
	}

	gcsBackend = &GoogleCloudStorageBackend{
		client: client,
		bucket: util.MustGetenv("GOOGLE_CLOUD_STORAGE_BUCKET"),
	}
}

// SaveToGCS saves a file named objectName in google cloud storage
//
// A successful SaveToGCS returns file's url and error == nil.
// In contrast, a failed one returns an empty string and corresponding error.
func SaveToGCS(r io.Reader, objectName string) (string, error) {
	// return objectName, nil
	ctx := context.Background()
	object := gcsBackend.client.Bucket(gcsBackend.bucket).Object(objectName)

	// Set a generation-match precondition to avoid potential race conditions and data
	// corruptions. The request to upload is aborted if the object's generation number
	// does not match your precondition.
	// For an object that does not yet exist, set the DoesNotExist precondition.
	// object = object.If(storage.Conditions{DoesNotExist: true})
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
	return attrs.MediaLink, nil
}

// DeleteFromGCS deletes a file in google cloud storage by its url
//
// A successful DeleteFromGCS returns error == nil. Otherwise, it returns corresponding error.
func DeleteFromGCS(objectName string) error {
	ctx := context.Background()
	object := gcsBackend.client.Bucket(gcsBackend.bucket).Object(objectName)

	// Set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to delete the file is aborted
	// if the object's generation number does not match your precondition.
	attrs, err := object.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("object.Attrs: %w", err)
	}
	object = object.If(storage.Conditions{GenerationMatch: attrs.Generation})

	return object.Delete(ctx)
}
