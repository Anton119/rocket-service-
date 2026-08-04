package redis_view

// SessionRedisView — представление сессии для Redis HashMap.
type SessionRedisView struct {
	UUID      string `redis:"uuid"`
	UserUUID  string `redis:"user_uuid"`
	Login     string `redis:"login"`
	CreatedAt string `redis:"created_at"`
	ExpiresAt string `redis:"expires_at"`
}
