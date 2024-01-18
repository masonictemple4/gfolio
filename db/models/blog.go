package models

import (
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
	Title       string    `gorm:"column:title;" json:"title"`
	Subtitle    string    `gorm:"column:subtitle;" json:"subtitle"`
	Description string    `gorm:"column:description;" json:"description"`
	Thumbnail   string    `gorm:"column:thumbnail;" json:"thumbnail"`
	ContentUrl  string    `gorm:"column:contenturl;" json:"contenturl"`
	Docpath     string    `gorm:"column:docpath;" json:"docpath"`
	Bucketname  string    `gorm:"column:bucketname;" json:"bucketname"`
	State       string    `gorm:"column:state;" json:"state"`
	Slug        string    `gorm:"column:slug;" json:"slug"`
	Tags        []Tag     `gorm:"many2many:post_tags;" json:"tags"`
	Media       []Media   `gorm:"many2many:post_media;" json:"media"`
	Authors     []User    `gorm:"many2many:post_authors;" json:"authors"`
	Comments    []Comment `json:"comments"`
}

/*
func (p *Post) FromPostInput(tx *gorm.DB, input *dtos.PostInput) error {
	p.Title = input.Title
	p.Subtitle = input.Subtitle
	p.Thumbnail = input.Thumbnail
	p.ContentUrl = input.ContentUrl
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
		return p.New(tx)
	} else {
		return p.update(tx)
	}

}

// made private because this is dangerous.
func (p *Post) update(tx *gorm.DB) error {
	return tx.Save(p).Error
}

// GenerateSlug will generate a slug for the post.
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
func (p *Post) GenerateSlug(input string) string {
	// TODO: Build library for this.
	if p.Slug != "" && input == "" {
		fmt.Println("The slug is already set: ", p.Slug)
		return p.Slug
	}

	if input != "" {
		return input
	}

	newSlug := slug.Make(p.Title)

	return newSlug
}

func (p *Post) generateFileName() string {
	return fmt.Sprintf("%s-%d.md", p.Slug, p.ID)
}

func (p *Post) generateBlogDir() (string, error) {
	datePath := p.CreatedAt.Format("2006/01/02")
	datePathParts := strings.Split(datePath, "/")
	if len(datePathParts) != 3 {
		return "", errors.New("post model: generatestorageobject: invalid date path")
	}
	return fmt.Sprintf("blogs/%s", datePath), nil
}

// blogs/{id}/{created}/{slug}.md
func (p *Post) GenerateDocPath() (string, error) {

	if p.Slug == "" {
		if slug := p.GenerateSlug(""); slug == "" {
			return "", errors.New("post model: generatestorageobject: invalid slug")
		}
	}

	blogDir, err := p.generateBlogDir()
	if err != nil {
		return "", err
	}

	obj := fmt.Sprintf("%s/%s", blogDir, p.generateFileName())

	return obj, nil
}

// Requires Bucketname
func (p *Post) GenerateContentUrl() string {
	baseUrl := os.Getenv("BUCKET_BASE_URL")
	return fmt.Sprintf("%s/%s/%s", baseUrl, p.Bucketname, p.Docpath)
}

func (p *Post) FindBySlug(tx *gorm.DB, slug string, opts *repository.RepositoryOpts) error {
	for name, opt := range opts.Preloads {
		tx = tx.Preload(name, opt)
	}
	return tx.Where("slug = ?", slug).Find(p).Error
}

// TODO: Opts (limits etc..) Preloads
func (p *Post) FindByID(tx *gorm.DB, id int, opts *repository.RepositoryOpts) error {
	return tx.First(p, id).Error
}
func (p *Post) Query(tx *gorm.DB, query map[string]any, opts *repository.RepositoryOpts, out any) error {
	return tx.Where(query).Find(out).Error
}

func (p *Post) New(tx *gorm.DB) error {
	return tx.Create(p).Error
}

func (p *Post) Update(tx *gorm.DB, id int, body map[string]any) error {
	return tx.Model(p).Where("id = ?", id).Updates(body).Error
}

// BEWARE!!! This method calls the method to make
// an update from the object istelf while we have
// the pointer.
//
// Update from the map is the recommended
// choice. However, if you understand the risks here
// I have decided to expose the UnsafeUpdate method
func (p *Post) UnsafeUpdate(tx *gorm.DB) error {
	return p.update(tx)
}

func (p *Post) Delete(tx *gorm.DB) error {
	return tx.Delete(p).Error
}

func (p *Post) DeleteById(tx *gorm.DB, id int) error {
	return tx.Delete(p, id).Error
}

func (p *Post) All(tx *gorm.DB, opts *repository.RepositoryOpts, out any) error {
	if opts != nil {
		for name, opt := range opts.Preloads {
			tx = tx.Preload(name, opt)
		}
	}
	return tx.Find(out).Error
}

func (p *Post) FindAssociations(tx *gorm.DB, assoc string, query map[string]any, out repository.AssociationEntity) error {
	if !out.ValidAssociation(p, assoc) {
		return errors.New("post model: findassociations: invalid association. Please verify the assoc value passed to this function and the out object.")
	}
	return tx.Model(p).Where(query).Association(assoc).Find(&out)
}
func (p *Post) DeleteAssociation(tx *gorm.DB, assoc string, del ...repository.AssociationEntity) error {
	return tx.Model(p).Association(assoc).Delete(del)
}

func (p *Post) AddAssociations(tx *gorm.DB, assoc string, inc ...repository.AssociationEntity) error {
	return tx.Model(p).Association(assoc).Append(inc)
}

func (p *Post) AssociationCount(tx *gorm.DB, assoc string, out repository.AssociationEntity) (int64, error) {
	return tx.Model(p).Association(assoc).Count(), nil
}

func (p *Post) ClearAssociations(tx *gorm.DB, assoc string) error {
	return tx.Model(p).Association(assoc).Clear()
}

func (p *Post) Raw(tx *gorm.DB, raw string, queryParams []any, opts *repository.RepositoryOpts) error {
	return tx.Raw(raw, queryParams...).Find(p).Error
}

// TODO: Fill out this method.
func (p *Post) AfterDelete(tx *gorm.DB) error {
	// Clean up Filestore
	// Clean up Media
	return nil
}
*/
