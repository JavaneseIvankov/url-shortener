package repository

type User struct {
	ID string
	Email string
	Password string
	Role string
}

type IUserRepository interface {
	CreateUser(user User) error 
	UpdateUser(user User) error
	GetUserByEmail(email string) error
	GetUserByID(id string) error
}

type UserRepositoryImpl struct {
}