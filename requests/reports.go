package requests

import "time"

type ReportsCreateRequests struct {
	ID          string    `json:"id" bson:"_id"`                   // ID
	ConId       string    `json:"conId" bson:"con_id"`             // 内容ID
	CreatedTime time.Time `json:"createdTime" bson:"created_time"` // 创建时间
	State       string    `json:"state" bson:"state"`              // 状态 默认0；0代表未处理，1代表已处理
	DataVersion string    `json:"dataVersion" bson:"data_version"` //
	RepType     string    `json:"repType" bson:"rep_type"`         // 请求类型: word，sentence
	UserId      string    `json:"userId" bson:"user_id"`           // 用户ID
	Tags        string    `json:"tags" bson:"tags"`                // 标签
	Desc        string    `json:"desc" bson:"desc"`                // 描述
	DataArea    string    `json:"dataArea" bson:"data_area"`       //
}

type ReportsFindRequests struct {
	ID          string    `json:"id" bson:"_id"`                   // ID
	ConId       string    `json:"conId" bson:"con_id"`             // 内容ID
	CreatedTime time.Time `json:"createdTime" bson:"created_time"` // 创建时间
	State       string    `json:"state" bson:"state"`              // 状态 默认0；0代表未处理，1代表已处理
	DataVersion string    `json:"dataVersion" bson:"data_version"` //
	RepType     string    `json:"repType" bson:"rep_type"`         // 请求类型: word，sentence
	UserId      string    `json:"userId" bson:"user_id"`           // 用户ID
	Tags        string    `json:"tags" bson:"tags"`                // 标签
	Desc        string    `json:"desc" bson:"desc"`                // 描述
	DataArea    string    `json:"dataArea" bson:"data_area"`       //
}

type ReportsFindOneRequests struct {
	ID          string    `json:"id" bson:"_id"`                   // ID
	ConId       string    `json:"conId" bson:"con_id"`             // 内容ID
	CreatedTime time.Time `json:"createdTime" bson:"created_time"` // 创建时间
	State       string    `json:"state" bson:"state"`              // 状态 默认0；0代表未处理，1代表已处理
	DataVersion string    `json:"dataVersion" bson:"data_version"` //
	RepType     string    `json:"repType" bson:"rep_type"`         // 请求类型: word，sentence
	UserId      string    `json:"userId" bson:"user_id"`           // 用户ID
	Tags        string    `json:"tags" bson:"tags"`                // 标签
	Desc        string    `json:"desc" bson:"desc"`                // 描述
	DataArea    string    `json:"dataArea" bson:"data_area"`       //
}

type ReportsListRequests struct {
	ID          string    `json:"id" bson:"_id"`                   // ID
	ConId       string    `json:"conId" bson:"con_id"`             // 内容ID
	CreatedTime time.Time `json:"createdTime" bson:"created_time"` // 创建时间
	State       string    `json:"state" bson:"state"`              // 状态 默认0；0代表未处理，1代表已处理
	DataVersion string    `json:"dataVersion" bson:"data_version"` //
	RepType     string    `json:"repType" bson:"rep_type"`         // 请求类型: word，sentence
	UserId      string    `json:"userId" bson:"user_id"`           // 用户ID
	Tags        string    `json:"tags" bson:"tags"`                // 标签
	Desc        string    `json:"desc" bson:"desc"`                // 描述
	DataArea    string    `json:"dataArea" bson:"data_area"`       //
	Page
}

type ReportsUpdateRequests struct {
	ID          string    `json:"id" bson:"_id"`                   // ID
	ConId       string    `json:"conId" bson:"con_id"`             // 内容ID
	CreatedTime time.Time `json:"createdTime" bson:"created_time"` // 创建时间
	State       int       `json:"state" bson:"state"`              // 状态 默认0；0代表未处理，1代表已处理
	DataVersion string    `json:"dataVersion" bson:"data_version"` //
	RepType     string    `json:"repType" bson:"rep_type"`         // 请求类型: word，sentence
	UserId      string    `json:"userId" bson:"user_id"`           // 用户ID
	Tags        string    `json:"tags" bson:"tags"`                // 标签
	Desc        string    `json:"desc" bson:"desc"`                // 描述
	DataArea    string    `json:"dataArea" bson:"data_area"`       //
}

type ReportsDeleteRequests struct {
	ID          string    `json:"id" bson:"_id"`                   // ID
	ConId       string    `json:"conId" bson:"con_id"`             // 内容ID
	CreatedTime time.Time `json:"createdTime" bson:"created_time"` // 创建时间
	State       int       `json:"state" bson:"state"`              // 状态 默认0；0代表未处理，1代表已处理
	DataVersion string    `json:"dataVersion" bson:"data_version"` //
	RepType     string    `json:"repType" bson:"rep_type"`         // 请求类型: word，sentence
	UserId      string    `json:"userId" bson:"user_id"`           // 用户ID
	Tags        string    `json:"tags" bson:"tags"`                // 标签
	Desc        string    `json:"desc" bson:"desc"`                // 描述
	DataArea    string    `json:"dataArea" bson:"data_area"`       //
}
