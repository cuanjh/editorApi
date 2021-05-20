package repository

import (
	"editorApi/commons"
	"editorApi/init/mgdb"
	"editorApi/mdbmodel/editor"
	"editorApi/requests"
	"editorApi/responses"
	"editorApi/tools/helpers"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type QRcode struct {
}

func (m *QRcode) Details(ctx *gin.Context, params requests.QRcodeDetailsRequests) (result responses.QRcodeResponses, err error) {
	collection := mgdb.MongoClient.Database(mgdb.DbEditor).Collection(editor.TbQRcode)
	var filter = bson.D{}
	if !helpers.Empty(params.UUID) {
		filter = append(filter, bson.E{"uuid", params.UUID})
	}
	dataResult := collection.FindOne(ctx, filter)
	dataResult.Decode(&result)
	return
}

func (m *QRcode) Add(ctx *gin.Context, params requests.QRcodeAddRequests) (inserted_id interface{}, err error) {
	collection := mgdb.MongoClient.Database(mgdb.DbEditor).Collection(editor.TbQRcode)
	var data = bson.M{
		"uuid":        uuid.NewV4().String(),
		"title":       params.Title,
		"info":        params.Info,
		"size":        params.Size,
		"is_del":      false,
		"created_on":  time.Now(),
		"update_time": time.Now(),
	}
	insertResult, err := collection.InsertOne(ctx, data)
	if err != nil {
		return
	}
	inserted_id = insertResult.InsertedID
	return
}

func (m *QRcode) Update(ctx *gin.Context, params requests.QRcodeUpdateRequests) (upserted_id interface{}, err error) {
	collection := mgdb.MongoClient.Database(mgdb.DbEditor).Collection(editor.TbQRcode)
	var filter = bson.D{}
	if !helpers.Empty(params.UUID) {
		filter = append(filter, bson.E{"uuid", params.UUID})
	}

	var update = bson.M{
		"title":       params.Title,
		"info":        params.Info,
		"size":        params.Size,
		"update_time": time.Now(),
	}

	updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return
	}
	upserted_id = updateResult.UpsertedID
	return
}

func (m *QRcode) List(ctx *gin.Context, params requests.QRcodeListRequests) (result []responses.QRcodeResponses, err error) {
	collection := mgdb.MongoClient.Database(mgdb.DbEditor).Collection(editor.TbQRcode)
	var filter = bson.D{}
	filter = append(filter, bson.E{"is_del", false})

	page := commons.DefaultPage()
	if !helpers.Empty(params.PageSize) {
		page.Limit = params.PageSize
	}

	if !helpers.Empty(params.PageIndex) && params.PageIndex > 0 {
		page.Skip = (params.PageIndex - 1) * page.Limit
	}

	var rank = bson.M{"created_on": -1}

	if !helpers.Empty(params.SortType) && !helpers.Empty(params.TextField) {
		rank = bson.M{params.TextField: params.SortType}
	}

	option := options.Find().SetSort(rank).SetLimit(page.Limit).SetSkip(page.Skip)

	cusor, err := collection.Find(
		ctx,
		filter,
		option,
	)

	defer cusor.Close(ctx)
	err = cusor.All(ctx, &result)
	return
}

func (s *QRcode) Delete(ctx *gin.Context, params requests.QRcodeDeleteRequests) (upserted_id interface{}, err error) {
	collection := mgdb.MongoClient.Database(mgdb.DbEditor).Collection(editor.TbQRcode)
	var filter = bson.D{}
	if !helpers.Empty(params.UUID) {
		filter = append(filter, bson.E{"uuid", params.UUID})
	}

	var update = bson.M{
		"is_del":      true,
		"update_time": time.Now(),
	}

	updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return
	}
	upserted_id = updateResult.UpsertedID
	return
}
