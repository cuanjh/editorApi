package repository

import (
	"context"
	"editorApi/init/mgdb"
	"editorApi/mdbmodel/editor"
	"editorApi/requests"
	"editorApi/responses"
	"editorApi/tools/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type UnlockInfosChapter struct {
	DataVersion int       `json:"data_version"`
	UserId      string    `json:"user_id"`
	CourseCode  string    `json:"course_code"`
	Level       string    `json:"level"`
	Chapter     string    `json:"chapter"`
	DataArea    string    `json:"data_area"`
	CreatedOn   time.Time `json:"created_on"`
}

func (m *UnlockInfosChapter) StatisticUnlockChapter(ctx context.Context, params requests.StatisticUnlockChapter) (result []responses.UnlockInfosChapter, err error) {
	collection := mgdb.OnlineClient.Database(mgdb.DbContent).Collection(editor.TbUnlockInfosChapter)

	var filter = bson.D{}
	if !helpers.Empty(params.CourseCode) {
		filter = append(filter, bson.E{"course_code", params.CourseCode})
	}

	if !helpers.Empty(params.StartDate) && !helpers.Empty(params.EndDate) {
		filter = append(filter, bson.E{"created_on", bson.M{"$gte": helpers.ParseInLocationDate(params.StartDate), "$lt": helpers.ParseInLocationDate(params.EndDate)}})
	}

	cusor, err := collection.Find(ctx, filter)

	defer cusor.Close(ctx)
	err = cusor.All(ctx, &result)
	return
}
