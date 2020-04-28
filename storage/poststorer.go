package storage

// PostStorer to store posts
type PostStorer interface {
	IsPostIDExists(id string) bool
	AddPostID(id string)
}
