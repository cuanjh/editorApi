package cmfunc

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
	"tkCommon/lib/rcluster"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/go-redis/redis"
)

var fileCache cache.Cache

func init() {
	var err error
	fileCache, err = cache.NewCache(
		"file",
		`{"CachePath":"./cache","FileSuffix":".cache"}`,
	)
	if err != nil {
		beego.Error(err)
	}
}

func RedisClient() *redis.ClusterClient {
	return rcluster.NewRedisClient()
}

func CacheId(data interface{}) (string, error) {
	j, err := json.Marshal(data)
	return fmt.Sprintf("%x", md5.Sum(j)), err
}

//保存缓存到文件中
func CacheSave(cacheId string, data interface{}, d time.Duration) bool {
	if fileCache != nil {
		cacheId = cacheId + ":v1data"
		bytes, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			return false
		}

		err = fileCache.Put(cacheId, bytes, d)
		if err != nil {
			fmt.Println(err)
			return false
		}
	}
	return true
}

//从文件中获取缓存
func CacheGet(cacheId string, rst interface{}) bool {
	fmt.Println("Cache Get")
	if fileCache != nil {
		cacheId = cacheId + ":v1data"
		dataBytes, ok := fileCache.Get(cacheId).([]byte)
		if ok {
			err := json.Unmarshal(dataBytes, rst)
			if err != nil {
				fmt.Println(err)
				return false
			}
			return true
		} else {
			fmt.Println(dataBytes)
		}
	} else {
		fmt.Println("FileCache is nil")
	}
	return false
}

func CacheSaveV2(cacheId string, data interface{}, ex time.Duration) error {
	cacheId = cacheId + ":v2data"

	rc := rcluster.NewRedisClient()
	if rc != nil {
		bytes, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			return err
		}

		ret := rc.Set(cacheId, bytes, ex)

		return ret.Err()
	}
	return errors.New("不能连接redis集群")
}

func CacheGetV2(cacheId string, rst interface{}) error {
	cacheId = cacheId + ":v2data"
	rc := rcluster.NewRedisClient()
	if rc != nil {
		ret := rc.Get(cacheId)
		dbytes, _ := ret.Bytes()

		if ret.Val() != "" {
			json.Unmarshal(dbytes, rst)
		}
	} else {
		return errors.New("不能连接redis集群")
	}
	return nil
}

//检测Cache是否存在
func CacheExist(cacheId string) bool {
	if fileCache != nil {
		return fileCache.IsExist(cacheId)
	}
	return false
}

//获取课程封面
//GetCourseCover
func GetCourseCover(courseCode string) string {
	return GetCourseAssets("course/covers/" + courseCode + "-2x.webp?v=4")
}

func GetCourseCoverV2(courseCode string) string {
	return GetCourseAssets("course/coversV2/" + courseCode + "-2x.webp")
}

//获取课程图标
func GetCourseFlag(lanCode string) string {
	return GetCourseAssets("course/icons/" + lanCode + "-3x.webp?v=4")
}

//获取国旗
func GetCountryFlag(country_code string) string {
	if country_code == "" {
		country_code = "CN"
	}
	return GetMobileAssets("country_flags/" + country_code + ".png")
}

//课程资源
func GetCourseAssets(url string) string {
	return beego.AppConfig.String("courseAssets") + url
}

//用户上传的资源
func GetUploadAssets(url string) string {
	return beego.AppConfig.String("uploadAssets") + url
}

//获取移动端静态文件资源
func GetMobileAssets(url string) string {
	return beego.AppConfig.String("mobileAssets") + url
}

func Photo(url, size string) string {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	} else {
		if url != "" {

			if strings.Contains(url, "app_image") {
				return GetUploadAssets(GetRandomPhoto())
			} else {

				if size == "" {
					if strings.Contains(url, "?") {
						return GetUploadAssets(url)
					} else {
						return GetUploadAssets(url) + "?v=3"
					}
				} else {
					if strings.Contains(url, "?") {
						return GetUploadAssets(url) + "&imageView2/0/w/" + size
					} else {
						return GetUploadAssets(url) + "?v=3&imageView2/0/w/" + size
					}
				}
			}
		} else {
			return GetUploadAssets(GetRandomPhoto())
		}
	}
	return ""
}

func GetRandomPhoto() string {
	rand.Seed(time.Now().Unix())
	return "uploadfiles/avatar/random/" + strconv.Itoa(rand.Intn(6)) + ".png"
}

//把值转换成json字符串
func ToJson(v interface{}) string {
	vbytes, e := json.Marshal(v)
	if e == nil {
		return string(vbytes)
	}
	return ""
}

//获取本地IP地址
func GetLocalIp() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		log.Fatal("获取IP失败" + err.Error())
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil &&
				(strings.HasPrefix(ipnet.IP.To4().String(), "192.") ||
					strings.HasPrefix(ipnet.IP.To4().String(), "10.") ||
					strings.HasPrefix(ipnet.IP.To4().String(), "172.")) {
				return ipnet.IP.To4().String()
			}
		}
	}
	return "0.0.0.0"
}

//获取远程IP地址
func GetRemoteIp() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		log.Fatal("获取IP失败" + err.Error())
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil &&
				!strings.HasPrefix(ipnet.IP.To4().String(), "192.") &&
				!strings.HasPrefix(ipnet.IP.To4().String(), "10.") &&
				!strings.HasPrefix(ipnet.IP.To4().String(), "172.") {
				return ipnet.IP.To4().String()
			}
		}
	}
	return "0.0.0.0"
}
