package storage

// repositories

type RedisRepository interface {
	RedisSet(key string, value interface{}) error
	RedisGet(key string) ([]byte, error)
	RedisDelete(key string) (int64, error)
}

// repositories
