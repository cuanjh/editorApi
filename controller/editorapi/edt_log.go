package editorapi

import (
	"editorApi/controller/servers"
	"editorApi/init/mgdb"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tblActionLogs = "action_logs"

// @Tags EditorLogsAPI（日志接口）
// @Summary 日志列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.ListsParams true "日志列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /editor/logs/list [post]
func EditorLogs(ctx *gin.Context) {
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

	var logs []bson.M
	logsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblActionLogs)

	if cusor, err = logsCollection.Find(
		ctx,
		bson.M{},
		options.Find().SetSkip(skip).SetSort(bson.M{
			"_id": -1,
		}).SetLimit(limit),
	); err != nil {
		checkErr(ctx, err)
		return
	}
	defer cusor.Close(ctx)
	cusor.All(ctx, &logs)
	servers.ReportFormat(ctx, true, "列表", gin.H{
		"logs": logs,
	})
}
