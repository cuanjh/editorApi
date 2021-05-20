package editor

type ClassWeixin struct {
	ID         string `bson:"_id" json:"id"`
	CourseCode string `bson:"course_code" json:"courseCode"`
	WeixinNo   string `bson:"weixin_no" json:"weixinNo"`
	WeixinCode string `bson:"weixin_code" json:"weixinCode"`
	IsShow     bool   `bson:"is_show" json:"isShow"`
}
