package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"

	// "editorApi/middleware"

	"github.com/gin-gonic/gin"
)

func InitContentTypeRouter(Router *gin.RouterGroup) {

	LangRouter := Router.Group("editor").Use(middleware.CORSMiddleware()).Use(middleware.JWTAuth())
	{
		LangRouter.POST("contenttype/list", editorapi.ContentTypeList)
		LangRouter.POST("contenttype/listorder", editorapi.ContentTypeListOrder)
		LangRouter.POST("contenttype/add", editorapi.ContentTypeAdd)
		LangRouter.POST("contenttype/edit", editorapi.ContentTypeEdit)
		LangRouter.POST("contenttype/del", editorapi.ContentTypeDel)

	}
}
