package config

type RedisConfiguration struct {
	Redishost string
	Redisport string
}

type MongodbConfiguration struct {
	Url      string
	Database string
}
