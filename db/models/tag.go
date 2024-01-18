package models

import (
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name string `gorm:"unique;index;" json:"name"`
}

/*
func (t *Tag) ValidAssociation(srcType any, assoc string) bool {
	switch srcType.(type) {
	case Post:
		if val, ok := validPostAssociationKey[assoc]; ok && reflect.TypeOf(t) == val {
			return true
		}
	}
	return false
}

func (t *Tag) Query(tx *gorm.DB, query map[string]any, opts *repository.RepositoryOpts, out any) error {
	return tx.Model(t).Where(query).Find(out).Error
}

func (t *Tag) FindByID(tx *gorm.DB, id int, opts *repository.RepositoryOpts) error {
	return tx.First(t, id).Error
}

func (t *Tag) New(tx *gorm.DB) error {
	return tx.Create(t).Error
}

func (t *Tag) Update(tx *gorm.DB, id int, body map[string]any) error {
	return tx.Model(t).Where("id = ?", id).Updates(body).Error
}

func (t *Tag) Delete(tx *gorm.DB) error {
	return tx.Delete(t).Error
}

func (t *Tag) DeleteById(tx *gorm.DB, id int) error {
	return tx.Delete(t, id).Error
}

func (t *Tag) All(tx *gorm.DB, opts *repository.RepositoryOpts, out any) error {
	if opts != nil {
		for name, opt := range opts.Preloads {
			tx = tx.Preload(name, opt)
		}
	}
	return tx.Find(t).Error
}

// TODO: Either need to break associations into their own interface or just redo the repository
// desing. I'm leaning towards using the repository package as the "repository" and registering
// all of these functions as generic.
func (t *Tag) FindAssociations(tx *gorm.DB, assoc string, query map[string]any, out repository.AssociationEntity) error {
	return errors.New("tag does not have any associations")
}

func (t *Tag) DeleteAssociation(tx *gorm.DB, assoc string, del ...repository.AssociationEntity) error {
	return tx.Model(t).Association(assoc).Delete(del)
}

func (t *Tag) AddAssociations(tx *gorm.DB, assoc string, inc ...repository.AssociationEntity) error {
	return tx.Model(t).Association(assoc).Append(inc)
}

func (t *Tag) AssociationCount(tx *gorm.DB, assoc string, out repository.AssociationEntity) (int64, error) {
	return tx.Model(t).Association(assoc).Count(), nil
}

func (t *Tag) ClearAssociations(tx *gorm.DB, assoc string) error {
	return tx.Model(t).Association(assoc).Clear()
}

func (t *Tag) Raw(tx *gorm.DB, raw string, queryParams []any, opts *repository.RepositoryOpts) error {
	return tx.Raw(raw, queryParams...).Find(t).Error
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
*/
