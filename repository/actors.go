package repository

import (
	"editorApi/commons"
	"editorApi/init/mgdb"
	"editorApi/requests"
	"editorApi/responses"
	"editorApi/tools/helpers"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const tb_actors = "actors"

type Actors struct {
	Uuid      string  `json:"uuid" bson:"uuid"`            // Uuid
	Name      string  `json:"name" bson:"name"`            // 声优名称
	Photo     string  `json:"photo" bson:"photo"`          // 头像
	Sound     string  `json:"sound" bson:"sound"`          // 声音地址
	Gender    int     `json:"gender" bson:"gender"`        // 性别 1,男；0，女；
	SoundTime float64 `json:"soundTime" bson:"sound_time"` // 时长
	Role      string  `json:"role" bson:"role"`            // 分组
	Country   string  `json:"country" bson:"country"`      // 国籍
	City      string  `json:"city" bson:"city"`            // 城市
	Feature   string  `json:"feature" bson:"feature"`      // 声音特点
	Lang      string  `json:"lang" bson:"lang"`            // 语言
	Desc      string  `json:"desc" bson:"desc"`            // 描述
	Status    int     `json:"status" bson:"status"`        // 状态 1，上线；2，下线
}

// 创建Actors
func (m *Actors) Create(ctx *gin.Context, params requests.ActorsCreateRequests) (inserted_id interface{}, err error) {
	db := mgdb.MongoClient.Database(mgdb.DbEditor).Collection(tb_actors)
	var add = bson.M{
		"uuid":       uuid.NewV4().String(),
		"name":       params.Name,
		"photo":      params.Photo,
		"sound":      params.Sound,
		"gender":     params.Gender,
		"sound_time": params.SoundTime,
		"role":       params.Role,
		"country":    params.Country,
		"city":       params.City,
		"feature":    params.Feature,
		"lang":       params.Lang,
		"desc":       params.Desc,
		"status":     params.Status,
	}

	insertResult, err := db.InsertOne(ctx, add)
	if err != nil {
		return
	}
	inserted_id = insertResult.InsertedID
	return
}

// 获取Actors更多数据
func (m *Actors) Find(ctx *gin.Context, params requests.ActorsFindRequests) (result []responses.ActorsResponses, err error) {
	db := mgdb.MongoClient.Database(mgdb.DbEditor).Collection(tb_actors)

	// 过滤条件
	var filter = bson.D{}
	if !helpers.Empty(params.Uuid) {
		filter = append(filter, bson.E{"uuid", params.Uuid})
	}

	if !helpers.Empty(params.Lang) {
		filter = append(filter, bson.E{"lang", params.Lang})
	}

	cursor, err := db.Find(ctx, filter)
	if err != nil {
		return
	}

	err = cursor.All(ctx, &result)
	return
}

// 获取Actors
func (m *Actors) FindOne(ctx *gin.Context, params requests.ActorsFindOneRequests) (result responses.ActorsResponses, err error) {
	db := mgdb.MongoClient.Database(mgdb.DbEditor).Collection(tb_actors)

	// 过滤条件
	var filter = bson.D{}
	if !helpers.Empty(params.Uuid) {
		filter = append(filter, bson.E{"uuid", params.Uuid})
	}

	if !helpers.Empty(params.Lang) {
		filter = append(filter, bson.E{"lang", params.Lang})
	}

	singleResult := db.FindOne(ctx, filter)
	if err != nil {
		return
	}

	err = singleResult.Decode(&result)
	return
}

// 获取Actors带分页
func (m *Actors) List(ctx *gin.Context, params requests.ActorsListRequests) (result []responses.ActorsResponses, err error) {
	db := mgdb.MongoClient.Database(mgdb.DbEditor).Collection(tb_actors)
	// 过滤条件
	var filter = bson.D{}
	if !helpers.Empty(params.Lang) {
		filter = append(filter, bson.E{"lang", params.Lang})
	}
	page := commons.DefaultPage()
	if !helpers.Empty(params.PageSize) {
		page.Limit = params.PageSize
	}

	if !helpers.Empty(params.PageIndex) && params.PageIndex > 0 {
		page.Skip = (params.PageIndex - 1) * page.Limit
	}

	option := options.Find().SetSkip(page.Skip).SetLimit(page.Limit)
	if !helpers.Empty(params.SortType) && !helpers.Empty(params.TextField) {
		option = option.SetSort(bson.M{
			params.TextField: params.SortType,
		})
	}
	cursor, err := db.Find(ctx, filter, option)
	defer cursor.Close(ctx)

	if err != nil {
		return
	}
	err = cursor.All(ctx, &result)
	return
}

// 更新Actors
func (m *Actors) Update(ctx *gin.Context, params requests.ActorsUpdateRequests) (updateResult interface{}, err error) {
	db := mgdb.MongoClient.Database(mgdb.DbEditor).Collection(tb_actors)

	var filter = bson.D{}
	if !helpers.Empty(params.Uuid) {
		filter = append(filter, bson.E{"uuid", params.Uuid})
	}
	if !helpers.Empty(params.Lang) {
		filter = append(filter, bson.E{"lang", params.Lang})
	}

	var update = bson.M{
		"name":       params.Name,
		"photo":      params.Photo,
		"sound":      params.Sound,
		"gender":     params.Gender,
		"sound_time": params.SoundTime,
		"role":       params.Role,
		"country":    params.Country,
		"city":       params.City,
		"feature":    params.Feature,
		"lang":       params.Lang,
		"desc":       params.Desc,
		"status":     params.Status,
	}

	updateResult, err = db.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return
	}

	return
}

// 删除Actors
func (m *Actors) Delete(ctx *gin.Context, params requests.ActorsDeleteRequests) (deleteResult interface{}, err error) {
	db := mgdb.MongoClient.Database(mgdb.DbEditor).Collection(tb_actors)
	var filter = bson.D{}
	if !helpers.Empty(params.Uuid) {
		filter = append(filter, bson.E{"uuid", params.Uuid})
	}
	deleteResult, err = db.DeleteOne(ctx, filter)
	return
}
