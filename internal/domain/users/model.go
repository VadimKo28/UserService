package users

type User struct {
	Name string `json:"name" required:"true"`
	Email    string `json:"email" required:"true"`
	Password string `json:"password" required:"true"`
}
