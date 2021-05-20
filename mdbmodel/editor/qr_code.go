package editor

import "time"

const TbQRcode = "qr_code"

type QRcode struct {
	UUID       string    `bson:"uuid" json:"uuid"`               //UUID
	Title      string    `bson:"title" json:"title"`             //标题
	Info       string    `bson:"info" json:"info"`               //json内容
	IsDel      bool      `bson:"is_del" json:"is_del"`           //是否删除
	CreatedOn  time.Time `bson:"created_on" json:"created_on"`   //创建时间
	UpdateTime time.Time `bson:"update_time" json:"update_time"` //更新时间
}
