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
