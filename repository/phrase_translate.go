package repository

import (
	"context"
	"editorApi/init/mgdb"
	"editorApi/requests"
	"editorApi/responses"
	"editorApi/tools/helpers"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

type PhraseTranslate struct {
	Parent    string `bson:"parent" json:"parent"`        // uuid
	ContentTr string `bson:"content_tr" json:"contentTr"` //翻译内容
}

func (m *PhraseTranslate) AddPhraseTranslate(params requests.PhraseTranslate) (id interface{}, err error) {
	if helpers.Empty(params.From) || helpers.Empty(params.To) {
		return id, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("phrase_" + strings.ToLower(params.From) + "_" + strings.ToLower(params.To))

	result, _ := m.FindOne(params)
	ctx := context.TODO()
	if !helpers.Empty(result.Parent) {
		var filter = bson.D{}
		if !helpers.Empty(params.Parent) {
			filter = append(filter, bson.E{"parent", params.Parent})
		}
		// 更新时间
		var update = bson.M{
			"content_tr": params.ContentTr,
		}
		updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": update})
		if err != nil {
			return 0, err
		}
		id = updateResult.UpsertedID
	} else {
		var add = bson.M{
			"parent":     params.Parent,
			"content_tr": params.ContentTr,
		}

		insertResult, err := collection.InsertOne(ctx, add)
		if err != nil {
			return 0, err
		}
		id = insertResult.InsertedID
	}
	return
}

func (m *PhraseTranslate) FindOne(params requests.PhraseTranslate) (result PhraseTranslate, err error) {
	if helpers.Empty(params.From) || helpers.Empty(params.To) {
		return result, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("phrase_" + strings.ToLower(params.From) + "_" + strings.ToLower(params.To))
	ctx := context.TODO()

	var filter = bson.D{}
	if !helpers.Empty(params.Parent) {
		filter = append(filter, bson.E{"parent", params.Parent})
	}

	singleResult := collection.FindOne(ctx, filter)
	err = singleResult.Decode(&result)
	return
}

func (m *PhraseTranslate) Detail(ctx *gin.Context, params requests.PhraseTranslateDetailRequests) (result responses.PhraseTranslateResponse, err error) {
	if helpers.Empty(params.From) || helpers.Empty(params.To) {
		return result, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("phrase_" + strings.ToLower(params.From) + "_" + strings.ToLower(params.To))
	var filter = bson.D{}
	if !helpers.Empty(params.Parent) {
		filter = append(filter, bson.E{"parent", params.Parent})
	}
	dataResult := collection.FindOne(ctx, filter)
	dataResult.Decode(&result)
	return
}

func (m *PhraseTranslate) FindAll(ctx *gin.Context, params requests.PhraseTranslateAllRequests) (result []responses.PhraseTranslateResponse, err error) {
	if helpers.Empty(params.From) || helpers.Empty(params.To) {
		return result, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("phrase_" + strings.ToLower(params.From) + "_" + strings.ToLower(params.To))

	var filter = bson.D{}
	if !helpers.Empty(params.Parents) {
		filter = append(filter, bson.E{"parent", bson.M{"$in": params.Parents}})
	}

	cusor, err := collection.Find(ctx, filter)
	defer cusor.Close(ctx)
	err = cusor.All(ctx, &result)

	if err != nil {
		return result, err
	}
	return
}
