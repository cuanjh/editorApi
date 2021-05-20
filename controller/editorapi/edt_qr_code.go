package editorapi

import (
	"editorApi/commons"
	"editorApi/requests"
	"editorApi/service"
	"editorApi/tools"
	"editorApi/tools/helpers"
	"fmt"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
)

// @Tags QRcodeAPI（二维码管理）
// @Summary 添加二维码
// @Description 添加二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.QRcodeAddRequests true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/qr_code/add [post]
func QRcodeAdd(ctx *gin.Context) {
	var request requests.QRcodeAddRequests
	ctx.BindJSON(&request)

	//2.验证
	UserValidations := commons.BaseValidations{}
	message, err := UserValidations.Check(request)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	obj := service.AppService()
	result, err := obj.AppService.QRcodeService.QRcodeAdd(ctx, request)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}
	commons.Success(ctx, result, "成功！", request)
}

// @Tags QRcodeAPI（二维码管理）
// @Summary 二维码列表
// @Description 二维码列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.QRcodeListRequests true "列表参数"
// @Success 200 object responses.TeacherResponse
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/qr_code/list [post]
func QRcodeList(ctx *gin.Context) {
	var request requests.QRcodeListRequests
	ctx.BindJSON(&request)

	//2.验证
	UserValidations := commons.BaseValidations{}
	message, err := UserValidations.Check(request)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	obj := service.AppService()
	result, err := obj.AppService.QRcodeService.QRcodeList(ctx, request)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}
	commons.Success(ctx, result, "成功！", request)
}

// @Tags QRcodeAPI（二维码管理）
// @Summary 更新二维码
// @Description 更新二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.QRcodeUpdateRequests true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/qr_code/update [post]
func QRcodeUpdate(ctx *gin.Context) {
	var request requests.QRcodeUpdateRequests
	ctx.BindJSON(&request)

	//2.验证
	UserValidations := commons.BaseValidations{}
	message, err := UserValidations.Check(request)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	obj := service.AppService()
	result, err := obj.AppService.QRcodeService.QRcodeUpdate(ctx, request)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}
	commons.Success(ctx, result, "成功！", request)
}

// @Tags QRcodeAPI（二维码管理）
// @Summary 删除二维码
// @Description 删除二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.QRcodeDeleteRequests true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/qr_code/delete [post]
func QRcodeDelete(ctx *gin.Context) {
	var request requests.QRcodeDeleteRequests
	ctx.BindJSON(&request)

	//2.验证
	UserValidations := commons.BaseValidations{}
	message, err := UserValidations.Check(request)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	obj := service.AppService()
	result, err := obj.AppService.QRcodeService.QRcodeDelete(ctx, request)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}
	commons.Success(ctx, result, "成功！", request)
}

// @Tags QRcodeAPI（二维码管理）
// @Summary 生成二维码
// @Description 生成二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.QRcodeImageRequests true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/qr_code/image [post]
func QRcodeImage(ctx *gin.Context) {
	var request requests.QRcodeImageRequests
	ctx.BindJSON(&request)

	//2.验证
	UserValidations := commons.BaseValidations{}
	message, err := UserValidations.Check(request)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	width, height := 400, 400
	if !helpers.Empty(request.Size) {
		width, height = request.Size, request.Size
	}

	qrCodeFilePath := "data/qrcode/"
	qrop := tools.NewQrCode(
		request.Url,
		width,
		height,
		qr.M,
		qr.Auto,
	)
	qrCodeFileName, _, err := qrop.Encode(qrCodeFilePath)
	fmt.Println(qrCodeFileName)
	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}
	commons.Success(ctx, "editor/"+qrCodeFilePath+qrCodeFileName, "成功！", nil)
}

// @Tags QRcodeAPI（二维码管理）
// @Summary 生成二维码
// @Description 生成二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body requests.QRcodeDetailsRequests true "列表参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/qr_code/details [get]
func QRcodeDetails(ctx *gin.Context) {
	var request requests.QRcodeDetailsRequests
	//ctx.BindJSON(&request)
	request.UUID = ctx.Query("uuid")

	//2.验证
	UserValidations := commons.BaseValidations{}
	message, err := UserValidations.Check(request)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	obj := service.AppService()
	result, err := obj.AppService.QRcodeService.QRcodeDetails(ctx, request)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}
	commons.SuccessJsonp(ctx, result, "成功！", request)
}
