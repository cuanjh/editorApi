package responses

type UnlockInfosPart struct {
	PartCode    string `bson:"partCode" json:"partCode"`
	UserId      string `bson:"userId" json:"userId"`
	CourseCode  string `bson:"courseCode" json:"courseCode"`
	Chapter     string `bson:"chapter" json:"chapter"`
	CorrectRate int64  `bson:"correctRate" json:"correctRate"`
	CreatedTime int64  `bson:"created_time" json:"createdTime"`
}
