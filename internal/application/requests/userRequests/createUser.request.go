package userRequests

import "time"

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type CreateUserResponse struct {
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	FullName         string    `json:"full_name"`
	PasswordChangeAt time.Time `json:"password_changed_at"`
	CreatedAt        time.Time `json:"created_at"`
}
