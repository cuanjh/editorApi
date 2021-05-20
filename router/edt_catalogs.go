package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"

	// "editorApi/middleware"

	"github.com/gin-gonic/gin"
)

func InitCatalogRouter(Router *gin.RouterGroup) {

	// LangRouter := Router.Group("editor").Use(middleware.CORSMiddleware()).Use(middleware.JWTAuth())
	LangRouter := Router.Group("editor").Use(middleware.CORSMiddleware())
	{
		LangRouter.POST("catalog/list", editorapi.CatalogList)
		LangRouter.POST("catalog/add", editorapi.CatalogAdd)
		LangRouter.POST("catalog/edit", editorapi.CatalogEdit)
		LangRouter.POST("catalog/del", editorapi.CatalogDel)
		LangRouter.POST("catalog/rename", editorapi.CatalogRename)
		LangRouter.POST("catalog/show", editorapi.CatalogShow)
		// LangRouter.GET("catalog/del/pro", editorapi.CatalogDelPro)
		LangRouter.POST("catalog/copy", editorapi.CatalogCopy)
		LangRouter.POST("catalog/move", editorapi.CatalogMove)
		LangRouter.POST("examin", editorapi.Examin)
		LangRouter.POST("examin/submit", editorapi.ExaminSubmit)
		LangRouter.POST("examin/reset", editorapi.ExaminReset)
	}
}
