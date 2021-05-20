package tencent

import (
	"fmt"
	"testing"
	"time"
)

func TestGetClient(t *testing.T)  {
	GetLiveClient()
	url := GetLivePushUrl("test",time.Now().Format("2006-01-01"))
	fmt.Println(url)
	fmt.Println(GetLivePullUrl("test",time.Now().Format("2006-01-01")))
}