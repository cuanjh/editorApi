package service

import (
	"context"
	"editorApi/repository"
	"editorApi/requests"
	"editorApi/responses"
)

type SentenceTranslateService struct {
}

func (s *SentenceTranslateService) AddSentenceTranslate(ctx context.Context, params requests.SentenceTranslate) (id interface{}, err error) {
	var sentence repository.SentenceTranslate
	return sentence.AddSentenceTranslate(params)
}

func (s *SentenceTranslateService) FindOne(ctx context.Context, params requests.SentenceTranslate) (result responses.SentenceTranslate, err error) {
	var sentence repository.SentenceTranslate
	return sentence.FindOne(params)
}
