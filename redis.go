package redis_client

import (
	"errors"
	"os"
	"strconv"
	"time"

	rds "github.com/go-redis/redis"
)

type redis struct {
	Connection *rds.Client
}

func newRedisConnection(database string) *redis {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	switch database {
	case "session":
		db := os.Getenv("REDIS_SESSION_DB")
		dbInt, _ := strconv.Atoi(db)
		return &redis{
			Connection: rds.NewClient(&rds.Options{
				Addr: host + ":" + port,
				DB:   dbInt,
			}),
		}
	default:
		return nil
	}
}

func NewRedisSessionClient() Redis {
	return newRedisConnection("session")
}

func (r redis) reDial() bool {
	for i := 0; i < 3; i++ {
		var err = r.Connection.Ping().Err()
		if err == nil {
			return true
		}
		time.Sleep(1 * time.Second)
	}
	return false
}

func (r redis) Set(key string, value interface{}, expiration time.Duration) error {
	if !r.reDial() {
		return errors.New("No redial")
	}
	result := r.Connection.Set(key, value, expiration)
	if result.Val() != "OK" {
		return errors.New(result.Val())
	}
	return nil
}

func (r redis) Get(key string) (string, error) {
	if !r.reDial() {
		return "", errors.New("No redial")
	}
	result := r.Connection.Get(key)
	iresult, err := result.Result()
	if err == rds.Nil {
		return "", nil
	}
	if err != nil && err != rds.Nil {
		return "", err
	}

	return iresult, nil
}

func (r redis) Exists(keys ...string) (bool, error) {
	if !r.reDial() {
		return false, errors.New("No redial")
	}
	result := r.Connection.Exists(keys...)
	if result.Val() == 1 {
		return true, nil
	}
	return false, nil
}

func (r redis) Del(keys ...string) (bool, error) {
	if !r.reDial() {
		return false, errors.New("No redial")
	}
	result := r.Connection.Del(keys...)
	if result.Val() == 1 {
		return true, nil
	}
	return false, nil
}
