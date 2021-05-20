package service

import (
	"editorApi/repository"
	"editorApi/requests"
	"editorApi/responses"
	"github.com/gin-gonic/gin"
)

type IActors interface {
	Create(ctx *gin.Context, params requests.ActorsCreateRequests) (inserted_id interface{}, err error)
	Find(ctx *gin.Context, params requests.ActorsFindRequests) (result []responses.ActorsResponses, err error)
	FindOne(ctx *gin.Context, params requests.ActorsFindOneRequests) (result responses.ActorsResponses, err error)
	List(ctx *gin.Context, params requests.ActorsListRequests) (result []responses.ActorsResponses, err error)
	Update(ctx *gin.Context, params requests.ActorsUpdateRequests) (upserted_id interface{}, err error)
	Delete(ctx *gin.Context, params requests.ActorsDeleteRequests) (deleteResult interface{}, err error)
}

type ActorsService struct {
}

// 添加数据
func (s *ActorsService) Create(ctx *gin.Context, params requests.ActorsCreateRequests) (inserted_id interface{}, err error) {
	var model repository.Actors
	inserted_id, err = model.Create(ctx, params)
	return
}

// 查询多条数据
func (s *ActorsService) Find(ctx *gin.Context, params requests.ActorsFindRequests) (result []responses.ActorsResponses, err error) {
	var model repository.Actors
	result, err = model.Find(ctx, params)
	return
}

//
func (s *ActorsService) FindOne(ctx *gin.Context, params requests.ActorsFindOneRequests) (result responses.ActorsResponses, err error) {
	var model repository.Actors
	result, err = model.FindOne(ctx, params)
	return
}

// 列表数据
func (s *ActorsService) List(ctx *gin.Context, params requests.ActorsListRequests) (result []responses.ActorsResponses, err error) {
	var model repository.Actors
	result, err = model.List(ctx, params)
	return
}

func (s *ActorsService) Update(ctx *gin.Context, params requests.ActorsUpdateRequests) (upserted_id interface{}, err error) {
	var model repository.Actors
	upserted_id, err = model.Update(ctx, params)
	return
}

func (s *ActorsService) Delete(ctx *gin.Context, params requests.ActorsDeleteRequests) (deleteResult interface{}, err error) {
	var model repository.Actors
	deleteResult, err = model.Delete(ctx, params)
	return
}
