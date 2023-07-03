package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/bujosa/xihe/env"
)

type Storage struct {
	bucketName string
	projectID  string
	subFolderPath string
	bucket  	*storage.BucketHandle
	context 	context.Context
}

func New() *Storage {
	projectId, err := env.GetString("GOOGLE_PROJECT_ID")
	if err != nil {
		panic(err)
	}
	bucketName, err := env.GetString("GOOGLE_BUCKET_NAME")
	if err != nil {
		panic(err)
	}
	_, err = env.GetString("GOOGLE_APPLICATION_CREDENTIALS")
	if err != nil {
		panic(err)
	}

	subFolderPath, err := env.GetString("GOOGLE_SUBFOLDER_PATH")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	return &Storage{
		bucketName: bucketName,
		projectID:  projectId,
		subFolderPath: subFolderPath,
		bucket:  	client.Bucket(bucketName),
		context: 	ctx,
	}
}

func (s Storage) Upload(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	parts := strings.Split(url, "/")
	urlFormat := parts[len(parts)-1]

	objectName := s.subFolderPath + urlFormat;
	object := s.bucket.Object(objectName)

	wc := object.NewWriter(s.context)

	if _, err = io.Copy(wc, resp.Body); err != nil {
		fmt.Println(err)
		return "", err
	}

	if err := wc.Close(); err != nil {
		fmt.Println(err)
		return "", err
	}

	// Return new url of the file in the bucket storage in gcp
	return "https://storage.googleapis.com/" + s.bucketName + "/" + objectName, nil
}

func (s Storage) AlreadyExist(url string) bool {
	isValid := s.validate(url)

	if !isValid {
		return false
	}

	object := s.bucket.Object(url)
	_, err := object.Attrs(s.context)
	if err == storage.ErrObjectNotExist {
		return false
	} else if err != nil {
		fmt.Println("Error al verificar la existencia del objeto:", err)
		return false
	} else {
		return true
	}
}

func (s Storage) validate(url string) bool {
	_, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}