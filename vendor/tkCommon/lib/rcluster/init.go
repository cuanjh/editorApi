package rcluster

import (
	"github.com/astaxie/beego/config"
)

var redisConfig config.Configer

func init() {
	redisConfig, _ = config.NewConfig(
		"ini",
		"/usr/local/consul-template/conf/goConf/redis-cluster.conf",
	)
}
