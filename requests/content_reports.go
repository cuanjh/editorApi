package requests

import "time"

type ContentReportsCreateRequests struct {
	ID          string    `json:"id" bson:"_id"`                   // ID
	DataVersion int64     `json:"dataVersion" bson:"data_version"` // 版本
	UserId      string    `json:"userId" bson:"user_id"`           // 用户ID
	Uuid        string    `json:"uuid" bson:"uuid"`                // uuid
	Code        string    `json:"code" bson:"code"`                // 课程编码
	Tags        string    `json:"tags" bson:"tags"`                // tags
	Agent       string    `json:"agent" bson:"agent"`              // agent
	DataArea    string    `json:"dataArea" bson:"data_area"`       //
	ParentUuids []string  `json:"parentUuids" bson:"parent_uuids"` // 所有父节点
	LangCode    string    `json:"langCode" bson:"lang_code"`       // 归属语言
	Desc        string    `json:"desc" bson:"desc"`                // 描述
	Img         string    `json:"img" bson:"img"`                  // 图片地址
	CreatedTime time.Time `json:"createdTime" bson:"created_time"` // 创建时间
	Status      int       `json:"status" bson:"status"`            // 状态 1，已处理；0，未处理；默认0
}

type ContentReportsFindRequests struct {
	ID          string    `json:"id" bson:"_id"`                   // ID
	DataVersion int64     `json:"dataVersion" bson:"data_version"` // 版本
	UserId      string    `json:"userId" bson:"user_id"`           // 用户ID
	Uuid        string    `json:"uuid" bson:"uuid"`                // uuid
	Code        string    `json:"code" bson:"code"`                // 课程编码
	Tags        string    `json:"tags" bson:"tags"`                // tags
	Agent       string    `json:"agent" bson:"agent"`              // agent
	DataArea    string    `json:"dataArea" bson:"data_area"`       //
	ParentUuids []string  `json:"parentUuids" bson:"parent_uuids"` // 所有父节点
	LangCode    string    `json:"langCode" bson:"lang_code"`       // 归属语言
	Desc        string    `json:"desc" bson:"desc"`                // 描述
	Img         string    `json:"img" bson:"img"`                  // 图片地址
	CreatedTime time.Time `json:"createdTime" bson:"created_time"` // 创建时间
	Status      int       `json:"status" bson:"status"`            // 状态 1，已处理；0，未处理；默认0
}

type ContentReportsFindOneRequests struct {
	ID          string    `json:"id" bson:"_id"`                   // ID
	DataVersion int64     `json:"dataVersion" bson:"data_version"` // 版本
	UserId      string    `json:"userId" bson:"user_id"`           // 用户ID
	Uuid        string    `json:"uuid" bson:"uuid"`                // uuid
	Code        string    `json:"code" bson:"code"`                // 课程编码
	Tags        string    `json:"tags" bson:"tags"`                // tags
	Agent       string    `json:"agent" bson:"agent"`              // agent
	DataArea    string    `json:"dataArea" bson:"data_area"`       //
	ParentUuids []string  `json:"parentUuids" bson:"parent_uuids"` // 所有父节点
	LangCode    string    `json:"langCode" bson:"lang_code"`       // 归属语言
	Desc        string    `json:"desc" bson:"desc"`                // 描述
	Img         string    `json:"img" bson:"img"`                  // 图片地址
	CreatedTime time.Time `json:"createdTime" bson:"created_time"` // 创建时间
	Status      int       `json:"status" bson:"status"`            // 状态 1，已处理；0，未处理；默认0
}

type ContentReportsListRequests struct {
	Page
	ID          string    `json:"id" bson:"_id"`                   // ID
	DataVersion int64     `json:"dataVersion" bson:"data_version"` // 版本
	UserId      string    `json:"userId" bson:"user_id"`           // 用户ID
	Uuid        string    `json:"uuid" bson:"uuid"`                // uuid
	Code        string    `json:"code" bson:"code"`                // 课程编码
	Tags        string    `json:"tags" bson:"tags"`                // tags
	Agent       string    `json:"agent" bson:"agent"`              // agent
	DataArea    string    `json:"dataArea" bson:"data_area"`       //
	ParentUuids []string  `json:"parentUuids" bson:"parent_uuids"` // 所有父节点
	LangCode    string    `json:"langCode" bson:"lang_code"`       // 归属语言
	Desc        string    `json:"desc" bson:"desc"`                // 描述
	Img         string    `json:"img" bson:"img"`                  // 图片地址
	CreatedTime time.Time `json:"createdTime" bson:"created_time"` // 创建时间
	Status      int       `json:"status" bson:"status"`            // 状态 1，已处理；0，未处理；默认0
}

type ContentReportsUpdateRequests struct {
	ID          string    `json:"id" bson:"_id"`                   // ID
	DataVersion int64     `json:"dataVersion" bson:"data_version"` // 版本
	UserId      string    `json:"userId" bson:"user_id"`           // 用户ID
	Uuid        string    `json:"uuid" bson:"uuid"`                // uuid
	Code        string    `json:"code" bson:"code"`                // 课程编码
	Tags        string    `json:"tags" bson:"tags"`                // tags
	Agent       string    `json:"agent" bson:"agent"`              // agent
	DataArea    string    `json:"dataArea" bson:"data_area"`       //
	ParentUuids []string  `json:"parentUuids" bson:"parent_uuids"` // 所有父节点
	LangCode    string    `json:"langCode" bson:"lang_code"`       // 归属语言
	Desc        string    `json:"desc" bson:"desc"`                // 描述
	Img         string    `json:"img" bson:"img"`                  // 图片地址
	CreatedTime time.Time `json:"createdTime" bson:"created_time"` // 创建时间
	Status      int       `json:"status" bson:"status"`            // 状态 1，已处理；0，未处理；默认0
}

type ContentReportsDeleteRequests struct {
	ID          string    `json:"id" bson:"_id"`                   // ID
	DataVersion int64     `json:"dataVersion" bson:"data_version"` // 版本
	UserId      string    `json:"userId" bson:"user_id"`           // 用户ID
	Uuid        string    `json:"uuid" bson:"uuid"`                // uuid
	Code        string    `json:"code" bson:"code"`                // 课程编码
	Tags        string    `json:"tags" bson:"tags"`                // tags
	Agent       string    `json:"agent" bson:"agent"`              // agent
	DataArea    string    `json:"dataArea" bson:"data_area"`       //
	ParentUuids []string  `json:"parentUuids" bson:"parent_uuids"` // 所有父节点
	LangCode    string    `json:"langCode" bson:"lang_code"`       // 归属语言
	Desc        string    `json:"desc" bson:"desc"`                // 描述
	Img         string    `json:"img" bson:"img"`                  // 图片地址
	CreatedTime time.Time `json:"createdTime" bson:"created_time"` // 创建时间
	Status      int       `json:"status" bson:"status"`            // 状态 1，已处理；0，未处理；默认0
}
