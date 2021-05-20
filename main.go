package main

import (
	"editorApi/commons"
	"editorApi/config"
	"editorApi/controller/editorapi"
	"editorApi/init/initNats"
	"editorApi/init/initRedis"
	"editorApi/init/initRouter"
	_ "editorApi/init/mgdb"
	"editorApi/init/qmlog"
	"editorApi/init/qmsql"
	"editorApi/init/registTable"
	"editorApi/requests"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func init() {
	commons.SetLogs(zap.DebugLevel, commons.LOGFORMAT_JSON, "")
}

// @title 全球说编辑器后台API服务文档
// @version 0.0.1
// @description 全球说编辑器后台API服务文档
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-token
// @BasePath /

func main() {
	qmlog.InitLog()                                            // 初始化日志
	db := qmsql.InitMysql(config.GinVueAdminconfig.MysqlAdmin) // 链接初始化数据库
	if config.GinVueAdminconfig.System.UseMultipoint {
		_ = initRedis.InitRedis() // 初始化redis服务
	}
	registTable.RegistTable(db)     //注册数据库表
	defer qmsql.DEFAULTDB.Close()   // 程序结束前关闭数据库链接
	natsConn := initNats.InitNats() //初始化Nats类库
	defer natsConn.Close()

	//上线目录或内容
	natsConn.QueueSubscribe(
		"content-online-job",
		"pushContentOnelineChan",
		func(msg *editorapi.PushOnlineMsg) {
			editorapi.PushContent(msg)
		},
	)

	//复制内容版本
	natsConn.QueueSubscribe(
		"content-copy-version",
		"CopyVersionChan",
		func(msg *editorapi.CatalogCopyParam) {
			editorapi.CopyContentVersion(msg)
		},
	)
	//上线内容信息和语种信息
	natsConn.QueueSubscribe(
		"courseinfo-online-job",
		"CourseInfoOnlineChan",
		func(msg *editorapi.PushOnlineCourseMsg) {
			editorapi.PushOnlineCourseInfos(msg)
		},
	)
	//上线和下线直播课程

	natsConn.QueueSubscribe(
		"liveCourseOnlineSub",
		"liveCourseOnlineSubChan",
		func(msg *editorapi.CourseOnlineMsg) {
			editorapi.LiveCourseOnlineSub(msg)
		},
	)

	natsConn.QueueSubscribe(
		"liveCourseOfflineSub",
		"liveCourseOfflineSubChan",
		func(msg *editorapi.CourseOfflineMsg) {
			editorapi.LiveCourseOfflineSub(msg)
		},
	)

	// 课程内容导入
	natsConn.QueueSubscribe(
		"ImportCourseContent",
		"ImportCourseContentChan",
		func(msg *editorapi.ImportContentMsg) {
			editorapi.HanderImportContent(msg)
		},
	)

	// 课程内容导出
	natsConn.QueueSubscribe(
		"ExportCourseContent",
		"ExportCourseContentChan",
		func(msg *editorapi.ContentExportParams) {
			editorapi.HanderContentExport(msg)
		},
	)

	// 课程解锁操作
	natsConn.QueueSubscribe(
		"ContentUnlock",
		"ContentUnlockChan",
		func(msg *editorapi.ContentUnlockMsg) {
			editorapi.HanderContentUnlock(msg)
		},
	)

	// 操作日志记录
	natsConn.QueueSubscribe(
		"OperateLog",
		"OperateLogChan",
		func(msg *editorapi.OperateLogParams) {
			editorapi.HanderOperateLog(msg)
		},
	)

	// 词典采集
	natsConn.QueueSubscribe(
		"CollectDict",
		"CollectDictChan",
		func(msg *requests.DictHanderParams) {
			editorapi.HanderDict(msg)
		},
	)
	// 把词典修改相关的内容导入到Es中
	natsConn.QueueSubscribe(
		"DictChangeMsg",
		"DictChangeMsgChan",
		func(msg *requests.DictESParams) {
			editorapi.PushContentToEs(msg)
		},
	)
	// 统计数据
	natsConn.QueueSubscribe(
		"StatisticUnlockChapter",
		"StatisticUnlockChapterChan",
		func(msg *requests.StatisticUnlockChapter) {
			editorapi.HanderStatisticUnlockChapter(msg)
		},
	)

	natsConn.QueueSubscribe(
		"StatisticUnlockPart",
		"StatisticUnlockPartChan",
		func(msg *requests.StatisticUnlockPart) {
			editorapi.HanderStatisticUnlockPart(msg)
		},
	)

	// CommentsUploadHander
	natsConn.QueueSubscribe(
		"CommentsUpload",
		"CommentsUploadChan",
		func(msg *requests.CommentsUploadHanderParams) {
			editorapi.CommentsUploadHander(msg)
		},
	)

	// CommentsUploadHander
	natsConn.QueueSubscribe(
		"HandleCourseData",
		"HandleCourseDataChan",
		func(msg *editorapi.CourseDataModle) {
			editorapi.HandleCourseData(msg)
		},
	)

	Router := initRouter.InitRouter() //注册路由

	qmlog.QMLog.Info("服务器开启") // 日志测试代码
	//Router.RunTLS(":443","ssl.pem", "ssl.key")  // https支持 需要添加中间件
	s := &http.Server{
		Addr:           ":8888",
		Handler:        Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	time.Sleep(10 * time.Microsecond)
	fmt.Printf(`
默认自动化文档地址:http://127.0.0.1%s/swagger/index.html
默认前端文件运行地址:http://127.0.0.1:8080
`, s.Addr)
	_ = s.ListenAndServe()
}
