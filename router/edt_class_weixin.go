package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"
	"github.com/gin-gonic/gin"
)

func InitClassWeixinRouter(Router *gin.RouterGroup) {

	LangRouter := Router.Group("editor").Use(middleware.CORSMiddleware())
	{
		LangRouter.POST("class_weixin/list", editorapi.ClassWeixinList)
		LangRouter.POST("class_weixin/add", editorapi.ClassWeixinAdd)
		LangRouter.PUT("class_weixin/edit", editorapi.ClassWeixinEdit)
		LangRouter.DELETE("class_weixin/del", editorapi.ClassWeixinDel)
	}
}