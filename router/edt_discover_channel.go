package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"

	// "editorApi/middleware"

	"github.com/gin-gonic/gin"
)

func InitDiscoverRouter(Router *gin.RouterGroup) {

	LangRouter := Router.Group("dis").Use(middleware.CORSMiddleware()).Use(middleware.JWTAuth())
	{
		LangRouter.GET("channel/list", editorapi.DisChannelList)
		LangRouter.POST("channel/add", editorapi.DisChannelAdd)
		LangRouter.POST("channel/edit", editorapi.DisChannelEdit)
		LangRouter.POST("channel/del", editorapi.DisChannelDel)
		LangRouter.POST("channel/listorder", editorapi.DisChannelListOrder)
	}
}
