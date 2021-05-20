package editorapi

import (
	"editorApi/controller/servers"
	"editorApi/init/initNats"
	"editorApi/init/mgdb"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/mongo"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tblOnline = "online_jobs"

type onlineJonParam struct {
	UUID       string `json:"online_uuid"`
	JobName    string `json:job_name`
	OnlineType string `json:"online_type"`
	DbEnv      string `json:"db_env"` //上线的环境：test ,online
}
type PushOnlineMsg struct {
	UUID       string `json:"uuid"`
	OnlineUUID string `json:"online_uuid"`
	OnlineType string `json:"online_type"`
	DbEnv      string `json:"db_env"` //上线的环境：test ,online
}

// @Tags EditorOnlineAPI（课程内容上线接口）
// @Summary 添加内容上线任务
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.onlineJonParam true "要发布的目录或课程版本的uuid;上线类型：catalog/content_version"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/online/job [post]
func OnlineJob(ctx *gin.Context) {
	defer rcv(ctx)
	var pushUUIDs onlineJonParam
	ctx.BindJSON(&pushUUIDs)
	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblOnline)
	uidStr := uuid.NewV4().String()
	collection.InsertOne(ctx, bson.M{
		"uuid":         uidStr,
		"job_name":     pushUUIDs.JobName,
		"created_time": time.Now().Unix(),
		"state":        0,
		"online_uuid":  pushUUIDs.UUID,
		"online_type":  pushUUIDs.OnlineType,
	})

	err := initNats.NatsConn.Publish(
		"content-online-job",
		PushOnlineMsg{
			UUID:       uidStr,
			OnlineUUID: pushUUIDs.UUID,
			OnlineType: pushUUIDs.OnlineType,
			DbEnv:      pushUUIDs.DbEnv,
		},
	)

	if err != nil {
		checkErr(ctx, err)
		return
	}
	servers.ReportFormat(ctx, true, "添加成功", gin.H{})
}

type onlineCourseCodes []string

type PushOnlineCourseMsg struct {
	UUID        string   `json:"uuid"`
	CourseCodes []string `json:"courseCodes"`
}

// @Tags EditorOnlineAPI（课程内容上线接口）
// @Summary 添加课程信息上线任务
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.onlineCourseCodes true "课程编码数组：例如：["ENG-Basic"],或者空[]"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/online/courseInfo [post]
func OnlineCourseInfoJob(ctx *gin.Context) {
	defer rcv(ctx)
	var courseCodes onlineCourseCodes
	ctx.BindJSON(&courseCodes)
	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblOnline)
	uidStr := uuid.NewV4().String()
	collection.InsertOne(ctx, bson.M{
		"uuid":         uidStr,
		"job_name":     "上线课程信息",
		"created_time": time.Now().Unix(),
		"state":        0,
		"courseCodes":  courseCodes,
		"online_type":  "coruseInfos",
	})

	err := initNats.NatsConn.Publish(
		"courseinfo-online-job",
		PushOnlineCourseMsg{
			UUID:        uidStr,
			CourseCodes: courseCodes,
		},
	)
	if err != nil {
		checkErr(ctx, err)
		return
	}
	servers.ReportFormat(ctx, true, "添加成功", gin.H{})
}

// @Tags EditorOnlineAPI（课程内容上线接口）
// @Summary 上线任务列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.ListsParams true "上线任务列表:state:0(还没执行)，1（正在执行），2（执行完成）"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /editor/online/list [post]
func OnlineList(ctx *gin.Context) {
	defer rcv(ctx)
	paras := ListsParams{}
	ctx.BindJSON(&paras)
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
	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblOnline)

	cusor, err = collection.Find(
		ctx,
		bson.M{},
		options.Find().SetProjection(map[string]int{
			"_id": 0,
		}).SetLimit(limit).SetSkip(skip),
	)
	defer cusor.Close(ctx)
	if err != nil {
		checkErr(ctx, err)
		return
	}
	var jobs []bson.M
	cusor.All(ctx, &jobs)
	servers.ReportFormat(ctx, true, "成功", gin.H{
		"jobs": jobs,
	})
}

type jobDelParam struct {
	UUID string `json:"job_uuid"`
}

// @Tags EditorOnlineAPI（课程内容上线接口）
// @Summary 删除上线任务
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.jobDelParam true "删除参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /editor/online/del [post]
func OnlineJobDel(ctx *gin.Context) {
	defer rcv(ctx)
	var paras *jobDelParam
	ctx.BindJSON(&paras)
	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblOnline)
	collection.DeleteOne(ctx, bson.M{
		"uuid": paras.UUID,
	})
	servers.ReportFormat(ctx, true, "成功", gin.H{})
}
