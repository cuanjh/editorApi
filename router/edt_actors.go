package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"
	"github.com/gin-gonic/gin"
)

func InitActorsRouter(Router *gin.RouterGroup) {
	ActorsRouter := Router.Group("editor").Use(middleware.CORSMiddleware(), middleware.JWTAuth())
	{
		ActorsRouter.POST("actors/create", editorapi.ActorsCreate)
		ActorsRouter.POST("actors/find", editorapi.ActorsFind)
		ActorsRouter.POST("actors/findone", editorapi.ActorsFindOne)
		ActorsRouter.POST("actors/list", editorapi.ActorsList)
		ActorsRouter.POST("actors/update", editorapi.ActorsUpdate)
		ActorsRouter.POST("actors/delete", editorapi.ActorsDelete)
	}
}
