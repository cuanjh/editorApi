package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"

	// "editorApi/middleware"

	"github.com/gin-gonic/gin"
)

func InitContentRouter(Router *gin.RouterGroup) {
	LangRouter := Router.Group("editor").Use(middleware.CORSMiddleware()).Use(middleware.JWTAuth())
	{
		LangRouter.POST("content", editorapi.Content)
		LangRouter.GET("content/showTypes", editorapi.ContentTypes)
		LangRouter.POST("content/edit", editorapi.ContentEdit)
		LangRouter.POST("content/del", editorapi.ContentDel)
		LangRouter.POST("content/search", editorapi.ContentSearch)
		LangRouter.POST("content/import", editorapi.ContentImport)
		LangRouter.POST("content/export", editorapi.ContentExport)
		LangRouter.POST("content/export_list", editorapi.ContentExportList)
	}
}
