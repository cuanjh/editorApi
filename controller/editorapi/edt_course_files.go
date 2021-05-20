package editorapi

import (
	"editorApi/commons"
	"editorApi/requests"
	"editorApi/service"
	"editorApi/tools/helpers"
	"github.com/gin-gonic/gin"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// @Tags 文档转码相关接口
// @Summary 文件转码列表
// @Description 文件转码列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.CourseFilesDeleteRequests true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/course_files/delete_file [post]
func DeleteFile(ctx *gin.Context) {
	var request requests.CourseFilesDeleteRequests
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		commons.Error(ctx, 500, err, "error")
	}

	//2.验证
	Validations := commons.BaseValidations{}
	message, err := Validations.Check(request)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	obj := service.AppService()
	result, err := obj.AppService.CourseFilesService.DeleteFile(ctx, request)

	if err != nil {
		commons.Error(ctx, 500, err, "error")
	}
	commons.Success(ctx, result, "success", request)
}

// @Tags 文档转码相关接口
// @Summary 文件转码列表
// @Description 文件转码列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.CourseFilesListRequests true "列表参数"
// @Success 200 object responses.CourseFilesListResponse
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/course_files/file_list [post]
func FileList(ctx *gin.Context) {
	var request requests.CourseFilesListRequests
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		commons.Error(ctx, 500, err, "error")
	}

	obj := service.AppService()
	result, err := obj.AppService.CourseFilesService.List(ctx, request)
	if err != nil {
		commons.Error(ctx, 500, err, "error")
	}
	commons.Success(ctx, result, "success", request)
}

// @Tags 文档转码相关接口
// @Summary 创建文档转码任务
// @Description 创建文档转码任务
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.CourseFilesCreateTranscodeRequests true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/course_files/create_transcode [post]
func CreateTranscode(ctx *gin.Context) {
	var params requests.CourseFilesRequests
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		commons.Error(ctx, 500, err, "error")
	}

	obj := service.AppService()
	response, err := obj.AppService.CourseFilesService.CreateTranscode(ctx, params)

	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		commons.Error(ctx, 500, err, err.Error())
	}
	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}
	commons.Success(ctx, response.Response, "success", params)
}

// @Tags 文档转码相关接口
// @Summary 查询文档转码任务
// @Description 查询文档转码任务
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.CourseFilesTranscodeRequests true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/course_files/describe_transcode [post]
func DescribeTranscode(ctx *gin.Context) {
	var params requests.CourseFilesTranscodeRequests
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		commons.Error(ctx, 500, err, "error")
	}

	obj := service.AppService()
	result, err := obj.AppService.CourseFilesService.DescribeTranscode(ctx, params)

	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		commons.Error(ctx, 500, err, err.Error())
		return
	}
	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}
	commons.Success(ctx, result.Response, "success", params)
}
