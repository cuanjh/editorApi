package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"

	// "editorApi/middleware"

	"github.com/gin-gonic/gin"
)

func InitModelRouter(Router *gin.RouterGroup) {

	LangRouter := Router.Group("editor").Use(middleware.CORSMiddleware()).Use(middleware.JWTAuth())
	{
		LangRouter.POST("model/list", editorapi.ModelLists)
		LangRouter.POST("model/add", editorapi.ModelAdd)
		LangRouter.POST("model/edit", editorapi.ModelEdit)
		LangRouter.POST("model/del", editorapi.ModelDel)
	}
}
