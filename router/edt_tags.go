package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"

	"github.com/gin-gonic/gin"
)

func InitTagRouter(Router *gin.RouterGroup) {

	LangRouter := Router.Group("editor").Use(middleware.CORSMiddleware()).Use(middleware.JWTAuth())
	{
		LangRouter.POST("tag/list", editorapi.TagsLists)
		LangRouter.GET("tag/types", editorapi.TagTypes)
		LangRouter.POST("tag/add", editorapi.TagAdd)
		LangRouter.POST("tag/edit", editorapi.TagEdit)
		LangRouter.POST("tag/del", editorapi.TagsDel)
		LangRouter.POST("tag/types/add", editorapi.TagTypesAdd)

	}
}
