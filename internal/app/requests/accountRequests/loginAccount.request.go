package accountRequests

import (
	"errors"
	"time"
)

type LoginAccountRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=16,passwd"`
}

type LoginAccountResponse struct {
	AccessToken string    `json:"access_token"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	IsValid     bool      `json:"is_valid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var (
	ErrEmailOrPasswordNotCorrect = errors.New("email or password not correct")
)
