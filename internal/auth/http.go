package auth

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8,max=72,bcrypt_max_bytes"`
}

type AuthResponse struct {
	UserID   int64  `json:"userId"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Token    string `json:"token"`
}
