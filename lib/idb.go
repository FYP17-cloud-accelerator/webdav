package lib

type IDb interface {
	AddUser(user User)
	GetUser(username string, c *Config) (*User, bool)
	GetUserCount() int
}
