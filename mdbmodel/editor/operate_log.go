package editor

type OperateLog struct {
	Language    string `bson:"language" json:"language"`                           //语种
	LangCode    string `bson:"lang_code" json:"lang_code"`                         //语种ID
	Course      string `bson:"course" json:"course"`                               //课题
	CourseID    string `bson:"course_id" json:"course_id"`                         //课题ID
	Version     string `bson:"version" json:"version"`                             //版本
	VersionId   string `bson:"version_id" json:"version_id"`                       //版本ID
	Catalogs    string `bson:"catalogs" json:"catalogs"`                           //目录
	Pattern     int    `bson:"pattern" json:"pattern"`                             //操作（分类、内容）
	Mold        int    `bson:"mold" json:"mold"`                                   //操作类型
	OperateDate string `bson:"operate_date" json:"operate_date"`                   //操作日期
	CreatedTime  int64  `bson:"created_time" json:"created_time"`                     //创建时间
	Uuid        string `bson:"uuid" json:"uuid"`                                   //日志ID
	ParentUuid  string `bson:"parent_uuid" json:"parent_uuid" validate:"required"` //ID
	Content     string `bson:"content" json:"content"`                             //日志内容
	UserId      string `bson:"user_id" json:"user_id"`                             //用户ID
	UserName    string `bson:"user_name" json:"user_name"`                         //操作人
	Ip          string `bson:"ip" json:"ip"`
}
