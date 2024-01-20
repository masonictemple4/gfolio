package models

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/gosimple/slug"
	"github.com/masonictemple4/masonictempl/internal/dtos"
	"gorm.io/gorm"
)

const (
	BlogStateDraft     = "draft"
	BlogStatePublished = "published"
	BlogStateArchived  = "archived"
)

func ValidBlogState(state string) bool {
	switch state {
	case BlogStateDraft, BlogStatePublished, BlogStateArchived:
		return true
	}
	return false
}

type Blog struct {
	gorm.Model
	Title       string `gorm:"column:title;" json:"title"`
	Subtitle    string `gorm:"column:subtitle;" json:"subtitle"`
	Description string `gorm:"column:description;" json:"description"`
	Thumbnail   string `gorm:"column:thumbnail;" json:"thumbnail"`
	ContentUrl  string `gorm:"column:contenturl;" json:"contenturl"`
	// Docpath should be the /asset_dir/blog_root/path-schema
	// will join the filestore root with this path when opening
	// local files. Thus eliminating our issues with storing
	// paths in db that may be shared on multiple devices.
	Docpath    string    `gorm:"column:docpath;" json:"docpath"`
	Bucketname string    `gorm:"column:bucketname;" json:"bucketname"`
	State      string    `gorm:"column:state;" json:"state"`
	Slug       string    `gorm:"column:slug;" json:"slug"`
	Tags       []Tag     `gorm:"many2many:blog_tags;" json:"tags"`
	Media      []Media   `gorm:"many2many:blog_media;" json:"media"`
	Authors    []User    `gorm:"many2many:blog_authors;" json:"authors"`
	Comments   []Comment `json:"comments"`
}

/*
dir := path.Join(rootPath, post.Date.Format("2006/01/02"), slug.Make(post.Title))
		if err := os.MkdirAll(dir, 0755); err != nil && err != os.ErrExist {
			log.Fatalf("failed to create dir %q: %v", dir, err)
		}

		// Create the output file.
		name := path.Join(dir, "index.html")
		f, err := os.Create(name)
*/

func (p *Blog) FromBlogInput(tx *gorm.DB, input *dtos.BlogInput) error {
	p.Title = input.Title
	p.Subtitle = input.Subtitle
	p.Thumbnail = input.Thumbnail
	p.ContentUrl = input.ContentUrl
	p.Description = input.Description

	if p.ID == 0 {
		if err := tx.Create(p).Error; err != nil {
			return err
		}
	} else {
		if err := tx.Save(p).Error; err != nil {
			return err
		}
	}

	if len(input.Tags) > 0 {
		err := p.ClearAssociations(tx, "Tags")
		if err != nil {
			return err
		}
		var tags []Tag
		err = TagFromStrings(tx, input.Tags, &tags)
		if err != nil {
			return err
		}
		p.Tags = tags
	}

	if len(input.Authors) > 0 {
		err := p.ClearAssociations(tx, "Authors")
		if err != nil {
			return err
		}
		var authors []User
		err = AuthorFromInput(tx, input.Authors, &authors)
		if err != nil {
			return err
		}
		p.Authors = authors
	}

	if len(input.Media) > 0 {
		err := p.ClearAssociations(tx, "Media")
		if err != nil {
			return err
		}
		var media []Media
		err = MediaFromStrings(tx, input.Media, &media)
		if err != nil {
			return err
		}

		p.Media = media
	}

	if p.ID == 0 {
		return tx.Create(p).Error
	} else {
		return tx.Save(p).Error
	}

}

// GenerateSlug will generate a slug for the blog.
// this method also takes in an optional input string
// to override the generated version for custom slugs.
//
// Leave the input empty if you would like to generate
// a slug from the title.
//
// IMPORTANT:
//
//	You must call Update on the object if you'd like to
//	persist this change in the database.
func (p *Blog) GenerateSlug(input string) string {
	// TODO: Build library for this.
	if p.Slug != "" && input == "" {
		fmt.Println("The slug is already set: ", p.Slug)
		return p.Slug
	}

	if input != "" {
		return input
	}

	// is title a single word?
	// if so use that as the slug.
	// for some reason slug.Make was returning
	// an empty string.
	pattern := `^\S+$`
	r, err := regexp.Compile(pattern)
	if err != nil {
		log.Printf("Error compiling regex: %v", err)
		return ""
	}
	println("The title we're checking: ", p.Title)

	if r.MatchString(p.Title) {
		return strings.ToLower(p.Title)
	}

	newSlug := slug.Make(p.Title)
	return newSlug
}

func (p *Blog) generateFileName() string {
	return fmt.Sprintf("%s-%d.md", p.Slug, p.ID)
}

func (p *Blog) generateBlogDir() (string, error) {
	datePath := p.CreatedAt.Format("2006/01/02")
	datePathParts := strings.Split(datePath, "/")
	if len(datePathParts) != 3 {
		return "", errors.New("blog model: generatestorageobject: invalid date path")
	}

	blogRoot := "blogs"
	if val := os.Getenv("BLOG_ROOT"); val != "" {
		blogRoot = os.Getenv("BLOG_ROOT")
	}

	blogRoot = strings.TrimSuffix(blogRoot, "/")

	return fmt.Sprintf("%s/%s", blogRoot, datePath), nil
}

// Returns fully qualified path for the blog's doc path
// This is the location on the filesystem. Use this when
// interacting directly with the file.
func (p *Blog) GenerateDocPath(assetRoot string) (string, error) {

	pwd := os.Getenv("PWD")

	if strings.Contains(assetRoot, pwd) {
		assetRoot = strings.Replace(assetRoot, pwd, "", 1)
	}

	if p.Slug == "" {
		if slug := p.GenerateSlug(""); slug == "" {
			return "", errors.New("blog model: generatestorageobject: invalid slug")
		}
	}

	blogDir, err := p.generateBlogDir()
	if err != nil {
		return "", err
	}

	obj := fmt.Sprintf("%s/%s/%s", assetRoot, blogDir, p.generateFileName())

	return obj, nil
}

// Requires Bucketname
// This is what you would request in href or src attributes.
// the public path for the file essentially.
// Locally this is served with the path handled with a local fileserver.
// Otherwise this will typically be a url for the bucket.
func (p *Blog) GenerateContentUrl() string {
	baseUrl := os.Getenv("BUCKET_BASE_URL")
	if baseUrl == "" {
		// With no bucket this will default to the
		// internal static file server path.
		wd := os.Getenv("PWD")
		return strings.Replace(p.Docpath, wd, "", 1)
	}

	return fmt.Sprintf("%s/%s/%s", baseUrl, p.Bucketname, p.Docpath)
}

// TODO: Fill out this method.
func (p *Blog) AfterDelete(tx *gorm.DB) error {
	// Clean up Filestore
	// Clean up Media
	return nil
}

func (p *Blog) ClearAssociations(tx *gorm.DB, assoc string) error {
	return tx.Model(p).Association(assoc).Clear()
}
