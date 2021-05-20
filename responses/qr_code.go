package responses

import "time"

type QRcodeResponses struct {
	UUID       string    `bson:"uuid" json:"uuid"`               //UUID
	Title      string    `bson:"title" json:"title"`             //标题
	Info       string    `bson:"info" json:"info"`               //json内容
	IsDel      bool      `bson:"is_del" json:"is_del"`           //是否删除
	Size       int       `bson:"size" json:"size"`               //二维码图片大小
	CreatedOn  time.Time `bson:"created_on" json:"created_on"`   //创建时间
	UpdateTime time.Time `bson:"update_time" json:"update_time"` //更新时间
}
