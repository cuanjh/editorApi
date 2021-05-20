package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"

	// "editorApi/middleware"

	"github.com/gin-gonic/gin"
)

func InitLangRouter(Router *gin.RouterGroup) {

	LangRouter := Router.Group("editor").Use(middleware.CORSMiddleware()).Use(middleware.JWTAuth())
	{
		LangRouter.POST("lang/list", editorapi.LangLists)
		LangRouter.POST("lang/add", editorapi.LangAdd)
		LangRouter.POST("lang/edit", editorapi.LangEdit)
		LangRouter.POST("lang/del", editorapi.LangDel)
		LangRouter.GET("info/token", editorapi.QiniuToken)
		LangRouter.GET("info/token/uploadfile", editorapi.QiniuUploadFileToken)
		LangRouter.GET("info/config", editorapi.ConfigInfo)
	}
}
