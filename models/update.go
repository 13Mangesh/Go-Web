package models

import (
	"fmt"
)

// Update : Model for a update
type Update struct {
	key string
}

// NewUpdate : Constructor for creating a new update
func NewUpdate(userID int64, body string) (*Update, error) {
	id, err := client.Incr(ctx, "update:next-id").Result()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("update:%d", id)
	pipe := client.Pipeline()
	pipe.HSet(ctx, key, "id", id)
	pipe.HSet(ctx, key, "user_id", userID)
	pipe.HSet(ctx, key, "body", body)
	pipe.LPush(ctx, "updates", id)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &Update{key}, nil
}

// GetBody : Helper method to get body of update
func (update *Update) GetBody() (string, error) {
	return client.HGet(ctx, update.key, "body").Result()
}

// GetUser : Hewlper method to get the user who posted the Update
func (update *Update) GetUser() (*User, error) {
	userID, err := client.HGet(ctx, update.key, "user_id").Int64()
	if err != nil {
		return nil, err
	}
	return GetUserByID(userID)
}

// Getupdates : To get updates slice from redis
func Getupdates() ([]*Update, error) {
	updateIDs, err := client.LRange(ctx,"updates", 0, 10).Result()
	if err != nil {
		return nil, err
	}
	updates := make([]*Update, len(updateIDs))
	for i, id := range updateIDs {
		key := "update:" + id
		updates[i] = &Update{key}
	}
	return updates, nil
}

// Postupdate : To add a update to redis
func Postupdate(userID int64, body string) (error) {
	_, err := NewUpdate(userID, body)
	return err
}