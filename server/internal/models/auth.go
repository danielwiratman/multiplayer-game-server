package models

import "time"

type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=32"`
	LastName  string `json:"last_name" validate:"required,min=2,max=32"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=128"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=128"`
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	UserID       int
	ExpiresIn    time.Time
}

type LoginResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresIn   time.Time `json:"expires_in"`
}

type RefreshTokenCacheVal struct {
	UserID    int
	ExpiresIn time.Time
}

type UserInfoFromToken struct {
	UserID int
}

type Account struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
