package auth

type AuthRequestDTO struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponseDTO struct {
	UserId int    `json:"userId"`
	Email  string `json:"email"`
	Token  string `json:"token"`
}
