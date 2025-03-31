package graph

import "bug-free/demo/graph/model"

// Resolver contains dependencies like database connections
// For this demo, we're using in-memory slices to store data
type Resolver struct {
	users []*model.User
	posts []*model.Post
}

// NewResolver creates a new resolver with empty slices
func NewResolver() *Resolver {
	return &Resolver{
		users: []*model.User{},
		posts: []*model.Post{},
	}
}