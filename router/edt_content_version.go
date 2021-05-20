package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"

	// "editorApi/middleware"

	"github.com/gin-gonic/gin"
)

func InitContentVersionRouter(Router *gin.RouterGroup) {

	LangRouter := Router.Group("editor").Use(middleware.CORSMiddleware()).Use(middleware.JWTAuth())
	{
		LangRouter.POST("content/version/list", editorapi.ContentVersionList)
		LangRouter.POST("content/version/add", editorapi.ContentVersionAdd)
		LangRouter.POST("content/version/edit", editorapi.ContentVersionEdit)
		LangRouter.POST("content/version/del", editorapi.ContentVersionDel)
	}
}
