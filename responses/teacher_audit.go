package responses

type TeacherAudit struct {
	UserId          string `bson:"user_id" json:"user_id"`                   //用户ID
	CreatedOn       string `bson:"created_on" json:"created_on"`             //创建时间、提交审核的时间
	Status          int    `bson:"status" json:"status"`                     //认证状态 1、待认证，2、认证未通过，3、认证通过
	Content         string `bson:"content" json:"content"`                   //审核内容
	Auditor         string `bson:"auditor" json:"auditor"`                   //审核人
	AuditorNickname string `bson:"auditor_nickname" json:"auditor_nickname"` //审核人昵称
	AuditTime       string `bson:"audit_time" json:"audit_time"`             //审核时间
}
