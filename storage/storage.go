package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/bujosa/xihe/utils"
)

type Storage struct {
	bucketName    string
	projectID     string
	subFolderPath string
	bucket        *storage.BucketHandle
	context       context.Context
}

func New(ctx context.Context) *Storage {
	projectId := ctx.Value(utils.ProjectIdKey).(string)
	bucketName := ctx.Value(utils.BucketNameKey).(string)
	subFolderPath := ctx.Value(utils.SubFolderPathKey).(string)

	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	return &Storage{
		bucketName:    bucketName,
		projectID:     projectId,
		subFolderPath: subFolderPath,
		bucket:        client.Bucket(bucketName),
		context:       ctx,
	}
}

func (s *Storage) Upload(url string) (string, error) {
	log.Print("Upload Picture with url: ", url)

	url = strings.Replace(url, "800x600", "500x500", 1)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error downloading picture: " + url)
		return "", err
	}
	defer resp.Body.Close()

	parts := strings.Split(url, "/")
	urlFormat := parts[len(parts)-1]

	objectName := s.subFolderPath + urlFormat
	object := s.bucket.Object(objectName)

	wc := object.NewWriter(s.context)

	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err = io.Copy(wc, resp.Body); err != nil {
		log.Println(err)
		return "", err
	}

	if err := wc.Close(); err != nil {
		log.Println(err)
		return "", err
	}

	formatUrl := "https://storage.googleapis.com/" + s.bucketName + "/" + objectName

	log.Println("Picture uploaded successfully: " + formatUrl)

	// Return new url of the file in the bucket storage in gcp
	return formatUrl, nil
}

func (s *Storage) AlreadyExist(url string) bool {
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

func (s *Storage) validate(url string) bool {
	_, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (s *Storage) RestartConnection() error {
	client, err := storage.NewClient(s.context)
	if err != nil {
		return err
	}

	s.bucket = client.Bucket(s.bucketName)

	return nil
}
