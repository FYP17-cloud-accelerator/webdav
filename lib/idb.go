package lib

type IDb interface {
	AddUser(user User)
	GetUser(username string) (*User, bool)
	GetUsers() map[string]*User
}
