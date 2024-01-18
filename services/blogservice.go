package services

import "github.com/masonictemple4/masonictempl/db"

type BlogService struct {
	Store db.BlogStore
}

func NewBlogService() *BlogService {
	return &BlogService{}
}
