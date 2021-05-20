package responses

import "time"

type UnlockInfosChapter struct {
	DataVersion int       `bson:"data_version" json:"dataVersion"`
	UserId      string    `bson:"user_id" json:"userId"`
	CourseCode  string    `bson:"course_code" json:"courseCode"`
	Level       string    `bson:"level" json:"level"`
	Chapter     string    `bson:"chapter" json:"chapter"`
	DataArea    string    `bson:"data_area" json:"dataArea"`
	CreatedOn   time.Time `bson:"created_on" json:"createdOn"`
}
