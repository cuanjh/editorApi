package editorapi

import (
	"editorApi/controller/servers"
	"editorApi/init/initNats"
	"editorApi/init/mgdb"
	"editorApi/init/qmlog"
	"editorApi/mdbmodel/editor"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tblContentInfo string = "course_content_infos"

type contentListParams struct {
	PageNo     int64  `json:"pageNo"` //注释
	PageSize   int64  `json:"pageSize"`
	ParentUuid string `json:"parent_uuid"`
}

// @Tags EditorContentVersionAPI（内容版本接口）
// @Summary 课程内容版本列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.contentListParams true "课程内容列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/content/version/list [post]
func ContentVersionList(c *gin.Context) {

	var param contentListParams
	c.BindJSON(&param)
	limit := param.PageSize
	if limit == 0 {
		limit = 40
	}
	offset := limit * param.PageNo
	filter := bson.M{
		"has_del": false,
	}

	if param.ParentUuid != "" {
		filter["parent_uuid"] = param.ParentUuid
	}
	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentInfo)
	cusor, err := collection.Find(c, filter, options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{
		"list_order": 1,
	}))
	defer cusor.Close(c)
	if err != nil {
		checkErr(c, err)
		return
	}
	var contentLists []*editor.Course_content_infos
	cusor.All(c, &contentLists)
	servers.ReportFormat(c, true, "内容版本列表", gin.H{
		"contents": contentLists,
	})
}

type contentAddParams struct {
	Cover        []string          `json:"cover"`
	Is_show      bool              `json:"is_show"`
	Title        map[string]string `json:"title"`
	Name         string            `json:"name"`
	Flag         []string          `json:"flag"`
	Has_changed  bool              `json:"has_changed"`
	Has_del      bool              `json:"has_del"`
	Desc         map[string]string `json:"desc"`
	Uuid         string            `json:"uuid"`
	Version      string            `json:"version"`
	Module       string            `json:"module"`
	Parent_uuid  string            `json:"parent_uuid"`
	Update_time  int64             `json:"update_time"`
	Tags         []string          `json:"tags"`
	CopyFromUuid string            `json:"copy_from_uuid"`
	CopySameLang bool              `json:"copy_same_lang"`
}

// @Tags EditorContentVersionAPI（内容版本接口）
// @Summary 课程内容版本添加
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.contentAddParams true "课程内容添加,module是内容类型，现在有两个值：basic（基础内容）和levelGrade（等级测试内容）"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/content/version/add [post]
func ContentVersionAdd(c *gin.Context) {

	var param contentAddParams
	c.BindJSON(&param)
	param.Uuid = uuid.NewV4().String()
	param.Has_changed = true
	param.Has_del = false
	param.Update_time = time.Now().Unix()
	param.Is_show = true
	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentInfo)
	if param.Version == "" {
		versionNum, _ := collection.CountDocuments(
			c,
			bson.M{
				"parent_uuid": param.Parent_uuid,
				"module":      param.Module,
			},
		)
		param.Version = "v" + strconv.Itoa(int(versionNum+1))
	}

	_, e := collection.InsertOne(c, param)
	if e != nil {
		checkErr(c, e)
		return
	}
	//发送消息，复制其他版本的内容结构
	if param.CopyFromUuid != "" {
		natsConn := initNats.NatsConn
		err := natsConn.Publish(
			"content-copy-version",
			CatalogCopyParam{
				Uuids: []string{
					param.CopyFromUuid,
				},
				ToUuid:   param.Uuid,
				SameLang: param.CopySameLang,
			},
		)
		qmlog.QMLog.Info("发布复制版本内容的消息")
		if err != nil {
			checkErr(c, err)
			return
		}
	}
	servers.ReportFormat(c, true, "添加成功", gin.H{
		"uuid":    param.Uuid,
		"version": param.Version,
	})
}

type contentInfo struct {
	Cover       []string          `bson:"cover" json:"cover"`
	Is_show     bool              `bson:"is_show" json:"is_show"`
	Title       map[string]string `bson:"title" json:"title"`
	Name        string            `bson:"name" json:"name"`
	Flag        []string          `bson:"flag" json:"flag"`
	Has_changed bool              `bson:"has_changed" json:"has_changed"`
	Desc        map[string]string `bson:"desc" json:"desc"`
	Update_time int64             `bson:"update_time" json:"update_time"`
	Tags        []string          `json:"tags"`
	Version     string            `bson:"version" json:"version"`
}

type contentEditParams struct {
	Uuid        string      `json:"uuid"`
	ContentInfo contentInfo `json:"content_info"`
}

// @Tags EditorContentVersionAPI（内容版本接口）
// @Summary 课程内容版本修改
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.contentEditParams true "课程内容修改"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/content/version/edit [post]
func ContentVersionEdit(c *gin.Context) {

	var param *contentEditParams
	c.BindJSON(&param)
	contentInfos := param.ContentInfo
	contentInfos.Update_time = time.Now().Unix()
	contentInfos.Has_changed = true
	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentInfo)
	r, e := collection.UpdateOne(c, bson.M{
		"uuid": param.Uuid,
	}, bson.M{
		"$set": contentInfos,
	})
	if e != nil || r.ModifiedCount == 0 {
		checkErr(c, e)
		return
	}
	servers.ReportFormat(c, true, "修改成功", gin.H{})
}

type contentDelParams struct {
	Uuid string `json:"uuid"`
}

// @Tags EditorContentVersionAPI（内容版本接口）
// @Summary 课程内容版本删除
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.contentDelParams true "课程内容删除"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/content/version/del [post]
func ContentVersionDel(c *gin.Context) {

	var param contentDelParams
	c.BindJSON(&param)

	dc, _ := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs).CountDocuments(c, bson.M{
		"parent_uuid": param.Uuid,
		"has_del":     false,
	})
	if dc > 0 {
		servers.ReportFormat(c, false, "此内容版本存在内容", gin.H{})
	} else {
		collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentInfo)
		_, e := collection.UpdateOne(c, bson.M{
			"uuid": param.Uuid,
		}, bson.M{
			"$set": bson.M{"has_del": true},
		})
		if e != nil {
			checkErr(c, e)
			return
		}
		servers.ReportFormat(c, true, "删除成功", gin.H{})
	}
}
