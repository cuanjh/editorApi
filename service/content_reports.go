package service

import (
	"editorApi/repository"
	"editorApi/requests"
	"editorApi/responses"
	"github.com/gin-gonic/gin"
)

type IContentReports interface {
	Create(ctx *gin.Context, params requests.ContentReportsCreateRequests) (inserted_id interface{}, err error)
	Find(ctx *gin.Context, params requests.ContentReportsFindRequests) (result []responses.ContentReportsResponses, err error)
	FindOne(ctx *gin.Context, params requests.ContentReportsFindOneRequests) (result responses.ContentReportsResponses, err error)
	List(ctx *gin.Context, params requests.ContentReportsListRequests) (result []responses.ContentReportsResponses, err error)
	Update(ctx *gin.Context, params requests.ContentReportsUpdateRequests) (upserted_id interface{}, err error)
	Delete(ctx *gin.Context, params requests.ContentReportsDeleteRequests) (deleteResult interface{}, err error)
}

type ContentReportsService struct {
}

// 添加数据
func (s *ContentReportsService) Create(ctx *gin.Context, params requests.ContentReportsCreateRequests) (inserted_id interface{}, err error) {
	var model repository.ContentReports
	inserted_id, err = model.Create(ctx, params)
	return
}

// 查询多条数据
func (s *ContentReportsService) Find(ctx *gin.Context, params requests.ContentReportsFindRequests) (result []responses.ContentReportsResponses, err error) {
	var model repository.ContentReports
	result, err = model.Find(ctx, params)
	return
}

//
func (s *ContentReportsService) FindOne(ctx *gin.Context, params requests.ContentReportsFindOneRequests) (result responses.ContentReportsResponses, err error) {
	var model repository.ContentReports
	result, err = model.FindOne(ctx, params)
	return
}

// 列表数据
func (s *ContentReportsService) List(ctx *gin.Context, params requests.ContentReportsListRequests) (result []responses.ContentReportsResponses, err error) {
	var model repository.ContentReports
	result, err = model.List(ctx, params)
	return
}

func (s *ContentReportsService) Update(ctx *gin.Context, params requests.ContentReportsUpdateRequests) (upserted_id interface{}, err error) {
	var model repository.ContentReports
	upserted_id, err = model.Update(ctx, params)
	return
}

func (s *ContentReportsService) Delete(ctx *gin.Context, params requests.ContentReportsDeleteRequests) (deleteResult interface{}, err error) {
	var model repository.ContentReports
	deleteResult, err = model.Delete(ctx, params)
	return
}
