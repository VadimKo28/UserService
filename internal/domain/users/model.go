package users

type UserCreateDTO struct {
	Name     string 
	Email    string 
	Password string 
}

type UserSignInDTO struct {
	Email    string 
	Password string 
}

type User struct {
	ID       int
	Name     string 
	Email    string 
}
