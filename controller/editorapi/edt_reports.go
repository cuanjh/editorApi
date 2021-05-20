package editorapi

import (
	"editorApi/commons"
	"editorApi/requests"
	"editorApi/service"
	"editorApi/tools/helpers"
	"github.com/gin-gonic/gin"
)

// @Tags 反馈接口
// @Summary 添加反馈
// @Description 添加反馈
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.ReportsCreateRequests true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"失败"}"
// @Router /editor/reports/create [post]
func ReportsCreate(ctx *gin.Context) {
	var params requests.ReportsCreateRequests
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
	response, err := obj.AppService.ReportsService.Create(ctx, params)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}

	commons.Success(ctx, response, "success", params)
}

// @Tags 反馈接口
// @Summary 查询反馈数据
// @Description 查询反馈数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.ReportsFindRequests true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"失败"}"
// @Router /editor/reports/find [post]
func ReportsFind(ctx *gin.Context) {
	var params requests.ReportsFindRequests
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
	response, err := obj.AppService.ReportsService.Find(ctx, params)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}

	commons.Success(ctx, response, "success", params)
}

// @Tags 反馈接口
// @Summary 查询一条反馈数据
// @Description 查询一条反馈数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.ReportsFindOneRequests true "列表参数"
// @Success 200 object responses.ReportsResponses
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"失败"}"
// @Router /editor/reports/findone [post]
func ReportsFindOne(ctx *gin.Context) {
	var params requests.ReportsFindOneRequests
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
	response, err := obj.AppService.ReportsService.FindOne(ctx, params)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}

	commons.Success(ctx, response, "success", params)
}

// @Tags 反馈接口
// @Summary 查询反馈列表分页数据
// @Description 查询反馈列表分页数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.ReportsListRequests true "列表参数"
// @Success 200 object responses.ReportsResponses
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"失败"}"
// @Router /editor/reports/list [post]
func ReportsList(ctx *gin.Context) {
	var params requests.ReportsListRequests
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
	response, err := obj.AppService.ReportsService.List(ctx, params)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}

	commons.Success(ctx, response, "success", params)
}

// @Tags 反馈接口
// @Summary 更新反馈数据
// @Description 更新反馈数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.ReportsUpdateRequests true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"失败"}"
// @Router /editor/reports/update [post]
func ReportsUpdate(ctx *gin.Context) {
	var params requests.ReportsUpdateRequests
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
	response, err := obj.AppService.ReportsService.Update(ctx, params)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}

	commons.Success(ctx, response, "success", params)
}

// @Tags 反馈接口
// @Summary 删除反馈数据
// @Description 删除反馈数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.ReportsDeleteRequests true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"失败"}"
// @Router /editor/reports/delete [post]
func ReportsDelete(ctx *gin.Context) {
	var params requests.ReportsDeleteRequests
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
	response, err := obj.AppService.ReportsService.Delete(ctx, params)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}

	commons.Success(ctx, response, "success", params)
}
