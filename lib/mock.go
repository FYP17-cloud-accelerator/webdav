package lib

import "errors"

type mockDB struct {
	Users map[string]*User
}

func NewMockDB() IDb {
	return &mockDB{
		Users: map[string]*User{},
	}
}

func (db *mockDB) AddUser(user User) {
	db.Users[user.Username] = &user
}

func (db *mockDB) GetUser(username string, c *Config) (*User, bool) {
	user, ok := db.Users[username]
	return user, ok
}

func (db *mockDB) GetUserCount() int {
	return len(db.Users)
}

func (db *mockDB) AddLog(logAccess *LogAccess) error {
	return errors.New("not implemented")
}

func (db *mockDB) UpdateAccess(logAccess *LogAccess) error {
	return errors.New("not implemented")
}
