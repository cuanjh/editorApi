package repository

import (
	"context"
	"editorApi/init/mgdb"
	"editorApi/requests"
	"editorApi/responses"
	"editorApi/tools/helpers"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

type DictTranslate struct {
	Parent    string      `bson:"parent" json:"parent"`        // uuid
	Expansion string      `bson:"expansion" json:"expansion"`  // 拓展
	ContentTr []ContentTr `bson:"content_tr" json:"contentTr"` // 词义
	Synonym   []Synonym   `bson:"synonym" json:"synonym"`      // 近义词
	Homonyms  []Homonyms  `bson:"homonyms" json:"homonyms"`    // 同词根
	Tags      []Tag       `bson:"tags" json:"tags"`            // 多语言标签
}

type Tag struct {
	Key  string `bson:"key" json:"key"`
	Name string `bson:"name"  json:"name"`
}

type ContentTr struct {
	Cx      string `json:"cx"`
	Content string `json:"content"`
}

type WordAttr struct {
	Content   string `json:"content"`
	ContentTr string `json:"contentTr"`
}

type Homonyms struct {
	Cx    string     `json:"cx"`
	Attrs []WordAttr `json:"attrs"`
}

type Synonym struct {
	Cx        string `json:"cx"`
	Content   string `json:"content"`
	ContentTr string `json:"contentTr"`
}

func (m *DictTranslate) AddDictTranslate(params requests.DictTranslate) (id interface{}, err error) {
	if helpers.Empty(params.From) || helpers.Empty(params.To) {
		return id, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("dict_" + strings.ToLower(params.From) + "_" + strings.ToLower(params.To))

	result, _ := m.FindOne(params)
	ctx := context.TODO()
	if !helpers.Empty(result.Parent) {
		var filter = bson.D{}
		if !helpers.Empty(params.Parent) {
			filter = append(filter, bson.E{"parent", params.Parent})
		}
		// 更新时间
		var update = bson.M{
			"expansion":  params.Expansion,
			"content_tr": params.ContentTr,
			"synonym":    params.Synonym,
			"homonyms":   params.Homonyms,
			"tags":       params.Tags,
		}
		updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": update})
		if err != nil {
			return 0, err
		}
		id = updateResult.UpsertedID
	} else {
		var add = bson.M{
			"parent":     params.Parent,
			"expansion":  params.Expansion,
			"content_tr": params.ContentTr,
			"synonym":    params.Synonym,
			"homonyms":   params.Homonyms,
			"tags":       params.Tags,
		}

		insertResult, err := collection.InsertOne(ctx, add)
		if err != nil {
			return 0, err
		}
		id = insertResult.InsertedID
	}
	return
}

func (m *DictTranslate) FindOne(params requests.DictTranslate) (result responses.DictTranslate, err error) {
	if helpers.Empty(params.From) || helpers.Empty(params.To) {
		return result, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("dict_" + strings.ToLower(params.From) + "_" + strings.ToLower(params.To))
	ctx := context.TODO()

	var filter = bson.D{}
	if !helpers.Empty(params.Parent) {
		filter = append(filter, bson.E{"parent", params.Parent})
	}

	singleResult := collection.FindOne(ctx, filter)
	err = singleResult.Decode(&result)
	return
}

func (m *DictTranslate) Detail(ctx *gin.Context, params requests.DictTranslateDetailRequests) (result responses.DictTranslate, err error) {
	if helpers.Empty(params.From) || helpers.Empty(params.To) {
		return result, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("dict_" + strings.ToLower(params.From) + "_" + strings.ToLower(params.To))
	var filter = bson.D{}
	if !helpers.Empty(params.Parent) {
		filter = append(filter, bson.E{"parent", params.Parent})
	}
	dataResult := collection.FindOne(ctx, filter)
	dataResult.Decode(&result)
	return
}

func (m *DictTranslate) Update(ctx context.Context, params requests.DictTranslate) (id interface{}, err error) {
	if helpers.Empty(params.From) || helpers.Empty(params.To) {
		return 0, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("dict_" + strings.ToLower(params.From) + "_" + strings.ToLower(params.To))
	var filter = bson.D{}
	if !helpers.Empty(params.Parent) {
		filter = append(filter, bson.E{"parent", params.Parent})
	}

	var update = bson.M{
		"expansion":  params.Expansion,
		"content_tr": params.ContentTr,
		"synonym":    params.Synonym,
		"homonyms":   params.Homonyms,
		"tags":       params.Tags,
	}

	updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return 0, err
	}
	id = updateResult.UpsertedID

	return
}

func (m *DictTranslate) AddTags(ctx *gin.Context, params requests.AddTags) (id interface{}, err error) {
	if helpers.Empty(params.From) || helpers.Empty(params.To) {
		return 0, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("dict_" + strings.ToLower(params.From) + "_" + strings.ToLower(params.To))

	// 老数据处理
	var param requests.DictTranslate
	param.Parent = params.Parent
	param.From = params.From
	param.To = params.To
	result, _ := m.FindOne(param)
	if !helpers.Empty(result.Tags) {
		for _, item := range result.Tags {
			var in bool
			for _, tag := range params.Tags {
				if tag.Key == item.Key {
					in = true
				}
			}
			if in == false {
				var tmp requests.Tag
				copier.Copy(tmp, item)
				params.Tags = append(params.Tags, tmp)
			}
		}
	}

	var dict = bson.M{
		"tags": params.Tags,
	}
	var filter = bson.D{}
	if !helpers.Empty(params.Parent) {
		filter = append(filter, bson.E{"parent", params.Parent})
	}
	updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": dict})
	if err != nil {
		return 0, err
	}
	id = updateResult.UpsertedID
	return
}
