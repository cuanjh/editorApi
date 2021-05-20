package tencent

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	live "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
)

const (
	LIVE_DOMAIN     = "92229.livepush.myqcloud.com"
	LIVE_DOMAIN_KEY = "42e30e51fc72571a5d7010c3d4dc450c"
)

func GetLiveClient() *live.Client {
	credential := common.NewCredential(
		"AKID0EjzjDEvud5SGoEZFfpkOmtlNrrcPBiy",
		"C8ofGMK2TmqdbKnxqTcV7mBLMUMxIMIy",
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 5
	cpf.SignMethod = "HmacSHA1"
	client, _ := live.NewClient(credential, "ap-beijing", cpf)

	return client
}

//获取推流地址
func GetLivePushUrl(streamName, txTime string) string {

	txTime = hex.EncodeToString([]byte(txTime))
	txt := LIVE_DOMAIN_KEY + streamName + txTime
	txtBytes := md5.Sum([]byte(txt))
	txSecret := hex.EncodeToString(txtBytes[0:])

	return "rtmp://" + LIVE_DOMAIN + "/live/" + streamName + "?txSecret=" + txSecret + "&txTime=" + txTime
}

//获取推流地址
func GetLivePullUrl(streamName, txTime string) string {

	txTime = hex.EncodeToString([]byte(txTime))
	txt := LIVE_DOMAIN_KEY + streamName + txTime
	txtBytes := md5.Sum([]byte(txt))
	txSecret := hex.EncodeToString(txtBytes[0:])

	return "http://" + LIVE_DOMAIN + "/live/" + streamName + ".m3u8?txSecret=" + txSecret + "&txTime=" + txTime
}
