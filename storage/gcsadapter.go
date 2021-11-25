package storage

import (
	"context"
	"io"
	"log"

	gcs "cloud.google.com/go/storage"
)

// GcsAdapter implements storer in google cloud storage
type GcsAdapter struct {
	bucketHandle *gcs.BucketHandle
}

// NewGcsAdapter creates a new instance. Accepts bucket name for bucket handle
func NewGcsAdapter(bucketName string) PostStorer {
	s := GcsAdapter{}
	client, err := gcs.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	s.bucketHandle = client.Bucket(bucketName)
	log.Printf("Bucket handle to bucket '%s' is created", bucketName)
	return &s
}

// AddPostID creates a new empty file. Name of the file is an id
func (s *GcsAdapter) AddPostID(id string) {
	fw := s.bucketHandle.Object(id).NewWriter(context.Background())

	if _, err := io.WriteString(fw, ""); err != nil {
		log.Panic(err)
	}

	if err := fw.Close(); err != nil {
		log.Panic(err)
	}
}

// IsPostIDExists lists all the objects in the bucket
// and checks if object with name id exists
func (s *GcsAdapter) IsPostIDExists(id string) bool {
	_, err := s.bucketHandle.Object(id).NewReader(context.Background())
	return err == nil
}
