package editor

import "time"

//直播间
type LiveRoom struct {
	ListOrder     int       `bson:"list_order" json:"list_order"`
	BuyNum        int64     `bson:"buy_num" json:"buy_num"`
	CoverV2       string    `bson:"cover_v2" json:"cover_v2"`
	Cover         string    `bson:"cover" json:"cover"`
	Published     string    `bson:"published" json:"published"`
	CourseType    int64     `bson:"course_type" json:"course_type"`
	Code          string    `bson:"code" json:"code"`
	FreeForMember float64   `bson:"free_for_member" json:"free_for_member"`
	Description   string    `bson:"description" json:"description"`
	TagKeys       []string  `bson:"tag_keys" json:"tag_keys"`
	MoneyType     string    `bson:"money_type" json:"money_type"`
	Money         float64   `bson:"money" json:"money"`
	MoneyDiscount float64   `bson:"moneyDiscount" json:"moneyDiscount"`
	LanCode       string    `bson:"lan_code" json:"lan_code"`
	ModuleName    string    `bson:"module_name" json:"module_name"`
	CreatedOn     time.Time `bson:"created_on" json:"created_on"`
	UpdateTime    time.Time `bson:"update_time" json:"update_time"`
	UserID        string    `bson:"user_id" json:"user_id"`
	LiveInfo      *LiveInfo `bson:"liveInfo" json:"liveInfo"`
	IsDel         int       `bson:"is_del" json:"is_del"`
}
type LiveEditRoom struct {
	CoverV2       string    `bson:"cover_v2" json:"cover_v2"`
	Cover         string    `bson:"cover" json:"cover"`
	Code          string    `bson:"code" json:"code"`
	LanCode       string    `bson:"lan_code" json:"lan_code"`
	FreeForMember float64   `bson:"free_for_member" json:"free_for_member"`
	Description   string    `bson:"description" json:"description"`
	TagKeys       []string  `bson:"tag_keys" json:"tag_keys"`
	MoneyType     string    `bson:"money_type" json:"money_type"`
	Money         float64   `bson:"money" json:"money"`
	UserId        string    `bson:"user_id" json:"user_id"`
	MoneyDiscount float64   `bson:"moneyDiscount" json:"moneyDiscount"`
	ModuleName    string    `bson:"module_name" json:"module_name"`
	UpdateTime    time.Time `bson:"update_time" json:"update_time"`
	LiveInfo      *LiveInfo `bson:"liveInfo" json:"liveInfo"`
}
type LiveInfo struct {
	CourseCode        string   `bson:"courseCode" json:"courseCode"`
	LiveUserUUID      string   `bson:"live_user_uuid" json:"liveUserUUID"`
	TechDesc          string   `bson:"tech_desc" json:"tech_desc"`
	Level             int      `bson:"level" json:"level"`        // 课程级别
	BaseRand          int64    `bson:"base_rand" json:"baseRand"` // 初始值
	TechName          string   `bson:"tech_name" json:"tech_name"`
	TechPhoto         string   `bson:"tech_photo" json:"tech_photo"`
	Online            bool     `bson:"online" json:"online"`
	RealStartTime     int      `bson:"realStartTime" json:"realStartTime"`
	RealEndTime       int      `bson:"realEndTime" json:"realEndTime"`
	UUID              string   `bson:"uuid" json:"uuid"`
	StartDate         string   `bson:"startDate" json:"startDate"`
	EndDate           string   `bson:"endDate" json:"endDate"`
	StartTime         string   `bson:"startTime" json:"startTime"`
	EndTime           string   `bson:"endTime" json:"endTime"`
	WeekDays          []string `bson:"weekDays" json:"weekDays"`
	VideoUrl          string   `bson:"videoUrl" json:"videoUrl"`
	VideoCoverUrl     string   `bson:"videoCoverUrl" json:"videoCoverUrl"`
	ShareBgUrl        string   `bson:"shareBgUrl" json:"shareBgUrl"`
	ShareTitle        string   `bson:"shareTitle" json:"shareTitle"`
	ShareDesc         string   `bson:"shareDesc" json:"shareDesc"`
	SharePoster       []string `bson:"sharePoster" json:"sharePoster"`
	Posters           []string `bson:"posters" json:"posters"`
	CourseNum         int      `bson:"courseNum" json:"courseNum"`
	ExcludeDates      []string `bson:"exclude_dates" json:"exclude_dates"`
	FinishInfo        string   `bson:"finishInfo" json:"finishInfo"`
	FinishTitle       string   `bson:"finishTitle" json:"finishTitle"`
	WeixinNo          string   `bson:"weixinNo" json:"weixinNo"`
	DateNotice        string   `bson:"date_notice" json:"dateNotice"`            // 直播日期描述
	BasicCourseCode   string   `bson:"basic_course_code" json:"basicCourseCode"` // 课程编码
	BasicContentLevel string   `bson:"basic_content_level" json:"basicContentLevel"`
	BasicChapterCover string   `bson:"basic_chapter_cover" json:"basicChapterCover"`
	BasicProfilePhoto string   `bson:"basic_profile_photo" json:"basicProfilePhoto"` // 头像
	DisTechPhoto      string   `bson:"dis_tech_photo" json:"disTechPhoto"`           // 直播课发现首页头像
}

//直播间课程
type LiveCourse struct {
	UUID          string `bson:"uuid" json:"uuid"`
	RoomId        int64  `bson:"roomId" json:"roomId"`
	ListOrder     int    `bson:"listOrder" json:"listOrder"`
	CourseCode    string `bson:"courseCode" json:"courseCode"`
	Title         string `bson:"title" json:"title"`
	LanCode       string `bson:"lanCode" json:"lanCode"`
	Date          string `bson:"date" json:"date"`
	StartTime     int    `bson:"startTime" json:"startTime"`
	RealStartTime int    `bson:"realStartTime" json:"realStartTime"`
	RealStEndTime int    `bson:"realEndTime" json:"realEndTime"`
	EndTime       int    `bson:"endTime" json:"EndTime"`
	Cover         string `bson:"cover" json:"cover"`
	State         int    `bson:"state" json:"state"`
	IsDel         bool   `bson:"isDel" json:"isDel"`
	LivePushUrl   string `bson:"livePushUrl" json:"livePushUrl"`
	LivePullUrl   string `bson:"livePullUrl" json:"livePullUrl"`
	VideoUrl      string `bson:"videoUrl" json:"videoUrl"`
	VideoCover    string `bson:"videoCover" json:"videoCover"`
	FinishTitle   string `bson:"finishTitle" json:"finishTitle"`
	FinishInfo    string `bson:"finishInfo" json:"finishInfo"`
	WeixinNo      string `bson:"weixinNo" json:"weixinNo"`
	OnlineNumber  int    `bson:"online_number" json:"online_number"`
}

type LiveCourseUserCount struct {
	CourseUuid string `bson:"course_uuid" json:"course_uuid"`
	CreatedOn  string `bson:"created_on" json:"created_on"`
	Number     int    `bson:"number" json:"number"`
}
