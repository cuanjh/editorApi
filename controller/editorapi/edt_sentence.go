package editorapi

import (
	"editorApi/commons"
	"editorApi/init/mgdb"
	"editorApi/requests"
	"editorApi/service"
	"editorApi/tools/helpers"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Tags Sentence（句子接口）
// @Summary 句子列表
// @Description 句子列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body requests.SentenceSearchRequests true "句子参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/sentence/list [post]
func SentenceList(ctx *gin.Context) {
	var request requests.SentenceSearchRequests
	ctx.BindJSON(&request)
	filter := bson.M{}
	if request.SearchType == 0 {
		filter["sentence"] = request.Sentence
	} else {
		filter["sentence"] = primitive.Regex{
			Pattern: request.Sentence,
			Options: "i",
		}
	}
	skip := (request.PageIndex - 1) * request.PageSize
	rsts := []*requests.Sentence{}

	mgdb.Find(
		mgdb.EnvEditor,
		mgdb.DbDict,
		"sentence_"+strings.ToLower(request.From),
		filter,
		bson.M{"created_on": -1},
		nil,
		skip,
		request.PageSize,
		&rsts,
	)
	commons.Success(ctx, rsts, "查询成功！", request)
}

// @Tags Sentence（句子接口）
// @Summary 句子详情
// @Description 句子详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body requests.DictDetailRequests true "辞典列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/sentence/detail [post]
func SentenceDetail(ctx *gin.Context) {
	var request requests.SentenceDetail
	ctx.BindJSON(&request)

	//2.验证
	UserValidations := commons.BaseValidations{}
	message, err := UserValidations.Check(request)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	obj := service.AppService()
	result, err := obj.AppService.SentenceService.Detail(ctx, request)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}
	commons.Success(ctx, result, "查询成功！", request)
}

// @Tags Sentence（句子接口）
// @Summary 更新句子
// @Description 更新句子
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body requests.SentenceUpdate true "更新句子"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/sentence/update [post]
func SentenceUpdate(ctx *gin.Context) {
	var request requests.SentenceUpdate
	ctx.BindJSON(&request)
	//2.验证
	UserValidations := commons.BaseValidations{}
	message, err := UserValidations.Check(request)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	obj := service.AppService()
	result, err := obj.AppService.SentenceService.Update(ctx, request)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}
	commons.Success(ctx, result, "更新成功！", request)
}
