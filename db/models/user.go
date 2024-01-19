package models

import (
	"github.com/masonictemple4/masonictempl/internal/dtos"
	"github.com/masonictemple4/masonictempl/internal/utils"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	// This is really the display name, not used to login it's just what's displayed on the blogs.
	Username       string `gorm:"column:username;uniqueIndex;" json:"username"`
	Password       string `gorm:"column:password;" json:"-"`
	Firstname      string `gorm:"column:firstname;" json:"firstname"`
	Lastname       string `gorm:"column:lastname;" json:"lastname"`
	Email          string `gorm:"column:email;uniqueIndex;" json:"email"`
	ProfilePicture string `gorm:"column:profilepicture;" json:"profilepicture"`
	Logintype      string `gorm:"column:logintype;" json:"logintype"`
}

func AuthorFromInput(tx *gorm.DB, input []dtos.BlogAuthorInput, out *[]User) error {
	var authors []User
	if err := utils.Convert(input, &authors); err != nil {
		return nil
	}
	for _, author := range authors {
		var user User
		err := tx.FirstOrCreate(&user, User{Username: author.Username, ProfilePicture: author.ProfilePicture}).Error
		if err != nil {
			return err
		}
		*out = append(*out, user)
	}
	return nil
}
