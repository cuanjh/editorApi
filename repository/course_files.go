package repository

import (
	"editorApi/init/mgdb"
	"editorApi/mdbmodel/editor"
	"editorApi/requests"
	"editorApi/responses"
	"editorApi/tools/helpers"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type CourseFiles struct {
}

func (m *CourseFiles) UpdateEventData(ctx *gin.Context, params requests.CourseFilesEventData) (upserted_id interface{}, err error) {
	collection := mgdb.OnlineClient.Database(mgdb.DbKuyu).Collection(editor.TbCourseFiles)

	var filter = bson.D{}
	if !helpers.Empty(params.TaskId) {
		filter = append(filter, bson.E{"task_id", params.TaskId})
	}

	var update = bson.M{"event_data": params, "state": 2}

	updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return
	}
	upserted_id = updateResult.UpsertedID
	return
}

// 创建
func (m *CourseFiles) Create(ctx *gin.Context, params requests.CourseFilesRequests) (inserted_id interface{}, err error) {
	collection := mgdb.OnlineClient.Database(mgdb.DbKuyu).Collection(editor.TbCourseFiles)
	var data = bson.M{"uuid": uuid.NewV4().String(), "file_url": params.FileUrl, "type": params.Type, "size": params.Size, "live_uuid": params.LiveUuid, "title": params.Title, "task_id": params.TaskId, "state": 1, "created_on": time.Now().Format("2006-01-02 15:04:05")}

	insertResult, err := collection.InsertOne(ctx, data)
	if err != nil {
		return
	}
	inserted_id = insertResult.InsertedID
	return
}

func (m *CourseFiles) DeleteFile(ctx *gin.Context, params requests.CourseFilesDeleteRequests) (upserted_id interface{}, err error) {
	collection := mgdb.OnlineClient.Database(mgdb.DbKuyu).Collection(editor.TbCourseFiles)

	var filter = bson.D{}
	if !helpers.Empty(params.Uuid) {
		filter = append(filter, bson.E{"uuid", params.Uuid})
	}

	if !helpers.Empty(params.LiveUuid) {
		filter = append(filter, bson.E{"live_uuid", params.LiveUuid})
	}

	var data = make(map[string]interface{})
	data["state"] = 4

	updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": data})
	if err != nil {
		return
	}
	upserted_id = updateResult.UpsertedID
	return
}

func (m *CourseFiles) List(ctx *gin.Context, params requests.CourseFilesListRequests) (result []responses.CourseFilesListResponse, err error) {
	collection := mgdb.OnlineClient.Database(mgdb.DbKuyu).Collection(editor.TbCourseFiles)
	var filter = bson.D{}
	filter = append(filter, bson.E{"state", bson.M{"$ne": 4}})
	if !helpers.Empty(params.LiveUuid) {
		filter = append(filter, bson.E{"live_uuid", params.LiveUuid})
	}

	cursor, err := collection.Find(ctx, filter)
	defer cursor.Close(ctx)
	if err != nil {
		return
	}
	err = cursor.All(ctx, &result)
	return
}
