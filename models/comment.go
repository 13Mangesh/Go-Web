package models

// GetComments : To get comments slice from redis
func GetComments() ([]string, error) {
	return client.LRange(ctx,"comments", 0, 10).Result()
}

// PostComment : To add a comment to redis
func PostComment(comment string) (error) {
	return client.LPush(ctx, "comments", comment).Err()
}