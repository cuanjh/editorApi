package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	MysqlAdmin       MysqlAdmin
	CourseConfig     CourseConfig
	Qiniu            Qiniu
	CasbinConfig     CasbinConfig
	RedisAdmin       RedisAdmin
	System           System
	Mongodb          Mongodb
	Tomongodb        Tomongodb
	Testmongodb      Testmongodb
	JWT              JWT
	NatsConfig       NatsConfig
	LiveCourseConfig LiveCourseConfig
	Tencent          Tencent
	UploadConfig     UploadConfig
}

type System struct {
	UseMultipoint bool
	Env           string
}

type JWT struct {
	SigningKey string
}

type CasbinConfig struct {
	ModelPath string // casbin model地址配置
}

type MysqlAdmin struct { // mysql admin 数据库配置
	Username string
	Password string
	Path     string
	Dbname   string
	Config   string
}

type RedisAdmin struct { // Redis admin 数据库配置
	Addr     string
	Password string
	DB       int
}
type Qiniu struct { // 七牛 密钥配置
	AccessKey string
	SecretKey string
}

type Tencent struct {
	SecretId                  string
	SecretKey                 string
	Region                    string
	SdkAppId                  int64
	SdkAppKey                 string
	AppId                     int64
	BizId                     int64
	CallbackKey               string
	TranscodeFileCallbackUrl  string
	TranscodeVideoCallbackUrl string
	OwnerAccount              string
}

//mongodb连接字符串
type Mongodb struct {
	Hosts       string
	User        string
	Passwd      string
	PoolLimit   uint64
	ReadReferer string
	ReplicaSet  string
}

//mongodb连接字符串
type Tomongodb struct {
	Hosts       string
	User        string
	Passwd      string
	PoolLimit   uint64
	ReadReferer string
	ReplicaSet  string
}
type Testmongodb struct {
	Hosts       string
	User        string
	Passwd      string
	PoolLimit   uint64
	ReadReferer string
	ReplicaSet  string
}

//课程配置相关
type CourseConfig struct {
	AssetsUrl string
}

type UploadConfig struct {
	UploadAssets string
}

//nats配置

type NatsConfig struct {
	Hosts string
}

//直播课程配置
type LiveCourseConfig struct {
	LivePushDomain    string
	LivePushDomainKey string
	LivePullDomain    string
	LivePullDomainKey string
}

var GinVueAdminconfig Config

func init() {
	v := viper.New()
	v.SetConfigName("config")           //  设置配置文件名 (不带后缀)
	v.AddConfigPath("./static/config/") // 第一个搜索路径
	v.SetConfigType("json")
	err := v.ReadInConfig() // 搜索路径，并读取配置数据
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	if err := v.Unmarshal(&GinVueAdminconfig); err != nil {
		fmt.Println(err)
	}
}
