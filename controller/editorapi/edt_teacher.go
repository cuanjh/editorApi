package editorapi

import (
	"editorApi/commons"
	"editorApi/requests"
	"editorApi/service"
	"editorApi/tools/helpers"
	"github.com/gin-gonic/gin"
)

const (
	tblTeacher      = "teacher"
	tblTeacherAudit = "teacher_audit"
)

// @Tags EditorTeacherAPI（教师接口）
// @Summary 获取教师信息列表
// @Description 获取教师信息表
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.TeacherListRequests true "列表参数"
// @Success 200 object responses.TeacherResponse
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/teacher/list [post]
func TeacherList(ctx *gin.Context) {
	request := requests.TeacherListRequests{}
	ctx.BindJSON(&request)

	//2.验证
	UserValidations := commons.BaseValidations{}
	message, err := UserValidations.Check(request)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	obj := service.AppService()
	result, err := obj.AppService.TeacherService.List(ctx, request)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}
	commons.Success(ctx, result, "查询成功！", request)
}

// @Tags EditorTeacherAPI（教师接口）
// @Summary 获取教师信息详情
// @Description 获取教师信息详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.TeacherEditRequests true "参数"
// @Success 200 object responses.TeacherResponse
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/teacher/details [post]
func TeacherDetails(ctx *gin.Context) {
	request := requests.TeacherEditRequests{}
	ctx.BindJSON(&request)

	//2.验证
	UserValidations := commons.BaseValidations{}
	message, err := UserValidations.Check(request)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	//obj := service.AppService()
	//result, err := obj.AppService.TeacherService.Edit(ctx, request)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}
	commons.Success(ctx, nil, "查询成功！", request)
}

// @Tags EditorTeacherAPI（教师接口）
// @Summary 审核教师信息
// @Description 审核教师信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.TeacherAuditRequests true "参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"失败"}"
// @Router /editor/teacher/audit [post]
func Audit(ctx *gin.Context) {
	request := requests.TeacherAuditRequests{}
	ctx.BindJSON(&request)

	//2.验证
	UserValidations := commons.BaseValidations{}
	message, err := UserValidations.Check(request)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	obj := service.AppService()
	_, err = obj.AppService.TeacherAuditService.AuditTeacher(ctx, request)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}
	commons.Success(ctx, nil, "提交成功！", request)
}
