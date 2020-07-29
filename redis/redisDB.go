package redis

import (
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

//Redis ...
type Redis interface {
	CreateToken(username string, lifeDuration time.Duration) (string, error)
	CheckToken(username, token string) (bool, error)
}

type dbRedis struct {
	db redis.Conn
}

//InitRedis ...
func InitRedis(network, host string) (Redis, error) {
	conn, err := redis.Dial(network, host)
	if err != nil {
		log.Fatal(err)
	}

	return &dbRedis{db: conn}, err
}

func (rds *dbRedis) CreateToken(username string, lifeDuration time.Duration) (string, error) {
	token := createUniqueToken()

	_, err := rds.db.Do("SET", username, token, "EX", strconv.Itoa(int(lifeDuration.Seconds())))
	if err != nil {
		fmt.Println(err.Error())
	}
	return token, err
}
func (rds *dbRedis) CheckToken(username, token string) (bool, error) {
	replay, err := rds.db.Do("GET", username)
	if err == nil && replay == nil {
		return false, nil
	}
	tkn, err := redis.String(replay, err)
	return tkn == token, err
}

func createUniqueToken() string {
	u := uuid.New()
	bt := u[:]
	return hex.EncodeToString(bt)
}
