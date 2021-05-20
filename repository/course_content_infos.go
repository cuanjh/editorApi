package repository

import (
	"context"
	"editorApi/init/mgdb"
	"editorApi/mdbmodel/editor"
	"editorApi/requests"
	"editorApi/responses"
	"editorApi/tools/helpers"
	"go.mongodb.org/mongo-driver/bson"
)

type CourseContentInfos struct {
}

func (m *CourseContentInfos) FindOne(ctx context.Context, params requests.CourseContentInfos) (result responses.CourseContentInfos, err error) {
	collection := mgdb.MongoClient.Database(mgdb.DbEditor).Collection(editor.TbCourseContentInfos)
	var filter = bson.D{}
	if !helpers.Empty(params.Course_code) {
		filter = append(filter, bson.E{"course_code", params.Course_code})
	}

	singleResult := collection.FindOne(ctx, filter)

	err = singleResult.Decode(&result)
	return
}
