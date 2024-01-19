package models

import (
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name string `gorm:"unique;index;" json:"name"`
}

func TagFromStrings(tx *gorm.DB, input []string, out *[]Tag) error {
	for _, t := range input {
		var tag Tag
		err := tx.FirstOrCreate(&tag, Tag{Name: t}).Error
		if err != nil {
			return err
		}
		*out = append(*out, tag)
	}

	return nil
}
