package requests

import "time"

type ContentExports struct {
	ID        string    `bson:"id" json:"id"`       // 下载数据ID
	UUID      string    `bson:"uuid" json:"uuid"`   // catalogs UUID
	Level     string    `bson:"level" json:"level"` //级别
	Name      string    `bson:"name" json:"name"`
	Code      string    `bson:"code" json:"code"`
	Url       string    `bson:"url" json:"url"`
	Status    int64     `bson:"status" json:"status"`         //1 代表正在处理，2；处理成功
	UserName  string    `bson:"user_name" json:"user_name"`   //操作人
	CreatedOn time.Time `bson:"created_on" json:"created_on"` //创建时间
}
