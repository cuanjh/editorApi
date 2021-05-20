package editorapi

import (
	"editorApi/commons"
	"editorApi/requests"
	"editorApi/service"
	"editorApi/tools/helpers"
	"github.com/gin-gonic/gin"
)

// @Tags Actors接口
// @Summary 添加Actors
// @Description 添加Actors
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.ActorsCreateRequests true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"失败"}"
// @Router /editor/actors/create [post]
func ActorsCreate(ctx *gin.Context) {
	var params requests.ActorsCreateRequests
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
	response, err := obj.AppService.ActorsService.Create(ctx, params)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}

	commons.Success(ctx, response, "success",params)
}

// @Tags Actors接口
// @Summary 查询Actors数据
// @Description 查询Actors数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.ActorsFindRequests true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"失败"}"
// @Router /editor/actors/find [post]
func ActorsFind(ctx *gin.Context) {
	var params requests.ActorsFindRequests
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
	response, err := obj.AppService.ActorsService.Find(ctx, params)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}

	commons.Success(ctx, response, "success",params)
}

// @Tags Actors接口
// @Summary 查询一条Actors数据
// @Description 查询一条Actors数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.ActorsFindOneRequests true "列表参数"
// @Success 200 object responses.ActorsResponses
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"失败"}"
// @Router /editor/actors/findone [post]
func ActorsFindOne(ctx *gin.Context) {
	var params requests.ActorsFindOneRequests
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
	response, err := obj.AppService.ActorsService.FindOne(ctx, params)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}

	commons.Success(ctx, response, "success",params)
}

// @Tags Actors接口
// @Summary 查询Actors列表分页数据
// @Description 查询Actors列表分页数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.ActorsListRequests true "列表参数"
// @Success 200 object responses.ActorsResponses
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"失败"}"
// @Router /editor/actors/list [post]
func ActorsList(ctx *gin.Context) {
	var params requests.ActorsListRequests
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
	response, err := obj.AppService.ActorsService.List(ctx, params)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}

	commons.Success(ctx, response, "success",params)
}

// @Tags Actors接口
// @Summary 更新Actors数据
// @Description 更新Actors数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.ActorsUpdateRequests true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"失败"}"
// @Router /editor/actors/update [post]
func ActorsUpdate(ctx *gin.Context) {
	var params requests.ActorsUpdateRequests
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
	response, err := obj.AppService.ActorsService.Update(ctx, params)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}

	commons.Success(ctx, response, "success",params)
}

// @Tags Actors接口
// @Summary 删除Actors数据
// @Description 删除Actors数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.ActorsDeleteRequests true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"失败"}"
// @Router /editor/actors/delete [post]
func ActorsDelete(ctx *gin.Context) {
	var params requests.ActorsDeleteRequests
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
	response, err := obj.AppService.ActorsService.Delete(ctx, params)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}

	commons.Success(ctx, response, "success",params)
}
