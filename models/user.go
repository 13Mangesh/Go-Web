package models

import (
	"fmt"
	"errors"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrUserNotFound : Custom error if user is not found
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidLogin : Custom error if credentials are wrong 
	ErrInvalidLogin = errors.New("invalid login")
)

// User : User model
type User struct {
	key string
}

// NewUser : Constructor for creating a new user
func NewUser(username string, hash []byte) (*User, error) {
	id, err := client.Incr(ctx, "user:next-id").Result()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("user:%d", id)
	pipe := client.Pipeline()
	pipe.HSet(ctx, key, "id", id)
	pipe.HSet(ctx, key, "username", username)
	pipe.HSet(ctx, key, "hash", hash)
	pipe.HSet(ctx, "user:by-username", username, id)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &User{key}, nil
}

// GetID : Helper method to get the Id
func (user *User) GetID() (int64, error) {
	return client.HGet(ctx, user.key, "id").Int64()
}

// GetUsername : Helper method to get the username
func (user *User) GetUsername() (string, error) {
	return client.HGet(ctx, user.key, "username").Result()
}

// GetHash : Helper method to get the hash
func (user *User) GetHash() ([]byte, error) {
	return client.HGet(ctx, user.key, "hash").Bytes()
}

// Authenticate : Authentication of user is performed 
func (user *User) Authenticate(password string) (error) {
	hash, err := user.GetHash()
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrInvalidLogin
	}
	return err
}

// GetUserByID : Helper ,ethod to get user by ID
func GetUserByID(id int64) (*User, error) {
	key := fmt.Sprintf("user:%d", id)
	return &User{key}, nil
}

// GetUserByUsername : Helper method to get user by  username
func GetUserByUsername(username string) (*User, error) {
	id, err := client.HGet(ctx, "user:by-username", username).Int64()
	if err == redis.Nil {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	return GetUserByID(id)
}

// AuthenticateUser : User login authentication is performed
func AuthenticateUser(username, password string) (*User, error) {
	user, err := GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, user.Authenticate(password)
}

// RegisterUser : Registers a user into the redis
func RegisterUser(username, password string) (error) {
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return err
	}
	_, err = NewUser(username, hash)
	return err
}