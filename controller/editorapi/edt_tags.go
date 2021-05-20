package editorapi

import (
	"editorApi/controller/servers"
	"editorApi/init/mgdb"
	"editorApi/init/qmlog"
	"editorApi/mdbmodel/editor"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/mongo"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tblContentTags string = "content_tags"
var tblContentTagTypes string = "content_tag_types"

// @Tags EditorTagAPI（标签接口）
// @Summary 标签类型列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/tag/types [get]
func TagTypes(c *gin.Context) {
	types := []map[string]string{}
	mgdb.Find(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		tblContentTagTypes,
		bson.M{},
		nil,
		map[string]int{
			"type": 1,
			"name": 1,
			"_id":  0,
		},
		0, 100,
		&types,
	)

	servers.ReportFormat(c, true, "tag类型列表", gin.H{
		"types": types,
	})

}

type tagType struct {
	Type string `json:"type" bson:"type"` //类型标示
	Name string `json:"name" bson:"name"` //类型名称
}

// @Tags EditorTagAPI（标签接口）
// @Summary 增加标签类型
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.tagType true "标签类型数据结构"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/tag/types/add [post]
func TagTypesAdd(c *gin.Context) {
	param := &tagType{}
	c.BindJSON(&param)
	mgdb.UpdateOne(
		mgdb.EnvEditor,
		mgdb.DbEditor,
		tblContentTagTypes,
		bson.M{
			"type": param.Type,
		},
		bson.M{
			"$set": bson.M{
				"name": param.Name,
			},
		},
		true,
	)

	servers.ReportFormat(c, true, "添加成功", gin.H{})

}

type tagsListsParams struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	PageNo   int64  `json:"pageNo"`
	PageSize int64  `json:"pageSize"`
}

// @Tags EditorTagAPI（标签接口）
// @Summary 标签列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.tagsListsParams true "标签列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /editor/tag/list [post]
func TagsLists(c *gin.Context) {

	paras := tagsListsParams{}
	c.BindJSON(&paras)
	var limit int64
	var skip int64

	limit = paras.PageSize
	if limit == 0 {
		limit = 1000
	}
	skip = paras.PageNo * limit

	var (
		err   error
		cusor *mongo.Cursor
	)
	tags := []*editor.Content_tags{}
	filter := bson.M{"has_del": false}
	if paras.Type != "" {
		filter["type"] = paras.Type
	}
	if paras.Name != "" {
		filter["name"] = primitive.Regex{
			Pattern: paras.Name,
			Options: "i",
		}
	}
	tagsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentTags)

	if cusor, err = tagsCollection.Find(
		c,
		filter,
		options.Find().SetSort(bson.M{"listorder": 1}),
		options.Find().SetLimit(limit),
		options.Find().SetSkip(skip),
	); err != nil {
		qmlog.QMLog.Error("查询标签列表出错：", err)
	}
	defer cusor.Close(c)
	cusor.All(c, &tags)

	servers.ReportFormat(c, true, "标签列表", gin.H{
		"tags": tags,
	})
}

// @Tags EditorTagAPI（标签接口）
// @Summary 添加标签
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editor.Content_tags true "标签信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/tag/add [post]
func TagAdd(c *gin.Context) {

	tags := editor.Content_tags{}
	c.BindJSON(&tags)

	tagsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentTags)
	check, _ := tagsCollection.CountDocuments(c, bson.M{
		"key":     tags.Key,
		"has_del": false,
	})
	if check > 0 {
		servers.ReportFormat(c, false, "标签key已经存在", gin.H{
			"tag": tags,
		})
	} else {
		tags.HasChanged = true
		tagsCollection.UpdateOne(
			c, bson.M{
				"key": tags.Key,
			},
			bson.M{"$set": tags},
			options.Update().SetUpsert(true),
		)

		servers.ReportFormat(c, true, "添加成功", gin.H{})
	}
}

// @Tags EditorTagAPI（标签接口）
// @Summary 添加标签
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editor.Content_tags true "标签信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/tag/edit [post]
func TagEdit(c *gin.Context) {

	tags := editor.Content_tags{}
	c.BindJSON(&tags)

	tagsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentTags)
	tags.HasChanged = true
	tagsCollection.UpdateOne(
		c, bson.M{
			"key": tags.Key,
		},
		bson.M{"$set": tags},
	)

	servers.ReportFormat(c, true, "编辑成功", gin.H{})

}

type tagsDelPara struct {
	Key string `json:"key"`
}

// @Tags EditorTagAPI（标签接口）
// @Summary 删除标签信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.tagsDelPara true "标签信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/tag/del [post]
func TagsDel(c *gin.Context) {

	paras := tagsDelPara{}
	c.BindJSON(&paras)

	tagsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentTags)
	r, err := tagsCollection.UpdateOne(
		c,
		bson.M{
			"key": paras.Key,
		},
		bson.M{
			"$set": bson.M{"has_del": true},
		},
	)
	if err != nil || r.ModifiedCount == 0 {
		servers.ReportFormat(c, false, "删除失败", gin.H{
			"modifiedCount": r.ModifiedCount,
		})
	}
	servers.ReportFormat(c, true, "删除成功", gin.H{
		"modifiedCount": r.ModifiedCount,
	})

}

var tblSearchWord = "search_word"
var courseContentDb = "courseContent"

type SearchWord struct {
	UUID      string `bson:"uuid" json:"uuid"`
	Word      string `bson:"word" json:"word"`
	ListOrder int64  `bson:"listOrder" json:"listOrder"`
}

// @Tags EditorSearchWordAPI（搜索关键词接口）
// @Summary 搜索关键词列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /editor/sword/list [post]
func SearchWordList(ctx *gin.Context) {
	searchWordCollection := mgdb.MongoClient.Database(courseContentDb).Collection(tblSearchWord)
	words := []*SearchWord{}
	cusor, _ := searchWordCollection.Find(
		ctx,
		bson.M{},
		options.Find().SetSort(map[string]int{
			"listOrder": 1,
		}),
	)
	defer cusor.Close(ctx)
	cusor.All(ctx, &words)

	servers.ReportFormat(ctx, true, "关键词列表", gin.H{
		"words": words,
	})
}

// @Tags EditorSearchWordAPI（搜索关键词接口）
// @Summary 添加关键词
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.SearchWord true "关键词"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/sword/add [post]
func SearchWordAdd(ctx *gin.Context) {
	var param *SearchWord
	ctx.BindJSON(&param)
	if param == nil {
		servers.ReportFormat(ctx, false, "关键词已经存在", gin.H{})
	}
	searchWordCollection := mgdb.MongoClient.Database(courseContentDb).Collection(tblSearchWord)
	c, _ := searchWordCollection.CountDocuments(ctx, bson.M{
		"word": param.Word,
	})
	if c > 0 {
		servers.ReportFormat(ctx, false, "关键词已经存在", gin.H{})
	} else {
		param.UUID = uuid.NewV4().String()
		param.ListOrder = time.Now().Unix()
		searchWordCollection.InsertOne(ctx, param)
		servers.ReportFormat(ctx, true, "添加成功", gin.H{})
	}

}

// @Tags EditorSearchWordAPI（搜索关键词接口）
// @Summary 编辑关键词
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.SearchWord true "关键词"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/sword/edit [post]
func SearchWordEdit(ctx *gin.Context) {
	var param *SearchWord
	ctx.BindJSON(&param)

	searchWordCollection := mgdb.MongoClient.Database(courseContentDb).Collection(tblSearchWord)
	c, _ := searchWordCollection.CountDocuments(ctx, bson.M{
		"word": param.Word,
		"uuid": bson.M{
			"$ne": param.UUID,
		},
	})
	if c > 0 {
		servers.ReportFormat(ctx, false, "关键词已经存在", gin.H{})
	} else {
		searchWordCollection.UpdateOne(ctx, bson.M{
			"uuid": param.UUID,
		}, bson.M{
			"$set": param,
		})
		servers.ReportFormat(ctx, true, "成功", gin.H{})
	}
}

type listOrdersPara []*SearchWord

// @Tags EditorSearchWordAPI（搜索关键词接口）
// @Summary 排序
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.listOrdersPara true "关键词"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/sword/listorder [post]
func SearchWordListOrder(ctx *gin.Context) {
	var params listOrdersPara
	ctx.BindJSON(&params)

	searchWordCollection := mgdb.MongoClient.Database(courseContentDb).Collection(tblSearchWord)

	for _, p := range params {
		searchWordCollection.UpdateOne(ctx, bson.M{
			"uuid": p.UUID,
		}, bson.M{
			"$set": bson.M{
				"listOrder": p.ListOrder,
			},
		})
	}

	servers.ReportFormat(ctx, true, "更新成功", gin.H{})

}

type searchWordDelParam struct {
	UUID string `bson:"uuid" json:"uuid"`
}

// @Tags EditorSearchWordAPI（搜索关键词接口）
// @Summary 添加关键词
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.searchWordDelParam true "关键词"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/sword/del [post]
func SearchWordDel(ctx *gin.Context) {
	var param *searchWordDelParam
	ctx.BindJSON(&param)

	searchWordCollection := mgdb.MongoClient.Database(courseContentDb).Collection(tblSearchWord)
	searchWordCollection.DeleteOne(ctx, bson.M{
		"uuid": param.UUID,
	})
	servers.ReportFormat(ctx, true, "成功", gin.H{})
}
