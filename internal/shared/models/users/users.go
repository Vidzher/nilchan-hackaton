package users

type User struct {
	ID       int    `json:"userId"`
	Email    string `json:"email"`
	Balance  int    `json:"-"`
	Password string `json:"-"`
}
