package requests

type TeacherAuditRequests struct {
	UserId  string `bson:"user_id" json:"user_id" validate:"required" label:"用户ID"` //用户ID
	Status  int    `bson:"status" json:"status" validate:"required" label:"认证状态"`   //认证状态 1、待认证，2、认证未通过，3、认证通过
	Content string `bson:"content" json:"content"`                                  //审核内容
}
