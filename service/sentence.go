package service

import (
	"context"
	"editorApi/repository"
	"editorApi/requests"
	"editorApi/responses"
	"editorApi/tools/helpers"

	"github.com/gin-gonic/gin"
)

type SentenceService struct {
}

func (s *SentenceService) List(ctx *gin.Context, params requests.SentenceSearchRequests) (result []responses.Sentence, err error) {
	var sentence repository.Sentence
	result, err = sentence.SentenceList(ctx, params)
	if err != nil {
		return nil, err
	}
	return
}

func (s *SentenceService) FindOne(ctx context.Context, params requests.Sentence) (result responses.Sentence, err error) {
	var sentence repository.Sentence
	return sentence.FindOne(params)
}

func (s *SentenceService) Detail(ctx *gin.Context, params requests.SentenceDetail) (result responses.Sentence, err error) {
	var sentence repository.Sentence
	result, err = sentence.Detail(ctx, params)

	var requestsSentenceTranslate requests.SentenceTranslate
	requestsSentenceTranslate.Parent = params.Uuid
	requestsSentenceTranslate.From = params.From
	requestsSentenceTranslate.To = params.To

	var sentenceTranslateModel repository.SentenceTranslate
	res, err := sentenceTranslateModel.FindOne(requestsSentenceTranslate)

	if err == nil {
		result.SentenceTr = res
	}
	return
}

func (s *SentenceService) Update(ctx *gin.Context, params requests.SentenceUpdate) (id interface{}, err error) {
	var sentence repository.Sentence
	id, err = sentence.Update(ctx, params)
	return
}

func (s *SentenceService) UpdateSoundInfos(ctx *gin.Context, params requests.SentenceUpdate) (id interface{}, err error) {
	var sentence repository.Sentence

	url := "https://fanyi.baidu.com/gettts?lan=en&text=The%20charity%20does%20a%20lot%20of%20good.&spd=3&source=web"
	filename := helpers.MD5("")
	dir := "dict/sentence/eng/"
	Download(url, "/opt/data/goPro/editorAPILinux/data/"+dir, filename+".mp3")

	id, err = sentence.Update(ctx, params)
	return
}

func (s *SentenceService) DeleteOne(ctx *gin.Context, params requests.SentenceDelete) (deleteResult interface{}, err error) {
	var sentence repository.Sentence
	sentence.DeleteOne(ctx, params)
	return
}

func (s *SentenceService) AddSentence(ctx context.Context, params requests.Sentence) (id interface{}, err error) {
	var sentence repository.Sentence
	return sentence.AddSentence(params)
}


func (s *SentenceService) SentenceAddCardId(ctx context.Context, params requests.SentenceCardId) (id interface{}, err error) {
	var sentence repository.Sentence
	return sentence.SentenceAddCardId(ctx, params)
}
