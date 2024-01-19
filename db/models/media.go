package models

import (
	"gorm.io/gorm"
)

const (
	MediaTypePhoto = "photo"
	MediaTypeVideo = "video"
	MediaTypeAudio = "audio"
)

func ValidMediaType(mt string) bool {
	switch mt {
	case MediaTypePhoto, MediaTypeVideo, MediaTypeAudio:
		return true
	}
	return false
}

// Might break these into separate models like video/photo
// Adding a unique index on url so we don't endup with 10000s
// of duplicate media objects.
type Media struct {
	gorm.Model
	MediaType string `gorm:"column:mediatype;" json:"mediatype"`
	Url       string `gorm:"column:url;uniqueIndex;" json:"url"`
	SmallUrl  string `gorm:"column:smallurl;" json:"smallurl"`
	MediumUrl string `gorm:"column:mediumurl;" json:"mediumurl"`
}

/*
func (m *Media) ValidAssociation(srcType any, assoc string) bool {

		switch srcType.(type) {
		case Blog:
			if val, ok := validBlogAssociationKey[assoc]; ok && val != nil && reflect.TypeOf(m) == val {
				return true
			}
		}
		return false
	}

	func (m *Media) All(tx *gorm.DB, opts *repository.RepositoryOpts, out any) error {
		if opts != nil {
			for name, opt := range opts.Preloads {
				tx = tx.Preload(name, opt)
			}
		}

		return tx.Find(out).Error
	}

// BEWARE!!! This method calls the method to make
// an update from the object istelf while we have
// the pointer.
//
// Update from the map is the recommended
// choice. However, if you understand the risks here
// I have decided to expose the UnsafeUpdate method

	func (m *Media) UnsafeUpdate(tx *gorm.DB) error {
		return m.update(tx)
	}

// WARNING!!! Because this uses the pointer to the object it could
// potentially cause weird behavior if something else is using that
// pointer and it changes etc..

	func (m *Media) update(tx *gorm.DB) error {
		return tx.Save(m).Error
	}
*/
func MediaFromStrings(tx *gorm.DB, input []string, out *[]Media) error {
	for _, url := range input {
		var media Media
		err := tx.FirstOrCreate(&media, Media{Url: url}).Error
		if err != nil {
			return err
		}
		*out = append(*out, media)
	}
	return nil
}
