package services

import (
	"context"

	"github.com/masonictemple4/masonictempl/db"
	"github.com/masonictemple4/masonictempl/db/models"
	"gorm.io/gorm"
)

type BlogService struct {
	Store *db.BlogStore
}

func NewBlogService(d *gorm.DB) *BlogService {

	if d == nil {
		panic("NewBlogService requires a gorm.DB")
	}

	store, err := db.NewBlogStore(db.WithDB(d))
	if err != nil {
		panic(err)
	}

	return &BlogService{Store: store}
}

func (b *BlogService) List(ctx context.Context) []models.Blog {
	blogs := make([]models.Blog, 0)
	order := "created_at desc"
	if err := b.Store.ListBlogs(&blogs, nil, nil, order, "Authors", "Tags"); err != nil {
		// TODO: Log the error
		return blogs
	}

	return blogs
}

func (b *BlogService) GetWithSlug(ctx context.Context, slug string, preloads ...string) (*models.Blog, error) {
	var blog models.Blog
	if err := b.Store.FindBySlug(&blog, slug, preloads...); err != nil {
		return nil, err
	}
	return &blog, nil
}
