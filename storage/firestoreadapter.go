package storage

import (
	"context"
	"log"

	firestore "cloud.google.com/go/firestore"
)

// FirestoreAdapter implements storer in firestore
type FirestoreAdapter struct {
	firestoreCollection *firestore.CollectionRef
}

// NewFirestoreAdapter constructor
func NewFirestoreAdapter(collectionName string, projectName string) PostStorer {
	fa := FirestoreAdapter{}
	client, err := firestore.NewClient(context.Background(), projectName)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	fa.firestoreCollection = client.Collection(collectionName)
	log.Printf("Client to firestore collection '%s' is created", collectionName)
	return &fa
}

// AddPostID creates a new empty file. Name of the file is an id
func (fa *FirestoreAdapter) AddPostID(id string) {
	_, err := fa.firestoreCollection.Doc(id).Set(context.Background(), "1")
	if err != nil {
		log.Panic(err)
	}
}

// IsPostIDExists lists all the objects in the bucket
// and checks if object with name id exists
func (fa *FirestoreAdapter) IsPostIDExists(id string) bool {
	_, err := fa.firestoreCollection.Doc(id).Get(context.Background())
	if err != nil {
		return false
	}
	return true
}
