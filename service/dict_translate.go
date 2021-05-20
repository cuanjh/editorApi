package service

import (
	"editorApi/repository"
	"editorApi/requests"
	"editorApi/responses"
	"github.com/gin-gonic/gin"
)

type DictTranslateService struct {
}

func (s *DictTranslateService) AddSentenceTranslate(ctx *gin.Context, params requests.DictTranslate) (id interface{}, err error) {
	var sentence repository.DictTranslate
	return sentence.AddDictTranslate(params)
}

func (s *DictTranslateService) FindOne(ctx *gin.Context, params requests.DictTranslate) (result responses.DictTranslate, err error) {
	var sentence repository.DictTranslate
	return sentence.FindOne(params)
}
