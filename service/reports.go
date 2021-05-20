package service

import (
	"editorApi/repository"
	"editorApi/requests"
	"editorApi/responses"
	"github.com/gin-gonic/gin"
)

type IReports interface {
	Create(ctx *gin.Context, params requests.ReportsCreateRequests) (inserted_id interface{}, err error)
	Find(ctx *gin.Context, params requests.ReportsFindRequests) (result []responses.ReportsResponses, err error)
	FindOne(ctx *gin.Context, params requests.ReportsFindOneRequests) (result responses.ReportsResponses, err error)
	List(ctx *gin.Context, params requests.ReportsListRequests) (result []responses.ReportsResponses, err error)
	Update(ctx *gin.Context, params requests.ReportsUpdateRequests) (upserted_id interface{}, err error)
	Delete(ctx *gin.Context, params requests.ReportsDeleteRequests) (deleteResult interface{}, err error)
}

type ReportsService struct {
}

// 添加数据
func (s *ReportsService) Create(ctx *gin.Context, params requests.ReportsCreateRequests) (inserted_id interface{}, err error) {
	var model repository.Reports
	inserted_id, err = model.Create(ctx, params)
	return
}

// 查询多条数据
func (s *ReportsService) Find(ctx *gin.Context, params requests.ReportsFindRequests) (result []responses.ReportsResponses, err error) {
	var model repository.Reports
	result, err = model.Find(ctx, params)
	return
}

//
func (s *ReportsService) FindOne(ctx *gin.Context, params requests.ReportsFindOneRequests) (result responses.ReportsResponses, err error) {
	var model repository.Reports
	result, err = model.FindOne(ctx, params)
	return
}

// 列表数据
func (s *ReportsService) List(ctx *gin.Context, params requests.ReportsListRequests) (result []responses.ReportsResponses, err error) {
	var model repository.Reports
	result, err = model.List(ctx, params)
	return
}

func (s *ReportsService) Update(ctx *gin.Context, params requests.ReportsUpdateRequests) (upserted_id interface{}, err error) {
	var model repository.Reports
	upserted_id, err = model.Update(ctx, params)
	return
}

func (s *ReportsService) Delete(ctx *gin.Context, params requests.ReportsDeleteRequests) (deleteResult interface{}, err error) {
	var model repository.Reports
	deleteResult, err = model.Delete(ctx, params)
	return
}
