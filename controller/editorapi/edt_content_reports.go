package editorapi

import (
	"editorApi/commons"
	"editorApi/requests"
	"editorApi/service"
	"editorApi/tools/helpers"
	"github.com/gin-gonic/gin"
)

// @Tags 课程内容反馈接口
// @Summary 添加课程内容反馈
// @Description 添加课程内容反馈
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/jsonf
// @Param data body requests.ContentReportsCreateRequests true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"失败"}"
// @Router /editor/content_reports/create [post]
func ContentReportsCreate(ctx *gin.Context) {
	var params requests.ContentReportsCreateRequests
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		commons.Error(ctx, 500, err, "error")
	}

	//2.验证
	validations := commons.BaseValidations{}
	message, err := validations.Check(params)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	obj := service.AppService()
	_, err = obj.AppService.ContentReportsService.Create(ctx, params)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}

	commons.Success(ctx, nil, "success", params)
}

// @Tags 课程内容反馈接口
// @Summary 查询课程内容反馈数据
// @Description 查询课程内容反馈数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.ContentReportsFindRequests true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"失败"}"
// @Router /editor/content_reports/find [post]
func ContentReportsFind(ctx *gin.Context) {
	var params requests.ContentReportsFindRequests
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		commons.Error(ctx, 500, err, "error")
	}

	//2.验证
	validations := commons.BaseValidations{}
	message, err := validations.Check(params)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	obj := service.AppService()
	response, err := obj.AppService.ContentReportsService.Find(ctx, params)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}

	commons.Success(ctx, response, "success", params)
}

// @Tags 课程内容反馈接口
// @Summary 查询一条课程内容反馈数据
// @Description 查询一条课程内容反馈数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.ContentReportsFindOneRequests true "列表参数"
// @Success 200 object responses.ContentReportsResponses
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"失败"}"
// @Router /editor/content_reports/findone [post]
func ContentReportsFindOne(ctx *gin.Context) {
	var params requests.ContentReportsFindOneRequests
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		commons.Error(ctx, 500, err, "error")
	}

	//2.验证
	validations := commons.BaseValidations{}
	message, err := validations.Check(params)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	obj := service.AppService()
	response, err := obj.AppService.ContentReportsService.FindOne(ctx, params)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}

	commons.Success(ctx, response, "success", params)
}

// @Tags 课程内容反馈接口
// @Summary 查询课程内容反馈列表分页数据
// @Description 查询课程内容反馈列表分页数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.ContentReportsListRequests true "列表参数"
// @Success 200 object responses.ContentReportsResponses
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"失败"}"
// @Router /editor/content_reports/list [post]
func ContentReportsList(ctx *gin.Context) {
	var params requests.ContentReportsListRequests
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		commons.Error(ctx, 500, err, "error")
	}

	//2.验证
	validations := commons.BaseValidations{}
	message, err := validations.Check(params)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	obj := service.AppService()
	response, err := obj.AppService.ContentReportsService.List(ctx, params)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}

	commons.Success(ctx, response, "success", params)
}

// @Tags 课程内容反馈接口
// @Summary 更新课程内容反馈数据
// @Description 更新课程内容反馈数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.ContentReportsUpdateRequests true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"失败"}"
// @Router /editor/content_reports/update [post]
func ContentReportsUpdate(ctx *gin.Context) {
	var params requests.ContentReportsUpdateRequests
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		commons.Error(ctx, 500, err, "error")
	}

	//2.验证
	validations := commons.BaseValidations{}
	message, err := validations.Check(params)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	obj := service.AppService()
	response, err := obj.AppService.ContentReportsService.Update(ctx, params)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}

	commons.Success(ctx, response, "success", params)
}

// @Tags 课程内容反馈接口
// @Summary 删除课程内容反馈数据
// @Description 删除课程内容反馈数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.ContentReportsDeleteRequests true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"失败"}"
// @Router /editor/content_reports/delete [post]
func ContentReportsDelete(ctx *gin.Context) {
	var params requests.ContentReportsDeleteRequests
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		commons.Error(ctx, 500, err, "error")
	}

	//2.验证
	validations := commons.BaseValidations{}
	message, err := validations.Check(params)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	obj := service.AppService()
	response, err := obj.AppService.ContentReportsService.Delete(ctx, params)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}

	commons.Success(ctx, response, "success", params)
}
