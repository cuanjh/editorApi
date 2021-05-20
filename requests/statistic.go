package requests

type StatisticUnlockChapter struct {
	StartDate  string `json:"startDate" validate:"required" label:"开始时间"`        //开始时间
	EndDate    string `json:"endDate" validate:"required" label:"结束时间"`          //结束时间
	CourseCode string `json:"courseCode" validate:"required" label:"courseCode"` //courseCode
	Code       string `json:"code" validate:"required" label:"code"`             //Code
}

type StatisticUnlockPart struct {
	StartDate  string `json:"startDate" validate:"required" label:"开始时间"`        //开始时间
	EndDate    string `json:"endDate" validate:"required" label:"结束时间"`          //结束时间
	CourseCode string `json:"courseCode" validate:"required" label:"courseCode"` //courseCode
	Chapter    string `json:"chapter"`                                           //courseCode
	Id         string `json:"id"`                                                //courseCode
	Code       string `json:"code" validate:"required" label:"code"`             //Code
}
