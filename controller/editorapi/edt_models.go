package editorapi

import (
	"editorApi/controller/servers"
	"editorApi/init/mgdb"
	"editorApi/mdbmodel/editor"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tblContentModels string = "content_models"

type modelListsParams struct {
	PageNo   int64 `json:"pageNo"`
	PageSize int64 `json:"pageSize"`
}

// @Tags EditorModelAPI（内容模型接口）
// @Summary 模型列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.modelListsParams true "模型列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /editor/model/list [post]
func ModelLists(c *gin.Context) {

	paras := modelListsParams{}
	c.BindJSON(&paras)
	var limit int64
	var skip int64

	limit = paras.PageSize
	if limit == 0 {
		limit = 40
	}
	skip = paras.PageNo * limit

	var (
		err   error
		cusor *mongo.Cursor
	)
	models := []*editor.Content_models{}

	modelsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentModels)

	if cusor, err = modelsCollection.Find(
		c,
		bson.M{"has_del": false},
		options.Find().SetSort(bson.M{"list_order": 1}),
		options.Find().SetLimit(limit).SetProjection(bson.M{"_id": 0}),
		options.Find().SetSkip(skip),
	); err != nil {
		checkErr(c, err)
		return
	}
	defer cusor.Close(c)
	cusor.All(c, &models)

	servers.ReportFormat(c, true, "模型列表", gin.H{
		"models": models,
	})
}

/**
模型添加信息
**/
type subFeildInfo struct {
	ListOrder int    `bson:"list_order" json:"list_order"`
	Feild     string `bson:"feild" json:"feild"`
	Type      string `bson:"type" json:"type"`
	Name      string `bson:"name" json:"name"`
	DataFrom  string `bson:"data_from" json:"data_from"`
	Desc      string `bson:"desc" json:"desc"`
	SubFeilds []subFeildInfoSon `bson:"sub_feilds" json:"sub_feilds"`
}

type subFeildInfoSon struct {
	ListOrder int    `bson:"list_order" json:"list_order"`
	Feild     string `bson:"feild" json:"feild"`
	Type      string `bson:"type" json:"type"`
	Name      string `bson:"name" json:"name"`
	DataFrom  string `bson:"data_from" json:"data_from"`
	Desc      string `bson:"desc" json:"desc"`
}

type FeildInfo struct {
	ListOrder int            `bson:"list_order" json:"list_order"`
	Feild     string         `bson:"feild" json:"feild"`
	Type      string         `bson:"type" json:"type"`
	Name      string         `bson:"name" json:"name"`
	DataFrom  string         `bson:"data_from" json:"data_from"`
	Desc      string         `bson:"desc" json:"desc"`
	SubFeilds []subFeildInfo `bson:"sub_feilds" json:"sub_feilds"`
}
type modelAddPara struct {
	Desc     string      `bson:"desc" json:"desc"`
	ModelKey string      `bson:"model_key" json:"model_key"`
	Feilds   []FeildInfo `bson:"feilds" json:"feilds"`
	Name     string      `bson:"name" json:"name"`
	Has_del  bool        `bson:"has_del" json:"has_del"`
}

// @Tags EditorModelAPI（内容模型接口）
// @Summary 添加模型
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.modelAddPara true "模型信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/model/add [post]
func ModelAdd(c *gin.Context) {

	model := modelAddPara{}
	c.BindJSON(&model)

	modelsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentModels)
	check, _ := modelsCollection.CountDocuments(c, bson.M{"model_key": model.ModelKey})
	if check > 0 {
		err := errors.New("模型Key已经存在")
		checkErr(c, err)
	} else {
		_, e := modelsCollection.InsertOne(c, model)
		checkErr(c, e)
		servers.ReportFormat(c, true, "添加成功", gin.H{})
	}
}

/**
编辑信息
**/
type modelEdit struct {
	Desc   string      `bson:"desc" json:"desc"`
	Feilds []FeildInfo `bson:"feilds" json:"feilds"`
	Name   string      `bson:"name" json:"name"`
}

type modelEditPara struct {
	ModelKey  string    `json:"model_key"`
	ModelInfo modelEdit `json:"model_info"`
}

// @Tags EditorModelAPI（内容模型接口）
// @Summary 编辑语种信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.modelEditPara true "语种信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"编辑成功"}"
// @Router /editor/model/edit [post]
func ModelEdit(c *gin.Context) {

	paras := modelEditPara{}
	c.BindJSON(&paras)

	modelsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentModels)
	r, e := modelsCollection.UpdateOne(
		c,
		bson.M{
			"model_key": paras.ModelKey,
			"has_del":   false,
		},
		bson.M{
			"$set": paras.ModelInfo,
		},
	)
	if e != nil {
		checkErr(c, e)
		return
	} else {
		servers.ReportFormat(c, true, "编辑成功", gin.H{
			"modifiedCount": r.ModifiedCount,
		})
	}
}

type modelDelPara struct {
	ModelKey string `json:"model_key"`
}

// @Tags EditorModelAPI（内容模型接口）
// @Summary 删除模型
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.modelDelPara true "模型信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/model/del [post]
func ModelDel(c *gin.Context) {

	paras := modelDelPara{}
	c.BindJSON(&paras)

	modelsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentModels)

	r, err := modelsCollection.UpdateOne(
		c,
		bson.M{
			"model_key": paras.ModelKey,
		},
		bson.M{
			"$set": bson.M{"has_del": true},
		},
	)

	if err != nil || r.ModifiedCount == 0 {
		if err == nil {
			err = errors.New("删除失败")
		}
		checkErr(c, err)
		return
	}

	servers.ReportFormat(c, true, "删除成功", gin.H{
		"modifiedCount": r.ModifiedCount,
	})

}
