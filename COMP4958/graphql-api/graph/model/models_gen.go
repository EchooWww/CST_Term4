// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Mutation struct {
}

type Post struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	AuthorID string `json:"authorID"`
	Author  *User  `json:"author"`
}



type User struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Email string  `json:"email"`
	Posts []*Post `json:"posts"`
}
