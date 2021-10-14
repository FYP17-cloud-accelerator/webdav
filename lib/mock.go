package lib

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

func (db *mockDB) GetUser(username string) (*User, bool) {
	user, ok := db.Users[username]
	return user, ok
}

func (db *mockDB) GetUsers() map[string]*User {
	return db.Users
}
