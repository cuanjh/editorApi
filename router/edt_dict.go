package router

import (
	"editorApi/controller/editorapi"
	"editorApi/middleware"

	"github.com/gin-gonic/gin"
)

func InitDiceRouter(Router *gin.RouterGroup) {
	DictRouter := Router.Group("editor").Use(middleware.CORSMiddleware()).Use(middleware.JWTAuth())
	{
		DictRouter.POST("dict/upload", editorapi.DictUpload)
		DictRouter.POST("dict/collect", editorapi.DictCollect)
		DictRouter.POST("dict/list", editorapi.DictList)
		DictRouter.POST("dict/detail", editorapi.DictDetail)
		DictRouter.POST("dict/update", editorapi.DictUpdate)
		DictRouter.POST("dict/tags", editorapi.DictTags)
		DictRouter.POST("dict/del", editorapi.DictDel)
		DictRouter.POST("dict/online", editorapi.DictOnline)
		DictRouter.POST("dict/offline", editorapi.DictOffline)
	}
}
