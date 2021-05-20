package editorapi

import (
	"editorApi/controller/servers"
	"editorApi/mdbmodel/editor"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
)

var tblDiscoverChannel = "discoverChannels"

// @Tags DisChannelAPI（发现频道接口）
// @Summary 获取频道列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /dis/channel/list [get]
func DisChannelList(ctx *gin.Context) {
	var channels []*editor.DiscoverChannel
	cusor, _ := toClient.Database(KUYU).Collection(tblDiscoverChannel).Find(
		ctx, bson.M{
			"isDel": false,
		},
	)
	defer cusor.Close(ctx)
	cusor.All(ctx, &channels)
	servers.ReportFormat(
		ctx,
		true,
		"频道列表",
		gin.H{
			"channels": channels,
		},
	)
}

// @Tags DisChannelAPI（发现频道接口）
// @Summary 增加频道
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editor.DiscoverChannel true "频道数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /dis/channel/add [post]
func DisChannelAdd(ctx *gin.Context) {
	var channel editor.DiscoverChannel
	ctx.BindJSON(&channel)
	channel.UUID = uuid.NewV4().String()
	channel.IsDel = false
	toClient.Database(KUYU).Collection(tblDiscoverChannel).InsertOne(ctx, channel)
	servers.ReportFormat(
		ctx,
		true,
		"添加成功",
		gin.H{},
	)
}

// @Tags DisChannelAPI（发现频道接口）
// @Summary 编辑频道
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editor.DiscoverChannel true "频道数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /dis/channel/edit [post]
func DisChannelEdit(ctx *gin.Context) {
	var channel editor.DiscoverChannel
	ctx.BindJSON(&channel)
	channel.IsDel = false
	toClient.Database(KUYU).Collection(tblDiscoverChannel).UpdateOne(
		ctx,
		bson.M{
			"uuid": channel.UUID,
		},
		bson.M{
			"$set": channel,
		},
	)
	servers.ReportFormat(
		ctx,
		true,
		"编辑成功",
		gin.H{},
	)
}

type radioChannelDelParam struct {
	UUID string `bson:"uuid" json:"uuid"`
}

// @Tags DisChannelAPI（发现频道接口）
// @Summary 删除频道
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.radioChannelDelParam true "频道数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /dis/channel/del [post]
func DisChannelDel(ctx *gin.Context) {
	var channel radioChannelDelParam
	ctx.BindJSON(&channel)

	toClient.Database(KUYU).Collection(tblDiscoverChannel).UpdateOne(
		ctx,
		bson.M{
			"uuid": channel.UUID,
		},
		bson.M{
			"$set": bson.M{
				"isDel": true,
			},
		},
	)

	servers.ReportFormat(
		ctx,
		true,
		"删除成功",
		gin.H{},
	)
}

type channelListOrder struct {
	UUID      string `json:"uuid"`
	ListOrder int    `json:"listOrder"`
}

type channelListOrders []*channelListOrder

// @Tags DisChannelAPI（发现频道接口）
// @Summary 频道排序
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body editorapi.channelListOrders true "频道数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /dis/channel/listorder [post]
func DisChannelListOrder(ctx *gin.Context) {
	var listOrders channelListOrders
	ctx.BindJSON(&listOrders)
	for _, l := range listOrders {
		toClient.Database(KUYU).Collection(tblDiscoverChannel).UpdateOne(
			ctx,
			bson.M{
				"uuid": l.UUID,
			},
			bson.M{
				"$set": bson.M{
					"listOrder": l.ListOrder,
				},
			},
		)
	}
	servers.ReportFormat(
		ctx,
		true,
		"排序成功",
		gin.H{},
	)
}
