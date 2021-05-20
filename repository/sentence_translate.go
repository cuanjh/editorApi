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

type SentenceTranslate struct {
	Parent      string        `bson:"parent"  json:"parent"`
	ContentTr   string        `bson:"content_tr"  json:"contentTr"`
	SentenceTag []SentenceTag `bson:"sentence_tag" json:"sentenceTag"` // 多语言标签
}

type SentenceTag struct {
	Key  string `bson:"key" json:"key"`
	Name string `bson:"name"  json:"name"`
}

func (m *SentenceTranslate) AddSentenceTranslate(params requests.SentenceTranslate) (id interface{}, err error) {
	if helpers.Empty(params.From) || helpers.Empty(params.To) {
		return id, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("sentence_" + strings.ToLower(params.From) + "_" + strings.ToLower(params.To))

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

func (m *SentenceTranslate) FindOne(params requests.SentenceTranslate) (result responses.SentenceTranslate, err error) {
	if helpers.Empty(params.From) || helpers.Empty(params.To) {
		return result, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("sentence_" + strings.ToLower(params.From) + "_" + strings.ToLower(params.To))
	ctx := context.TODO()

	var filter = bson.D{}
	if !helpers.Empty(params.Parent) {
		filter = append(filter, bson.E{"parent", params.Parent})
	}

	singleResult := collection.FindOne(ctx, filter)
	err = singleResult.Decode(&result)
	return
}

func (m *SentenceTranslate) Update(ctx context.Context, params requests.SentenceTranslate) (id interface{}, err error) {
	if helpers.Empty(params.From) || helpers.Empty(params.To) {
		return 0, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("sentence_" + strings.ToLower(params.From) + "_" + strings.ToLower(params.To))
	var filter = bson.D{}
	if !helpers.Empty(params.Parent) {
		filter = append(filter, bson.E{"parent", params.Parent})
	}

	var update = bson.M{
		"content_tr": params.ContentTr,
	}

	updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return 0, err
	}
	id = updateResult.UpsertedID

	return
}

func (m *SentenceTranslate) DeleteOne(ctx *gin.Context, params requests.SentenceTranslate) (deleteResult interface{}, err error) {
	if helpers.Empty(params.From) || helpers.Empty(params.To) {
		return 0, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("sentence_" + strings.ToLower(params.From) + "_" + strings.ToLower(params.To))
	var filter = bson.D{}
	if !helpers.Empty(params.Parent) {
		filter = append(filter, bson.E{"parent", params.Parent})
	}
	deleteResult, err = collection.DeleteOne(ctx, filter)
	return
}
