package repository

import (
	"context"
	"editorApi/init/mgdb"
	"editorApi/requests"
	"editorApi/responses"
	"editorApi/tools/helpers"
	"errors"
	array "github.com/chenhg5/collection"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"time"
)

// 短语 一对多的形式
type Phrase struct {
	Content   string    `bson:"content" json:"content"`      //短语单词
	Uuid      string    `bson:"uuid" json:"uuid"`            //uuid
	DictUuid  []string  `bson:"dict_uuid" json:"dict_uuid"`  //多个字典UUID
	IsDel     bool      `bson:"is_del" json:"isDel"`         //是否删除
	CreatedOn time.Time `bson:"created_on" json:"createdOn"` //创建时间
	UpdatedOn time.Time `bson:"updated_on" json:"updatedOn"` //更新时间
}

func (m *Phrase) AddPhrase(params requests.Phrase) (id interface{}, err error) {
	if helpers.Empty(params.From) {
		return nil, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("phrase_" + strings.ToLower(params.From))

	ctx := context.TODO()
	result, _ := m.FindOne(params)
	if !helpers.Empty(result.Uuid) {
		var filter = bson.D{}
		if !helpers.Empty(params.Uuid) {
			filter = append(filter, bson.E{"uuid", params.Uuid})
		}

		for _, item := range result.DictUuid {
			params.DictUuid = append(params.DictUuid, item)
		}

		params.DictUuid = array.Collect(params.DictUuid).Unique().ToStringArray()

		// 更新时间
		var update = bson.M{
			"content":    params.Content,
			"dict_uuid":  params.DictUuid,
			"is_del":     false,
			"updated_on": time.Now(),
		}

		updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": update})
		if err != nil {
			return 0, err
		}
		id = updateResult.UpsertedID
	} else {
		var add = bson.M{
			"content":    params.Content,
			"uuid":       params.Uuid,
			"dict_uuid":  params.DictUuid,
			"is_del":     false,
			"created_on": time.Now(),
			"updated_on": time.Now(),
		}

		insertResult, err := collection.InsertOne(ctx, add)
		if err != nil {
			return 0, err
		}
		id = insertResult.InsertedID
	}
	return
}

func (m *Phrase) FindOne(params requests.Phrase) (result Phrase, err error) {
	if helpers.Empty(params.From) {
		return result, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("phrase_" + strings.ToLower(params.From))
	ctx := context.TODO()

	var filter = bson.D{}
	if !helpers.Empty(params.Uuid) {
		filter = append(filter, bson.E{"uuid", params.Uuid})
	}

	singleResult := collection.FindOne(ctx, filter)
	err = singleResult.Decode(&result)
	return
}

func (m *Phrase) FindAll(ctx *gin.Context, params requests.PhraseAllRequests) (result []responses.PhraseResponse, err error) {
	if helpers.Empty(params.From) {
		return result, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("phrase_" + strings.ToLower(params.From))

	var filter = bson.D{}
	if !helpers.Empty(params.DictUuid) {
		filter = append(filter, bson.E{"dict_uuid", bson.M{"$in": params.DictUuid}})
	}

	cusor, err := collection.Find(ctx, filter)
	defer cusor.Close(ctx)
	err = cusor.All(ctx, &result)

	if err != nil {
		return result, err
	}

	var parents []string
	for _, phrase := range result {
		parents = append(parents, phrase.Uuid)
	}

	var phraseTranslateAllRequests requests.PhraseTranslateAllRequests
	phraseTranslateAllRequests.From = params.From
	phraseTranslateAllRequests.To = params.To
	phraseTranslateAllRequests.Parents = parents

	var phraseTranslateModel PhraseTranslate
	phraseTranslates, err := phraseTranslateModel.FindAll(ctx, phraseTranslateAllRequests)
	if err != nil {
		return result, err
	}

	for key, phrase := range result {
		for _, phraseTranslate := range phraseTranslates {
			if phraseTranslate.Parent == phrase.Uuid {
				result[key].ContentTr = phraseTranslate.ContentTr
			}
		}
	}

	return
}
