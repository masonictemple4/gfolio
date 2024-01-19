package dtos

import "github.com/masonictemple4/masonictempl/internal/customdate"

type MediaReturn struct {
	ID        uint                   `json:"id"`
	CreatedAt customdate.DefaultDate `json:"createdat"`
	UpdatedAt customdate.DefaultDate `json:"updatedat"`
	MediaType string                 `json:"mediatype"`
	Url       string                 `json:"url"`
	SmallUrl  string                 `json:"smallurl"`
	MediumUrl string                 `json:"mediumurl"`
}

// TODO: This will probably be more of a mpfd
type MediaInput struct {
}

type UpdateMediaInput struct {
}
