package editor

const TbTeacher = "teacher"

const (
	STATUS_CHECK_PENDING = 1
	STATUS_NOT_PASS      = 2
	STATUS_PASS          = 3
	STATUS_FREEZE        = 4
)

type Teacher struct {
	UserId           string   `bson:"user_id" json:"user_id"`                     //用户ID
	RealName         string   `bson:"real_name" json:"real_name"`                 //真实姓名
	Nationality      string   `bson:"nationality" json:"nationality"`             //国籍
	Gender           string   `bson:"gender" json:"gender"`                       //性别
	IdentityCard     string   `bson:"identity_card" json:"identity_card"`         //身份证、护照
	CertificateType  string   `bson:"certificate_type" json:"certificate_type"`   //证件类型
	CertificateFront string   `bson:"certificate_front" json:"certificate_front"` //证件正面
	CertificateBack  string   `bson:"certificate_back" json:"certificate_back"`   //证件反面
	BirthDate        string   `bson:"birth_date" json:"birth_date"`               //出生年月日
	Address          string   `bson:"address" json:"address"`                     //现居住地址
	LiveNickname     string   `bson:"live_nickname" json:"live_nickname"`         //直播昵称
	LanCode          []string `bson:"lan_code" json:"lan_code"`                   //语种
	Status           int      `bson:"status" json:"status"`                       //认证状态 1、待认证，2、认证未通过，3、认证通过，4、已删除、已冻结
	Introduction     string   `bson:"introduction" json:"introduction"`           //简介
	ApproveDate      string   `bson:"approve_date" json:"approve_date"`           //认证日期
	CourseNumber     int      `bson:"course_number" json:"course_number"`         //课程数量
	CreatedOn        string   `bson:"created_on" json:"created_on"`               //创建时间
	AuditTime        string   `bson:"audit_time" json:"audit_time"`               //审核时间
}
