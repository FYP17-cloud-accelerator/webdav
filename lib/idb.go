package lib

type IDb interface {
	AddUser(user User)
	GetUser(username string, c *Config) (*User, bool)
	GetUserCount() int
	AddLog(logAccess *LogAccess) error
	UpdateAccess(logAccess *LogAccess) error
}

type DbUser struct {
	Username string
	Password string
	Scope    string
}
