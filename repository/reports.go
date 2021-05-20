package repository

import (
	"editorApi/commons"
	"editorApi/init/mgdb"
	"editorApi/requests"
	"editorApi/responses"
	"editorApi/tools/helpers"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const tb_reports = "reports"

type Reports struct {
	ID          string    `json:"id" bson:"_id"`                   // ID
	ConId       string    `json:"conId" bson:"con_id"`             // 内容ID
	CreatedTime time.Time `json:"createdTime" bson:"created_time"` // 创建时间
	State       int       `json:"state" bson:"state"`              // 状态 默认0；0代表未处理，1代表已处理
	DataVersion string    `json:"dataVersion" bson:"data_version"` //
	RepType     string    `json:"repType" bson:"rep_type"`         // 请求类型: word，sentence
	UserId      string    `json:"userId" bson:"user_id"`           // 用户ID
	Tags        string    `json:"tags" bson:"tags"`                // 标签
	Desc        string    `json:"desc" bson:"desc"`                // 描述
	DataArea    string    `json:"dataArea" bson:"data_area"`       //
}

// 创建Reports
func (m *Reports) Create(ctx *gin.Context, params requests.ReportsCreateRequests) (inserted_id interface{}, err error) {
	db := mgdb.MongoClient.Database(mgdb.DbDict).Collection(tb_reports)

	insertResult, err := db.InsertOne(ctx, params)
	if err != nil {
		return
	}
	inserted_id = insertResult.InsertedID
	return
}

// 获取Reports更多数据
func (m *Reports) Find(ctx *gin.Context, params requests.ReportsFindRequests) (result []responses.ReportsResponses, err error) {
	db := mgdb.MongoClient.Database(mgdb.DbDict).Collection(tb_reports)

	// 过滤条件
	var filter = bson.D{}
	if !helpers.Empty(params.UserId) {
		filter = append(filter, bson.E{"user_id", params.UserId})
	}

	if params.State != "" {
		filter = append(filter, bson.E{"state", params.State})
	}

	if !helpers.Empty(params.RepType) {
		filter = append(filter, bson.E{"rep_type", params.RepType})
	}

	if !helpers.Empty(params.Desc) {
		filter = append(filter, bson.E{"desc", primitive.Regex{Pattern: params.Desc, Options: "i"}})
	}

	cursor, err := db.Find(ctx, filter)
	if err != nil {
		return
	}

	err = cursor.All(ctx, &result)
	return
}

// 获取Reports
func (m *Reports) FindOne(ctx *gin.Context, params requests.ReportsFindOneRequests) (result responses.ReportsResponses, err error) {
	db := mgdb.MongoClient.Database(mgdb.DbDict).Collection(tb_reports)

	// 过滤条件
	var filter = bson.D{}
	if !helpers.Empty(params.UserId) {
		filter = append(filter, bson.E{"user_id", params.UserId})
	}

	if params.State != "" {
		filter = append(filter, bson.E{"state", params.State})
	}

	if !helpers.Empty(params.RepType) {
		filter = append(filter, bson.E{"rep_type", params.RepType})
	}

	if !helpers.Empty(params.Desc) {
		filter = append(filter, bson.E{"desc", primitive.Regex{Pattern: params.Desc, Options: "i"}})
	}

	singleResult := db.FindOne(ctx, filter)
	if err != nil {
		return
	}

	err = singleResult.Decode(&result)
	return
}

// 获取Reports带分页
func (m *Reports) List(ctx *gin.Context, params requests.ReportsListRequests) (result []responses.ReportsResponses, err error) {
	db := mgdb.MongoClient.Database(mgdb.DbDict).Collection(tb_reports)
	// 过滤条件
	var filter = bson.D{}

	if !helpers.Empty(params.UserId) {
		filter = append(filter, bson.E{"user_id", params.UserId})
	}

	if params.State != "" {
		filter = append(filter, bson.E{"state", params.State})
	}

	if !helpers.Empty(params.RepType) {
		filter = append(filter, bson.E{"rep_type", params.RepType})
	}

	if !helpers.Empty(params.Desc) {
		filter = append(filter, bson.E{"desc", primitive.Regex{Pattern: params.Desc, Options: "i"}})
	}

	page := commons.DefaultPage()
	if !helpers.Empty(params.PageSize) {
		page.Limit = params.PageSize
	}

	if !helpers.Empty(params.PageIndex) && params.PageIndex > 0 {
		page.Skip = (params.PageIndex - 1) * page.Limit
	}

	option := options.Find().SetSkip(page.Skip).SetLimit(page.Limit)
	if !helpers.Empty(params.SortType) && !helpers.Empty(params.TextField) {
		option = option.SetSort(bson.M{
			params.TextField: params.SortType,
		})
	}
	cursor, err := db.Find(ctx, filter, option)
	defer cursor.Close(ctx)

	if err != nil {
		return
	}
	err = cursor.All(ctx, &result)
	dbClient := mgdb.MongoClient.Database(mgdb.DbDict)

	for k, r := range result {
		con := map[string]string{}
		if r.RepType == "word" {
			rst := dbClient.Collection("dict_"+strings.ToLower(r.FromLang)).FindOne(
				ctx,
				bson.M{"uuid": r.ConId},
				options.FindOne().SetProjection(bson.M{"content": 1, "_id": 0}),
			)
			rst.Decode(&con)
			r.Content = con["content"]
		}
		if r.RepType == "sentence" {
			rst := dbClient.Collection("sentence_"+strings.ToLower(r.FromLang)).FindOne(
				ctx,
				bson.M{"uuid": r.ConId},
				options.FindOne().SetProjection(bson.M{"sentence": 1, "_id": 0}),
			)
			rst.Decode(&con)
			r.Content = con["sentence"]
		}

		result[k] = r
	}
	return
}

// 更新Reports
func (m *Reports) Update(ctx *gin.Context, params requests.ReportsUpdateRequests) (updateResult interface{}, err error) {
	db := mgdb.MongoClient.Database(mgdb.DbDict).Collection(tb_reports)

	var filter = bson.D{}
	if !helpers.Empty(params.ID) {
		id, _ := primitive.ObjectIDFromHex(params.ID)
		filter = append(filter, bson.E{"_id", id})
	}
	var update = bson.M{
		"state": params.State,
	}

	updateResult, err = db.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return
	}

	return
}

// 删除Reports
func (m *Reports) Delete(ctx *gin.Context, params requests.ReportsDeleteRequests) (deleteResult interface{}, err error) {
	db := mgdb.MongoClient.Database(mgdb.DbDict).Collection(tb_reports)
	var filter = bson.D{}
	if !helpers.Empty(params.ID) {
		id, _ := primitive.ObjectIDFromHex(params.ID)
		filter = append(filter, bson.E{"_id", id})
	}

	deleteResult, err = db.DeleteOne(ctx, filter)
	return
}
