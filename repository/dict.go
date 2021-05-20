package repository

import (
	"context"
	"editorApi/commons"
	"editorApi/init/mgdb"
	"editorApi/requests"
	"editorApi/responses"
	"editorApi/tools/helpers"
	"errors"
	"github.com/jinzhu/copier"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Dict struct {
	Uuid       string       `bson:"uuid" json:"uuid"`              // uuid
	CardId     string       `bson:"card_id" json:"cardId"`         // card_id
	ListOrder  string       `bson:"list_order" json:"listOrder"`   // 排序
	Content    string       `bson:"content" json:"content"`        // 单词
	Images     []Image      `bson:"images" json:"images"`          // 图片
	SoundInfos []SoundInfos `bson:"sound_infos" json:"soundInfos"` // 音标
	IsDel      bool         `bson:"is_del" json:"isDel"`           // 是否删除
	CreatedOn  time.Time    `bson:"created_on" json:"createdOn"`   // 创建时间
	UpdatedOn  time.Time    `bson:"updated_on" json:"updatedOn"`   // 更新时间
}

const (
//DICT_ENG string = "eng" //英语
//DICT_CHI string = "chi" //汉语（简）
//DICT_CHO string = "cho" //汉语（繁）
//DICT_KHM string = "khm" //高棉语
//DICT_MOG string = "mog" //蒙古语
//DICT_JPN string = "jpn" //日语
//DICT_KOR string = "kor" //韩语
//DICT_FRE string = "fre" //法语
//DICT_GER string = "ger" //德语
//DICT_RUS string = "rus" //俄语
//DICT_SPA string = "spa" //西班牙语
//DICT_POR string = "por" //葡萄牙语
//DICT_ARA string = "ara" //阿拉伯语
//DICT_ITA string = "ita" //意大利语
//DICT_DAN string = "dan" //丹麦语
)

type Image struct {
	Url  string `bson:"url" json:"url"`
	Name string `bson:"name" json:"name"`
}

type SoundInfos struct {
	Ct     string `bson:"ct" json:"ct"`         // 类型 en：英； ：美;
	Ps     string `bson:"ps" json:"ps"`         // 音标
	Sound  string `bson:"sound" json:"sound"`   // 声音
	Gender string `bson:"gender" json:"gender"` // male: 男音 female：女音
}

func (m *Dict) DictList(ctx *gin.Context, params requests.DictListRequests) (result []responses.DictResponse, err error) {
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("dict_" + strings.ToLower(params.From))

	var filter = bson.D{}
	if params.OnLine == "Y" {
		filter = append(filter, bson.E{"done", true})
	}

	if params.OnLine == "N" {
		filter = append(filter, bson.E{"done", nil})
	}

	if !helpers.Empty(params.IsDel) {
		filter = append(filter, bson.E{"is_del", params.IsDel})
	}

	// 支持模糊查询
	if !helpers.Empty(params.Content) {
		if params.SearchType == 1 {
			filter = append(filter, bson.E{"content", primitive.Regex{Pattern: params.Content, Options: "i"}})
		} else if params.SearchType == 0 {
			filter = append(filter, bson.E{"content", params.Content})
		}
	}

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

func (m *Dict) Detail(ctx context.Context, params requests.DictDetailRequests) (result responses.DictResponse, err error) {
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("dict_" + strings.ToLower(params.From))
	var filter = bson.D{}
	if !helpers.Empty(params.Uuid) {
		filter = append(filter, bson.E{"uuid", params.Uuid})
	}
	dataResult := collection.FindOne(ctx, filter)
	dataResult.Decode(&result)
	return
}

/**
查询sound_infos数据中数组只有2个数据
*/
func (m *Dict) SoundInfosSize(ctx *gin.Context, params requests.DictDetailRequests) (result []responses.DictResponse, err error) {
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("dict_" + strings.ToLower(params.From))
	var filter = bson.D{}
	filter = append(filter, bson.E{"sound_infos", bson.M{"$size": 2}})

	cusor, err := collection.Find(
		ctx,
		filter,
	)

	defer cusor.Close(ctx)
	err = cusor.All(ctx, &result)
	return
}

func (m *Dict) AddDict(params requests.Dict) (id interface{}, err error) {
	if helpers.Empty(params.From) || helpers.Empty(params.Content) {
		return 0, errors.New("参数不能为空！")
	}

	if helpers.Empty(params.Uuid) {
		params.Uuid = helpers.MD5(strings.TrimRight(params.Content, "."))
	}

	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("dict_" + strings.ToLower(params.From))
	ctx := context.TODO()
	result, _ := m.DictFindOne(params)
	if !helpers.Empty(result.Uuid) {

		/**
		if !helpers.Empty(result.SoundInfos) {
			for _, soundInfo := range result.SoundInfos {
				for key, item := range params.SoundInfos {
					if strings.ToLower(soundInfo.Ct) == strings.ToLower(item.Ct) {
						params.SoundInfos[key].Sound = soundInfo.Sound
						params.SoundInfos[key].Gender = soundInfo.Gender
					}
				}
			}
		}
		**/
		if !helpers.Empty(result.SoundInfos) {
			for _, value := range result.SoundInfos {
				var tmp requests.SoundInfos
				err := copier.Copy(&tmp, &value)
				if err != nil {
					continue
				}
				isExist := false
				for _, item := range params.SoundInfos {
					if item.Sound == tmp.Sound {
						isExist = true
						continue
					}
				}
				if isExist == false {
					params.SoundInfos = append(params.SoundInfos, tmp)
				}
			}
		}

		if !helpers.Empty(result.Images) {
			for _, value := range result.Images {
				var tmp requests.Image
				err := copier.Copy(&tmp, &value)
				if err != nil {
					continue
				}
				isExist := false
				for _, item := range params.Images {
					if item.Url == tmp.Url {
						isExist = true
						continue
					}
				}
				if isExist == false {
					params.Images = append(params.Images, tmp)
				}
			}
		}

		var filter = bson.D{}
		if !helpers.Empty(params.Uuid) {
			filter = append(filter, bson.E{"uuid", params.Uuid})
		}

		var update = bson.M{
			"content":     params.Content,
			"list_order":  strings.ToLower(params.Content),
			"images":      params.Images,
			"sound_infos": params.SoundInfos,
			"is_del":      params.IsDel,
			"updated_on":  time.Now(),
			"created_on":  time.Now(),
			"done":        false,
		}

		updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": update})
		if err != nil {
			return 0, err
		}
		id = updateResult.UpsertedID
	} else {
		var add = bson.M{
			"uuid":        params.Uuid,
			"list_order":  strings.ToLower(params.Content),
			"content":     params.Content,
			"images":      params.Images,
			"sound_infos": params.SoundInfos,
			"is_del":      false,
			"updated_on":  time.Now(),
			"created_on":  time.Now(),
			"done":        false,
		}

		insertResult, err := collection.InsertOne(ctx, add)
		if err != nil {
			return 0, err
		}
		id = insertResult.InsertedID
	}
	return
}

func (m *Dict) DictDel(ctx *gin.Context, params requests.DictDelRequests) {

	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("dict_" + strings.ToLower(params.From))

	var filter = bson.M{
		"uuid": bson.M{
			"$in": params.Uuids,
		},
	}

	var update = bson.M{
		"is_del": true,
		"done":   false,
	}

	collection.UpdateMany(ctx, filter, bson.M{"$set": update})

	return
}
func (m *Dict) DictUpdate(ctx *gin.Context, params requests.DictUpdateRequests) (id interface{}, err error) {
	if helpers.Empty(params.From) {
		return 0, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("dict_" + strings.ToLower(params.From))

	var filter = bson.D{}
	if !helpers.Empty(params.Uuid) {
		filter = append(filter, bson.E{"uuid", params.Uuid})
	}

	var update = bson.M{
		"list_order":  strings.ToLower(params.Content),
		"content":     params.Content,
		"images":      params.Images,
		"sound_infos": params.SoundInfos,
		"updated_on":  time.Now(),
		"created_on":  time.Now(),
		"done":        false,
	}

	updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return 0, err
	}
	id = updateResult.UpsertedID

	//更新翻译
	var paramsDictTranslate requests.DictTranslate
	paramsDictTranslate.From = params.From
	paramsDictTranslate.To = params.To
	paramsDictTranslate.Parent = params.Uuid
	paramsDictTranslate.Expansion = params.DictTranslate.Expansion
	paramsDictTranslate.ContentTr = params.DictTranslate.ContentTr
	paramsDictTranslate.Synonym = params.DictTranslate.Synonym
	paramsDictTranslate.Homonyms = params.DictTranslate.Homonyms
	paramsDictTranslate.Tags = params.DictTranslate.Tags

	var dictTranslateModel DictTranslate
	id, err = dictTranslateModel.Update(ctx, paramsDictTranslate)
	return
}

func (m *Dict) DictFindOne(params requests.Dict) (result Dict, err error) {
	if helpers.Empty(params.From) {
		return result, errors.New("参数不能为空！")
	}

	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("dict_" + strings.ToLower(params.From))
	ctx := context.TODO()

	var filter = bson.D{}
	if !helpers.Empty(params.Uuid) {
		filter = append(filter, bson.E{"uuid", params.Uuid})
	}

	if !helpers.Empty(params.ListOrder) {
		filter = append(filter, bson.E{"list_order", params.ListOrder})
	}

	singleResult := collection.FindOne(ctx, filter)
	err = singleResult.Decode(&result)
	return
}

func (m *Dict) DictFindOneByUuid(from, uuid string) (result Dict, err error) {
	if helpers.Empty(uuid) || helpers.Empty(from) {
		return result, errors.New("参数不能为空！")
	}

	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("dict_" + strings.ToLower(from))
	ctx := context.TODO()

	var filter = bson.D{}
	filter = append(filter, bson.E{"uuid", uuid})

	singleResult := collection.FindOne(ctx, filter)
	err = singleResult.Decode(&result)
	return
}

func (m *Dict) DictFindOneByCardId(from, card_id string) (result Dict, err error) {
	if helpers.Empty(card_id) || helpers.Empty(from) {
		return result, errors.New("参数不能为空！")
	}

	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("dict_" + strings.ToLower(from))
	ctx := context.TODO()

	var filter = bson.D{}
	filter = append(filter, bson.E{"card_id", card_id})

	singleResult := collection.FindOne(ctx, filter)
	err = singleResult.Decode(&result)
	return
}

func (m *Dict) DictUpdateCardId(ctx *gin.Context, params requests.DictCardId) (id interface{}, err error) {
	if helpers.Empty(params.From) {
		return 0, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("dict_" + strings.ToLower(params.From))
	var filter = bson.D{}
	if !helpers.Empty(params.CardId) {
		filter = append(filter, bson.E{"card_id", params.CardId})
	}

	var update = bson.M{
		"uuid":        params.Uuid,
		"content":     params.Content,
		"sound_infos": params.SoundInfos,
	}

	updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return 0, err
	}
	id = updateResult.UpsertedID

	return
}

func (m *Dict) DictAddCardId(ctx *gin.Context, params requests.DictCardId) (id interface{}, err error) {
	if helpers.Empty(params.From) {
		return 0, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("dict_" + strings.ToLower(params.From))
	var dict = bson.M{
		"card_id":     params.CardId,
		"uuid":        params.Uuid,
		"list_order":  strings.ToLower(params.Content),
		"content":     params.Content,
		"sound_infos": params.SoundInfos,
	}
	updateResult, err := collection.InsertOne(ctx, dict)
	if err != nil {
		return 0, err
	}
	id = updateResult.InsertedID
	return
}

func (m *Dict) FindAll(ctx context.Context, params requests.DictFindAll) (result []responses.DictResponse, err error) {
	if helpers.Empty(params.From) {
		return result, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("dict_" + strings.ToLower(params.From))

	var filter = bson.D{}
	if !helpers.Empty(params.CreatedOn) {
		filter = append(filter, bson.E{"created_on", bson.M{"$lt": params.CreatedOn}})
	}

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

func (m *Dict) DictUpdateSound(ctx *gin.Context, params requests.DictSound) (id interface{}, err error) {
	if helpers.Empty(params.From) {
		return 0, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("dict_" + strings.ToLower(params.From))
	var filter = bson.D{}
	if !helpers.Empty(params.Uuid) {
		filter = append(filter, bson.E{"uuid", params.Uuid})
	}

	var update = bson.M{
		"sound_infos": params.SoundInfos,
	}

	updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return 0, err
	}
	id = updateResult.UpsertedID

	return
}
