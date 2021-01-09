package models

import (
	"strconv"
	"fmt"
)

// Update : Model for a update
type Update struct {
	id int64
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
	pipe.LPush(ctx, fmt.Sprintf("user:%d:updates", userID), id)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &Update{id}, nil
}

// GetBody : Helper method to get body of update
func (update *Update) GetBody() (string, error) {
	key := fmt.Sprintf("update:%d", update.id)
	return client.HGet(ctx, key, "body").Result()
}

// GetUser : Hewlper method to get the user who posted the Update
func (update *Update) GetUser() (*User, error) {
	key := fmt.Sprintf("update:%d", update.id)
	userID, err := client.HGet(ctx, key, "user_id").Int64()
	if err != nil {
		return nil, err
	}
	return GetUserByID(userID)
}

func queryUpdates(key string) ([]*Update, error) {
	updateIDs, err := client.LRange(ctx, key, 0, 10).Result()
	if err != nil {
		return nil, err
	}
	updates := make([]*Update, len(updateIDs))
	for i, strID := range updateIDs {
		id, err := strconv.Atoi(strID)
		if err != nil {
			return nil, err
		}
		updates[i] = &Update{int64(id)}
	}
	return updates, nil
}

// GetAllupdates : To get updates slice from redis
func GetAllupdates() ([]*Update, error) {
	return queryUpdates("updates")
}

// Getupdates : To get updates slice from redis for a particular user
func Getupdates(userID int64) ([]*Update, error) {
	key := fmt.Sprintf("user:%d:updates", userID)
	return queryUpdates(key)
}

// Postupdate : To add a update to redis
func Postupdate(userID int64, body string) (error) {
	_, err := NewUpdate(userID, body)
	return err
}