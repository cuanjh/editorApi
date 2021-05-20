package rcluster

import (
	"strings"

	"github.com/go-redis/redis"
)

var rcClient *redis.ClusterClient

func init() {

	rcClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: strings.Split(
			redisConfig.String("redisClusterHosts"),
			",",
		),
		PoolSize: 100,
	})
}

func NewRedisClient() *redis.ClusterClient {
	if rcClient == nil {
		rcClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: strings.Split(
				redisConfig.String("redisClusterHosts"),
				",",
			),
			PoolSize: 1000,
		})
	}
	return rcClient
}
