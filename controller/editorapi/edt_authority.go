package editorapi

import (
	"fmt"
	"editorApi/controller/servers"
	"editorApi/init/mgdb"
	"editorApi/mdbmodel/editor"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var tblAuthority string = "catalogs_authority"

type authority struct {
	UserUUID  string `json:"user_uuid"` //目录分配给编辑的UUID
	Authority string `json:"authority"` //权限：r/rw
}

type authoritySetParam struct {
	Type        string      `json:"type"`
	CatalogUUID string      `json:"uuid"`
	Authorities []authority `json:"authorities"`
}

// @Tags EditorAuthorityAPI(课程编辑权限接口)
// @Summary 设置课程编辑权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.authoritySetParam true "权限设置:authority的值为r(读)\rw(读写);type值：catalog/content_version"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"设置成功"}"
// @Router /editor/authority/set [post]
func AuthoritySet(ctx *gin.Context) {
	defer rcv(ctx)
	var param authoritySetParam
	ctx.BindJSON(&param)

	collection := mgdb.MongoClient.Database(EDITOR_DB).Collection(tblCatalogs)
	if param.Type == "content_version" {
		collection = mgdb.MongoClient.Database(EDITOR_DB).Collection(tblContentInfo)
	}
	fmt.Println("%v", param)

	// authorities := []editor.CatalogAuthority{}
	for _, a := range param.Authorities {

		collection.UpdateOne(
			ctx,
			bson.M{"uuid": param.CatalogUUID},
			bson.M{
				"$pull": bson.M{
					"authorities": bson.M{"user_uuid": a.UserUUID},
				},
			},
		)
		collection.UpdateOne(
			ctx,
			bson.M{
				"uuid": param.CatalogUUID,
			},
			bson.M{
				"$push": bson.M{
					"authorities": editor.CatalogAuthority{
						UserUUID:        a.UserUUID,
						Auth:            a.Authority,
						ExaminState:     0,
						ExaminStateInfo: editor.ExaminStateInfo{},
					},
				},
			},
		)
	}

	servers.ReportFormat(ctx, true, "成功", gin.H{})
}
