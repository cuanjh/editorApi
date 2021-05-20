package editorapi

import (
	"editorApi/controller/servers"
	"editorApi/mdbmodel/editor"
	"editorApi/tools/helpers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	tbClassWeixin = "class_weixin"
)

// @Tags LiveAPI（课程微信接口）
// @Summary 课程微信列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.classWeixinParams true "分页参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /editor/class_weixin/list [post]
func ClassWeixinList(ctx *gin.Context) {
	param := classWeixinParams{}
	ctx.BindJSON(&param)
	var limit int64
	var skip int64

	limit = param.PageSize
	if helpers.Empty(limit) {
		limit = 40
	}
	skip = param.PageNo * limit

	if helpers.Empty(param.CourseCode) {
		servers.ReportFormat(ctx, false, "课程编码不能为空！", gin.H{})
		return
	}

	class_weixin := []*editor.ClassWeixin{}

	roomsCollection := toClient.Database(KUYU).Collection(tbClassWeixin)
	cusor, err := roomsCollection.Find(
		ctx,
		bson.M{
			"course_code": param.CourseCode,
		},
		options.Find().SetSort(bson.M{"created_time": -1}).SetLimit(limit).SetSkip(skip),
	)
	if !helpers.Empty(err) {
		checkErr(ctx, err)
		return
	}
	defer cusor.Close(ctx)

	cusor.All(ctx, &class_weixin)

	servers.ReportFormat(ctx, true, "课程微信列表", gin.H{
		"class_weixin": class_weixin,
	})
	return
}

type classWeixinParams struct {
	PageNo     int64  `json:"pageNo"`
	PageSize   int64  `json:"pageSize"`
	CourseCode string `bson:"course_code" json:"courseCode"`
}

// @Tags LiveAPI（课程微信接口）
// @Summary 添加课程微信
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.classWeixinParam true "课程微信数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /editor/class_weixin/add [post]
func ClassWeixinAdd(ctx *gin.Context) {
	var param classWeixinParam
	ctx.BindJSON(&param)

	if helpers.Empty(param.CourseCode) {
		servers.ReportFormat(ctx, false, "课程编码不能为空！", gin.H{})
		return
	}

	classWeixinCollection := toClient.Database(KUYU).Collection(tbClassWeixin)
	// 添加时，需要保证当前只有一个显示
	if param.IsShow {
		classWeixinCollection.UpdateMany(
			ctx,
			bson.M{
				"course_code": param.CourseCode,
			},
			bson.M{"$set": bson.M{
				"is_show": false,
			}},
		)
	}

	param.CreatedTime = time.Now().Unix()
	insertResult, err := classWeixinCollection.InsertOne(ctx, param)

	if !helpers.Empty(err) || helpers.Empty(insertResult.InsertedID) {
		servers.ReportFormat(ctx, false, "添加失败", gin.H{})
		return
	}
	servers.ReportFormat(ctx, true, "添加成功", gin.H{})
}

type classWeixinParam struct {
	CourseCode  string `bson:"course_code" json:"courseCode"`
	WeixinNo    string `bson:"weixin_no" json:"weixinNo"`
	WeixinCode  string `bson:"weixin_code" json:"weixinCode"`
	IsShow      bool   `bson:"is_show" json:"isShow"`
	CreatedTime int64  `bson:"created_time" json:"createdTime"`
}

// @Tags LiveAPI（课程微信接口）
// @Summary 编辑课程微信
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.classWeixinEditPara true "课程微信数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /editor/class_weixin/edit [put]
func ClassWeixinEdit(ctx *gin.Context) {
	var paras = classWeixinEditPara{}
	ctx.BindJSON(&paras)

	if helpers.Empty(paras.Id) {
		servers.ReportFormat(ctx, false, "ID不能为空！", gin.H{})
		return
	}

	classWeixinCollection := toClient.Database(KUYU).Collection(tbClassWeixin)
	_id, _ := primitive.ObjectIDFromHex(paras.Id)

	// 编辑时，需要保证当前只有一个显示
	if paras.IsShow {
		classWeixinCollection.UpdateMany(
			ctx,
			bson.M{
				"course_code": paras.CourseCode,
			},
			bson.M{"$set": bson.M{
				"is_show": false,
			}},
		)
	}

	classWeixinCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": _id,
		},
		bson.M{"$set": bson.M{
			"course_code": paras.CourseCode,
			"weixin_no":   paras.WeixinNo,
			"weixin_code": paras.WeixinCode,
			"is_show":     paras.IsShow,
		}},
	)

	servers.ReportFormat(ctx, true, "编辑成功", gin.H{})
}

type classWeixinEditPara struct {
	Id         string `bson:"_id" json:"id"`
	CourseCode string `bson:"course_code" json:"courseCode"`
	WeixinNo   string `bson:"weixin_no" json:"weixinNo"`
	WeixinCode string `bson:"weixin_code" json:"weixinCode"`
	IsShow     bool   `bson:"is_show" json:"isShow"`
}

// @Tags LiveAPI（课程微信接口）
// @Summary 删除课程微信
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.classWeixinDelPara true "课程微信数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /editor/class_weixin/del [delete]
func ClassWeixinDel(ctx *gin.Context) {
	paras := classWeixinDelPara{}
	ctx.BindJSON(&paras)

	if helpers.Empty(paras.Id) {
		servers.ReportFormat(ctx, false, "ID不能为空！", gin.H{})
		return
	}

	classWeixinCollection := toClient.Database(KUYU).Collection(tbClassWeixin)
	_id, _ := primitive.ObjectIDFromHex(paras.Id)
	filter := bson.M{"_id": _id}
	res, err := classWeixinCollection.DeleteOne(
		ctx,
		filter,
		nil,
	)

	if !helpers.Empty(err) || helpers.Empty(res.DeletedCount) {
		servers.ReportFormat(ctx, false, "数据不存在或删除失败", gin.H{})
		return
	}
	servers.ReportFormat(ctx, true, "删除成功", gin.H{})
}

type classWeixinDelPara struct {
	Id string `json:"id"`
}
