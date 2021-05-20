package editorapi

import (
	"editorApi/commons"
	"editorApi/init/initNats"
	"editorApi/repository"
	"editorApi/requests"
	"editorApi/service"
	"editorApi/tools/helpers"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// @Tags StatisticAPI（统计管理）
// @Summary 添加统计Chapter
// @Description 添加统计Chapter
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.StatisticUnlockChapter true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/statistic/unlock_chapter [post]
func StatisticUnlockChapter(ctx *gin.Context) {
	var request requests.StatisticUnlockChapter
	ctx.BindJSON(&request)

	//2.验证
	UserValidations := commons.BaseValidations{}
	message, err := UserValidations.Check(request)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	// 异步操作
	initNats.NatsConn.Publish("StatisticUnlockChapter", &request)

	//obj := service.AppService()
	//result, err := obj.AppService.StatisticService.StatisticUnlockChapter(ctx, request)
	//if err != nil {
	//	commons.Error(ctx, 500, err, err.Error())
	//}

	commons.Success(ctx, nil, "提交成功！", request)
}

// @Tags StatisticAPI（统计管理）
// @Summary 添加统计Part
// @Description 添加统计Part
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.StatisticUnlockPart true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/statistic/unlock_part [post]
func StatisticUnlockPart(ctx *gin.Context) {
	var request requests.StatisticUnlockPart
	ctx.BindJSON(&request)

	//2.验证
	UserValidations := commons.BaseValidations{}
	message, err := UserValidations.Check(request)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	var contentExportsModel repository.ContentExports
	var paramsContentExports requests.ContentExports
	paramsContentExports.Code = request.Code
	paramsContentExports.Name = request.CourseCode
	id := uuid.NewV4().String()
	paramsContentExports.ID = id

	contentExportsModel.AddContentExports(ctx, paramsContentExports)
	request.Id = id
	// 异步操作
	initNats.NatsConn.Publish("StatisticUnlockPart", &request)

	//obj := service.AppService()
	//result, err := obj.AppService.StatisticService.StatisticUnlockChapter(ctx, request)
	//if err != nil {
	//	commons.Error(ctx, 500, err, err.Error())
	//}

	commons.Success(ctx, nil, "提交成功！", request)
}

func HanderStatisticUnlockPart(request *requests.StatisticUnlockPart) {
	obj := service.AppService()
	obj.AppService.StatisticService.HanderStatisticUnlockPart(*request)
}

func HanderStatisticUnlockChapter(request *requests.StatisticUnlockChapter) {
	obj := service.AppService()
	obj.AppService.StatisticService.HanderStatisticUnlockChapter(*request)
}
