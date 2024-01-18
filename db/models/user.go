package models

import (
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

/*
func (u *User) All(tx *gorm.DB, opts *repository.RepositoryOpts, out any) error {
	if opts != nil {
		for name, opt := range opts.Preloads {
			tx = tx.Preload(name, opt)
		}
	}
	return tx.Find(out).Error
}
func (u *User) FindByUsername(tx *gorm.DB, uname string, opts *repository.RepositoryOpts) error {
	return tx.First(u, "username = ?", uname).Error
}

// BEWARE!!! This method calls the method to make
// an update from the object istelf while we have
// the pointer.
//
// Update from the map is the recommended
// choice. However, if you understand the risks here
// I have decided to expose the UnsafeUpdate method
func (u *User) UnsafeUpdate(tx *gorm.DB) error {
	return u.update(tx)
}

// WARNING!!! Because this uses the pointer to the object it could
// potentially cause weird behavior if something else is using that
// pointer and it changes etc..
func (u *User) update(tx *gorm.DB) error {
	return tx.Save(u).Error
}

// Updates an object from map, this ensures we don't have any stale state
// in our object and instead just updates the fields that exist in body.
func (u *User) Update(tx *gorm.DB, id int, body map[string]any) error {
	return tx.Model(u).Where("id = ?", id).Updates(body).Error
}

func (u *User) ValidAssociation(srcType any, assoc string) bool {
	switch srcType.(type) {
	case Blog:
		if val, ok := validBlogAssociationKey[assoc]; ok && reflect.TypeOf(u) == val {
			return true
		}
	}
	return false
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

func (u *User) GenerateProfilePicturePath(ext string) string {
	switch ext {
	case "jpg", "jpeg", "png", "gif":
	default:
		return ""
	}
	return fmt.Sprintf("users/%s-%d/%s.%s", u.Username, u.ID, "profilepic", ext)
}

func (u *User) GenerateProfilePictureURL(path string) string {
	baseurl := os.Getenv("BUCKET_BASE_URL")
	bucketName := os.Getenv("STORAGE_BUCKET")
	return fmt.Sprintf("%s/%s/%s", baseurl, bucketName, path)
}
*/
