package dtos

import "github.com/masonictemple4/masonictempl/internal/customdate"

type UserReturn struct {
	ID             uint                   `json:"id"`
	CreatedAt      customdate.DefaultDate `json:"createdat"`
	UpdatedAt      customdate.DefaultDate `json:"updatedat"`
	Username       string                 `json:"username"`
	Firstname      string                 `json:"firstname"`
	Lastname       string                 `json:"lastname"`
	ProfilePicture string                 `json:"profilepicture"`
}

// TODO: We can expand on this in future or just replace it with UserReturn
type BlogDetailAuthorReturn struct {
	Username       string `json:"username"`
	ProfilePicture string `json:"profilepicture"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required" example:"test@example.com"`
	Password string `json:"password" validate:"required" example:"password123"`
}

type RegisterInput struct {
	Email           string `json:"email" validate:"required" example:"test@example.com"`
	Password        string `json:"password" validate:"required" example:"password123"`
	ConfirmPassword string `json:"confirmpassword" validate:"required" example:"password123"`
}

type UpdateUserInput struct {
	// TODO: Populate update user input.
}

type GithubAuthInput struct {
	Token          string `json:"token"`
	Email          string `json:"email"`
	AccountID      string `json:"accountid"`
	ProfilePicture string `json:"profilepicture"`
}

type GoogleAuthInput struct {
	IdToken        string `json:"idtoken"`
	Token          string `json:"token"`
	Email          string `json:"email"`
	AccountID      string `json:"accountid"`
	ProfilePicture string `json:"profilepicture"`
	ExpiresAt      int64  `json:"expires"`
}
