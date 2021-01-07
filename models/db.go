package models

import (
	"github.com/go-redis/redis"
	"context"
)


var ctx = context.Background()
var client *redis.Client

// Init : Initialize the redis client
func Init() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}