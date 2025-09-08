package dto

type UserInput struct {
	Name  string `json:"name" binding:"required" validate:"required,min=2,max=255"`
	Email string `json:"email" binding:"required" validate:"required,email,max=255"`
}

type UserRegistrationInput struct {
	Name            string `json:"name" binding:"required" validate:"required,min=2,max=255"`
	Email           string `json:"email" binding:"required" validate:"required,email,max=255"`
	Password        string `json:"password" binding:"required" validate:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" binding:"required" validate:"required,eqfield=Password"`
}

type UserLoginInput struct {
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required,min=6"`
}

type UserChangePasswordInput struct {
	CurrentPassword string `json:"current_password" binding:"required" validate:"required"`
	NewPassword     string `json:"new_password" binding:"required" validate:"required,min=8,max=72,nefield=CurrentPassword"`
	ConfirmPassword string `json:"confirm_password" binding:"required" validate:"required,eqfield=NewPassword"`
}
