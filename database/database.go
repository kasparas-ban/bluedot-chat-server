package database

import (
	"time"

	"github.com/go-redis/redis"
)

type Connection struct {
	ConnectionId string `json:"sessionId" gorm:"primarykey"`
	UserId       uint   `json:"userId"`
}

type database struct {
	Client *redis.Client
}

var ConnectionsDB = &database{}

func (c *database) Connect(url string) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		panic("Failed to connect to Redis cache")
	}

	c.Client = redis.NewClient(opt)

	if err := c.Client.Ping().Err(); err != nil {
		panic(err)
	}
}

func AddConn(username string) error {
	return ConnectionsDB.Client.SetNX(username, "serverID?", 15*time.Minute).Err()
}

func ReadConn(username string) (string, error) {
	val, err := ConnectionsDB.Client.Get(username).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
