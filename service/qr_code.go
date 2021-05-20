package service

import (
	"editorApi/repository"
	"editorApi/requests"
	"editorApi/responses"
	"github.com/gin-gonic/gin"
)

type QRcodeService struct {
}

func (s *QRcodeService) QRcodeAdd(ctx *gin.Context, params requests.QRcodeAddRequests) (inserted_id interface{}, err error) {
	var teacher = repository.QRcode{}
	return teacher.Add(ctx, params)
}

func (s *QRcodeService) QRcodeDelete(ctx *gin.Context, params requests.QRcodeDeleteRequests) (upserted_id interface{}, err error) {
	var teacher = repository.QRcode{}
	return teacher.Delete(ctx, params)
}

func (s *QRcodeService) QRcodeUpdate(ctx *gin.Context, params requests.QRcodeUpdateRequests) (upserted_id interface{}, err error) {
	var teacher = repository.QRcode{}
	return teacher.Update(ctx, params)
}

func (s *QRcodeService) QRcodeList(ctx *gin.Context, params requests.QRcodeListRequests) (result []responses.QRcodeResponses, err error) {
	var teacher = repository.QRcode{}
	return teacher.List(ctx, params)
}

func (s *QRcodeService) QRcodeDetails(ctx *gin.Context, params requests.QRcodeDetailsRequests) (result responses.QRcodeResponses, err error) {
	var teacher = repository.QRcode{}
	return teacher.Details(ctx, params)
}
