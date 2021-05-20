package repository

import (
	"editorApi/commons"
	"editorApi/init/mgdb"
	"editorApi/requests"
	"editorApi/responses"
	"editorApi/tools/helpers"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const tb_content_reports = "content_reports"

type ContentReports struct {
	ID          string    `json:"id" bson:"_id"`                   // ID
	DataVersion int64     `json:"dataVersion" bson:"data_version"` // 版本
	UserId      string    `json:"userId" bson:"user_id"`           // 用户ID
	Uuid        string    `json:"uuid" bson:"uuid"`                // uuid
	Code        string    `json:"code" bson:"code"`                // 课程编码
	Tags        string    `json:"tags" bson:"tags"`                // tags
	Agent       string    `json:"agent" bson:"agent"`              // agent
	DataArea    string    `json:"dataArea" bson:"data_area"`       //
	ParentUuids []string  `json:"parentUuids" bson:"parent_uuids"` // 所有父节点
	LangCode    string    `json:"langCode" bson:"lang_code"`       // 归属语言
	Desc        string    `json:"desc" bson:"desc"`                // 描述
	Img         string    `json:"img" bson:"img"`                  // 图片地址
	CreatedTime time.Time `json:"createdTime" bson:"created_time"` // 创建时间
	Status      int       `json:"status" bson:"status"`            // 状态 1，已处理；0，未处理；默认0
}

// 创建课程内容反馈
func (m *ContentReports) Create(ctx *gin.Context, params requests.ContentReportsCreateRequests) (inserted_id interface{}, err error) {
	db := mgdb.MongoClient.Database(mgdb.DbContent).Collection(tb_content_reports)

	var add = bson.M{
		"data_version": params.DataVersion,
		"user_id":      params.UserId,
		"uuid":         uuid.NewV4().String(),
		"code":         params.Code,
		"tags":         params.Tags,
		"agent":        params.Agent,
		"data_area":    params.DataArea,
		"parent_uuids": params.ParentUuids,
		"lang_code":    params.LangCode,
		"desc":         params.Desc,
		"img":          params.Img,
	}

	insertResult, err := db.InsertOne(ctx, add)
	if err != nil {
		return
	}
	inserted_id = insertResult.InsertedID
	return
}

// 获取课程内容反馈更多数据
func (m *ContentReports) Find(ctx *gin.Context, params requests.ContentReportsFindRequests) (result []responses.ContentReportsResponses, err error) {
	db := mgdb.MongoClient.Database(mgdb.DbContent).Collection(tb_content_reports)

	// 过滤条件
	var filter = bson.D{}

	if !helpers.Empty(params.ID) {
		id, _ := primitive.ObjectIDFromHex(params.ID)
		filter = append(filter, bson.E{"_id", id})
	}

	if !helpers.Empty(params.UserId) {
		filter = append(filter, bson.E{"user_id", params.UserId})
	}
	if !helpers.Empty(params.Uuid) {
		filter = append(filter, bson.E{"uuid", params.Uuid})
	}
	cursor, err := db.Find(ctx, filter)
	if err != nil {
		return
	}

	err = cursor.All(ctx, &result)
	return
}

// 获取课程内容反馈
func (m *ContentReports) FindOne(ctx *gin.Context, params requests.ContentReportsFindOneRequests) (result responses.ContentReportsResponses, err error) {
	db := mgdb.MongoClient.Database(mgdb.DbContent).Collection(tb_content_reports)

	// 过滤条件
	var filter = bson.D{}
	if !helpers.Empty(params.Code) {
		filter = append(filter, bson.E{"code", params.Code})
	}
	if !helpers.Empty(params.Uuid) {
		filter = append(filter, bson.E{"uuid", params.Uuid})
	}

	singleResult := db.FindOne(ctx, filter)
	if err != nil {
		return
	}

	err = singleResult.Decode(&result)
	return
}

// 获取课程内容反馈带分页
func (m *ContentReports) List(ctx *gin.Context, params requests.ContentReportsListRequests) (result []responses.ContentReportsResponses, err error) {
	db := mgdb.MongoClient.Database(mgdb.DbContent).Collection(tb_content_reports)
	// 过滤条件
	var filter = bson.D{}

	if !helpers.Empty(params.ID) {
		id, _ := primitive.ObjectIDFromHex(params.ID)
		filter = append(filter, bson.E{"_id", id})
	}

	if !helpers.Empty(params.Uuid) {
		filter = append(filter, bson.E{"uuid", params.Uuid})
	}

	//code
	if !helpers.Empty(params.Code) {
		filter = append(filter, bson.E{"code", params.Code})
	}

	//lang_code
	if !helpers.Empty(params.LangCode) {
		filter = append(filter, bson.E{"lang_code", params.LangCode})
	}

	//desc
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
			helpers.CamelToCase(params.TextField): params.SortType,
		})
	}
	cursor, err := db.Find(ctx, filter, option)
	defer cursor.Close(ctx)

	if err != nil {
		return
	}
	err = cursor.All(ctx, &result)
	return
}

// 更新课程内容反馈
func (m *ContentReports) Update(ctx *gin.Context, params requests.ContentReportsUpdateRequests) (updateResult interface{}, err error) {
	db := mgdb.MongoClient.Database(mgdb.DbContent).Collection(tb_content_reports)

	var filter = bson.D{}
	if !helpers.Empty(params.Uuid) {
		filter = append(filter, bson.E{"uuid", params.Uuid})
	}

	if !helpers.Empty(params.ID) {
		id, _ := primitive.ObjectIDFromHex(params.ID)
		filter = append(filter, bson.E{"_id", id})
	}

	var update = bson.M{
		"data_version": params.DataVersion,
		"user_id":      params.UserId,
		"uuid":         params.Uuid,
		"code":         params.Code,
		"tags":         params.Tags,
		"agent":        params.Agent,
		"data_area":    params.DataArea,
		"parent_uuids": params.ParentUuids,
		"lang_code":    params.LangCode,
		"desc":         params.Desc,
		"img":          params.Img,
		"status":       params.Status,
	}

	updateResult, err = db.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return
	}

	return
}

// 删除课程内容反馈
func (m *ContentReports) Delete(ctx *gin.Context, params requests.ContentReportsDeleteRequests) (deleteResult interface{}, err error) {
	db := mgdb.MongoClient.Database(mgdb.DbContent).Collection(tb_content_reports)
	var filter = bson.D{}
	if !helpers.Empty(params.Uuid) {
		filter = append(filter, bson.E{"uuid", params.Uuid})
	}

	if !helpers.Empty(params.ID) {
		id, _ := primitive.ObjectIDFromHex(params.ID)
		filter = append(filter, bson.E{"_id", id})
	}

	deleteResult, err = db.DeleteOne(ctx, filter)
	return
}
