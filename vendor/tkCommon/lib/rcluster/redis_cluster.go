package rcluster

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/gitstliu/go-redis-cluster"
)

func createConn() (*redis.Cluster, error) {

	hosts := redisConfig.String("redisClusterHosts")

	return redis.NewCluster(&redis.Options{
		StartNodes:   strings.Split(hosts, ","),
		ConnTimeout:  50 * time.Millisecond,
		ReadTimeout:  50 * time.Millisecond,
		WriteTimeout: 50 * time.Millisecond,
		KeepAlive:    16,
		AliveTime:    60 * time.Second,
	})
}

type Rcluster struct {
	cluster *redis.Cluster
	closed  bool
}

func NewRcluster() *Rcluster {
	conn, err := createConn()
	if err != nil {
		log.Println("连接redis集群错误：", err)
		return nil
	}
	return &Rcluster{
		cluster: conn,
		closed:  false,
	}
}

func (this *Rcluster) Do(action string, params ...interface{}) (interface{}, error) {
	if this.closed {
		return nil, errors.New("此类库已经关闭")
	}
	return this.cluster.Do(action, params...)
}

func (this *Rcluster) Close() {
	this.closed = true
	this.cluster.Close()
}

//设置key value

func (this *Rcluster) Set(key, val string) (string, error) {
	return redis.String(this.Do("SET", key, val))
}

//设置带过期时间的key value
func (this *Rcluster) SetEx(key, val string, exp int) (string, error) {
	return redis.String(this.Do("SETEX", key, exp, val))
}

//设置带过期时间的key value
func (this *Rcluster) SetNx(key, val string) (string, error) {

	return redis.String(this.Do("SETNX", key, val))
}

func (this *Rcluster) Get(key string) (string, error) {

	return redis.String(this.Do("GET", key))
}

func (this *Rcluster) Exists(key string) (bool, error) {
	return redis.Bool(this.Do("EXISTS", key))
}

func (this *Rcluster) Del(key ...interface{}) {
	this.Do("DEL", key...)
}

func (this *Rcluster) Incr(key string) (int, error) {

	return redis.Int(this.Do("INCR", key))

}

func (this *Rcluster) Incrby(key string, inc int) (int, error) {

	return redis.Int(this.Do("INCRBY", key, inc))

}

func (this *Rcluster) Decr(key string) (int, error) {
	return redis.Int(this.Do("Decr", key))
}

func (this *Rcluster) Decrby(key string, inc int) (int, error) {
	return redis.Int(this.Do("DECRBY", key, inc))
}

func (this *Rcluster) HSet(hashName, key string, val interface{}) {
	this.Do("HSET", hashName, key, val)
}

func (this *Rcluster) HGet(hashName, key string) (interface{}, error) {
	return this.Do("HGET", hashName, key)
}

func (this *Rcluster) HMSet(params ...interface{}) (interface{}, error) {
	return this.Do("HMSET", params...)
}

func (this *Rcluster) HMGet(params ...interface{}) (interface{}, error) {
	return this.Do("HMGET", params...)
}

func (this *Rcluster) HGetAll(hashName string) (interface{}, error) {
	return this.Do("HMGETALL", hashName)
}

func (this *Rcluster) HDel(hashName string, key string) (interface{}, error) {
	return this.Do("HDEL", hashName, key)
}

func (this *Rcluster) HKeys(hashName string) ([]interface{}, error) {
	return redis.Values(this.Do("HKeys", hashName))
}

func (this *Rcluster) HVals(hashName string) ([]interface{}, error) {
	return redis.Values(this.Do("HVALS", hashName))
}

func (this *Rcluster) HLen(hashName string) (int, error) {
	return redis.Int(this.Do("HMGETALL", hashName))
}

func (this *Rcluster) LLPush(params ...interface{}) (interface{}, error) {
	return this.Do("LPUSH", params...)
}

func (this *Rcluster) LRPush(params ...interface{}) (interface{}, error) {
	return this.Do("RPUSH", params...)
}

func (this *Rcluster) LLRange(listName string, startIndex, endIndex int) ([]string, error) {
	return redis.Strings(this.Do("LRANGE", listName, startIndex, endIndex))
}

func (this *Rcluster) LLRem(listName string, count int, val interface{}) (int, error) {
	return redis.Int(this.Do("LREM", listName, count, val))
}

func (this *Rcluster) LLPop(listName string) (interface{}, error) {
	return this.Do("LPOP", listName)
}

func (this *Rcluster) LRPop(listName string) (interface{}, error) {
	return this.Do("RPOP", listName)
}

func (this *Rcluster) LLTrim(listName string, startIndex, endIndex int) ([]string, error) {
	return redis.Strings(this.Do("LTRIM", listName, startIndex, endIndex))
}

func (this *Rcluster) LRPopLPush(listName1, listName2 string) (interface{}, error) {
	return this.Do("RPOPLPUSH", listName1, listName2)
}

func (this *Rcluster) LLen(key string) (int, error) {
	return redis.Int(this.Do("LLEN", key))
}

func (this *Rcluster) SADD(params ...interface{}) (int, error) {
	return redis.Int(this.Do("SADD", params...))
}

func (this *Rcluster) SMembers(setName string) ([]string, error) {
	return redis.Strings(this.Do("SMEMEBERS", setName))
}

func (this *Rcluster) SRem(params ...interface{}) (interface{}, error) {
	return this.Do("SREM", params...)
}
func (this *Rcluster) SPop(setName string) (string, error) {
	return redis.String(this.Do("SPOP", setName))
}
func (this *Rcluster) SMemberNum(setName string) (int, error) {
	return redis.Int(this.Do("SCARD", setName))
}

func (this *Rcluster) SIsMember(setName string, mem interface{}) (bool, error) {
	return redis.Bool(this.Do("SISMEMBER", setName, mem))
}

func (this *Rcluster) ZAdd(params ...interface{}) (interface{}, error) {
	return this.Do("ZADD", params...)
}

func (this *Rcluster) ZRange(zSetName string, start, end int) ([]string, error) {
	return redis.Strings(this.Do("ZRANGE", zSetName, start, end, "WITHSCORES"))
}

func (this *Rcluster) ZRevRange(zSetName string, start, end int) ([]string, error) {
	return redis.Strings(this.Do("ZREVRANGE", zSetName, start, end, "WITHSCORES"))
}

func (this *Rcluster) ZRangeByScore(zSetName string, score1, score2 int) ([]string, error) {
	return redis.Strings(this.Do("ZRANGEBYSCORE", zSetName, score1, score2, "WITHSCORES"))
}

func (this *Rcluster) ZRevRangeByScore(zSetName string, score1, score2 int) ([]string, error) {
	return redis.Strings(this.Do("ZREVRANGEBYSCORE", zSetName, score1, score2, "WITHSCORES"))
}

func (this *Rcluster) ZScore(zSetName string, mem interface{}) (string, error) {
	return redis.String(this.Do("ZSCORE", zSetName, mem))
}

func (this *Rcluster) ZMemberNum(zSetName string) (int, error) {
	return redis.Int(this.Do("ZCARD", zSetName))
}

func (this *Rcluster) ZCount(zSetName string, min, max int) (int, error) {
	return redis.Int(this.Do("ZCOUNT", zSetName, min, max))
}

func (this *Rcluster) ZRank(zSetName string, member interface{}) (int, error) {
	return redis.Int(this.Do("ZRANK", zSetName, member))
}

func (this *Rcluster) ZRevRank(zSetName string, member interface{}) (int, error) {
	return redis.Int(this.Do("ZREVRANK", zSetName, member))
}

func (this *Rcluster) ZRem(zSetName string, member interface{}) (int, error) {
	return redis.Int(this.Do("ZREM", zSetName, member))
}

func (this *Rcluster) ZRemRangeByRank(zSetName string, startIndex, endIndex int) (int, error) {
	return redis.Int(this.Do("ZREMRANGEBYRANK", zSetName, startIndex, endIndex))
}

func (this *Rcluster) ZRemRangeByScore(zSetName string, min, max int) (int, error) {
	return redis.Int(this.Do("ZREMRANGEBYSCORE", zSetName, min, max))
}
