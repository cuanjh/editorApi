package editorapi

import (
	"editorApi/controller/servers"
	"editorApi/init/mgdb"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tblContentType = "content_types"

type contenType struct {
	ListOrder int      `bson:"list_order" json:"list_order"`
	Name      string   `bson:"name" json:"name"`
	ModelKey  string   `bson:"model_key" json:"model_key"`
	ModelKeys []string `bson:"model_keys" json:"model_keys"`
	Type      string   `bson:"type" json:"type"`
	Desc      string   `bson:"desc" json:"desc"`
	HasDel    bool     `bson:"has_del" json:"has_del"`
}

// @Tags EditorContentTypeAPI(课程内容类型接口)
// @Summary 获取内容类型列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.ListsParams true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /editor/contenttype/list [post]
func ContentTypeList(ctx *gin.Context) {

	paras := ListsParams{}
	ctx.BindJSON(&paras)
	var limit int64
	var skip int64

	limit = paras.PageSize
	if limit == 0 {
		limit = 40
	}
	skip = paras.PageNo * limit

	var types []*contenType

	mgdb.Find(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		tblContentType,
		bson.M{
			"has_del": false,
		},
		bson.M{
			"list_order": 1,
			"_id":        -1,
		},
		nil,
		skip,
		limit,
		&types,
	)

	servers.ReportFormat(ctx, true, "内容类型列表", gin.H{
		"types": types,
	})
}

// @Tags EditorContentTypeAPI(课程内容类型接口)
// @Summary 添加内容类型
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.contenType true "添加类型参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /editor/contenttype/add [post]
func ContentTypeAdd(ctx *gin.Context) {
	var para *contenType
	ctx.BindJSON(&para)
	para.HasDel = false
	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentType)
	if c, _ := collection.CountDocuments(ctx, bson.M{"type": para.Type, "has_del": false}); c > 0 {
		servers.ReportFormat(ctx, false, "该题型已存在", gin.H{})
		return
	}
	para.HasDel = false

	collection.UpdateOne(
		ctx,
		bson.M{
			"type": para.Type,
		},
		bson.M{"$set": para},
		options.Update().SetUpsert(true),
	)
	servers.ReportFormat(ctx, true, "添加成功", gin.H{})
}

// @Tags EditorContentTypeAPI(课程内容类型接口)
// @Summary 编辑内容类型
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.contenType true "添加类型参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /editor/contenttype/edit [post]
func ContentTypeEdit(ctx *gin.Context) {
	var para *contenType
	ctx.BindJSON(&para)
	para.HasDel = false
	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentType)
	collection.UpdateOne(
		ctx,
		bson.M{
			"type": para.Type,
		},
		bson.M{"$set": para},
		options.Update().SetUpsert(true),
	)
	servers.ReportFormat(ctx, true, "编辑成功", gin.H{})
}

type contentTypeDelParam struct {
	Type string `json:"type"`
}

// @Tags EditorContentTypeAPI(课程内容类型接口)
// @Summary 删除内容类型
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.contentTypeDelParam true "删除内容类型参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/contenttype/del [post]
func ContentTypeDel(ctx *gin.Context) {

	var para *contentTypeDelParam
	ctx.BindJSON(&para)
	mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentType).DeleteOne(ctx, bson.M{
		"type": para.Type,
	})

	servers.ReportFormat(ctx, true, "删除成功", gin.H{})
}

type contentTypeListOrder struct {
	ListOrder int64  `json:"list_order"`
	Type      string `json:"type"`
}

type contentTypeListOrderInfos []*contentTypeListOrder

// @Tags EditorContentTypeAPI(课程内容类型接口)
// @Summary 内容类型排序
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.contentTypeListOrderInfos true "排序参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/contenttype/listorder [post]
func ContentTypeListOrder(ctx *gin.Context) {

	var paras contentTypeListOrderInfos
	ctx.BindJSON(&paras)
	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentType)

	for _, t := range paras {
		collection.UpdateOne(
			ctx,
			bson.M{
				"type": t.Type,
			},
			bson.M{
				"$set": bson.M{
					"list_order": t.ListOrder,
				},
			},
		)
	}
	servers.ReportFormat(ctx, true, "成功", gin.H{})
}
