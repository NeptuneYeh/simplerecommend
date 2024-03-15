package accountRequests

import "time"

type CreateAccountRequest struct {
	Name     string `json:"name" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6,max=16,passwd"`
	Email    string `json:"email" binding:"required,email"`
}

type AccountResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	IsValid   bool      `json:"is_valid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
