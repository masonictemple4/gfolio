package handlers

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
)

func NewPostsHandler() PostsHandler {
	// Replace this in-memory function with a call to a database.
	postsGetter := func() (posts []Post, err error) {
		return []Post{{Name: "templ", Author: "author"}}, nil
	}
	return PostsHandler{
		GetPosts: postsGetter,
		Log:      log.Default(),
	}
}

type PostsHandler struct {
	Log      *log.Logger
	GetPosts func() ([]Post, error)
}

func (ph PostsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ps, err := ph.GetPosts()
	if err != nil {
		ph.Log.Printf("failed to get posts: %v", err)
		http.Error(w, "failed to retrieve posts", http.StatusInternalServerError)
		return
	}
	templ.Handler(posts(ps)).ServeHTTP(w, r)
}

type Post struct {
	Name   string
	Author string
}
