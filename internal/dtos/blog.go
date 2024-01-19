package dtos

import "github.com/masonictemple4/masonictempl/internal/customdate"

type BlogReturn struct {
	ID          uint                     `json:"id" yaml:"id"`
	CreatedAt   customdate.DefaultDate   `json:"createdat" yaml:"createdat"`
	UpdatedAt   customdate.DefaultDate   `json:"updatedat" yaml:"updatedat"`
	Title       string                   `json:"title" yaml:"title"`
	Description string                   `json:"description" yaml:"description"`
	Subtitle    string                   `json:"subtitle" yaml:"subtitle"`
	Thumbnail   string                   `json:"thumbnail" yaml:"thumbnail"`
	ContentUrl  string                   `json:"contenturl" yaml:"contenturl"`
	State       string                   `json:"state" yaml:"state"`
	Slug        string                   `json:"slug" yaml:"slug"`
	Tags        []TagReturn              `json:"tags" yaml:"tags"`
	Media       []MediaReturn            `json:"media" yaml:"media"`
	Authors     []BlogDetailAuthorReturn `json:"authors" yaml:"authors"`
	// The comments aren't getting included in the blog blog frontmatter so it does
	// not have to include the yaml tag.
	Comments []CommentReturn `json:"comments"`
}

type BlogInput struct {
	Title       string            `json:"title" yaml:"title" validate:"required" example:"This is a title."`
	Description string            `json:"description" yaml:"description" example:"This is a description."`
	Subtitle    string            `json:"subtitle" yaml:"subtitle" validate:"" example:"This is a subtitle."`
	ContentUrl  string            `json:"contenturl" yaml:"contenturl" validate:"required" example:"This is the source/path to the markdown file for the blog content."`
	Thumbnail   string            `json:"thumbnail" yaml:"thumbnail" validate:"" example:"https://storage.googleapis.com/masonictemple4-pub/images/test-profile-picture.jpeg"`
	Tags        []string          `json:"tags" yaml:"tags" validate:"" example:"['tag1', 'tag2']"`
	Media       []string          `json:"media" yaml:"media" validate:"" example:"['https://storage.google.com/...']"`
	Authors     []BlogAuthorInput `json:"authors" yaml:"authors" validate:""`
}

type UpdateBlogInput struct {
	Title       string   `json:"title" validate:"" example:"This is a title."`
	Description string   `json:"description" example:"This is a description."`
	Subtitle    string   `json:"subtitle" validate:"" example:"This is a subtitle."`
	ContentUrl  string   `json:"contenturl" validate:"" example:"This is the source/path to the markdown file for the blog content."`
	State       string   `json:"state" validate:"" example:"draft"`
	Tags        []string `json:"tags" validate:"" example:"['tag1', 'tag2']"`
}

type BlogAuthorInput struct {
	Username       string `json:"username" yaml:"username"`
	ProfilePicture string `json:"profilepicture" yaml:"profilepicture"`
}
