package lib

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/webdav"
)

type sqlDb struct {
	host     string
	db       string
	user     string
	password string
}

func NewSqlDb(host string, db string, user string, password string) IDb {
	return &sqlDb{
		user:     user,
		password: password,
		host:     host,
		db:       db,
	}
}

func (db *sqlDb) getClient() (*sql.DB, error) {
	client, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", db.user, db.password, db.host, db.db))
	return client, err
}

func (db *sqlDb) AddUser(user User) {
	client, err := db.getClient()
	if err != nil {
		return
	}
	defer client.Close()

	_, err = client.Exec("INSERT INTO users(username, password, scope) VALUES(?, ?, ?)", user.Username, user.Password, user.Scope)
	if err != nil {
		return
	}
	log.Println(user)
}

func (db *sqlDb) GetUser(username string, c *Config) (*User, bool) {
	client, err := db.getClient()
	if err != nil {
		return nil, true
	}
	defer client.Close()

	dbuser := DbUser{}
	err = client.QueryRow("SELECT username, password, scope FROM user WHERE username = ?", username).Scan(&dbuser.Username, &dbuser.Password, &dbuser.Scope)
	if err != nil {
		return nil, true
	}
	user := User{
		Username: dbuser.Username,
		Password: dbuser.Password,
		Scope:    dbuser.Scope,
		Modify:   c.User.Modify,
		Rules:    c.User.Rules,
		Handler: &webdav.Handler{
			Prefix: c.User.Handler.Prefix,
			FileSystem: WebDavDir{
				Dir:     webdav.Dir(dbuser.Scope),
				NoSniff: c.NoSniff,
			},
			LockSystem: webdav.NewMemLS(),
		},
	}
	return &user, true
}

func (db *sqlDb) GetUserCount() int {
	client, err := db.getClient()
	if err != nil {
		return 0
	}
	defer client.Close()

	var count int
	err = client.QueryRow("SELECT COUNT(*) FROM user").Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

func (db *sqlDb) AddLog(logAccess *LogAccess) error {
	client, err := db.getClient()
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.Exec(
		"INSERT INTO logs(user_id, filename, path, extension, access_time, mod_time, size) VALUES((SELECT user_id FROM user WHERE username = ?), ?, ?, ?, ?, ?, ?)",
		logAccess.Username, logAccess.FileName, logAccess.Path, logAccess.Extension, logAccess.AccessTime, logAccess.ModTime, logAccess.Size)
	return err
}

func (db *sqlDb) UpdateAccess(logAccess *LogAccess) error {
	client, err := db.getClient()
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.Exec(
		"REPLACE INTO logs(user_id, full_path, access_time) VALUES((SELECT user_id FROM user WHERE username = ?), ?, ?)",
		logAccess.Username, logAccess.FullPath, logAccess.AccessTime)
	return err
}
