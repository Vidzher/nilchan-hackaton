package auth

type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type AuthResponseDTO struct {
	UserID   int    `json:"userId"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Token    string `json:"token"`
}
