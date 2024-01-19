package dtos

type TagReturn struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type TagInput struct {
	Name string `json:"name" validate:"required" example:"tagname"`
}
