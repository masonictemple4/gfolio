package dtos

import "github.com/masonictemple4/masonictempl/internal/customdate"

type CommentReturn struct {
	ID        uint                   `json:"id"`
	CreatedAt customdate.DefaultDate `json:"createdat"`
	UpdatedAt customdate.DefaultDate `json:"updatedat"`
	User      UserReturn             `json:"user"`
	Text      string                 `json:"text"`
}

type CommentInput struct {
	BlogID uint   `json:"blogid" validate:"required" example:"1"`
	Text   string `json:"text" validate:"required" example:"This is a comment."`
}

// TODO: Not sure we'll need an update comment just yet.
