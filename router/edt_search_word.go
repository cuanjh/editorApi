package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"

	"github.com/gin-gonic/gin"
)

func InitSearcWordRouter(Router *gin.RouterGroup) {

	LangRouter := Router.Group("editor").Use(middleware.CORSMiddleware())
	{
		LangRouter.POST("sword/list", editorapi.SearchWordList)
		LangRouter.POST("sword/add", editorapi.SearchWordAdd)
		LangRouter.POST("sword/edit", editorapi.SearchWordEdit)
		LangRouter.POST("sword/listorder", editorapi.SearchWordListOrder)
		LangRouter.POST("sword/del", editorapi.SearchWordDel)
	}
}
