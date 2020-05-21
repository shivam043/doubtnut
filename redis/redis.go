package redis

import (
	"errors"
	"github.com/doubtnut/logging"
	"github.com/go-redis/redis"
	"os"
	"time"
)

var client *redis.Client

var Logger = logging.NewLogger()

func Init() {
	SetUpRedisClient()
}

// create connection with redis
func SetUpRedisClient() {
	// close existing connection
	if client != nil {
		client.Close()
	}

	host := os.Getenv("REDIS_HOST")

	if host == "" {
		host = "127.0.0.1:6379"
	}

	Logger.Printf("redis-connect - Connecting to server : %s", host)

	// create client for redis
	client = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "",
		DB:       0,
	})

	// check client connection using ping-pong api
	_, err := client.Ping().Result()

	// if error in connecting, then wait for 3 seconds and retry connection
	if err != nil {
		Logger.Errorf("redis-connect - Redis connection failed with reason : %s, trying again in 3 seconds", err.Error())
		time.Sleep(time.Duration(3))

	}

}

// find value for a given key
func GetValue(key string) (string, error) {
	value, err := client.Get(key).Result()

	if err == redis.Nil {
		return "", errors.New("key (" + key + ") does not exist")
	} else if err != nil {
		return "", errors.New("Error while finding redis key  (" + key + ") : " + err.Error())
	} else {
		return value, nil
	}
}

// set value for a given key
func SetValue(key string, value string, timeout int) bool {
	err := client.Set(key, value, time.Duration(timeout)).Err()
	if err != nil {
		Logger.Errorf("Error while setting redis value : %s key (%s), value (%s)", err.Error(), key, value)
		return false
	}

	return true
}

func Keys(key string) *redis.StringSliceCmd {
	return client.Keys(key)
}

func Del(key string) {
	client.Del(key)
}
