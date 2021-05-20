package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"

	// "editorApi/middleware"

	"github.com/gin-gonic/gin"
)

func InitImageRouter(Router *gin.RouterGroup) {

	LangRouter := Router.Group("editor").Use(middleware.CORSMiddleware())
	{
		LangRouter.POST("image/search", editorapi.ImageSearch)
		LangRouter.GET("image/download", editorapi.ImageDownload)
		LangRouter.POST("image/add", editorapi.ImageAdd)
		LangRouter.POST("image/del", editorapi.ImageDel)
		LangRouter.POST("image/add/more", editorapi.ImageAddMore)
		LangRouter.POST("image/edit", editorapi.ImageEdit)
		LangRouter.POST("image/tags", editorapi.ImageTags)
		LangRouter.POST("image/tag/add", editorapi.ImageTagAdd)
		LangRouter.POST("image/tag/del", editorapi.ImageTagDel)
	}
}
