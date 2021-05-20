package responses

import "time"

type ReportsResponses struct {
	ID          string    `json:"id" bson:"_id"`                   // ID
	ConId       string    `json:"conId" bson:"con_id"`             // 内容ID
	Content     string    `json:"content" bson:"content"`          // 内容
	CreatedTime time.Time `json:"createdTime" bson:"created_time"` // 创建时间
	State       int       `json:"state" bson:"state"`              // 状态 默认0；0代表未处理，1代表已处理
	RepType     string    `json:"repType" bson:"rep_type"`         // 请求类型: word，sentence
	UserId      string    `json:"userId" bson:"user_id"`           // 用户ID
	Tags        string    `json:"tags" bson:"tags"`                // 标签
	Desc        string    `json:"desc" bson:"desc"`                // 描述
	FromLang    string    `json:"fromLang" bson:"from_lang"`       // from
	ToLang      string    `json:"toLang" bson:"to_lang"`           // to
	DataArea    string    `json:"dataArea" bson:"data_area"`       //
}
