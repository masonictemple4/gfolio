package db

import (
	"errors"

	"github.com/masonictemple4/masonictempl/db/models"
	"gorm.io/gorm"
)

var (
	ErrNoDB = errors.New("db: no database provided")
)

type BlogStore struct {
	db *gorm.DB
}

type OptionsFn func(*BlogStore)

/*
TODO: Was going to set this up with a pattern like:

	func WithSQLite(name string) func(*BlogStore) {
		return func(bs *BlogStore) {
			bs.db = NewSqliteDB(name, nil)
		}
	}

However, that would require a new options fn type
to then pass into the NewBlogStore function. This
felt cleaner.

Usage:
store := NewSqliteDB("blog.db", nil)
newStore := NewBlogStore(store)
*/
func WithDB(gDB *gorm.DB) func(*BlogStore) {
	return func(bs *BlogStore) {
		bs.db = gDB
	}
}

func NewBlogStore(opts ...OptionsFn) (store *BlogStore, err error) {
	store = &BlogStore{}
	for _, o := range opts {
		o(store)
	}

	if store.db == nil {
		err = ErrNoDB
		store = nil
	}

	err = store.migrate()

	return
}

func (bs *BlogStore) Close() error {
	db, err := bs.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (bs *BlogStore) migrate() error {
	return bs.db.AutoMigrate(
		&models.Blog{},
		&models.Tag{},
		&models.Media{},
		&models.User{},
		&models.Comment{},
	)
}

func (bs *BlogStore) ListBlogs(b *[]models.Blog, query map[string]any, limits map[string]int, order string, preloads ...string) error {
	tx := bs.db.Model(b)

	if len(preloads) > 0 {
		for _, preload := range preloads {
			tx = tx.Preload(preload)
		}
	}

	if limits != nil {
		limit, limitOk := limits["limit"]
		if limitOk {
			tx = tx.Limit(limit)
		}
		offset, offsetOk := limits["offset"]
		if offsetOk {
			tx = tx.Offset(offset)
		}
	}

	if query != nil {
		for k, v := range query {
			tx = tx.Where(k, v)
		}
	}

	if order != "" {
		tx.Order(order)
	}

	return tx.Find(b).Error
}

func (bs *BlogStore) DB() *gorm.DB {
	return bs.db
}

func (bs *BlogStore) UpdateFromMap(tx *gorm.DB, b *models.Blog, body map[string]any, bid int) error {
	return bs.db.Model(b).Where("id = ?", bid).Updates(body).Error
}

func (bs *BlogStore) FindByID(tx *gorm.DB, b *models.Blog, bid int) error {
	return tx.First(b, bid).Error
}

func (bs *BlogStore) FindBySlug(b *models.Blog, slug string, preloads ...string) error {
	tx := bs.db

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	return tx.Where("slug = ?", slug).First(b).Error
}
