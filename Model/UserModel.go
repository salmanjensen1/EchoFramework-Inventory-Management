package Model

type User struct {
	ID       int    `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
	Address  string `json:"address,omitempty"`
	IsAdmin  bool   `json:"isAdmin,omitempty"`
}
