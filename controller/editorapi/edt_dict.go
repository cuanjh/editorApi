package editorapi

import (
	"editorApi/commons"
	"editorApi/init/initNats"
	"editorApi/init/mgdb"
	"editorApi/requests"
	"editorApi/service"
	"editorApi/tools/helpers"
	"strings"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func DictCollect(ctx *gin.Context) {
	header, err := ctx.FormFile("filename")
	if err != nil {
		commons.Error(ctx, 500, err, "文件名不能为空！")
	}
	dst := "data/dict/" + uuid.NewV4().String() + ".txt"
	// gin 简单做了封装,拷贝了文件流
	if err := ctx.SaveUploadedFile(header, dst); err != nil {
		commons.Error(ctx, 500, err, "文件保存失败！")
	}

	var param = requests.DictHanderParams{
		FilePath: dst,
		From:     "eng",
		To:       "chi",
	}

	obj := service.AppService()
	obj.AppService.DictService.ReadFile(&param)

	commons.Success(ctx, nil, "成功！", nil)
}

func DictUpload(ctx *gin.Context) {
	header, err := ctx.FormFile("filename")
	if err != nil {
		commons.Error(ctx, 500, err, "文件名不能为空！")
	}
	dst := "data/dict/" + uuid.NewV4().String() + ".txt"
	// gin 简单做了封装,拷贝了文件流
	if err := ctx.SaveUploadedFile(header, dst); err != nil {
		commons.Error(ctx, 500, err, "文件保存失败！")
	}

	// 异步操作
	initNats.NatsConn.Publish("CollectDict",
		&requests.DictHanderParams{
			FilePath: dst,
			From:     "eng",
			To:       "chi",
		},
	)

	commons.Success(ctx, nil, "成功！", nil)
}

func HanderDict(param *requests.DictHanderParams) {
	obj := service.AppService()
	obj.AppService.DictService.ReadFile(param)
}

// @Tags Dict（辞典接口）
// @Summary 辞典列表
// @Description 辞典列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body requests.DictListRequests true "辞典列表"
// @Success 200 object responses.DictResponse
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/dict/list [post]
func DictList(ctx *gin.Context) {
	var request requests.DictListRequests
	ctx.BindJSON(&request)

	//2.验证
	UserValidations := commons.BaseValidations{}
	message, err := UserValidations.Check(request)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	obj := service.AppService()
	result, err := obj.AppService.DictService.List(ctx, request)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}
	commons.Success(ctx, result, "查询成功！", request)
}

// @Tags Dict（辞典接口）
// @Summary 辞典详情
// @Description 辞典详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body requests.DictDetailRequests true "辞典列表"
// @Success 200 object responses.DictResponse
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/dict/detail [post]
func DictDetail(ctx *gin.Context) {
	var request requests.DictDetailRequests
	ctx.BindJSON(&request)

	//2.验证
	UserValidations := commons.BaseValidations{}
	message, err := UserValidations.Check(request)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	obj := service.AppService()
	result, err := obj.AppService.DictService.Detail(ctx, request)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}
	commons.Success(ctx, result, "查询成功！", request)
}

// @Tags Dict（辞典接口）
// @Summary 辞典更新
// @Description 辞典更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body requests.DictUpdateRequests true "辞典列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/dict/update [post]
func DictUpdate(ctx *gin.Context) {
	var request requests.DictUpdateRequests
	ctx.BindJSON(&request)

	//2.验证
	UserValidations := commons.BaseValidations{}
	message, err := UserValidations.Check(request)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}

	obj := service.AppService()
	result, err := obj.AppService.DictService.Update(ctx, request)

	if err != nil {
		commons.Error(ctx, 500, err, err.Error())
	}
	commons.Success(ctx, result, "查询成功！", request)
}

// @Tags Dict（辞典接口）
// @Summary 单词删除
// @Description 单词删除
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body requests.DictDelRequests true "单词删除参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/dict/del [post]
func DictDel(ctx *gin.Context) {
	var request requests.DictDelRequests
	ctx.BindJSON(&request)

	//2.验证
	UserValidations := commons.BaseValidations{}
	message, err := UserValidations.Check(request)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}
	if request.CType == "word" {
		mgdb.UpdateMany(
			mgdb.EnvEditor,
			mgdb.DbDict,
			"dict_"+strings.ToLower(request.From),
			bson.M{
				"uuid": bson.M{
					"$in": request.Uuids,
				},
			},
			bson.M{
				"$set": bson.M{
					"is_del": true,
				},
			},
			false,
		)
	} else {
		mgdb.UpdateMany(
			mgdb.EnvEditor,
			mgdb.DbDict,
			"sentence_"+strings.ToLower(request.From),
			bson.M{
				"uuid": bson.M{
					"$in": request.Uuids,
				},
			},
			bson.M{
				"$set": bson.M{
					"is_del": true,
				},
			},
			false,
		)
	}

	initNats.NatsConn.Publish(
		"DictChangeMsg",
		requests.DictESParams{
			CType:   request.CType,
			From:    strings.ToLower(request.From),
			To:      strings.ToLower(request.To),
			Uuids:   request.Uuids,
			Operate: "delete",
		},
	)

	commons.Success(ctx, nil, "删除成功！", request)
}

// @Tags Dict（辞典接口）
// @Summary 辞典上线
// @Description 辞典上线
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body requests.DictOnlineRequests true "辞典上线参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/dict/online [post]
func DictOnline(ctx *gin.Context) {
	var request requests.DictOnlineRequests
	ctx.BindJSON(&request)

	//2.验证
	UserValidations := commons.BaseValidations{}
	message, err := UserValidations.Check(request)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}
	if request.CType == "word" {
		mgdb.UpdateMany(
			mgdb.EnvEditor,
			mgdb.DbDict,
			"dict_"+strings.ToLower(request.From),
			bson.M{
				"uuid": bson.M{
					"$in": request.Uuids,
				},
			},
			bson.M{
				"$set": bson.M{
					"done": true,
				},
			},
			false,
		)
	} else {
		mgdb.UpdateMany(
			mgdb.EnvEditor,
			mgdb.DbDict,
			"sentence_"+strings.ToLower(request.From),
			bson.M{
				"uuid": bson.M{
					"$in": request.Uuids,
				},
			},
			bson.M{
				"$set": bson.M{
					"done": true,
				},
			},
			false,
		)
	}
	initNats.NatsConn.Publish(
		"DictChangeMsg",
		requests.DictESParams{
			CType:   request.CType,
			From:    strings.ToLower(request.From),
			To:      strings.ToLower(request.To),
			Uuids:   request.Uuids,
			Operate: "online",
		},
	)
	commons.Success(ctx, nil, "上线成功！", request)
}

// @Tags Dict（辞典接口）
// @Summary 辞典下线
// @Description 辞典下线
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body requests.DictOnlineRequests true "辞典下线参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/dict/offline [post]
func DictOffline(ctx *gin.Context) {
	var request requests.DictOnlineRequests
	ctx.BindJSON(&request)

	//2.验证
	UserValidations := commons.BaseValidations{}
	message, err := UserValidations.Check(request)
	if err != nil && !helpers.Empty(message) {
		commons.Error(ctx, 500, err, message)
	}
	if request.CType == "word" {
		mgdb.UpdateMany(
			mgdb.EnvEditor,
			mgdb.DbDict,
			"dict_"+strings.ToLower(request.From),
			bson.M{
				"uuid": bson.M{
					"$in": request.Uuids,
				},
			},
			bson.M{
				"$set": bson.M{
					"done": false,
				},
			},
			false,
		)
	} else {
		mgdb.UpdateMany(
			mgdb.EnvEditor,
			mgdb.DbDict,
			"sentence_"+strings.ToLower(request.From),
			bson.M{
				"uuid": bson.M{
					"$in": request.Uuids,
				},
			},
			bson.M{
				"$set": bson.M{
					"done": false,
				},
			},
			false,
		)
	}
	initNats.NatsConn.Publish(
		"DictChangeMsg",
		requests.DictESParams{
			CType:   request.CType,
			From:    strings.ToLower(request.From),
			To:      strings.ToLower(request.To),
			Uuids:   request.Uuids,
			Operate: "offline",
		},
	)
	commons.Success(ctx, nil, "下线成功！", request)
}

// @Tags Dict（辞典接口）
// @Summary 辞典Tags
// @Description 辞典Tags
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Failure 400 {string} string "{"code":500,"data":{},"msg":"获取失败信息"}"
// @Router /editor/dict/tags [post]
func DictTags(ctx *gin.Context) {
	var tags []requests.Tag
	tags = append(tags, requests.Tag{Key: "ENG-Basic", Name: "Pro官方课程"})
	tags = append(tags, requests.Tag{Key: "KEN-Basic", Name: "Kid官方课程"})
	tags = append(tags, requests.Tag{Key: "ENG-TextbookVersion", Name: "小学英语（通用版）"})
	tags = append(tags, requests.Tag{Key: "ENG-PrimaryPEP", Name: "小学英语（人教版）"})
	tags = append(tags, requests.Tag{Key: "ENG-SEPOxford", Name: "小学英语（沪教版）"})
	commons.Success(ctx, tags, "查询成功！", nil)
}

func DictAddTag(ctx *gin.Context) {
	header, err := ctx.FormFile("filename")
	if err != nil {
		commons.Error(ctx, 500, err, "文件名不能为空！")
	}
	dst := "data/dict/" + uuid.NewV4().String() + ".xlsx"
	// gin 简单做了封装,拷贝了文件流
	if err := ctx.SaveUploadedFile(header, dst); err != nil {
		commons.Error(ctx, 500, err, "文件保存失败！")
	}

	var request requests.DictAddTag
	ctx.BindJSON(&request)
	request.FilePath = dst
	request.From = "eng"
	request.To = "chi"

	obj := service.AppService()
	obj.AppService.DictService.DictAddTag(ctx, request)

	commons.Success(ctx, nil, "成功！", nil)
}
