package requests

type TeacherListRequests struct {
	Page
	RealName string `bson:"real_name" json:"real_name"` //教师真实姓名
	LanCode  string `json:"lan_code"`                   //课程code
	Status   int    `bson:"status" json:"status"`       //认证状态 1、待认证，2、认证未通过，3、认证通过，4、已删除、已冻结
}

type TeacherEditRequests struct {
	UserId string `bson:"user_id" json:"user_id" validate:"required" label:"用户ID"` //用户ID(必填)
}

type TeacherRequests struct {
	UserId string `bson:"user_id" json:"user_id" validate:"required" label:"用户ID"` //用户ID(必填)
	Status   int    `bson:"status" json:"status"`       //认证状态 1、待认证，2、认证未通过，3、认证通过，4、已删除、已冻结
}