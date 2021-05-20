package repository

import (
	"context"
	"editorApi/init/mgdb"
	"editorApi/mdbmodel/editor"
	"editorApi/requests"
	"editorApi/responses"
	"editorApi/tools/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"sync"
)

type UnlockInfosPart struct {
	PartCode    string `json:"partCode"`
	UserId      string `json:"userId"`
	CourseCode  string `json:"courseCode"`
	Chapter     string `json:"chapter"`
	CorrectRate int64  `json:"correctRate"`
	CreatedTime int64  `json:"createdTime"`
}

func (m *UnlockInfosPart) StatisticUnlockPart(ctx context.Context, params requests.StatisticUnlockPart) (result []responses.UnlockInfosPart, err error) {
	collection := mgdb.OnlineClient.Database(mgdb.DbContent).Collection(editor.TbUnlockInfosPart)

	var filter = bson.D{}
	if !helpers.Empty(params.CourseCode) {
		filter = append(filter, bson.E{"courseCode", params.CourseCode})
	}

	if !helpers.Empty(params.Chapter) {
		filter = append(filter, bson.E{"chapter", params.Chapter})
	}

	if !helpers.Empty(params.StartDate) && !helpers.Empty(params.EndDate) {
		filter = append(filter, bson.E{"created_time", bson.M{"$gte": helpers.ParseInTimestamp(params.StartDate), "$lt": helpers.ParseInTimestamp(params.EndDate)}})
	}

	cusor, err := collection.Find(ctx, filter)

	defer cusor.Close(ctx)
	err = cusor.All(ctx, &result)
	return
}

func (m *UnlockInfosPart) CountPart(ctx context.Context, params requests.StatisticUnlockPart) (result int64, err error) {
	collection := mgdb.OnlineClient.Database(mgdb.DbContent).Collection(editor.TbUnlockInfosPart)

	var filter = bson.D{}
	filter = append(filter, bson.E{"correctRate", bson.M{"$gt": 0}})

	if !helpers.Empty(params.CourseCode) {
		filter = append(filter, bson.E{"courseCode", params.CourseCode})
	}

	if !helpers.Empty(params.Chapter) {
		filter = append(filter, bson.E{"chapter", params.Chapter})
	}

	if !helpers.Empty(params.StartDate) && !helpers.Empty(params.EndDate) {
		filter = append(filter, bson.E{"created_time", bson.M{"$gte": helpers.ParseInTimestamp(params.StartDate), "$lt": helpers.ParseInTimestamp(params.EndDate)}})
	}

	result, err = collection.CountDocuments(ctx, filter)
	return
}

var mutex sync.Mutex
func (m *UnlockInfosPart) CountPartUnique(ctx context.Context, params requests.StatisticUnlockPart) (num int, err error) {
	mutex.Lock()
	defer mutex.Unlock()

	collection := mgdb.OnlineClient.Database(mgdb.DbContent).Collection(editor.TbUnlockInfosPart)

	var filter = bson.D{}
	filter = append(filter, bson.E{"correctRate", bson.M{"$gt": 0}})

	if !helpers.Empty(params.CourseCode) {
		filter = append(filter, bson.E{"courseCode", params.CourseCode})
	}

	if !helpers.Empty(params.Chapter) {
		filter = append(filter, bson.E{"chapter", params.Chapter})
	}

	if !helpers.Empty(params.StartDate) && !helpers.Empty(params.EndDate) {
		filter = append(filter, bson.E{"created_time", bson.M{"$gte": helpers.ParseInTimestamp(params.StartDate), "$lt": helpers.ParseInTimestamp(params.EndDate)}})
	}

	pipeline := []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": "$userId", "count": bson.M{"$sum": 1}}},
		{"$project": bson.M{"count": 1, "_id": 0}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)

	type resp struct {
		Count int `bson:"count"`
	}
	var result []resp

	defer cursor.Close(ctx)
	cursor.All(ctx, &result)
	for _, item := range result {
		num = num + item.Count
	}
	return
}
