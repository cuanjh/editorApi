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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 短语 一对多的形式
type Sentence struct {
	Uuid        string               `bson:"uuid" json:"uuid"`                //uuid
	CardId      string               `bson:"card_id" json:"cardId"`           // card_id
	Mold        int                  `bson:"mold" json:"mold"`                //类型 1：口语；2：书面语；
	Sentence    string               `bson:"sentence" json:"sentence"`        //短语单词
	Image       []string             `bson:"image" json:"image"`              //图片
	SoundInfos  []SentenceSoundInfos `bson:"sound_infos" json:"soundInfos"`   //声音
	CourseInfos []CourseInfos        `bson:"course_infos" json:"courseInfos"` //课程
	Source      string               `bson:"source" json:"source"`            //来源
	IsDel       bool                 `bson:"is_del" json:"isDel"`             //是否删除
	CreatedOn   time.Time            `bson:"created_on" json:"createdOn"`     //创建时间
	UpdatedOn   time.Time            `bson:"updated_on" json:"updatedOn"`     //更新时间
}

type SentenceSoundInfos struct {
	Sound  string `json:"sound"`  //声音
	Gender string `json:"gender"` //male: 男音 female：女音
}

type CourseInfos struct {
	Uuid        string               `bson:"uuid" json:"uuid"` //uuid
	CourseCode  string               `bson:"course_code" json:"courseCode"`
	ChapterCode string               `bson:"chapter_code" json:"chapterCode"`
	Image       []string             `bson:"image" json:"image"`            //图片
	SoundInfos  []SentenceSoundInfos `bson:"sound_infos" json:"soundInfos"` //声音
}

func (m *Sentence) AddSentence(params requests.Sentence) (id interface{}, err error) {
	if helpers.Empty(params.From) || helpers.Empty(params.Sentence) {
		return id, errors.New("参数不能为空！")
	}

	if helpers.Empty(params.Uuid) {
		params.Uuid = helpers.MD5(params.Sentence)
	}

	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("sentence_" + strings.ToLower(params.From))
	ctx := context.TODO()
	result, _ := m.FindOne(params)
	if !helpers.Empty(result.Uuid) {
		var filter = bson.D{}
		if !helpers.Empty(params.Uuid) {
			filter = append(filter, bson.E{"uuid", params.Uuid})
		}

		if !helpers.Empty(result.SoundInfos) {
			for _, value := range result.SoundInfos {
				var tmp requests.SentenceSoundInfos
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

		if !helpers.Empty(result.Image) {
			for _, value := range result.Image {
				isExist := false
				for _, item := range params.Image {
					if item == value {
						isExist = true
						continue
					}
				}
				if isExist == false {
					params.Image = append(params.Image, value)
				}
			}
		}

		var update = bson.M{
			"card_id":      params.CardId,
			"mold":         params.Mold,
			"sentence":     params.Sentence,
			"image":        params.Image,
			"sound_infos":  params.SoundInfos,
			"course_infos": params.CourseInfos,
			"source":       params.Source,
			"is_del":       false,
			"updated_on":   time.Now(),
		}
		// 更新时间
		updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": update})
		if err != nil {
			return 0, err
		}
		id = updateResult.UpsertedID
	} else {
		var add = bson.M{
			"card_id":      params.CardId,
			"uuid":         params.Uuid,
			"mold":         params.Mold,
			"sentence":     params.Sentence,
			"image":        params.Image,
			"sound_infos":  params.SoundInfos,
			"course_infos": params.CourseInfos,
			"source":       params.Source,
			"is_del":       false,
			"updated_on":   time.Now(),
			"created_on":   time.Now(),
		}
		insertResult, err := collection.InsertOne(ctx, add)
		if err != nil {
			return 0, err
		}
		id = insertResult.InsertedID
	}
	return
}

func (m *Sentence) FindOne(params requests.Sentence) (result responses.Sentence, err error) {
	if helpers.Empty(params.From) {
		return result, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("sentence_" + strings.ToLower(params.From))
	ctx := context.TODO()

	var filter = bson.D{}
	if !helpers.Empty(params.Uuid) {
		filter = append(filter, bson.E{"uuid", params.Uuid})
	}

	singleResult := collection.FindOne(ctx, filter)
	err = singleResult.Decode(&result)
	return
}

func (m *Sentence) SentenceList(ctx *gin.Context, params requests.SentenceSearchRequests) (result []responses.Sentence, err error) {
	if helpers.Empty(params.From) {
		return result, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("sentence_" + strings.ToLower(params.From))

	var filter = bson.M{}

	// 支持模糊查询
	if !helpers.Empty(params.Sentence) {

		if params.SearchType == 1 {
			filter["$text"] = bson.M{
				"$search": params.Sentence,
			}
		} else if params.SearchType == 0 {
			filter["sentence"] = params.Sentence
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

func (m *Sentence) Detail(ctx context.Context, params requests.SentenceDetail) (result responses.Sentence, err error) {
	if helpers.Empty(params.From) {
		return result, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("sentence_" + strings.ToLower(params.From))
	var filter = bson.D{}
	if !helpers.Empty(params.Uuid) {
		filter = append(filter, bson.E{"uuid", params.Uuid})
	}
	dataResult := collection.FindOne(ctx, filter)
	dataResult.Decode(&result)
	return
}

func (m *Sentence) Update(ctx context.Context, params requests.SentenceUpdate) (id interface{}, err error) {
	if helpers.Empty(params.From) {
		return 0, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("sentence_" + strings.ToLower(params.From))
	var filter = bson.D{}
	if !helpers.Empty(params.Uuid) {
		filter = append(filter, bson.E{"uuid", params.Uuid})
	}

	var update = bson.M{
		"uuid":         params.Uuid,
		"mold":         params.Mold,
		"sentence":     params.Sentence,
		"image":        params.Image,
		"sound_infos":  params.SoundInfos,
		"source":       params.Source,
		"course_infos": params.CourseInfos,
		"done":         false,
	}

	updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return 0, err
	}
	id = updateResult.UpsertedID

	var paramsSentenceTranslate requests.SentenceTranslate
	paramsSentenceTranslate.From = params.From
	paramsSentenceTranslate.To = params.To
	paramsSentenceTranslate.Parent = params.SentenceTr.Parent
	paramsSentenceTranslate.ContentTr = params.SentenceTr.ContentTr
	paramsSentenceTranslate.Tags = params.SentenceTr.Tags

	var sentenceTranslateModel SentenceTranslate
	id, err = sentenceTranslateModel.Update(ctx, paramsSentenceTranslate)
	return
}

func (m *Sentence) UpdateSoundInfos(ctx context.Context, params requests.SentenceUpdate) (id interface{}, err error) {
	if helpers.Empty(params.From) {
		return 0, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("sentence_" + strings.ToLower(params.From))
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

func (m *Sentence) FindAll(ctx context.Context, params requests.SentenceFindAll) (result []responses.SentenceFindAll, err error) {
	if helpers.Empty(params.From) {
		return result, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("sentence_" + strings.ToLower(params.From))

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

func (m *Sentence) UpdateImage(ctx context.Context, params requests.SentenceUpdate) (id interface{}, err error) {
	if helpers.Empty(params.From) {
		return 0, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("sentence_" + strings.ToLower(params.From))
	var filter = bson.D{}
	if !helpers.Empty(params.Uuid) {
		filter = append(filter, bson.E{"uuid", params.Uuid})
	}

	var update = bson.M{
		"sound": params.SoundInfos,
	}

	updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return 0, err
	}
	id = updateResult.UpsertedID
	return
}

func (m *Sentence) DeleteOne(ctx *gin.Context, params requests.SentenceDelete) (deleteResult interface{}, err error) {
	if helpers.Empty(params.From) {
		return 0, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("sentence_" + strings.ToLower(params.From))
	var filter = bson.D{}
	if !helpers.Empty(params.Uuid) {
		filter = append(filter, bson.E{"uuid", params.Uuid})
	}
	deleteResult, err = collection.DeleteOne(ctx, filter)
	return
}

func (m *Sentence) SentenceAddCardId(ctx context.Context, params requests.SentenceCardId) (id interface{}, err error) {
	if helpers.Empty(params.From) || helpers.Empty(params.Uuid) {
		return 0, errors.New("参数不能为空！")
	}
	collection := mgdb.MongoClient.Database(mgdb.DbDict).Collection("sentence_" + strings.ToLower(params.From))
	var filter = bson.D{}
	if !helpers.Empty(params.Uuid) {
		filter = append(filter, bson.E{"uuid", params.Uuid})
	}

	var update = bson.M{
		"card_id": params.CardId,
	}

	updateResult, err := collection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return 0, err
	}
	id = updateResult.UpsertedID
	return
}
