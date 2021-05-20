package repository

import (
	"editorApi/commons"
	"editorApi/init/mgdb"
	"editorApi/mdbmodel/editor"
	"editorApi/requests"
	"editorApi/responses"
	"editorApi/tools/helpers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Teacher struct {
}

func (m *Teacher) TeacherList(ctx *gin.Context, params requests.TeacherListRequests) (result []responses.TeacherResponse, err error) {
	collection := mgdb.OnlineClient.Database(mgdb.DbKuyu).Collection(editor.TbTeacher)

	var filter = bson.D{}
	if !helpers.Empty(params.RealName) {
		filter = append(filter, bson.E{"real_name", params.RealName})
	}

	if !helpers.Empty(params.LanCode) {
		filter = append(filter, bson.E{"lan_code", params.LanCode})
	}

	if !helpers.Empty(params.Status) {
		filter = append(filter, bson.E{"status", params.Status})
	}

	page := commons.DefaultPage()
	if !helpers.Empty(params.PageSize) {
		page.Limit = params.PageSize
	}

	if !helpers.Empty(params.PageIndex) && params.PageIndex > 0 {
		page.Skip = (params.PageIndex - 1) * page.Limit
	}

	var rank = bson.M{"created_on": -1}

	if !helpers.Empty(params.SortType) && !helpers.Empty(params.TextField) {
		rank = bson.M{params.TextField: params.SortType}
	}

	option := options.Find().SetSort(rank).SetLimit(page.Limit).SetSkip(page.Skip)

	cusor, err := collection.Find(
		ctx,
		filter,
		option,
	)

	defer cusor.Close(ctx)
	err = cusor.All(ctx, &result)
	return
}

func (m *Teacher) FindOne(ctx *gin.Context, params requests.TeacherEditRequests) (result responses.TeacherResponse, err error) {
	collection := mgdb.OnlineClient.Database(mgdb.DbKuyu).Collection(editor.TbTeacher)
	var filter = bson.D{}
	if !helpers.Empty(params.UserId) {
		filter = append(filter, bson.E{"user_id", params.UserId})
	}
	singleResult := collection.FindOne(
		ctx,
		filter,
	)
	err = singleResult.Decode(&result)
	return
}

// 更新
func (m *Teacher) Audit(ctx *gin.Context, params requests.TeacherRequests) (upserted_id interface{}, err error) {
	collection := mgdb.OnlineClient.Database(mgdb.DbKuyu).Collection(editor.TbTeacher)

	var filter = bson.D{}
	if !helpers.Empty(params.UserId) {
		filter = append(filter, bson.E{"user_id", params.UserId})
	}

	var data = make(map[string]interface{})
	data["status"] = params.Status
	data["audit_time"] = time.Now().Format("2006-01-02 15:04:05")

	upserted_id, err = collection.UpdateOne(
		ctx,
		filter,
		bson.M{
			"$set": data,
		},
	)
	return
}
