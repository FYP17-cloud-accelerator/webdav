package lib

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/webdav"
)

type redisDb struct {
	host    string
	db      int
	expires time.Duration
}

type DbUser struct {
	Username string
	Password string
	Scope    string
}

// Check if cache is connected
func ping(client *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	return err
}

func NewRedisDb(host string, db int, exp time.Duration) IDb {
	return &redisDb{
		host:    host,
		db:      db,
		expires: exp,
	}
}

func (db *redisDb) getClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:        db.host,
		Password:    "",
		DB:          db.db,
		ReadTimeout: 10 * time.Second,
	})
	err := ping(client)
	return client, err

}
func (db *redisDb) AddUser(user User) {
	client, err := db.getClient()
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jsn, err := json.Marshal(DbUser{
		Username: user.Username,
		Password: user.Password,
		Scope:    user.Scope,
	})
	if err != nil {
		log.Print(err.Error())
		return
	}
	log.Print(string(jsn))
	client.HSet(ctx, "users", user.Username, jsn)
}

func (db *redisDb) GetUser(username string, c *Config) (*User, bool) {
	client, err := db.getClient()
	if err != nil {
		log.Print(err.Error())
		return nil, true
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := client.HGet(ctx, "users", username).Result()
	if err != nil {
		log.Print(err.Error())
		return nil, false
	}
	dbuser := DbUser{}
	json.Unmarshal([]byte(result), &dbuser)
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

func (db *redisDb) GetUserCount() int {
	client, err := db.getClient()
	if err != nil {
		log.Print(err.Error())
		return 0
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := client.HGetAll(ctx, "users").Result()
	if err != nil {
		log.Print(err.Error())
		return 0
	}
	return len(result)
}

func (db *redisDb) AddLog(logAccess *LogAccess) error {
	return errors.New("not implemented")
}
