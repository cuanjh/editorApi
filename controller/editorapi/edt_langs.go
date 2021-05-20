package editorapi

import (
	"errors"
	"editorApi/config"
	"editorApi/controller/servers"
	"editorApi/init/mgdb"
	"editorApi/mdbmodel/editor"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tblCourseLangs string = "course_langs"

type ListsParams struct {
	PageNo   int64 `json:"pageNo"`
	PageSize int64 `json:"pageSize"`
}

// @Tags EditorLangAPI（语言种类接口）
// @Summary 编辑器语种列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.ListsParams true "语种列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /editor/lang/list [post]
func LangLists(c *gin.Context) {

	paras := ListsParams{}
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
	langs := []*editor.Course_langs{}

	langsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCourseLangs)

	if cusor, err = langsCollection.Find(
		c,
		bson.M{"has_del": false},
		options.Find().SetSort(bson.M{"list_order": 1}),
		options.Find().SetLimit(limit),
		options.Find().SetSkip(skip),
	); err != nil {
		checkErr(c, err)
		return
	}
	defer cusor.Close(c)
	cusor.All(c, &langs)

	servers.ReportFormat(c, true, "语种列表", gin.H{
		"langs":     langs,
		"assetsUrl": config.GinVueAdminconfig.CourseConfig.AssetsUrl,
	})
}

/**
语言信息
**/
type langPara struct {
	Is_show        bool              `json:"is_show"`
	Title          map[string]string `json:"title"`
	Flag           []string          `json:"flag"`
	Desc           map[string]string `json:"desc"`
	List_order     int64             `json:"list_order"`
	Lan_code       string            `json:"lan_code"`
	Word_direction string            `json:"word_direction"`
	Is_hot         bool              `json:"is_hot"`
	Has_del        bool              `json:"has_del"`
}

// @Tags EditorLangAPI（语言种类接口）
// @Summary 添加语种
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.langPara true "语种信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /editor/lang/add [post]
func LangAdd(c *gin.Context) {

	lang := langPara{}
	c.BindJSON(&lang)
	lang.Has_del = false
	lang.Is_show = true
	langsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCourseLangs)
	check, _ := langsCollection.CountDocuments(c, bson.M{
		"lan_code": lang.Lan_code,
		"has_del":  false,
	})
	if check > 0 {
		err := errors.New("语言种类已经存在")
		checkErr(c, err)
	} else {
		r, e := langsCollection.InsertOne(c, lang)
		if e != nil {
			checkErr(c, e)
			return
		}
		servers.ReportFormat(c, true, "添加成功", gin.H{
			"assetsUrl": config.GinVueAdminconfig.CourseConfig.AssetsUrl,
			"id":        r.InsertedID,
		})
	}
}

/**
语言编辑信息
**/
type langEdit struct {
	Is_show        bool              `json:"is_show"`
	Title          map[string]string `json:"title"`
	Flag           []string          `json:"flag"`
	Desc           map[string]string `json:"desc"`
	List_order     int64             `json:"list_order"`
	Word_direction string            `json:"word_direction"`
	Is_hot         bool              `json:"is_hot"`
}

type langEditPara struct {
	LanCode  string   `json:"lan_code"`
	LangInfo langEdit `json:"lang_info"`
}

// @Tags EditorLangAPI（语言种类接口）
// @Summary 编辑语种信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.langEditPara true "语种信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"编辑成功"}"
// @Router /editor/lang/edit [post]
func LangEdit(c *gin.Context) {

	paras := langEditPara{}
	c.BindJSON(&paras)

	langsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCourseLangs)
	r, e := langsCollection.UpdateOne(
		c,
		bson.M{
			"lan_code": paras.LanCode,
			"has_del":  false,
		},
		bson.M{
			"$set": paras.LangInfo,
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

type langDelPara struct {
	LanCode string `json:"lan_code"`
}

// @Tags EditorLangAPI（语言种类接口）
// @Summary 删除语种信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.langDelPara true "语种信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/lang/del [post]
func LangDel(c *gin.Context) {

	paras := langDelPara{}
	c.BindJSON(&paras)
	dc, _ := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCourseInfos).CountDocuments(c, bson.M{
		"lan_code": paras.LanCode,
		"has_del":  false,
	})

	if dc > 0 {
		servers.ReportFormat(c, false, "此语种存在课程", gin.H{})
	} else {
		langsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCourseLangs)
		r, err := langsCollection.UpdateOne(
			c,
			bson.M{
				"lan_code": paras.LanCode,
				"has_del":  false,
			},
			bson.M{
				"$set": bson.M{"has_del": true},
			},
		)
		if err != nil {
			checkErr(c, err)
			return
		}
		servers.ReportFormat(c, true, "删除成功", gin.H{
			"modifiedCount": r.ModifiedCount,
		})
	}

}
func rcv(ctx *gin.Context) {
	if r := recover(); r != nil {
		var ok bool
		err, ok := r.(error)
		if !ok {
			debug.PrintStack()
			checkErr(ctx, err)
		}
	}
}

func checkErr(ctx *gin.Context, err error) {
	if err != nil {
		servers.ReportFormat(ctx, false, "失败", gin.H{
			"err": err.Error(),
		})
	}
}
