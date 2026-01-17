package users

type UserCreateDTO struct {
	Name     string 
	Email    string 
	Password string 
}

type GetUserDTO struct {
	ID       int
	Name     string 
	Email    string 
}
