package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"

	"github.com/gin-gonic/gin"
)

func InitSentenceRouter(Router *gin.RouterGroup) {
	SentenceRouter := Router.Group("editor").Use(middleware.CORSMiddleware(), middleware.JWTAuth())
	{
		SentenceRouter.POST("sentence/list", editorapi.SentenceList)
		SentenceRouter.POST("sentence/detail", editorapi.SentenceDetail)
		SentenceRouter.POST("sentence/update", editorapi.SentenceUpdate)
	}
}
