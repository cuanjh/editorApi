package editorapi

import (
	"editorApi/controller/servers"
	"editorApi/init/mgdb"
	"editorApi/mdbmodel/editor"
	"editorApi/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type examinParam struct {
	CatalogUUID string `json:"catalog_uuid"` //目录UUID
	UserUUID    string `json:"user_uuid"`    //分配的权限UUID
	ExaminState int    `json:"examin_state"` //审核状态，0正在编辑，1提交审核，2审核通过，3审核没通过
	Comment     string `json:"comment"`      //审核评语
}

// @Tags EditorCatalogAPI(目录接口)
// @Summary 审核接口
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.examinParam true "内容审核: 审核状态，0正在编辑，1提交审核，2审核通过，3审核没通过"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/examin [post]
func Examin(ctx *gin.Context) {
	defer rcv(ctx)
	var param examinParam
	ctx.BindJSON(&param)
	claims, _ := ctx.Get("claims")
	waitUse := claims.(*middleware.CustomClaims)
	userUUID := waitUse.UUID.String()
	catalogsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs)
	catalogsCollection.UpdateOne(
		ctx,
		bson.M{
			"uuid":                  param.CatalogUUID,
			"authorities.user_uuid": param.UserUUID,
		},
		bson.M{
			"$set": bson.M{
				"authorities.$.examinState": param.ExaminState,
				"authorities.$.examinStateInfo": editor.ExaminStateInfo{
					UserUUID: userUUID,
					Comment:  param.Comment,
				},
			},
		},
	)
	servers.ReportFormat(ctx, true, "成功", gin.H{})
}

type examinSubmitParam struct {
	CatalogUUID string `json:"catalog_uuid"` //目录UUID
}

// @Tags EditorCatalogAPI(目录接口)
// @Summary 提交审核
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.examinSubmitParam true "提交审核"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/examin/submit [post]
func ExaminSubmit(ctx *gin.Context) {
	defer rcv(ctx)
	var param examinSubmitParam
	ctx.BindJSON(&param)
	claims, _ := ctx.Get("claims")
	waitUse := claims.(*middleware.CustomClaims)
	userUUID := waitUse.UUID.String()
	catalogsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs)
	catalogsCollection.UpdateOne(
		ctx,
		bson.M{
			"uuid": param.CatalogUUID,
			"authorities": bson.M{
				"$elemMatch": bson.M{
					"user_uuid": userUUID,
				},
			},
		},
		bson.M{
			"$set": bson.M{
				"authorities.$.examinState": 1,
			},
		},
	)
	servers.ReportFormat(ctx, true, "提交成功", gin.H{})
}

// @Tags EditorCatalogAPI(目录接口)
// @Summary 审核状态重置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.examinSubmitParam true "审核状态重置"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /editor/examin/reset [post]
func ExaminReset(ctx *gin.Context) {
	defer rcv(ctx)
	var param examinSubmitParam
	ctx.BindJSON(&param)
	claims, _ := ctx.Get("claims")
	waitUse := claims.(*middleware.CustomClaims)
	userUUID := waitUse.UUID.String()
	catalogsCollection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs)
	catalogsCollection.UpdateOne(
		ctx,
		bson.M{
			"uuid": param.CatalogUUID,
			"authorities": bson.M{
				"$elemMatch": bson.M{
					"user_uuid": userUUID,
				},
			},
		},
		bson.M{
			"$set": bson.M{
				"authorities.$.examinState": 0,
			},
		},
	)
	servers.ReportFormat(ctx, true, "重置成功", gin.H{})
}
